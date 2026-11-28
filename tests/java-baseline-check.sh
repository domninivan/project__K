#!/usr/bin/env bash
# Dedicated PostgreSQL test server only; requires CREATEDB.
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
: "${BASELINE_TEST_ADMIN_URL:?Set an admin URL for a dedicated test server}"
test_database="baseline_test_${$}_${RANDOM}"
export PGUSER="${DB_USERNAME:-postgres}"
export PGPASSWORD="${DB_PASSWORD:?Set DB_PASSWORD for the test server}"
psql_admin() { psql -X -qAt --no-password -v ON_ERROR_STOP=1 "$BASELINE_TEST_ADMIN_URL" "$@"; }
cleanup() { psql_admin -c "DROP DATABASE IF EXISTS $test_database"; }
trap cleanup EXIT
psql_admin -c "CREATE DATABASE $test_database"
export DB_URL="jdbc:${BASELINE_TEST_ADMIN_URL%/*}/$test_database"
sql() { psql -X -qAt --no-password -v ON_ERROR_STOP=1 "${DB_URL#jdbc:}" "$@"; }
check() { bash "$ROOT_DIR/Scripts/java-baseline-check.sh"; }

[[ "$(check)" == false ]] || { echo 'Empty DB must not baseline'; exit 1; }
sql -f "$ROOT_DIR/JavaBackend/production-java-backend/src/main/resources/db/migration/V1__initial_schema.sql"
[[ "$(check)" == true ]] || { echo 'V1 must baseline'; exit 1; }
sql -c "INSERT INTO product(title) VALUES ('Keep this data')"
[[ "$(check)" == true ]] || { echo 'Data must not affect schema detection'; exit 1; }
[[ "$(sql -c "SELECT count(*) FROM product")" == 1 ]] || exit 1
[[ "$(sql -c "SELECT to_regclass('public.flyway_schema_history') IS NULL")" == t ]] || exit 1
sql -c "ALTER TABLE product ADD COLUMN unexpected text"
if check; then echo 'Changed schema must be rejected'; exit 1; fi
sql -c "CREATE TABLE flyway_schema_history(version text)"
[[ "$(check)" == false ]] || { echo 'Managed DB must not baseline again'; exit 1; }
export DB_URL="${DB_URL}?currentSchema=public"
if check; then echo 'Unsupported JDBC options must be rejected'; exit 1; fi
echo '[OK] Empty, V1, preserved data, schema mismatch, managed DB, URL validation.'
