import os
import json
import requests
import tiktoken
import re
from dotenv import load_dotenv
from confluent_kafka import Consumer, Producer, KafkaError
import json
import psycopg2
from psycopg2 import sql, OperationalError
from ollama import Client


def truncate_text_for_llm(text: str, max_tokens: int = 3500) -> str:
    """
    Обрезает текст до заданного количества токенов, используя токенизацию OpenAI (cl100k_base).

    Функция используется для подготовки текста к передаче в модели Ollama или другие
    языковые модели, которые имеют ограничения на количество токенов во входных данных.
    Сначала текст кодируется в токены, затем берутся только первые max_tokens токенов,
    и они декодируются обратно в строку.

    Параметры:
    ----------
    text : str
        Исходный текст, который нужно обрезать.
    max_tokens : int
        Максимальное количество токенов, которое можно оставить в тексте.
        По умолчанию — 3500.

    Возвращает:
    -----------
    str
        Обрезанный текст, не превышающий заданное количество токенов.
    """

    # Получаем энкодер для модели, совместимой с cl100k_base
    enc = tiktoken.get_encoding("cl100k_base")

    # Преобразуем текст в список токенов (числовых идентификаторов)
    tokens = enc.encode(text)

    # Обрезаем список токенов до заданного количества
    truncated_tokens = tokens[:max_tokens]

    # Декодируем токены обратно в текст
    truncated_text = enc.decode(truncated_tokens)

    return truncated_text


def get_text_by_content_guid(config: dict, content_guid: str) -> str | None:
    """
    Извлекает текстовое содержимое из таблицы PostgreSQL по значению CONTENT_GUID.

    Функция подключается к базе данных PostgreSQL, выполняет безопасный параметризованный
    SQL-запрос и возвращает содержимое текстового поля TEXT из таблицы ml_table
    для указанного идентификатора `CONTENT_GUID`.

    Параметры:
    ----------
    config : dict
        Словарь с параметрами подключения к базе данных. Пример структуры:
        {
            "database": "db_name",
            "user": "username",
            "password": "password",
            "host": "localhost",
            "port": "5432"
        }
    content_guid : str
        Уникальный идентификатор контента (значение поля CONTENT_GUID в таблице).

    Возвращает:
    -----------
    str | None
        Возвращает значение поля `TEXT`, если запись найдена.
        Возвращает `None`, если запись отсутствует или произошла ошибка при подключении.

    """

    connection = None
    cursor = None

    try:
        # Подключаемся к базе PostgreSQL, используя параметры из config
        connection = psycopg2.connect(
            dbname=config['database'],
            user=config['user'],
            password=config['password'],
            host=config['host'],
            port=config['port']
        )

        # Создаём курсор для выполнения SQL-запросов
        cursor = connection.cursor()

        # Формируем безопасный параметризованный SQL-запрос
        query = sql.SQL("SELECT TEXT FROM ml_table WHERE CONTENT_GUID = %s")

        # Выполняем запрос, передавая параметр content_guid
        cursor.execute(query, (content_guid,))

        # Извлекаем первую (и единственную) строку результата
        result = cursor.fetchone()

        if result:
            return result[0]  # Возвращаем содержимое поля TEXT
        else:
            return None       # Если запись не найдена — возвращаем None

    except OperationalError as e:
        # Обрабатываем ошибки соединения с PostgreSQL
        print(f"Ошибка при работе с PostgreSQL: {e}")
        return None

    finally:
        # Закрываем курсор и соединение, даже при ошибке
        if cursor:
            cursor.close()
        if connection:
            connection.close()


