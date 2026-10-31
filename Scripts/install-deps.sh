#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=== Установка зависимостей проекта ==="

# Resolve the installed JDK when the shell inherited a stale JAVA_HOME.
if [[ ! -x "${JAVA_HOME:-}/bin/java" || ! -x "${JAVA_HOME:-}/bin/javac" ]]; then
    if [[ -x /usr/local/sdkman/candidates/java/current/bin/javac ]]; then
        export JAVA_HOME=/usr/local/sdkman/candidates/java/current
    elif command -v javac >/dev/null 2>&1; then
        JAVA_HOME="$(dirname "$(dirname "$(readlink -f "$(command -v javac)")")")"
        export JAVA_HOME
    else
        echo "ОШИБКА: JDK не найден. Пересоберите Dev Container с Java 21." >&2
        exit 1
    fi
fi
export PATH="$JAVA_HOME/bin:$PATH"

# Go backend
if [[ -f "$ROOT/GoBackend/go.mod" ]]; then
    echo "=== GoBackend ==="
    (
        cd "$ROOT/GoBackend"
        go mod download
    )
fi

# React + Vite
if [[ -f "$ROOT/Frontend/package-lock.json" ]]; then
    echo "=== Frontend: npm ci ==="
    (
        cd "$ROOT/Frontend"
        npm ci
    )
elif [[ -f "$ROOT/Frontend/package.json" ]]; then
    echo "=== Frontend: npm install ==="
    (
        cd "$ROOT/Frontend"
        npm install
    )
fi


# Java backend: Maven
JAVA_DIR="$ROOT/JavaBackend/production-java-backend"

if [[ -f "$JAVA_DIR/pom.xml" ]]; then
    echo "=== JavaBackend: Maven ==="
    (
        cd "$JAVA_DIR"

        if [[ -f ./mvnw ]]; then
            sh ./mvnw -B dependency:go-offline
        elif command -v mvn >/dev/null 2>&1; then
            mvn -B dependency:go-offline
        else
            echo "ОШИБКА: Maven не установлен, а mvnw отсутствует." >&2
            exit 1
        fi
    )
else
    echo "ОШИБКА: не найден файл $JAVA_DIR/pom.xml" >&2
    exit 1
fi


echo "=== Установка зависимостей завершена ==="
