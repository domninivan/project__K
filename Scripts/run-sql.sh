
#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Использование: $0 <путь-к-sql-файлу>"
    exit 1
fi

SQL_FILE="$1"

if [[ ! -f "$SQL_FILE" ]]; then
    echo "SQL-файл не найден: $SQL_FILE"
    exit 1
fi

PGPASSWORD="${DB_PASSWORD:-postgres}" psql \
    -h "${DB_HOST:-postgres}" \
    -p "${DB_PORT:-5432}" \
    -U "${DB_USERNAME:-postgres}" \
    -d "${DB_NAME:-postgres}" \
    -v ON_ERROR_STOP=1 \
    -f "$SQL_FILE"
