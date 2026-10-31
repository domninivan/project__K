# Case 22 — План производства по заказам

Учебный проект по кейсу №22 «План производства по заказам (мастер-план)».

## Архитектура

```text
React
   |
   v
Java Spring Boot
  /            \
 v              v
Go              1С
```

- React взаимодействует только с Java.
- Java — центральный API и оркестратор бизнес-процесса заказа.
- Go отвечает за производственное планирование.
- 1С отвечает за материалы, склад, закупки и себестоимость.
- Go и 1С напрямую друг с другом не взаимодействуют.

---

# Основной способ запуска — Dev Container

Dev Container используется как единая среда разработки для команды.

Внутри `app` доступны Java, Go, Node.js, Python и инструменты разработки.
PostgreSQL запускается отдельным сервисом Docker Compose.

## Что нужно установить

На компьютере достаточно:

- Git
- Docker Desktop или Podman
- Visual Studio Code
- расширение **Dev Containers**

Отдельно устанавливать Java, Maven, Tomcat, Go, Node.js и Python для работы через Dev Container не требуется.

---

## 1. Получить проект

Если проект ещё не скачан:

```bash
git clone <URL_РЕПОЗИТОРИЯ>
cd project__K
```

Если проект уже есть:

```bash
git pull
```

---

## 2. Открыть Dev Container

Открыть в VS Code именно корень:

```text
project__K
```

Затем:

```text
Ctrl + Shift + P
→ Dev Containers: Rebuild and Reopen in Container
```

VS Code использует:

```text
.devcontainer/devcontainer.json
docker-compose.yml
```

---

## 3. Что запускается

```text
app
└── среда разработки
    ├── Java
    ├── Go
    ├── Node.js
    ├── Python
    └── /workspaces/project__K

postgres
└── PostgreSQL 17
```

Java внутри Dev Container подключается к PostgreSQL через:

```text
postgres:5432
```

Переменные окружения задаются в `docker-compose.yml`:

```text
DB_URL=jdbc:postgresql://postgres:5432/postgres
DB_USERNAME=postgres
DB_PASSWORD=postgres
```

---

## 4. Проверить окружение

В терминале Dev Container:

```bash
java -version
```

Для Java-backend используется Java 21.

Проверить подключение к БД:

```bash
echo $DB_URL
```

Ожидается:

```text
jdbc:postgresql://postgres:5432/postgres
```

Проверить сервисы:

```bash
docker compose ps
```

---

## 5. Запустить Java-backend

```bash
cd JavaBackend/production-java-backend
./mvnw spring-boot:run
```

После успешного запуска:

```text
Tomcat started on port 8080
Started ProductionApplication
```

Backend:

```text
http://localhost:8080
```

Отдельный Tomcat не нужен — Spring Boot использует встроенный.

---

## 6. Проверить API

```text
GET http://localhost:8080/api/orders/hello
```

Ожидаемый ответ:

```text
Hello from Case 22
```

Список заказов:

```text
GET http://localhost:8080/api/orders
```

Заказ по ID:

```text
GET http://localhost:8080/api/orders/{orderId}
```

---

## 7. Создать тестовый заказ

```text
POST http://localhost:8080/api/orders
Content-Type: application/json
```

```json
{
  "dueDate": "2026-11-15",
  "comment": "Тестовый заказ",
  "items": [
    {
      "productId": 1,
      "quantity": 100
    },
    {
      "productId": 2,
      "quantity": 50
    }
  ]
}
```

---

# PostgreSQL

Внутри Docker-сети:

```text
host: postgres
port: 5432
database: postgres
user: postgres
password: postgres
```

С компьютера:

```text
host: 127.0.0.1
port: 5433
database: postgres
user: postgres
password: postgres
```

---

## Инициализация БД

При первом создании PostgreSQL volume автоматически выполняются:

```text
database/01-schema.sql
database/02-seed.sql
```

`01-schema.sql` создаёт:

```text
product
product_operation
product_material_requirement
production_order
production_order_item
```

`02-seed.sql` добавляет тестовые данные.

Если скрипты были изменены после первого запуска:

```bash
docker compose down -v
docker compose up -d postgres
```

> `docker compose down -v` удаляет локальные данные PostgreSQL.

---

## Проверка БД

Список таблиц:

```bash
docker compose exec postgres psql -U postgres -d postgres -c "\dt"
```

Детали:

```bash
docker compose exec postgres psql -U postgres -d postgres -c "SELECT * FROM product;"
```

Заказы:

```bash
docker compose exec postgres psql -U postgres -d postgres -c "SELECT * FROM production_order;"
```

---

# Локальный запуск Java без Dev Container

Также можно запускать Java непосредственно на компьютере.

Нужны:

- JDK 21
- Git
- Docker Desktop / Podman

Maven отдельно не нужен: используется Maven Wrapper.

Запустить PostgreSQL:

```bash
docker compose up -d postgres
```

Перейти в Java-проект:

### Windows

```powershell
cd JavaBackend\production-java-backend
```

### Linux / macOS

```bash
cd JavaBackend/production-java-backend
```

Запустить:

### Windows PowerShell

```powershell
.\mvnw.cmd spring-boot:run
```

### Windows cmd

```bat
mvnw.cmd spring-boot:run
```

### Linux / macOS

```bash
./mvnw spring-boot:run
```

---

# Конфигурация подключения к БД

`database.properties`:

```properties
db.driver=org.postgresql.Driver
db.url=${DB_URL:jdbc:postgresql://127.0.0.1:5433/postgres}
db.username=${DB_USERNAME:postgres}
db.password=${DB_PASSWORD:postgres}
```

Локально:

```text
127.0.0.1:5433
```

В Dev Container:

```text
postgres:5432
```

---

# Сборка Java

Dev Container / Linux:

```bash
./mvnw clean package
```

Windows:

```powershell
.\mvnw.cmd clean package
```

JAR:

```text
JavaBackend/production-java-backend/target/
```

Запуск:

```bash
java -jar target/production-java-backend.jar
```

---

# Полезные команды Docker

```bash
docker compose ps
docker compose logs postgres
docker compose restart postgres
docker compose stop postgres
docker compose down
```

Полностью удалить локальную БД:

```bash
docker compose down -v
```

---

# Технологии Java-backend

- Java 21
- Spring Boot 3.5
- Spring MVC
- Hibernate 5.6
- native `SessionFactory`
- Spring Transactions
- PostgreSQL 17
- c3p0
- MapStruct
- Bean Validation
- Maven Wrapper

Spring Boot отвечает за запуск приложения и web-слой.

Persistence-слой остаётся на DAO + Hibernate `SessionFactory`.

---

# Структура репозитория

```text
project__K/
├── .devcontainer/
│   ├── devcontainer.json
│   └── devcontainer-lock.json
│
├── database/
│   ├── 01-schema.sql
│   └── 02-seed.sql
│
├── JavaBackend/
│   └── production-java-backend/
│       ├── pom.xml
│       ├── mvnw
│       ├── mvnw.cmd
│       ├── .mvn/
│       └── src/
│
├── GoBackend/
├── Frontend/
├── docker-compose.yml
└── README.md
```

---

# Быстрый старт для участника команды

```text
1. git pull / git clone
2. открыть project__K в VS Code
3. Ctrl+Shift+P
4. Dev Containers: Rebuild and Reopen in Container
5. открыть терминал
6. cd JavaBackend/production-java-backend
7. ./mvnw spring-boot:run
8. открыть http://localhost:8080/api/orders/hello
```

Если сервер вернул:

```text
Hello from Case 22
```

среда готова к работе.
