# CRM — бэкенд (Go + Gin + PostgreSQL)

Бэкенд CRM-системы для предприятия оптово-розничной торговли: приём заказов на закупку,
управление продуктами, категориями, поставщиками, счетами-фактурами и заявками на закупку.

Реализация по ТЗ ВКР (замена исходного стека PHP/Laravel/MySQL на Go/Gin/PostgreSQL).

## Технологии

- **Go 1.24** + **Gin** — HTTP-фреймворк
- **PostgreSQL 16** — СУБД
- **sqlx** + **pgx** — доступ к данным (параметризованные запросы → защита от SQL-инъекций)
- **golang-migrate** — миграции схемы (встроены в бинарник)
- **JWT** (golang-jwt) — аутентификация
- **bcrypt** — хеширование паролей
- **Docker / Docker Compose**

## Архитектура (трёхслойная)

```
cmd/api/                 — точка входа, сборка зависимостей, graceful shutdown
internal/
  config/                — конфигурация из переменных окружения
  model/                 — доменные модели, DTO, доменные ошибки
  handler/  (transport)  — Gin-обработчики, роутинг, middleware (JWT, RBAC)
  service/  (business)   — бизнес-логика, правила доступа, хеширование
  repository/(data)      — работа с PostgreSQL через sqlx
  pkg/
    hash/                — bcrypt
    jwtutil/             — выпуск/проверка JWT
migrations/              — SQL-миграции (встроены через go:embed)
```

Поток запроса: **handler → service → repository → PostgreSQL**.

## Быстрый старт (Docker)

```bash
cp .env.example .env          # при необходимости поправьте значения
docker compose up --build -d  # поднимет PostgreSQL + API
```

API будет доступен на `http://localhost:8080`. Миграции применяются автоматически
при старте, при первом запуске создаётся администратор (`admin` / `admin123` по умолчанию).

Проверка: `curl http://localhost:8080/health`

## Локальный запуск (без Docker)

Нужен запущенный PostgreSQL. Заполните `.env`, затем:

```bash
go mod tidy
go run ./cmd/api
```

## Роли и доступ

- **admin** — полный доступ ко всем сущностям и управлению пользователями.
- **user** — чтение справочников (продукты/категории/поставщики/счета), создание и
  просмотр **своих** заявок на закупку.

## Swagger / OpenAPI

Интерактивная документация доступна после запуска:

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **OpenAPI JSON:** `http://localhost:8080/swagger/doc.json`

Авторизация в UI: нажмите **Authorize** и введите `Bearer <token>` (токен из `/auth/login`).

Перегенерация документации после изменения аннотаций:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger   # или: swag init -g cmd/api/main.go -o docs --parseInternal --parseDepth 2
```

Сгенерированные файлы в `docs/` коммитятся в репозиторий (нужны для сборки образа).

## API

Базовый префикс: `/api/v1`. Защищённые маршруты требуют заголовок
`Authorization: Bearer <token>`.

### Аутентификация

| Метод | Путь | Доступ | Описание |
|-------|------|--------|----------|
| POST | `/auth/register` | публично | Регистрация (роль `user`) |
| POST | `/auth/login` | публично | Вход, возвращает JWT |
| GET | `/me` | авторизован | Текущий пользователь |

### Пользователи (только admin)

`GET/POST /users`, `GET/PUT/DELETE /users/:id`

### Справочники

Чтение — любой авторизованный; создание/изменение/удаление — только **admin**.

- `GET/POST /categories`, `GET/PUT/DELETE /categories/:id`
- `GET/POST /suppliers`, `GET/PUT/DELETE /suppliers/:id`
- `GET/POST /products`, `GET/PUT/DELETE /products/:id` (фильтр: `GET /products?category_id=1`)
- `GET/POST /invoices`, `GET/PUT/DELETE /invoices/:id`

### Заявки на закупку

| Метод | Путь | Доступ | Описание |
|-------|------|--------|----------|
| GET | `/purchase-requests` | user/admin | admin — все, user — свои |
| POST | `/purchase-requests` | user/admin | Создать заявку |
| GET | `/purchase-requests/:id` | владелец/admin | Просмотр |
| DELETE | `/purchase-requests/:id` | владелец/admin | Удалить |
| PATCH | `/purchase-requests/:id/status` | admin | Сменить статус |

Статусы заявки: `new`, `in_progress`, `approved`, `rejected`, `completed`.

### Аналитика (любой авторизованный)

Реальные агрегаты, вычисляемые в PostgreSQL.

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/analytics/summary` | Сводные KPI (кол-во сущностей, выручка) |
| GET | `/analytics/sales?from=&to=` | Финансовая аналитика по счетам за период |
| GET | `/analytics/purchase-requests` | Количество заявок в разрезе статусов |
| GET | `/analytics/top-products?limit=10` | Самые востребованные продукты |

