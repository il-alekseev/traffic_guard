#!/bin/bash

# Простой скрипт для создания топиков через healthcheck
BOOTSTRAP_SERVER="localhost:9092"

# Функция создания топика если не существует
create_topic_if_not_exists() {
    local topic=$1
    local partitions=$2
    local retention=$3
    local cleanup_policy=$4
    
    if ! kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --list | grep -q "$topic"; then
        echo "Creating topic: $topic"
        kafka-topics --bootstrap-server $BOOTSTRAP_SERVER --create \
            --topic "$topic" \
            --partitions "$partitions" \
            --replication-factor 1 \
            --config "retention.ms=$retention" \
            --config "cleanup.policy=$cleanup_policy"
    fi
}

# Создаем топики если их нет
create_topic_if_not_exists "ml-analysis-results" 3 604800000 "delete"
create_topic_if_not_exists "url-content-analysis" 3 604800000 "delete"
create_topic_if_not_exists "url-metadata-results" 3 604800000 "delete"
create_topic_if_not_exists "url-processing-requests" 3 604800000 "delete"

exit 0