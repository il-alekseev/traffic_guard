#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import json
import numpy as np
import logging
from dotenv import load_dotenv
from confluent_kafka import Consumer, Producer, KafkaError
import onnxruntime as ort
from transformers import AutoTokenizer
import psycopg2
from psycopg2 import sql, OperationalError


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
)
logger = logging.getLogger("traffic-classifier")

def softmax(x: np.ndarray, axis: int = -1) -> np.ndarray:
    x = x - np.max(x, axis=axis, keepdims=True)
    e = np.exp(x)
    return e / np.sum(e, axis=axis, keepdims=True)


def delivery_callback(err, msg):
    if err:
        logger.error("Ошибка доставки: %s", err)
    else:
        logger.info("Отправлено в %s [partition %d] @ offset %d",
                    msg.topic(), msg.partition(), msg.offset())


def insert_classified_content(conn, request_id, content_id, url, source_content, identified_class) -> bool:
    """
    Пишем в PostgreSQL так же, как в исходном скрипте:
    таблица classified_content(request_id, content_id, url, source_content, identified_class)
    """
    try:
        with conn.cursor() as cur:
            cur.execute(
                sql.SQL("""
                    INSERT INTO classified_content (
                        request_id, content_id, url, source_content, identified_class
                    )
                    VALUES (%s, %s, %s, %s, %s)
                """),
                (request_id, content_id, url, source_content, identified_class)
            )
        conn.commit()
        return True
    except psycopg2.Error as e:
        conn.rollback()
        logger.warning("Не удалось записать в БД: request_id=%s", request_id, exc_info=True)
        return False