def get_llama_cpp_response(model_name: str, full_prompt: str, server_url: str) -> str:
    """
    Отправляет запрос к серверу LLaMA.cpp и возвращает распознанную категорию текста.

    Функция взаимодействует с сервером LLaMA.cpp через REST API.  
    Она формирует POST-запрос на эндпоинт /completion, передаёт промпт и параметры генерации,
    а затем пытается интерпретировать ответ модели как JSON, чтобы извлечь поле category.

    Параметры:
    ----------
    model_name : str
        Имя модели, загруженной на сервере LLaMA.cpp (например, "llama-3-8b", "mistral-7b").
    full_prompt : str
        Полный текст запроса, который будет передан в модель (включая системный контекст и пользовательский ввод).
    server_url : str
        URL-адрес сервера LLaMA.cpp, например: "http://localhost:8080".

    Возвращает:
    -----------
    str
        Значение поля "category" из JSON-ответа модели.
        Если ответ невозможно разобрать как JSON — возвращает строку с текстом ошибки.

    """

    # Формируем тело запроса для API LLaMA.cpp
    payload = {
        "model": model_name,       # название модели
        "prompt": full_prompt,     # передаваемый текст (инструкция + контент)
        # параметр случайности (чем ниже, тем стабильнее)
        "temperature": 0.62,
        "max_tokens": 85,          # ограничение длины генерируемого ответа
        "stop": ["\n\n"],          # сигнал для остановки генерации
    }

    # Отправляем POST-запрос на сервер LLaMA.cpp
    response = requests.post(f"{server_url}/completion", json=payload)

    # Проверяем успешность ответа (выбросит исключение при HTTP-ошибке)
    response.raise_for_status()

    # Извлекаем сгенерированный текст из JSON-ответа
    result_text = response.json().get("content", "").strip()

    # Пробуем интерпретировать ответ как JSON, чтобы получить категорию
    try:
        result_json = json.loads(result_text)
        return result_json.get("category")

    # Если модель вернула невалидный JSON — возвращаем необработанный текст
    except json.JSONDecodeError:
        return "Неизвестный класс"


def is_recognizable_category(category: str, categories_file: str = "categories.txt") -> str:
    """
    Модель может галлюцинировать и выдавать категории, которых нет в списке рассматриваемых. 
    Функция проверяет, является ли переданная категория одной из распознаваемых.

    Аргументы:
    category (str): категория, которую нужно проверить.
    categories_file (str): путь к файлу с распознаваемыми категориями (по умолчанию "categories.txt").

    Возвращает:
    str: Возвращает входящую категорию, если категория есть в списке, иначе Положительная категория.
    """
    try:
        # Открываем файл с категориями и считываем их в список, убирая лишние пробелы и символы переноса строки
        with open(categories_file, "r", encoding="utf-8") as f:
            recognizable_categories = [line.strip()
                                       for line in f if line.strip()]

        # Проверяем наличие категории в списке
        if category in recognizable_categories:
            return category
        else:
            return "Положительная категория"
    except FileNotFoundError:
        print(f"Файл {categories_file} не найден.")
        return False
    except Exception as e:
        print(f"Произошла ошибка при проверке категории: {e}")
        return False


def get_ollama_response(model_name: str, prompt: str, server_url: str) -> str:
    """
    Отправляет запрос к локальной или удалённой Ollama и возвращает категорию,
    полученную в ответе модели.

    Функция формирует запрос к языковой модели через Ollama API, получает текстовый ответ
    и пытается интерпретировать его как JSON. Если модель возвращает корректный JSON,
    из него извлекается значение по ключу "category". В противном случае возвращается
    строка "Неизвестный класс".

    Параметры:
    ----------
    model_name : str
        Имя модели, которую необходимо вызвать (например, "llama3", "mistral", "gpt-4").
    prompt : str
        Текст запроса (вопрос или инструкция), передаваемый модели.
    server_url : str
        URL сервера Ollama, к которому подключается клиент (например, "http://localhost:11434").

    Возвращает:
    -----------
    str
        Значение поля "category" из ответа модели, если удалось успешно разобрать JSON.
        Если модель не вернула корректный JSON, возвращает строку "Неизвестный класс".
    """

    # Инициализируем клиент для взаимодействия с Ollama
    client = Client(host=server_url)

    # Отправляем запрос к модели с заданными параметрами генерации
    response = client.generate(
        model=model_name,
        prompt=prompt,
        options={
            "num_predict": 80,    # ограничиваем длину генерируемого ответа
            # делаем ответ более детерминированным (меньше случайности)
            "temperature": 0.62,
            "format": "text"      # просим модель вернуть чистый текст
        }
    )

    # Удаляем внутренние теги <think>...</think>, если они присутствуют в ответе
    llm_output = re.sub(r'<think>.*?</think>', '',
                        response['response'], flags=re.DOTALL).strip()

    try:
        # Пытаемся интерпретировать ответ как JSON
        json_llm_output = json.loads(llm_output)
        # Возвращаем значение ключа "category", если оно есть
        category = json_llm_output.get("category")
        print("Первоначальная распознанная категория: ",category)
        return is_recognizable_category(category)
    except json.JSONDecodeError:
        # Если JSON некорректный — возвращаем стандартную фразу
        return "Неизвестный класс"


