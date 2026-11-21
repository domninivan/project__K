# Локальная разработка

Инструкция сверена с конфигурацией ветки 21.11.2026; запуск в рамках переноса документации не выполнялся.

Первый запуск через Dev Container, установка программ, автоматическая подготовка БД и пересборка описаны в [первоначальной настройке](initial_setup.md). Здесь приведены команды повседневной работы и настройки запуска с хоста.

Нужны Java 21, Go версии из `GoBackend/go.mod`, Node/npm совместимые с `Frontend/package.json`, Bash, psql. Для общего скрипта запуска нужен `air`, для Go-миграций — `migrate`. Настройки и установка среды описаны в `.devcontainer/post-create.sh`; при ошибке установки проверяйте доступность этих инструментов.

После успешного post-create зависимости и Go-миграции уже установлены. Для обычного запуска в контейнере из корня:

```bash
bash Scripts/app-start.sh
```

После изменения зависимостей выполните `bash Scripts/install-dependencies.sh`, после добавления Go-миграций — `bash Scripts/migrate-go-helper.sh up`.

Go-миграции выполняются в отдельной `go_db`. Начальная Java-схема и seed монтируются из `database/` и применяются PostgreSQL только при первой инициализации пустого тома. Изменение SQL не обновляет уже существующий том автоматически.

| Параметры | Внутри Compose | Запуск с Windows-хоста |
| --- | --- | --- |
| PostgreSQL | `postgres:5432` | `localhost:5433` |
| Java | DB_URL, DB_USERNAME, DB_PASSWORD из Compose | Задать DB_URL на `jdbc:postgresql://localhost:5433/postgres` и свои учётные параметры |
| Go | GO_DB_* из Compose, БД `go_db` | По `GoBackend/.env.example`, заменить GO_DB_HOST на localhost и GO_DB_PORT на 5433 |
| Java → Go | `go.base-url=http://localhost:8081`, когда оба процесса в app | Адрес процесса Go; по умолчанию localhost:8081 |

Java читает параметры из ресурсов и окружения; Go — из окружения и локального `.env` (через godotenv). Для Bash-помощника миграций можно создать локальный `GoBackend/.env` по примеру. Не публикуйте локальные секреты. При запуске с хоста создайте также отдельную БД `go_db`, затем примените Go-миграции; в Dev Container это делает post-create. С хоста сначала можно запустить только БД командой `docker compose up -d postgres` из корня.

Команды запуска по отдельности из соответствующего каталога:

| Каталог | Команда |
| --- | --- |
| `GoBackend` | `go run ./cmd/api` (переменные окружения уже заданы) |
| `JavaBackend/production-java-backend` | Windows: `.\mvnw.cmd spring-boot:run`; Linux: `bash ./mvnw spring-boot:run` |
| `Frontend` | `npm install`, затем `npm run dev` |

Java: `http://localhost:8080`; Go: `http://localhost:8081`; Vite (фиксированный порт): `http://localhost:5173` (точный адрес в выводе Vite). Проверки Go: `/api/v1/live` и `/api/v1/ready`. Готовность порта не подтверждает сквозной сценарий. React пока работает локально, маршруты Java–Go расходятся: [известные ограничения](03-architecture.md).

Для доступа с хоста к приложениям контейнера используйте вкладку Ports в VS Code: [проброс портов](initial_setup.md#запуск-приложений). В Compose явно опубликован только PostgreSQL.
