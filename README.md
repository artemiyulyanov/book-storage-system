# book-storage-system

Пет-проект на Go, демонстрирующий переход от монолитного CRUD-сервиса к микросервисной архитектуре: разделение на независимые сервисы, единая точка входа через API Gateway, межсервисное взаимодействие по gRPC, асинхронная событийная шина на Kafka и авторизация через JWT.

## Содержание

- [Архитектура](#архитектура)
- [Сервисы](#сервисы)
- [Стек технологий](#стек-технологий)
- [Структура репозитория](#структура-репозитория)
- [Как это работает](#как-это-работает)
- [Запуск проекта](#запуск-проекта)
- [Работа с миграциями](#работа-с-миграциями)
- [API](#api)
- [Планы на будущее](#планы-на-будущее)

## Архитектура

### Синхронный слой (запросы клиента)

```
                     ┌──────────────┐
                     │    Клиент    │
                     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  api-gateway │   единственная точка входа наружу
                     │  (роутинг +  │
                     │ JWT-проверка)│
                     └──┬────────┬──┘
                        │        │
              ┌─────────▼──┐  ┌──▼──────────┐
              │book-service│  │ auth-service│
              └─────┬──────┘  └──────┬──────┘
                    │                │ gRPC
              ┌─────▼──────┐  ┌──────▼──────┐
              │  books_db  │  │ user-service│
              │ (Postgres) │  └──────┬──────┘
              └────────────┘         │
                                ┌─────▼──────┐
                                │  users_db  │
                                │ (Postgres) │
                                └────────────┘
```

### Асинхронный слой (события → аналитика/аудит)

```
┌────────────┐
│book-service│──┐
└────────────┘  │      ┌───────┐      ┌───────────────┐      ┌────────────┐
                 ├─────▶│ Kafka │─────▶│ event-service │─────▶│ ClickHouse │
┌────────────┐   │      │(topics│      │ (batch insert │      │  (events)  │
│auth-service│──┘       │ по    │      │  по 100 msg   │      └────────────┘
└────────────┘          │ event)│      │  или 5 сек)   │
                         └───────┘      └───────────────┘
```

`Kafka` и `ZooKeeper` — общая инфраструктура, поднимается отдельно из `infra/` и используется несколькими сервисами (не привязана к одному владельцу данных, в отличие от `books_db`/`users_db`).

Ключевые архитектурные решения:

- **Database-per-service** — у каждого сервиса своя база данных, никаких общих таблиц и прямых обращений в чужую БД. Это обратная сторона монолита: связь между сущностями (например, книга и её автор) реализована через `author_id` на уровне приложения, а не через SQL `FOREIGN KEY`.
- **Единственная точка входа** — только `api-gateway` слушает порт наружу. `book-service`, `user-service` и `auth-service` доступны только внутри docker-сети, что и позволяет им доверять данным, проброшенным gateway'ем (например, заголовку `X-User-Id`).
- **Аутентификация на границе системы** — JWT проверяется один раз, на `api-gateway`. Внутренние сервисы не валидируют токен повторно — они получают уже проверенную личность пользователя через заголовок.
- **Синхронное межсервисное взаимодействие по gRPC** — там, где сервисам нужно напрямую спросить друг у друга данные (`auth-service` → `user-service` за хешем пароля), используется gRPC с protobuf-контрактом, а не HTTP/JSON.
- **Асинхронная событийная шина для аналитики/аудита (Kafka → ClickHouse)** — `book-service` и `auth-service` публикуют доменные события (`book.created`/`updated`/`deleted`, `user.registered`, `user.loggedIn`) в Kafka асинхронно, отдельной горутиной, не блокируя ответ клиенту. Это read-only контур: если Kafka или `event-service` недоступны, основной CRUD-флоу продолжает работать, теряются только записи в аналитике.
- **event-service изолирован от остальной системы** — у него нет HTTP-эндпоинтов и он не зарегистрирован в `api-gateway`. Единственная роль — consumer Kafka-топиков и writer в ClickHouse. Это намеренно: он не источник истины о сущностях, а журнал произошедшего.
- **Общий код — только инфраструктура** — переиспользуемая логика (JWT, хеширование паролей, HTTP-хелперы, подключение к БД и запуск миграций, типы событий и Kafka-producer) вынесена в модуль `common`. Доменные модели каждого сервиса (`Book`, `User`) остаются приватными для своего сервиса.

## Сервисы

| Сервис | Отвечает за | Хранилище | Протоколы |
|---|---|---|---|
| **api-gateway** | Единая точка входа, роутинг запросов, проверка JWT, проброс идентичности пользователя вниз | — (stateless) | HTTP (снаружи), HTTP (к book-service/user-service/auth-service) |
| **auth-service** | Логин: сверка пароля, выдача JWT | — (stateless) | HTTP (снаружи), gRPC (к user-service), Kafka (producer: `user.registered`, `user.loggedIn`) |
| **user-service** | Пользователи: регистрация, хранение профиля и хеша пароля | `users_db` (PostgreSQL) | HTTP (CRUD через gateway), gRPC-сервер (для auth-service) |
| **book-service** | Книги: CRUD, привязка книги к автору (`author_id`) | `books_db` (PostgreSQL) | HTTP (через gateway), Kafka (producer: `book.created`/`updated`/`deleted`) |
| **event-service** | Consumer доменных событий из Kafka, батчевая запись в ClickHouse для аналитики/аудита. HTTP-API не имеет | ClickHouse | Kafka (consumer, своя горутина на топик) |
| **common** | Общий Go-модуль: JWT, хеширование паролей, HTTP-хелперы, подключение к БД, миграции (Postgres + ClickHouse), типы событий и Kafka-producer | — | — |

Отдельно от сервисов — `infra`: общая инфраструктура (Kafka + ZooKeeper), которая не принадлежит ни одному сервису и поднимается независимо.

## Стек технологий

- **Язык**: Go
- **HTTP-роутинг**: [gorilla/mux](https://github.com/gorilla/mux)
- **База данных**: PostgreSQL (по инстансу на сервис)
- **Драйвер БД**: pgx / pgxpool
- **Миграции**: [golang-migrate](https://github.com/golang-migrate/migrate) — версионируемые up/down SQL-миграции, применяются отдельным раннером перед стартом сервиса (для Postgres- и ClickHouse-сервисов — разные функции-обёртки в `common/database`)
- **Межсервисное взаимодействие**: gRPC + Protocol Buffers
- **Событийная шина**: Apache Kafka (`confluentinc/cp-kafka`) + ZooKeeper (`confluentinc/cp-zookeeper`), Go-клиент — [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
- **Аналитическое хранилище событий**: ClickHouse (`clickhouse-go/v2`) — append-only таблица `events` (`MergeTree`, партиционирование по месяцу, TTL 180 дней)
- **Аутентификация**: JWT (HS256), хеширование паролей — bcrypt
- **Валидация**: [go-playground/validator](https://github.com/go-playground/validator)
- **Контейнеризация**: Docker, multi-stage сборка, Docker Compose (по одному compose-файлу на сервис + отдельный `infra/docker-compose.yml`)
- **Общий код**: локальный Go-модуль `common`, подключаемый через `replace` в `go.mod`

## Структура репозитория

```
book-storage-system/
├── api-gateway/
│   ├── cmd/                    # точка входа
│   ├── internal/
│   │   ├── proxy/              # роутинг и проксирование запросов
│   │   └── middleware/         # JWT-проверка
│   └── go.mod
├── auth-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── grpcclients/         # gRPC-клиент к user-service
│   │   └── handlers/            # логин, публикация user.loggedIn / user.registered
│   └── go.mod
├── book-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── database/            # repository-слой
│   │   ├── handlers/            # CRUD, публикация book.created/updated/deleted
│   │   └── models/
│   ├── migrations/
│   └── go.mod
├── user-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── database/
│   │   ├── grpcserver/          # gRPC-сервер (internal API для auth-service)
│   │   └── handlers/
│   ├── migrations/
│   └── go.mod
├── event-service/
│   ├── cmd/                     # точка входа: миграции ClickHouse + запуск consumer'ов
│   ├── internal/
│   │   ├── consumer/            # Kafka consumer, отдельная горутина и батчинг на топик
│   │   └── clickhouse/          # клиент ClickHouse, батчевая вставка в events
│   ├── migrations/
│   │   └── clickhouse/          # DDL таблицы events
│   └── go.mod
├── proto/
│   └── user/
│       └── user.proto           # контракт gRPC-сервиса пользователей
├── common/
│   ├── auth/                    # JWT, bcrypt
│   ├── database/                # подключение к БД, запуск миграций (Postgres + ClickHouse)
│   ├── events/                  # типы событий, Envelope, Kafka producer
│   ├── network/                 # HTTP-хелперы, валидация, парсинг запросов
│   ├── proto/                   # сгенерированный код из proto/
│   └── go.mod
├── infra/
│   └── docker-compose.yml       # общая инфраструктура: Kafka + ZooKeeper
└── LICENSE
```

## Как это работает

### Регистрация и логин

1. Клиент отправляет `POST /api/auth/register` (регистрация) в `auth-service` через gateway — пароль хешируется через bcrypt перед сохранением, в открытом виде нигде не хранится.
2. Для логина клиент идёт в `auth-service` (`POST /api/auth/login`).
3. `auth-service` запрашивает у `user-service` данные пользователя по email — **по gRPC**, а не по HTTP, и не напрямую в БД (`auth-service` не имеет доступа к `users_db`).
4. `auth-service` сверяет присланный пароль с полученным хешем (`bcrypt.CompareHashAndPassword`) и, если всё совпало, выпускает JWT.

### Авторизованные запросы

1. Клиент прикладывает токен в заголовке `Authorization: Bearer <token>`.
2. `api-gateway` проверяет подпись и срок действия токена. Проверка настраивается **по HTTP-методу** — например, `GET /api/books` может быть публичным, а `POST/PUT/DELETE /api/books` требовать авторизации.
3. При успешной проверке gateway извлекает `user_id` из токена и прокидывает его дальше сервису в заголовке `X-User-Id`.
4. `book-service` доверяет этому заголовку — он физически недоступен извне и получает трафик только от gateway, — и использует его, например, чтобы не дать пользователю редактировать чужие книги (`UPDATE ... WHERE id = $1 AND author_id = $2`).

### Проксирование в gateway

`api-gateway` не содержит бизнес-логики — только маршрутизацию. Запрос `GET /api/books/42` попадает на прокси, зарегистрированный под префиксом `/api/books`, префикс срезается (`http.StripPrefix`), и `book-service` получает уже чистый путь `/42`.

### Событийный поток: Kafka → event-service → ClickHouse

1. Когда `book-service` создаёт/обновляет/удаляет книгу или `auth-service` регистрирует/логинит пользователя, обработчик асинхронно (в отдельной горутине, не блокируя HTTP-ответ клиенту) публикует событие в Kafka через общий `events.Producer` из `common`. Каждый тип события — свой топик (`book.created`, `book.updated`, `book.deleted`, `user.registered`, `user.loggedIn`).
2. `event-service` не имеет HTTP-API и не зарегистрирован в `api-gateway` — это чистый consumer. На старте он прогоняет ClickHouse-миграции, затем поднимает отдельную горутину-читателя на каждый топик.
3. Каждый consumer копит входящие события в батч и сбрасывает его в ClickHouse либо когда батч набрал 100 сообщений, либо раз в 5 секунд по таймеру — в зависимости от того, что наступит раньше.
4. Offset в Kafka коммитится только после успешной вставки батча в ClickHouse. Если вставка не удалась, коммит не происходит, и при следующем запуске consumer перечитает те же сообщения.
5. Таблица `events` в ClickHouse — append-only журнал (`MergeTree`, партиционирование по месяцу, TTL 180 дней), полностью отдельный от «боевых» Postgres-баз `books_db`/`users_db`. Он не участвует в бизнес-логике сервисов — это источник для аналитики и аудита, а не текущее состояние сущностей.

> В `common/events` уже описаны типы `user.updated` и `user.deleted`, но пока ни один сервис их не публикует — это задел на будущее, в том числе под доработку профиля пользователя.

## Запуск проекта

Все сервисы общаются через единую внешнюю docker-сеть `book-storage-system-backend` — её нужно создать один раз перед первым запуском:

```bash
docker network create book-storage-system-backend
```

Каждый сервис (кроме `api-gateway` и `auth-service`, которые самостоятельны и без БД) поднимается собственным `docker-compose.yml`, лежащим внутри его директории, с контекстом сборки, поднятым до корня репозитория — это нужно, чтобы Docker видел соседний модуль `common` при сборке. `infra/docker-compose.yml` — общая инфраструктура (Kafka + ZooKeeper), не привязанная к конкретному сервису: её нужно поднять до `book-service`, `auth-service` и `event-service`, так как они читают/пишут в Kafka.

В каждой директории сервиса лежит `.env.template` с перечнем нужных переменных — скопируйте его в `.env` и заполните перед запуском.

```bash
# 1. Настройте .env с JWT_SECRET (в корне каждого сервиса, где он требуется)
echo "JWT_SECRET=$(openssl rand -base64 32)" > api-gateway/.env
echo "JWT_SECRET=$(openssl rand -base64 32)" > auth-service/.env   # тот же секрет, что и выше

# 2. Поднимите общую инфраструктуру (Kafka + ZooKeeper)
cd infra && docker compose up -d

# 3. Поднимите сервисы
cd ../book-service && docker compose up -d --build
cd ../user-service && docker compose up -d --build
cd ../event-service && docker compose up -d --build
cd ../auth-service && docker compose up -d --build
cd ../api-gateway && docker compose up -d --build
```

Требования: Docker, Docker Compose v2.20+ (используется `condition: service_completed_successfully` в `depends_on` для контейнеров с миграциями).

## Работа с миграциями

Схема каждой Postgres-БД версионируется через пары `.up.sql`/`.down.sql` в директории `migrations/` соответствующего сервиса. У `event-service` — свой каталог `migrations/clickhouse` с DDL таблицы `events` (для ClickHouse пары down-миграций не используются). Миграции применяются отдельным раннером до старта самого сервиса — сервис никогда не стартует со «старой» схемой.

Создание новой миграции:

```bash
migrate create -ext sql -dir book-service/migrations -seq add_book_status
```

## API

Все запросы идут через `api-gateway` (`http://localhost:8080`).

| Метод | Путь | Авторизация | Описание |
|---|---|---|---|
| `POST` | `/api/auth/register` | нет | Регистрация пользователя |
| `POST` | `/api/auth/login` | нет | Логин, выдаёт JWT |
| `GET` | `/api/books` | нет | Список книг |
| `GET` | `/api/books/{id}` | нет | Книга по ID |
| `POST` | `/api/books` | да | Создать книгу (автор — текущий пользователь) |
| `PUT` | `/api/books/{id}` | да | Обновить книгу (только своя) |
| `DELETE` | `/api/books/{id}` | да | Удалить книгу (только свою) |
| `GET` | `/api/users/` | да | Список пользователей |
| `GET` | `/api/users/{id}` | да | Пользователь по ID |

`event-service` наружу не смотрит: у него нет HTTP/gRPC API, он не зарегистрирован в `api-gateway`, и данные в ClickHouse сейчас можно посмотреть только напрямую, в обход системы.

## Планы на будущее

- Манифесты Kubernetes (Deployment / Service / Ingress) поверх текущей архитектуры
- Health-check эндпоинты для gRPC (`grpc_health_v1`)
- Распараллеливание агрегирующих запросов в gateway через `errgroup`
- CI-пайплайн для автоматической сборки и тестов каждого сервиса
- Публикация `user.updated` / `user.deleted` (типы уже описаны в `common/events`, продюсер пока не вызывается)
- Хотя бы read-only HTTP/gRPC API поверх ClickHouse в `event-service`, чтобы отдавать агрегированную аналитику наружу, а не только копить её

## Лицензия

[MIT](LICENSE)