def get_llm_studio_response(model_name: str, prompt: str, server_url: str) -> str:
    """
    Отправляет запрос к LLM Studio API и возвращает категорию, сгенерированную моделью.

    Функция формирует POST-запрос к серверу LLM Studio (или совместимому API),
    передавая параметры модели, промпт и ожидаемую JSON-схему ответа.  
    В ответ получает JSON-объект, из которого извлекает текстовое поле category.

    Параметры:
    ----------
    model_name : str
        Название модели, которую необходимо вызвать (например, "mistral", "llama3", "gpt-4").
    prompt : str
        Текст запроса или инструкция, передаваемая языковой модели.
    server_url : str
        URL сервера LLM Studio (например, "http://localhost:8000").

    Возвращает:
    -----------
    str
        Категория текста, возвращённая моделью в поле `text`.
        Если сервер вернул неожиданный формат данных — возбуждается исключение.

    """

    # Формируем тело POST-запроса с настройками генерации и схемой ожидаемого JSON-ответа
    payload = {
        "model": model_name,
        "prompt": prompt,
        "parameters": {
            "max_new_tokens": 85,   # ограничиваем длину ответа
            # делаем ответ детерминированным (меньше случайности)
            "temperature": 0.62
        },
        "json_schema": {  # ожидаемая структура ответа от модели
            "type": "object",
            "properties": {
                "category": {
                    "type": "string",
                    "description": "Категория текста"
                }
            },
            "required": ["category"],           # поле category обязательно
            "additionalProperties": False       # запрещаем другие поля
        }
    }

    # Отправляем POST-запрос на сервер LLM Studio
    response = requests.post(f"{server_url}/v1/generate", json=payload)

    # Проверяем, что запрос выполнен успешно
    response.raise_for_status()

    # Преобразуем JSON-ответ от сервера в Python-объект
    data = response.json()

    # Извлекаем категорию из первого результата
    return data["results"][0]["text"]


def delivery_callback(err, msg):
    """Callback для обработки результата отправки сообщения в Kafka"""
    if err:
        print(f'✗ Ошибка доставки: {err}')
    else:
        print(
            f'✓ Отправлено в {msg.topic()} [partition {msg.partition()}] @ offset {msg.offset()}')


def send_to_kafka(producer, topic, request_id, content_guid, recognised_class):
    """Отправка сообщения в Kafka"""
    message = {
        'REQUEST_ID': request_id,
        'CONTENT_GUID': content_guid,
        'RECOGNISED_CLASS': recognised_class
    }

    try:
        producer.produce(
            topic=topic,
            key=request_id.encode('utf-8'),
            value=json.dumps(message).encode('utf-8'),
            callback=delivery_callback
        )
        producer.poll(0)
    except Exception as e:
        print(f"✗ Ошибка: {e}")


