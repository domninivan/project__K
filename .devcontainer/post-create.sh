
#!/usr/bin/env bash

set -Eeuo pipefail

PROJECT="/workspaces/project__K"
JAVA_PROJECT="$PROJECT/JavaBackend/production-java-backend"
GO_BACKEND="$PROJECT/GoBackend"

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
# 4. Make scripts executable
# ----------------------------------------

# if [ -d "$PROJECT/Scripts" ]; then
#     find "$PROJECT/Scripts" \
#         -type f \
#         -name "*.sh" \
#         -exec chmod +x {} +
# fi


if [[ ! -f "$GO_BACKEND/.env" ]]; then
    if [[ -f "$GO_BACKEND/.env.example" ]]; then
        cp "$GO_BACKEND/.env.example" "$GO_BACKEND/.env"
        echo "[OK] Created GoBackend/.env from GoBackend/.env.example"
    else
        echo "[ERROR] GoBackend/.env.example not found"
        exit 1
    fi
else
    echo "[OK] GoBackend/.env already exists, skipping"
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