class OnnxUnifiedClassifier:
    """
    Content (11):
      0 Алкоголь
      1 Интернет магазины
      2 Компьютерные игры
      3 Криптомайнинг
      4 Порнография и секс
      5 Прокси и анонимайзеры
      6 Разрешенный ресурс
      7 Табак
      8 Фильмы и видео онлайн
      9 Азартные игры
      10 Наркотики

    Sentiment (4):
      0 normal
      1 suicide
      2 depression
      3 toxic

    Unified (13):
      0  Алкоголь и табак    = pc0 + pc7
      1  Интернет магазины   = pc1
      2  Компьютерные игры   = pc2
      3  Криптомайнинг       = pc3
      4  Порнография и секс  = pc4
      5  Прокси и анонимайзеры = pc5
      6  Фильмы и видео онлайн = pc8
      7  Азартные игры       = pc9
      8  Наркотики           = pc10
      9  Суицид              = ps1
      10 Депрессия           = ps2
      11 Агрессия            = ps3
      12 Разрешенный ресурс  = pc6 + ps0
    """
    def __init__(self, content_onnx: str, sent_onnx: str, tokenizer_dir: str, provider: str = "auto"):
        self.labels = [
            "Алкоголь и табак",
            "Интернет магазины",
            "Компьютерные игры",
            "Криптомайнинг",
            "Порнография и секс",
            "Прокси и анонимайзеры",
            "Фильмы и видео онлайн",
            "Азартные игры",
            "Наркотики",
            "Суицид",
            "Депрессия",
            "Агрессия",
            "Разрешенный ресурс",
        ]

        self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_dir, use_fast=True)

        providers = self._pick_providers(provider)
        self.sess_content = ort.InferenceSession(content_onnx, providers=providers)
        self.sess_sent = ort.InferenceSession(sent_onnx, providers=providers)

        # какие входы реально ждут модели
        self.content_inputs = {i.name for i in self.sess_content.get_inputs()}
        self.sent_inputs = {i.name for i in self.sess_sent.get_inputs()}

    @staticmethod
    def _pick_providers(provider: str):
        provider = provider.lower().strip()
        avail = ort.get_available_providers()

        if provider == "auto":
            if "CUDAExecutionProvider" in avail:
                return ["CUDAExecutionProvider", "CPUExecutionProvider"]
            return ["CPUExecutionProvider"]
        if provider in ("cuda", "gpu"):
            if "CUDAExecutionProvider" not in avail:
                raise RuntimeError(f"CUDAExecutionProvider недоступен. Доступно: {avail}")
            return ["CUDAExecutionProvider", "CPUExecutionProvider"]
        if provider == "cpu":
            return ["CPUExecutionProvider"]
        raise ValueError("ORT_PROVIDER must be: auto|cpu|cuda")

    def _encode(self, text: str, max_length: int):
        enc = self.tokenizer(
            text,
            truncation=True,
            max_length=max_length,
            padding=True,
            return_tensors="np",
        )
        return {k: v.astype(np.int64) for k, v in enc.items()}

    @staticmethod
    def _filter_inputs(encoded: dict, needed: set):
        ort_inputs = {}
        if "input_ids" in needed:
            ort_inputs["input_ids"] = encoded["input_ids"]
        if "attention_mask" in needed:
            ort_inputs["attention_mask"] = encoded["attention_mask"]
        if "token_type_ids" in needed:
            ort_inputs["token_type_ids"] = encoded.get(
                "token_type_ids", np.zeros_like(encoded["input_ids"], dtype=np.int64)
            )
        return ort_inputs

    def predict(self, text: str, max_length: int = 256, top_k: int = 5):
        enc = self._encode(text, max_length=max_length)

        in_content = self._filter_inputs(enc, self.content_inputs)
        in_sent = self._filter_inputs(enc, self.sent_inputs)

        logits_c = self.sess_content.run(None, in_content)[0]  # (1,11)
        logits_s = self.sess_sent.run(None, in_sent)[0]        # (1,4)

        pc = softmax(logits_c, axis=-1)[0]  # (11,)
        ps = softmax(logits_s, axis=-1)[0]  # (4,)

        v = np.array([
            (pc[0] + pc[7])*0.7,   # Алкоголь + Табак
            pc[1],
            pc[2],
            pc[3],
            pc[4],
            pc[5],
            pc[8],
            pc[9],
            pc[10],
            ps[1],           # suicide -> Суицид
            ps[2],           # depression -> Депрессия
            ps[3],           # toxic -> Агрессия
            (pc[6] + ps[0])*0.5,   # Разрешенный ресурс + normal
        ], dtype=np.float32)

        # нормализация до суммы 1
        s = float(v.sum())
        if s > 0:
            v = v / s

        k = min(top_k, len(self.labels))
        idxs = np.argsort(-v)[:k]
        top = [{"label": self.labels[int(i)], "score": float(v[int(i)])} for i in idxs]

        return {"vector": v.tolist(), "labels": self.labels, "top": top}


