# Case 22 — План производства по заказам

Учебный проект по кейсу №22 «План производства по заказам (мастер-план)».

На текущем этапе Java-backend отвечает за:
- справочник деталей;
- технологические операции деталей;
- нормативы материалов;
- производственные заказы и позиции заказов;
- REST API для frontend;
- дальнейшую оркестрацию Go-сервиса и 1С.

## Что нужно установить

Для запуска Java-части нужны:

- Git
- JDK 21
- Docker Desktop или Podman с поддержкой Docker Compose

Отдельно устанавливать Maven и Tomcat не нужно.

Maven запускается через Maven Wrapper, а Spring Boot использует встроенный Tomcat.

---

## Структура запуска

Java запускается локально через Spring Boot.

PostgreSQL запускается в Docker.

```text
Java / Spring Boot
        |
        v
127.0.0.1:5433
        |
        v
PostgreSQL container
```

Dev Container для обычного запуска Java-backend не требуется.

---

## 1. Клонировать проект

```bash
git clone <URL_РЕПОЗИТОРИЯ>
cd project__K
```

---

## 2. Запустить PostgreSQL

Из корня `project__K`:

```bash
docker compose up -d postgres
```

Проверить, что контейнер запущен:

```bash
docker ps
```

PostgreSQL должен быть доступен на:

```text
host: 127.0.0.1
port: 5433
database: postgres
user: postgres
password: postgres
```

При первом создании контейнера автоматически выполняются:

```text
database/01-schema.sql
database/02-seed.sql
```

`01-schema.sql` создаёт таблицы.

`02-seed.sql` добавляет тестовые данные.

---

## 3. Запустить Java-backend

Перейти в Maven-проект:

```bash
cd JavaBackend/production-java-backend
```

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

После успешного запуска в консоли должно появиться сообщение примерно:

```text
Tomcat started on port 8080
Started ProductionApplication
```

Backend будет доступен по адресу:

```text
http://localhost:8080
```

---

## 4. Проверить, что сервер работает

Открыть:

```text
GET http://localhost:8080/api/orders/hello
```

Ожидаемый ответ:

```text
Hello from Case 22
```

---

## 5. Получить список заказов

```text
GET http://localhost:8080/api/orders
```

---

## 6. Создать тестовый заказ

```text
POST http://localhost:8080/api/orders
Content-Type: application/json
```

Пример JSON:

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

## Если база пустая и SQL-скрипты не выполнились

Init-скрипты PostgreSQL выполняются только при создании нового volume.

Если контейнер уже запускался раньше, а SQL-скрипты были добавлены позже, нужно удалить старый volume:

```bash
docker compose down -v
```

После этого снова создать PostgreSQL:

```bash
docker compose up -d postgres
```

> Внимание: `docker compose down -v` удаляет данные текущей локальной БД.

---

## Полезные команды Docker

Запустить PostgreSQL:

```bash
docker compose up -d postgres
```

Остановить:

```bash
docker compose stop postgres
```

Посмотреть логи:

```bash
docker compose logs postgres
```

Проверить таблицы:

```bash
docker compose exec postgres psql -U postgres -d postgres -c "\dt"
```

Проверить тестовые детали:

```bash
docker compose exec postgres psql -U postgres -d postgres -c "SELECT * FROM product;"
```

Полностью пересоздать БД:

```bash
docker compose down -v
docker compose up -d postgres
```

---

## Сборка Java без запуска

### Windows

```powershell
.\mvnw.cmd clean package
```

### Linux / macOS

```bash
./mvnw clean package
```

JAR появится в:

```text
JavaBackend/production-java-backend/target/
```

Запустить собранный JAR можно так:

```bash
java -jar target/production-java-backend.jar
```

---

## Используемые технологии

Java backend:

- Java 21
- Spring Boot 3.5
- Spring MVC
- Hibernate 5.6
- native `SessionFactory`
- Spring Transactions
- PostgreSQL
- c3p0
- MapStruct
- Maven Wrapper

База данных:

- PostgreSQL 17
- Docker Compose
- автоматическая инициализация через SQL

---

## Текущая архитектура

```text
React
   |
   v
Java Spring Boot
  /            \
 v              v
Go              1С
```

Правила взаимодействия:

- React общается только с Java.
- Java является центральным API и оркестратором.
- Go отвечает за производственное планирование.
- 1С отвечает за материалы, склад, закупки и себестоимость.
- Go и 1С напрямую друг с другом не взаимодействуют.

---

## Если проект не запускается

Сначала проверить PostgreSQL:

```bash
docker ps
```

Если контейнер не запущен:

```bash
docker compose up -d postgres
```

Потом проверить, что порт `5433` опубликован.

Java ожидает БД по адресу:

```text
jdbc:postgresql://127.0.0.1:5433/postgres
```

Если в логе есть:

```text
Connection refused
```

то Java работает, но PostgreSQL недоступен.

---

## Быстрый запуск

### Windows

```powershell
git clone <URL_РЕПОЗИТОРИЯ>
cd project__K

docker compose up -d postgres

cd JavaBackend\production-java-backend
.\mvnw.cmd spring-boot:run
```

После этого проверить:

```text
http://localhost:8080/api/orders/hello
```
