# Локальная разработка

Инструкция сверена с конфигурацией ветки 21.11.2026; запуск в рамках переноса документации не выполнялся.

Первый запуск через Dev Container, установка программ, автоматическая подготовка БД и пересборка описаны в [первоначальной настройке](initial_setup.md). Здесь приведены команды повседневной работы и настройки запуска с хоста.

Нужны Java 21, Go версии из `GoBackend/go.mod`, Node/npm совместимые с `Frontend/package.json`, Bash, psql. Для общего скрипта запуска нужен `air`, для Go-миграций — `migrate`. Настройки и установка среды описаны в `.devcontainer/post-create.sh`; при ошибке установки проверяйте доступность этих инструментов.

После успешного post-create зависимости и Go-миграции уже установлены. Для обычного запуска в контейнере из корня:

```bash
bash Scripts/app-start.sh
```

После изменения зависимостей выполните `bash Scripts/install-dependencies.sh`, после добавления Go-миграций — `bash Scripts/migrate-go-helper.sh up`.

Перед созданием Dev Container `initializeCommand` автоматически подготавливает
`JavaBackend/production-java-backend/.env` из шаблона и генерирует JWT_SECRET_BASE64 из 32 случайных байт.
Непустой существующий ключ и остальные настройки сохраняются. Подготовка выполняется
одноразовым Docker-контейнером с тем же базовым образом, до чтения env_file в Compose;
Node.js, Python и Bash на хосте не нужны. Файл исключён из Git, ключ в лог не выводится.
После ручного изменения переменных пересоздайте Dev Container.

При прямом запуске Compose без Dev Containers сначала выполните подготовку из корня
проекта (PowerShell; из Bash замените `${PWD}` на `$(pwd)`):

```powershell
docker run --rm --user 0:0 --mount "type=bind,source=${PWD},target=/workspace" --workdir /workspace --entrypoint bash mcr.microsoft.com/devcontainers/base:trixie .devcontainer/init-java-env.sh
```

Проверка подготовки на временных файлах: в таком же контейнере вместо
`.devcontainer/init-java-env.sh` запустите `tests/init-java-env.sh`.

Go-миграции выполняются в отдельной `go_db`. Java-схему и демонстрационные данные создаёт только Flyway при запуске JavaBackend, до Hibernate validate. Скрипты находятся в `JavaBackend/production-java-backend/src/main/resources/db/migration`; Docker создаёт только саму БД. Применённые миграции не редактировать — добавлять новые версии. Для запуска Java нужен `JWT_SECRET_BASE64`; пользователи регистрируются без токена через POST /api/users с выбором роли.

Общий `app-start.sh` автоматически выбирает baseline только для проверенной исходной V1 без истории Flyway. Пустая БД и БД с историей запускаются без baseline. При несовпадении структуры скрипт останавливается до запуска сервисов. Удалять Docker-том для перехода не требуется.

### Регистрация пользователей

Регистрация пользователей: `POST /api/users` доступен без авторизации с любой
из ролей CLIENT, PLANNER, SHOP_MANAGER, WORKER, REPORTING. Это выбранный режим
учебного проекта: регистрирующийся сам выбирает свои полномочия. Другие API
по-прежнему требуют JWT и соответствующую роль; доступ к чужим заказам закрыт.
Первого планировщика можно создать этим же запросом, SQL не требуется.

В Postman выберите No Auth, Body → raw → JSON, Content-Type application/json:

```json
{
  "fullName": "Планировщик",
  "username": "planner1",
  "password": "PlannerTest123!",
  "role": "PLANNER"
}
```

Ответ — 201 с id, fullName, username, role, active. Пароль сохраняется как BCrypt-хеш
и не возвращается. Затем `POST /api/auth/login` с username и password выдаёт JWT.
Пароль — от 8 символов, максимум 72 байта UTF-8. Ошибки полей/неизвестная роль — 400,
повтор логина — 409 USERNAME_TAKEN. Токен на регистрацию не передавайте; невалидный
Bearer-токен в заголовке всё равно будет отклонён фильтром.

