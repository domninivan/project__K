
#!/usr/bin/env bash

set -Eeuo pipefail

PROJECT="/workspaces/project__K"
JAVA_PROJECT="$PROJECT/JavaBackend/production-java-backend"
GO_BACKEND="$PROJECT/GoBackend"

git config --global --add safe.directory "$PROJECT"

echo "========================================"
echo "Preparing development environment..."
echo "========================================"

# ----------------------------------------
# 3. Fix ownership of named volumes
# ----------------------------------------

sudo chown -R vscode:vscode \
    /home/vscode/.m2 \
    /home/vscode/.cache/go-build \
    /go/bin \
    /go/pkg/mod \
    "$PROJECT/Frontend/node_modules" \
    "$JAVA_PROJECT/target"

sudo mkdir -p /go/pkg/sumdb
sudo chown -R vscode:vscode /go/pkg/sumdb
# ----------------------------------------
# 4. Make scripts executable (optional)
# ----------------------------------------

if [[ -d "$PROJECT/Scripts" ]]; then
    echo "Setting executable permissions for shell scripts..."

    if find "$PROJECT/Scripts" \
        -type f \
        -name "*.sh" \
        -exec chmod +x {} +; then
        echo "[OK] Executable permissions configured"
    else
        echo "[WARNING] Could not set executable permissions."
        echo "[WARNING] Scripts will be launched through bash."
    fi
fi

# ----------------------------------------
# 5. Install PostgreSQL client
# ----------------------------------------

if ! command -v psql >/dev/null 2>&1; then
    sudo apt-get update
    sudo apt-get install -y postgresql-client
fi

# ----------------------------------------
# 6. Install Go development tools
# ----------------------------------------

if ! command -v migrate >/dev/null 2>&1; then
    echo "Installing golang-migrate CLI..."

    if GOBIN=/go/bin go install \
        -tags 'postgres' \
        github.com/golang-migrate/migrate/v4/cmd/migrate@latest; then
        echo "[OK] golang-migrate installed"
    else
        echo "[ERROR] Failed to install golang-migrate CLI"
        exit 1
    fi
else
    echo "[OK] golang-migrate already installed"
fi


if [[ -f "$GO_BACKEND/.env.example" ]]; then
    cp "$GO_BACKEND/.env.example" "$GO_BACKEND/.env"
    echo "[OK] GoBackend/.env updated from GoBackend/.env.example"
else
    echo "[ERROR] GoBackend/.env.example not found"
    exit 1
fi

if ! command -v air >/dev/null 2>&1; then
    go install github.com/air-verse/air@latest
fi

# ----------------------------------------
# 7. Install Bruno CLI
# ----------------------------------------

if ! command -v bru >/dev/null 2>&1; then
    npm install -g @usebruno/cli
fi

# ----------------------------------------
# 8. Install project dependencies
# ----------------------------------------

if [ -f "$PROJECT/Scripts/install-dependencies.sh" ]; then
    echo "Installing project dependencies..."

    bash "$PROJECT/Scripts/install-dependencies.sh"
else
    echo "WARNING: install-dependencies.sh not found"
fi

# ----------------------------------------
# Create GoBackend database if not exists
# ----------------------------------------

echo ""
echo "========================================"
echo "Preparing GoBackend database..."
echo "========================================"

: "${GO_DB_HOST:?GO_DB_HOST is not set}"
: "${GO_DB_PORT:?GO_DB_PORT is not set}"
: "${GO_DB_USER:?GO_DB_USER is not set}"
: "${GO_DB_PASSWORD:?GO_DB_PASSWORD is not set}"
: "${GO_DB_NAME:?GO_DB_NAME is not set}"

export PGPASSWORD="$GO_DB_PASSWORD"

cleanup() {
    unset PGPASSWORD
}

trap cleanup EXIT

if ! pg_isready \
    -h "$GO_DB_HOST" \
    -p "$GO_DB_PORT" \
    -U "$GO_DB_USER" >/dev/null 2>&1; then

    echo "[ERROR] PostgreSQL is not available"
    exit 1
fi

DB_EXISTS=$(psql \
    -h "$GO_DB_HOST" \
    -p "$GO_DB_PORT" \
    -U "$GO_DB_USER" \
    -d postgres \
    -tAc "SELECT EXISTS (
        SELECT 1
        FROM pg_database
        WHERE datname = '$GO_DB_NAME'
    );")

if [[ "$DB_EXISTS" == "t" ]]; then
    echo "[OK] Database '$GO_DB_NAME' already exists. Skipping."
else
    psql \
        -h "$GO_DB_HOST" \
        -p "$GO_DB_PORT" \
        -U "$GO_DB_USER" \
        -d postgres \
        -v ON_ERROR_STOP=1 \
        -c "CREATE DATABASE \"$GO_DB_NAME\" OWNER \"$GO_DB_USER\";"

    echo "[OK] Database '$GO_DB_NAME' created."
fi

echo "[OK] GoBackend database preparation completed."

# ----------------------------------------
# Apply GoBackend database migrations
# ----------------------------------------

echo ""
echo "========================================"
echo "Applying GoBackend migrations..."
echo "========================================"

bash "$PROJECT/Scripts/migrate-go-helper.sh" up

echo "[OK] GoBackend migrations applied."

# ----------------------------------------
# 9. Verify installed tools
# ----------------------------------------

echo ""
echo "========================================"
echo "Installed tools:"
echo "========================================"

git --version
go version
air -v
java -version
mvn -version
node --version
npm --version
bru --version
python3 --version
psql --version
docker --version

echo ""
echo "========================================"
echo "Development environment is ready!"
echo "========================================"
