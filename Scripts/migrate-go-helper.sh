
#!/usr/bin/env bash

set -euo pipefail

# Использование migrate-go-helper.sh:
#
# ./Scripts/migrate-go-helper.sh up                       - Применить все новые миграции
# ./Scripts/migrate-go-helper.sh down                     - Откатить последнюю применённую миграцию
# ./Scripts/migrate-go-helper.sh down-all                 - Полностью откатить все миграции
# ./Scripts/migrate-go-helper.sh version                  - Показать текущую версию миграций
# ./Scripts/migrate-go-helper.sh force 1                  - Принудительно установить указанную версию
# ./Scripts/migrate-go-helper.sh create add_active_column - Создать файлы новой миграции

# Корень проекта
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Пути GoBackend
GO_BACKEND_DIR="$PROJECT_ROOT/GoBackend"
MIGRATIONS_DIR="$GO_BACKEND_DIR/migrations"
ENV_FILE="$GO_BACKEND_DIR/.env"

# Загружаем .env только если переменные не переданы извне
if [[ -z "${GO_DB_HOST:-}" ||
      -z "${GO_DB_PORT:-}" ||
      -z "${GO_DB_USER:-}" ||
      -z "${GO_DB_PASSWORD:-}" ||
      -z "${GO_DB_NAME:-}" ]]; then

    if [[ ! -f "$ENV_FILE" ]]; then
        echo "[ERROR] GoBackend/.env not found."
        exit 1
    fi

    echo "[INFO] Loading configuration from GoBackend/.env"

    set -a
    source "$ENV_FILE"
    set +a
else
    echo "[INFO] Using environment variables provided externally."
fi

# Проверяем необходимые переменные
: "${GO_DB_HOST:?Не задан GO_DB_HOST}"
: "${GO_DB_PORT:?Не задан GO_DB_PORT}"
: "${GO_DB_USER:?Не задан GO_DB_USER}"
: "${GO_DB_PASSWORD:?Не задан GO_DB_PASSWORD}"
: "${GO_DB_NAME:?Не задан GO_DB_NAME}"

# Формируем PostgreSQL connection string
DATABASE_URL="postgres://${GO_DB_USER}:${GO_DB_PASSWORD}@${GO_DB_HOST}:${GO_DB_PORT}/${GO_DB_NAME}?sslmode=disable"

# Проверяем наличие golang-migrate
if ! command -v migrate &>/dev/null; then
    echo "Ошибка: golang-migrate не установлен."
    echo "Установка:"
    echo "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Проверяем наличие директории миграций
if [[ ! -d "$MIGRATIONS_DIR" ]]; then
    echo "Ошибка: директория migrations не найдена."
    exit 1
fi

# Создание новой миграции
if [[ "${1:-}" == "create" ]]; then
    if [[ -z "${2:-}" ]]; then
        echo "Использование: ./Scripts/migrate-go-helper.sh create migration_name"
        exit 1
    fi

    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$2"
    exit 0
fi

case "${1:-}" in
    up)
        echo "Применение миграций..."
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
        ;;

    down)
        echo "Откат последней миграции..."
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down 1
        ;;

    down-all)
        read -r -p "Откатить ВСЕ миграции? (yes/no): " confirm

        if [[ "$confirm" == "yes" ]]; then
            migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down -all
        else
            echo "Отменено."
        fi
        ;;

    version)
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
        ;;

    force)
        if [[ -z "${2:-}" ]]; then
            echo "Использование: ./Scripts/migrate-go-helper.sh force VERSION"
            exit 1
        fi

        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$2"
        ;;

    *)
        echo "Использование: ./Scripts/migrate-go-helper.sh {up|down|down-all|version|force VERSION|create NAME}"
        exit 1
        ;;
esac
