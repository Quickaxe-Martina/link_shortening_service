#!/bin/bash

# Кол-во запросов
REQUESTS=1000
# Кол-во параллельных воркеров
CONCURRENCY=10
# URL сервиса
BASE_URL="http://localhost:8080/api/shorten"

# Генерация случайной строки
random_string() {
  head /dev/urandom | tr -dc a-z0-9 | head -c 6
}

# Создаём временный файл с телами запросов
TMPFILE=$(mktemp)
for i in $(seq 1 $REQUESTS); do
  RAND_STR=$(random_string)
  echo "{\"URL\": \"https://hey-example-$RAND_STR.com\"}" >> $TMPFILE
done

# Запускаем нагрузку через hey
./hey_linux_amd64 -n $REQUESTS -c $CONCURRENCY -m POST -T "application/json" -D $TMPFILE $BASE_URL

# Удаляем временный файл
rm $TMPFILE