Даты — в формате `YYYY-MM-DD`, обе границы опциональны.

### Отчёты

Чтение/экспорт — любой авторизованный; создание/генерация/удаление — **admin**.

| Метод | Путь | Доступ | Описание |
|-------|------|--------|----------|
| GET | `/reports` | авторизован | Список отчётов |
| GET | `/reports/:id` | авторизован | Отчёт |
| GET | `/reports/:id/export` | авторизован | Скачать как CSV (UTF-8 BOM, для Excel) |
| POST | `/reports` | admin | Создать отчёт вручную |
| POST | `/reports/generate-sales` | admin | Сформировать отчёт о продажах за период |
| DELETE | `/reports/:id` | admin | Удалить отчёт |

`POST /reports/generate-sales` принимает `{"from":"2026-01-01","to":"2026-12-31","title":"..."}`
(все поля опциональны) и сохраняет отчёт с готовым CSV-содержимым на основе реальной аналитики.

## Демо-данные (сидер)

При первом запуске на пустую БД автоматически создаются демо-данные: пользователь,
3 категории, 2 поставщика, 6 продуктов, 4 заявки на закупку и 2 счёта-фактуры.
Наполнение **идемпотентно** (при повторном старте пропускается) и отключается через
`SEED_DEMO=false`.

Демо-пользователь: `ivanov` / `user123` (роль `user`).

## Пример использования

```bash
# Вход администратором
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | jq -r .token)

# Создать категорию
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Вентиляционные трубы","description":"Основная категория"}'

# Создать продукт
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Труба d100","price":350.00,"category_id":1}'
```

## Переменные окружения

См. `.env.example`. Ключевые: `JWT_SECRET` (обязательно), `DB_*`,
`ADMIN_USERNAME`/`ADMIN_PASSWORD`.

## CI/CD (GitHub Actions)

Один пайплайн `.github/workflows/ci.yml` (CI/CD), запускается на push/PR в `master`:

- **test** — `gofmt`, `go vet`, `go build`, `go test -race`.
- **docker** — пробная сборка Docker-образа (валидация Dockerfile).
- **deploy** — только при push в `master` и после успешных `test` + `docker`:
  заходит на сервер по SSH и выполняет
  `cd ~/crm && git fetch && git reset --hard origin/master && docker compose up --build -d`.

Образ собирается прямо на сервере (тем же `docker-compose.yml`, что и при ручном
запуске). GHCR не используется. Файл `.env` живёт на сервере и в Git не попадает —
поэтому секреты приложения в GitHub не нужны, только доступ по SSH.

### Что нужно настроить (один раз)

1. **SSH-ключ для GitHub Actions.** На своём ПК:
   ```bash
   ssh-keygen -t ed25519 -C "github-deploy" -f deploy_key
   ```
   Публичный ключ добавить на сервер:
   ```bash
   ssh-copy-id -i deploy_key.pub root@<сервер>
   # или вручную дописать deploy_key.pub в ~/.ssh/authorized_keys на сервере
   ```

2. **Секреты репозитория** (Settings → Secrets and variables → Actions → New secret):

   | Секрет | Значение |
   |--------|----------|
   | `SSH_HOST` | адрес сервера (например `165.245.241.121`) |
   | `SSH_USER` | пользователь SSH (например `root`) |
   | `SSH_KEY`  | приватный ключ `deploy_key` (весь файл целиком) |
   | `SSH_PORT` | порт SSH (опц., по умолчанию 22) |

3. На сервере проект уже склонирован в `~/crm`, есть `.env` и запущенные контейнеры.

После этого любой `git push` в `master` автоматически прогоняет проверки и,
если они зелёные, обновляет приложение на сервере.

## Что дальше (вне текущей итерации)

- Экспорт в нативный XLSX/PDF (сейчас — CSV, открывается в Excel)
- Уведомления ответственных лиц по заявкам
- Интеграция внешних API
```