### Автоматический baseline Java

Перед запуском сервисов `Scripts/java-baseline-check.sh` выполняет только чтение БД через psql.
Он сравнивает структуру public с эталоном `Scripts/java-v1.signature`, полученным из
неизменённой V1 на PostgreSQL 17: таблицы, колонки, ограничения, индексы,
последовательности и дополнительные объекты. Данные и текущие значения счётчиков
не влияют на сравнение. Только точное совпадение без истории разрешает baseline версии 1.
Результат передаётся в JAVA_DB_BASELINE_ON_MIGRATE только процессу Java, без записи в .env.

Используются DB_URL, DB_USERNAME, DB_PASSWORD с теми же значениями по умолчанию,
что в database.properties. Поддерживается JDBC URL вида jdbc:postgresql://host:port/database
без query-параметров и встроенных учётных данных. При иных переопределениях настроек Spring
нужно согласовать подключение с этими переменными. Ошибка подключения или неизвестная
структура прекращает общий запуск. Наличие истории не скрывает ошибки Flyway: их проверяет
сам Flyway при старте Java. Формат каталога другой версии PostgreSQL может привести
к консервативному отказу; не включайте baseline без проверки схемы.

При прямом запуске Maven автоматической проверки нет. Для заведомо исходной V1 без
истории после проверки БД можно однократно выполнить
`JAVA_DB_BASELINE_ON_MIGRATE=true ./mvnw spring-boot:run`. Для пустой БД и БД с историей
это не требуется. SQL V1 и эталон не менять при добавлении новых миграций.

Проверка helper на отдельном тестовом PostgreSQL с правом CREATEDB:
задайте BASELINE_TEST_ADMIN_URL, DB_USERNAME, DB_PASSWORD и выполните
`bash tests/java-baseline-check.sh`. Тест создаёт и удаляет только свою тестовую БД.

| Параметры | Внутри Compose | Запуск с Windows-хоста |
| --- | --- | --- |
| PostgreSQL | `postgres:5432` | `localhost:5433` |
| Java | DB_URL, DB_USERNAME, DB_PASSWORD из Compose | Задать DB_URL на `jdbc:postgresql://localhost:5433/postgres` и свои учётные параметры |
| Go | GO_DB_* из Compose, БД `go_db` | По `GoBackend/.env.example`, заменить GO_DB_HOST на localhost и GO_DB_PORT на 5433 |
| Java → Go | `go.base-url=http://localhost:8081`, когда оба процесса в app | Адрес процесса Go; по умолчанию localhost:8081 |

Java читает параметры из ресурсов и окружения; Go — из окружения и локального `.env` (через godotenv). Для Bash-помощника миграций можно создать локальный `GoBackend/.env` по примеру. Не публикуйте локальные секреты. При запуске с хоста создайте также отдельную БД `go_db`, затем примените Go-миграции; в Dev Container это делает post-create. С хоста сначала можно запустить только БД командой `docker compose up -d postgres` из корня.

Команды запуска по отдельности из соответствующего каталога:

| Каталог | Команда |
| --- | --- |
| `GoBackend` | `go run ./cmd/api` (переменные окружения уже заданы) |
| `JavaBackend/production-java-backend` | Windows: `.\mvnw.cmd spring-boot:run`; Linux: `bash ./mvnw spring-boot:run` |
| `Frontend` | `npm install`, затем `npm run dev` |

Java: `http://localhost:8080`; Go: `http://localhost:8081`; Vite (фиксированный порт): `http://localhost:5173` (точный адрес в выводе Vite). Проверки Go: `/api/v1/live` и `/api/v1/ready`. Готовность порта не подтверждает сквозной сценарий. React пока работает локально, маршруты Java–Go расходятся: [известные ограничения](03-architecture.md).

Для доступа с хоста к приложениям контейнера используйте вкладку Ports в VS Code: [проброс портов](initial_setup.md#запуск-приложений). В Compose явно опубликован только PostgreSQL.
