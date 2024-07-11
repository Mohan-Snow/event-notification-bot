#!/bin/bash

# Загружаем переменные из .env файла
set -o allexport
source ../.env
set -o allexport

# Заменяем переменные окружения в строке подключения
DB_CONNECTION_STRING="postgres://${DATA_SOURCE_USERNAME}:${DATA_SOURCE_PASSWORD}@${DATA_SOURCE_HOST}:${DATA_SOURCE_PORT}/${DATA_SOURCE_NAME}?sslmode=disable"

# Выполняем команды goose
goose -dir ../migration postgres "$DB_CONNECTION_STRING" status
goose -dir ../migration postgres "$DB_CONNECTION_STRING" up
