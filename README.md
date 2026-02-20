# BIOCAD Solution

Go-сервис для фоновой обработки `.tsv` файлов.
Сервис отслеживает директорию, отправляет новые файлы в очередь, обрабатывает их через **processor**, сохраняет данные в PostgreSQL, логирует и фиксирует ошибки, генерирует `.docx` по `unit_guid` и предоставляет HTTP API для получения обработанных данных.

Проект построен как асинхронный pipeline с использованием очереди сообщений.

---

## Технологии

* **Go**
* **PostgreSQL**
* **RabbitMQ**
* REST API
* Асинхронная обработка (producer → processor)
* Генерация DOCX

---

## Зависимости

* `github.com/caarlos0/env/v10` — загрузка конфигурации из env
* `github.com/joho/godotenv` — чтение `.env` файла
* `github.com/jackc/pgx/v5` — драйвер PostgreSQL
* `github.com/rabbitmq/amqp091-go` — клиент RabbitMQ
* `github.com/gorilla/mux` — HTTP роутер
* `github.com/mmonterroca/docxgo/v2` — генерация `.docx`

Установка зависимостей:

```bash
cd directory-viewing-service/
go mod tidy
```

---

## Как работает

* Сервис сканирует `WORK_DIR`
* Новый `.tsv` отправляется в RabbitMQ
* **processor**:

    * парсит файл
    * сохраняет данные в PostgreSQL
    * логирует ошибки
    * сохраняет ошибки в файл
    * генерирует `<unit_guid>.docx`
* Результат сохраняется в `OUT_DIR`

Ошибки парсинга:

* логируются
* сохраняются в БД
* записываются в выходной файл `logs.log`

---

## Быстрый старт

1. Скопировать конфиг:

```bash
cp .env.example .env
```

2. Указать свои настройки в `.env`:

* Пути не должны содержать слеша в конце: `/`
* Путь до директории должен существовать. Не успел залоггировать/вынести в отдельную функцию создание папки по несуществующему пути.
* В случае, если сервис не обрабатывает директорию, то нужно проверить путь до директории (o/i).

```env
WATCH_INTERVAL=3m
WORK_DIR=путь_к_tsv
OUT_DIR=путь_для_docx

DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=processing_files_db
DB_PORT=5432
DB_SSLMODE=disable

RABBITMQ_HOST=localhost
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_PORT=5672
```

3. Запустить сервис:

```bash
docker-compose up -d

cd directory-viewing-service/cmd/
go run .
```

Сервис начнёт отслеживать директорию и принимать HTTP-запросы.

---

## API

```
GET /{uid}?limit=N
```

Пример:

```
GET http://localhost:8080/123e4567?limit=10
```

* `uid` — `unit_guid`
* `limit` — необязательный параметр (количество записей)
* Если `limit` не указан, то возвращаются все записи по `uid`

Если `limit` не число — `400 Bad Request`.

---