def main():
    """
    Основной цикл обработки сообщений из Kafka и взаимодействия с LLM (Ollama).

    Скрипт выполняет следующие действия:
    1. Загружает конфигурацию из .env файла (модель, Kafka, база данных и т.д.).
    2. Настраивает потребителя и производителя Kafka.
    3. В бесконечном цикле читает сообщения из топика CONSUMER_TOPIC.
    4. Обрабатывает полученный контент, сокращая текст и формируя промпт.
    5. Отправляет промпт в LLM и получает распознанную категорию.
    6. Отправляет результат обратно в Kafka в PRODUCER_TOPIC.
    7. Корректно завершает работу при нажатии Ctrl + C.

    Используемые переменные окружения:
    ---------------------------------
    MODEL_NAME, LLM_URL, KAFKA_URL, CONSUMER_TOPIC, PRODUCER_TOPIC,
    KAFKA_CONSUMER_GROUP_ID, KAFKA_PRODUCER_GROUP_ID, KAFKA_TIMEOUT,
    DATABASE_NAME, DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST, DATABASE_PORT

    Пример запуска:
    ---------------
    >>> python main.py
    """

    # Загружаем системный промпт из файла (контекст для модели)
    with open("system_prompt.txt", "r", encoding="utf-8") as f:
        system_prompt = f.read()

    # Загружаем переменные окружения из файла .env
    load_dotenv()

    # Конфигурация модели и подключения
    model_name = os.getenv("MODEL_NAME")
    llm_url = os.getenv("LLM_URL")
    kafka_url = os.getenv("KAFKA_URL")
    consumer_topic = os.getenv("CONSUMER_TOPIC")
    producer_topic = os.getenv("PRODUCER_TOPIC")
    kafka_consumer_group_id = os.getenv("KAFKA_CONSUMER_GROUP_ID")
    kafka_producer_group_id = os.getenv("KAFKA_PRODUCER_GROUP_ID")
    kafka_timeout = float(os.getenv("KAFKA_TIMEOUT"))

    # Параметры базы данных (если они понадобятся другим функциям)
    database_config = {
        'user': os.getenv("DATABASE_USER"),
        'password': os.getenv("DATABASE_PASSWORD"),
        'host': os.getenv("DATABASE_HOST"),
        'database': os.getenv("DATABASE_NAME"),
        'port': os.getenv("DATABASE_PORT"),
    }

    # Конфигурация Kafka consumer (читает сообщения)
    kafka_consumer_config = {
        'bootstrap.servers': kafka_url,         # адрес брокера Kafka
        'group.id': kafka_consumer_group_id,    # идентификатор группы потребителей
        'auto.offset.reset': 'earliest',        # читаем с начала, если нет смещений
    }

    # Конфигурация Kafka producer (отправляет результаты)
    kafka_producer_config = {
        'bootstrap.servers': kafka_url,         # адрес брокера Kafka
        'client.id': kafka_producer_group_id,   # идентификатор клиента
        'acks': 'all',                          # подтверждение записи
        'retries': 3                            # попытки повторной отправки
    }

    # Создаём Kafka-потребителя и подписываемся на нужный топик
    consumer = Consumer(kafka_consumer_config)
    consumer.subscribe([consumer_topic])

    # Создаём Kafka-производителя для отправки результатов
    producer = Producer(kafka_producer_config)

    try:
        while True:
            # Ожидаем новое сообщение с таймаутом
            msg = consumer.poll(kafka_timeout)

            if msg is None:
                # Сообщений пока нет — продолжаем ожидание
                continue

            if msg.error():
                # Обработка ошибок сообщений Kafka
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    print(f"Конец партиции: {msg.topic()} [{msg.partition()}]")
                else:
                    print(f"Ошибка Kafka: {msg.error()}")
                continue

            # Декодируем JSON-сообщение из Kafka
            message = json.loads(msg.value().decode('utf-8'))

            # Извлекаем данные из сообщения
            content_id = message.get('content_id')
            request_id = message.get('request_id')
            content = message.get('content')

            print("Полученное сообщение: ", content)
            # Сокращаем текст, чтобы он поместился в лимит LLM
            resized_content = truncate_text_for_llm(content)

            # Формируем итоговый промпт для LLM
            final_prompt = f"{system_prompt}\n\n{resized_content}"

            # Получаем ответ от LLM (категорию текста)
            recognised_class = get_ollama_response(
                model_name, final_prompt, llm_url)
            print("Распознанный класс:", recognised_class)

            # Отправляем результат обратно в Kafka
            send_to_kafka(producer, producer_topic, request_id,
                          content_id, recognised_class)
            print("------------------------")

    except KeyboardInterrupt:
        # Корректное завершение работы при остановке вручную
        print("Остановка потребителя Kafka")

    finally:
        # Закрываем соединение с Kafka при завершении программы
        consumer.close()
        print("Соединение с Kafka закрыто")


if __name__ == "__main__":
    main()
