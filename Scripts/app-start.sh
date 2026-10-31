
#!/bin/bash

# Переходим в корень проекта
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "=================================="
echo " Starting project..."
echo "=================================="

# Запуск Go Backend
echo "[1/3] Starting Go Backend..."
(
    cd "$ROOT_DIR/GoBackend"
    go run ./cmd/main.go
) &

GO_PID=$!

# Запуск Java Backend
echo "[2/3] Starting Java Backend..."
(
    cd "$ROOT_DIR/JavaBackend/production-java-backend"
    ./mvnw spring-boot:run
) &

JAVA_PID=$!

# Запуск Frontend
echo "[3/3] Starting Frontend..."
(
    cd "$ROOT_DIR/Frontend"
    npm run dev
) &

FRONT_PID=$!

echo ""
echo "=================================="
echo " All services started!"
echo "=================================="
echo "Go PID: $GO_PID"
echo "Java PID: $JAVA_PID"
echo "Frontend PID: $FRONT_PID"
echo ""
echo "Press CTRL+C to stop all services."

# При остановке завершаем все процессы
cleanup() {
    echo ""
    echo "Stopping services..."
    kill "$GO_PID" "$JAVA_PID" "$FRONT_PID" 2>/dev/null
    wait
    echo "All services stopped."
}

trap cleanup SIGINT SIGTERM

wait