def main():
    load_dotenv()

    # Kafka
    kafka_url = os.getenv("KAFKA_URL")
    consumer_topic = os.getenv("CONSUMER_TOPIC")
    producer_topic = os.getenv("PRODUCER_TOPIC")
    kafka_consumer_group_id = os.getenv("KAFKA_CONSUMER_GROUP_ID", "onnx_consumer_group")
    kafka_producer_group_id = os.getenv("KAFKA_PRODUCER_GROUP_ID", "onnx_producer")
    kafka_timeout = float(os.getenv("KAFKA_TIMEOUT", "1.0"))

    if not kafka_url or not consumer_topic or not producer_topic:
        raise RuntimeError("Нужны KAFKA_URL, CONSUMER_TOPIC, PRODUCER_TOPIC в окружении/.env")

    # ONNX
    content_onnx = os.getenv("CONTENT_ONNX_PATH", "content_classifier.onnx")
    sent_onnx = os.getenv("SENTIMENT_ONNX_PATH", "sentiment_classifier.onnx")
    tokenizer_dir = os.getenv("TOKENIZER_DIR", "models/content")
    provider = os.getenv("ORT_PROVIDER", "auto")
    max_length = int(os.getenv("MAX_LENGTH", "256"))
    top_k = int(os.getenv("TOP_K", "5"))

    db_config = {
        "dbname": os.getenv("DATABASE_NAME"),
        "user": os.getenv("DATABASE_USER"),
        "password": os.getenv("DATABASE_PASSWORD"),
        "host": os.getenv("DATABASE_HOST"),
        "port": os.getenv("DATABASE_PORT"),
    }
    if not all(db_config.values()):
        raise RuntimeError("Нужны DATABASE_NAME/USER/PASSWORD/HOST/PORT в окружении/.env")

    try:
        db_conn = psycopg2.connect(**db_config)
        logger.info("PostgreSQL connected")
    except OperationalError as e:
        logger.error("Не удалось подключиться к БД: %s", e)
        raise

    consumer = Consumer({
        "bootstrap.servers": kafka_url,
        "group.id": kafka_consumer_group_id,
        "auto.offset.reset": "earliest",
        "enable.auto.commit": False,  # коммитим вручную после успеха
    })
    consumer.subscribe([consumer_topic])

    producer = Producer({
        "bootstrap.servers": kafka_url,
        "client.id": kafka_producer_group_id,
        "acks": "all",
        "retries": 3,
    })

    clf = OnnxUnifiedClassifier(content_onnx, sent_onnx, tokenizer_dir, provider=provider)

    logger.info("ONNX Kafka+PostgreSQL worker started.")
    try:
        while True:
            msg = consumer.poll(kafka_timeout)

            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    logger.info("Конец партиции: %s [%d]", msg.topic(), msg.partition())
                else:
                    logger.error("Ошибка Kafka: %s", msg.error())
                continue

            try:
                message = json.loads(msg.value().decode("utf-8"))

                # совместимость со старым и новым форматом
                content_id = message.get("content_id") or message.get("CONTENT_GUID") or message.get("content_guid")
                request_id = message.get("request_id") or message.get("REQUEST_ID") or message.get("request_guid")
                content = message.get("content") or message.get("text") or message.get("CONTENT") or message.get("TEXT")
                url = message.get("url") or message.get("URL", "")

                if not content_id or not request_id or not content:
                    logger.warning("Пропуск сообщения: нет content_id/request_id/content")
                    consumer.commit(message=msg, asynchronous=False)
                    continue

                pred = clf.predict(content, max_length=max_length, top_k=top_k)

                top1_label = pred["top"][0]["label"]
                top1_score = pred["top"][0]["score"]

                # 1) Запись в PostgreSQL (как раньше)
                db_ok = insert_classified_content(
                    conn=db_conn,
                    request_id=str(request_id),
                    content_id=str(content_id),
                    url=str(url),
                    source_content=str(content),
                    identified_class=str(top1_label),
                )
                if not db_ok:
                    # если БД критична — можно raise и НЕ коммитить offset
                    logger.warning("Не удалось записать в БД: request_id=%s", request_id)

                # 2) Отправка в Kafka результата (Top-5 + вектор)
                out_message = {
                    "REQUEST_ID": str(request_id),
                    "CONTENT_GUID": str(content_id),

                    # обратно-совместимые поля
                    "RECOGNISED_CLASS": top1_label,
                    "RECOGNISED_SCORE": round(float(top1_score), 6),

                    # новые поля
                    "TOP5": pred["top"],
                    "VECTOR13": pred["vector"],
                    "LABELS13": pred["labels"],
                    "URL": url,
                }

                producer.produce(
                    topic=producer_topic,
                    key=str(request_id).encode("utf-8"),
                    value=json.dumps(out_message, ensure_ascii=False).encode("utf-8"),
                    callback=delivery_callback,
                )
                producer.poll(0)

                # коммитим offset только после успешной обработки
                consumer.commit(message=msg, asynchronous=False)

                logger.info("Request-id: %s, top1: %s (%.4f)", request_id, top1_label, top1_score)

            except Exception as e:
                logger.error("Ошибка обработки сообщения: %s", e, exc_info=True)
                # consumer.commit(message=msg, asynchronous=False)

    except KeyboardInterrupt:
        logger.info("Остановка потребителя Kafka")

    finally:
        try:
            consumer.close()
        except Exception:
            pass
        try:
            producer.flush()
        except Exception:
            pass
        try:
            db_conn.close()
        except Exception:
            pass
        logger.info("Завершение. Kafka/DB закрыты.")


if __name__ == "__main__":
    main()
