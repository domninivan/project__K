#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
command -v psql >/dev/null 2>&1 || { echo '[ERROR] Для проверки Java-БД нужен psql.' >&2; exit 1; }

# Same defaults as database.properties. Do not guess a database from Go settings.
java_db_url="${DB_URL:-jdbc:postgresql://127.0.0.1:5433/postgres}"
if [[ "$java_db_url" != jdbc:postgresql://* || "$java_db_url" == *'?'* || "$java_db_url" == *'@'* || "$java_db_url" == *'#'* ]]; then
    echo '[ERROR] Автопроверка поддерживает DB_URL вида jdbc:postgresql://host:port/database без параметров и учётных данных.' >&2
    exit 1
fi

state="$(PGPASSWORD="${DB_PASSWORD:-postgres}" PGCONNECT_TIMEOUT=10 \
    PGOPTIONS='-c default_transaction_read_only=on' \
    psql -X -qAt -v ON_ERROR_STOP=1 --no-password \
        --dbname="${java_db_url#jdbc:}" --username="${DB_USERNAME:-postgres}" \
        --file="$ROOT_DIR/Scripts/java-schema-state.sql")" || {
    echo '[ERROR] Не удалось проверить Java-БД. Сервисы не запущены.' >&2
    exit 1
}
expected="$(tr -d '\r\n' < "$ROOT_DIR/Scripts/java-v1.signature")"
case "$state" in
    EMPTY|MANAGED)
        echo "[INFO] Java-БД: $state. Baseline не требуется." >&2
        echo false
        ;;
    "$expected")
        echo '[INFO] Java-БД соответствует исходной V1, истории нет. Включён однократный baseline версии 1.' >&2
        echo true
        ;;
    *)
        echo '[ERROR] Java-БД не имеет истории Flyway и отличается от исходной V1. Автоматический baseline запрещён.' >&2
        echo '[ERROR] Проверьте DB_URL и схему БД; инструкция: docs/15-local-development.md.' >&2
        exit 1
        ;;
esac
