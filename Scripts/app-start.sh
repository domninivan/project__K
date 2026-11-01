
#!/usr/bin/env bash

set -u

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

GO_PORT=8081
JAVA_PORT=8080
FRONT_PORT=5173

fail() {
    echo "[ERROR] $*" >&2
    exit 1
}

# Проверяем необходимые команды.
for cmd in air go npm setsid; do
    command -v "$cmd" >/dev/null 2>&1 || fail "Не найдена команда: $cmd"
done

[[ -x "$ROOT_DIR/JavaBackend/production-java-backend/mvnw" ]] \
    || fail "Не найден Maven Wrapper или отсутствует право на выполнение"

if [[ ! -f "$ROOT_DIR/GoBackend/.air.toml" ]]; then
    echo "[WARNING] GoBackend/.air.toml не найден. Air может не запуститься."
fi

GO_PID=""
JAVA_PID=""
FRONT_PID=""

PIDS=()
declare -A SERVICE_NAMES

cleanup() {
    local exit_code=$?

    # Отключаем обработчики, чтобы cleanup не вызвался повторно.
    trap - INT TERM EXIT

    echo
    echo "[STOP] Завершение всех сервисов..."

    # Отправляем SIGTERM каждой группе процессов.
    for pid in "${PIDS[@]}"; do
        echo "[STOP] ${SERVICE_NAMES[$pid]}"
        kill -TERM -- "-$pid" 2>/dev/null || true
    done

    # Ждём до 5 секунд.
    for i in {1..5}; do
        local running=0

        for pid in "${PIDS[@]}"; do
            if kill -0 -- "-$pid" 2>/dev/null; then
                running=1
            fi
        done

        (( running == 0 )) && break
        sleep 1
    done

    # Принудительно завершаем оставшиеся группы.
    for pid in "${PIDS[@]}"; do
        kill -KILL -- "-$pid" 2>/dev/null || true
    done

    # Собираем завершившиеся фоновые процессы.
    for pid in "${PIDS[@]}"; do
        wait "$pid" 2>/dev/null || true
    done

    echo "[STOP] Все сервисы остановлены."
    exit "$exit_code"
}

trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

start_service() {
    local name="$1"
    local dir="$2"
    shift 2

    echo "[START] $name"

    setsid bash -c '
        cd "$1" || exit 1
        shift
        exec "$@"
    ' _ "$dir" "$@" &

    local pid=$!
    PIDS+=("$pid")
    SERVICE_NAMES["$pid"]="$name"

    echo "[INFO] $name запущен, PID=$pid, PGID=$pid"
}

wait_for_port() {
    local name="$1"
    local port="$2"
    local timeout="$3"

    echo "[WAIT] Ожидание $name на порту $port..."

    for ((i = 0; i < timeout; i++)); do
        if (echo > "/dev/tcp/127.0.0.1/$port") >/dev/null 2>&1; then
            echo "[OK] $name доступен на порту $port"
            return 0
        fi

        sleep 1
    done

    echo "[WARNING] $name не открыл порт $port за ${timeout} секунд."
    return 1
}

echo
echo "========================================"
echo "Проверка окружения"
echo "========================================"

[[ -d "$ROOT_DIR/GoBackend" ]] || fail "GoBackend не найден"
[[ -d "$ROOT_DIR/Frontend" ]] || fail "Frontend не найден"
[[ -d "$ROOT_DIR/JavaBackend/production-java-backend" ]] \
    || fail "JavaBackend не найден"

[[ -f "$ROOT_DIR/GoBackend/go.mod" ]] || fail "Не найден go.mod"
[[ -f "$ROOT_DIR/Frontend/package.json" ]] || fail "Не найден package.json"
[[ -d "$ROOT_DIR/Frontend/node_modules" ]] \
    || fail "Не найдена Frontend/node_modules. Выполни npm install."

echo "[OK] Окружение проверено."

echo
echo "========================================"
echo "Запуск сервисов"
echo "========================================"

# Go Backend.
start_service "Go Backend" "$ROOT_DIR/GoBackend" air
GO_PID="${PIDS[-1]}"

# Java Backend.
start_service "Java Backend" \
    "$ROOT_DIR/JavaBackend/production-java-backend" \
    ./mvnw spring-boot:run
JAVA_PID="${PIDS[-1]}"

# Frontend.
start_service "Frontend" "$ROOT_DIR/Frontend" npm run dev -- --host 0.0.0.0
FRONT_PID="${PIDS[-1]}"

echo
echo "========================================"
echo "Ожидание запуска сервисов"
echo "========================================"

wait_for_port "Go Backend" "$GO_PORT" 30 || true
wait_for_port "Java Backend" "$JAVA_PORT" 30 || true
wait_for_port "Frontend" "$FRONT_PORT" 30 || true

echo
echo "========================================"
echo "Сервисы запущены"
echo "========================================"
echo "Java:     http://localhost:$JAVA_PORT"
echo "Go:       http://localhost:$GO_PORT"
echo "Frontend: http://localhost:$FRONT_PORT"
echo "========================================"
echo
echo "Нажми Ctrl+C для остановки всех сервисов."
echo

# Если любой сервис завершится, останавливаем остальные.
EXITED_PID=""

if wait -n -p EXITED_PID "${PIDS[@]}"; then
    EXIT_CODE=0
else
    EXIT_CODE=$?
fi

EXITED_NAME="${SERVICE_NAMES[$EXITED_PID]:-Неизвестный сервис}"

echo
echo "[ERROR] Сервис '$EXITED_NAME' завершился (код: $EXIT_CODE)."
echo "[INFO] Останавливаем остальные сервисы..."

exit "$EXIT_CODE"
