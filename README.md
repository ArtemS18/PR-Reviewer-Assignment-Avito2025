# PR-Reviewer-Assignment-Service

Сервис для назначения ревьюеров на pull request’ы.  
Состоит из Go‑backend’а (`reviewer-api`), PostgreSQL и вспомогательных скриптов для нагрузки и заполнения БД.

---

## Структура проекта

- `reviewer-api/` — Go‑сервис.
  - `cmd/reviewer-api/main.go` — входная точка API.
  - `cmd/migrate/main.go` — миграции (GORM AutoMigrate).
  - `internal/app/ds` — модели БД (`Team`, `User`, `PullRequest`, `Reviewer`).
  - `internal/app/http-server/handlers` — HTTP‑хендлеры.
  - `Dockerfile` — сборка сервиса.
- `docker-compose.yaml` — Postgres + backend.
- `seed.sql` — SQL‑скрипт для заполнения БД тестовыми данными.
- `testing/` — Python‑скрипты для сидинга и нагрузочного тестирования.
- `README.md` — этот файл.

---

## Запуск через Docker

1. Создайте файл `.env` в корне репозитория (либо используйте имеющийся):

   ```env
   HTTP_HOST=0.0.0.0
   HTTP_PORT=8080

   DB_HOST=postgres
   DB_PORT=5432
   DB_NAME=postgres
   DB_USER=postgres
   DB_PASS=postgres
   ```


### 2. Запуск через Docker Compose

```docker-compose up --build ```
Таблицы создадутся автоматически при запуске контейнеров. Для заполнения их тестовыми данными (200 users и 20 teams) воспользуйтесь python-скриптами (```Нагрузочное тестирование и заполнен данными (Python)```)

- Приложение: [http://localhost:8080](http://localhost:8080)


### 3. Миграции БД

**Локально (создают только сами таблицы, без заполнения данных):**

```cd reviewer-api
go run ./cmd/migrate
```

## Тестирование

### Unit-тесты (Go)

```cd reviewer-api
go test ./internal/app/http-server/handlers/...
```

### Нагрузочное тестирование и заполнен данными (Python)

Установите зависимости:

```
cd testing
pip install -r req.txt
```

Запуск стресс-теста (длится порядка 10 секунд) ПЕРЕЗАПИСЫВАЕТ users(200) и teams(20):

```
python main.py test
```

Наполнить БД тестовыми данными ПЕРЕЗАПИСЫВАЕТ users(200) и teams(20):

Python-скрипт

```
python main.py fill_db
```
---

Пример вывода теста:
```
--- Stress test results ---
Total requests: 300
Success: 300 (100.0000%)
Latency p50: 80.49 ms
Latency p95: 117.35 ms
Latency p99: 127.35 ms
Errors: 0

SLI check:
- success_rate >= 99.9%: OK
- p95 <= 300 ms: OK
```
