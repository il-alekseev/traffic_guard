# Описание компонента “Системные метрики хоста”

## Общее описание

Системные метрики хоста состоят из двух компонентов:

- **Prometheus**- система мониторинга и оповещения с открытым исходным кодом, которая в реальном времени собирает и анализирует метрики работы приложений и серверов. Инструмент написан на Golang.
    
- **Node Exporter** - это инструмент для сбора метрик с хостов(серверов) Linux и предоставления их системе мониторинга Prometheus. Инструмент написан на Golang.
    
## Взаимосвязь с другими компонетами


![](Components/Prometheus/img/1.png)  
  

В рамках системы системные метрики хоста выступают как компоненты собирающие метрики хоста, на котором развернута система а также хоста на котором развернута LLM. Визуализирования собранных метрик производится в Grafana (компонент системы), которая позволяет создать дашборды.


# Описание настройки

Инструмент будет развернут в системе в виде Docker- контейнера в системе управления контейнерами Docker Compose.

## Настройка контейнера в системе

### Node Exporter


Перед запуском контейнера необходимо указать переменные окружения в .env файле.


`./deploy/.env` :

```
...
NODE_TAG=v1.9.1

NODE_PORT=9100

NODE_VOLUME_PORT=9100
...
```
  

`docker-compose.yml`:

```yaml
  node_exporter:
    image: prom/node-exporter:${NODE_TAG}
    container_name: node_exporter
    ports:
      - ${NODE_VOLUME_PORT}:${NODE_PORT}
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    restart: always
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:${NODE_PORT}"]
      interval: 10s
      timeout: 5s
      retries: 6
      start_period: 10s
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
    networks:
      - continent
```
  

### Prometheus


Перед запуском контейнера необходимо указать переменные окружения в .env файле.

`./deploy/.env` :

```
PROMETHEUS_TAG=v3.4.1

PROMETHEUS_PORT=9090

PROMETHEUS_VOLUME_PORT=9090

PROMETHEUS_PROMETHEUS_TARGETS=prometheus:9090

PROMETHEUS_NODEEXPORTER_TARGETS=node_exporter:9100
```  

`docker-compose.yml`:

```yaml
  prometheus:
    image: prom/prometheus:${PROMETHEUS_TAG}
    container_name: prometheus
    ports:
      - ${PROMETHEUS_VOLUME_PORT}:${PROMETHEUS_PORT}
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./prometheus/data:/prometheus
    environment:
      - PROMETHEUS_PROMETHEUS_TARGETS=${PROMETHEUS_PROMETHEUS_TARGETS}
      - PROMETHEUS_NODEEXPORTER_TARGETS=${PROMETHEUS_NODEEXPORTER_TARGETS}
    restart: always
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:${PROMETHEUS_PORT}"]
      interval: 10s
      timeout: 5s
      retries: 6
      start_period: 10s
    depends_on:
      node_exporter:
        condition: service_healthy
    networks:
      - continent
```  

Также помимо этого необходимо настроить конфигурационный файл.

`./deploy/prometheus/prometheus.yml`:

  
```yaml
global:
  scrape_interval: 5s
  evaluation_interval: 5s
scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ['prometheus:9090']
  - job_name: "node"
    static_configs:
      - targets: ['node_exporter:9100']
  - job_name: "node_ml"
    static_configs:
      - targets: ['XX.XX.XX.XX:9100']
```