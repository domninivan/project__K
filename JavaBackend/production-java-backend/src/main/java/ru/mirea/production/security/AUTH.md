# JWT и доступ к заказам

Вся конфигурация безопасности находится в `src/main/resources/security-context.xml`:
провайдер проверки пароля, BCrypt, JWT, CORS, stateless filter chain и правила URL.
XML подключается через существующий `@ImportResource` приложения; Java SecurityConfig нет.
Используется Spring Security OAuth2 Resource Server для проверки Bearer JWT.

## Настройка и запуск

Для Dev Container храните `JWT_SECRET_BASE64` в `JavaBackend/production-java-backend/.env`: Base64 от
минимум 32 случайных байт. Compose передаёт переменные из этого файла сервису
`app` через `env_file`. Файл исключён из Git; шаблон — `JavaBackend/production-java-backend/.env.example`.
Dev Containers initializeCommand автоматически создаёт .env и ключ до запуска Compose. Существующий непустой ключ сохраняется. Для прямого запуска Compose порядок подготовки описан в docs/15-local-development.md.
Используйте один и тот же ключ при перезапусках: смена ключа делает прежние
токены недействительными. Не публикуйте `.env` и не добавляйте его в Git.

После изменения `.env` пересоздайте Dev Container (Dev Containers: Rebuild Container),
затем запускайте Java как обычно: `bash ./mvnw spring-boot:run` из каталога backend.
Обычный restart уже созданного контейнера не перечитывает env_file.
Java сама `.env` не читает: при запуске с Windows передайте ключ через окружение.

Пример разового запуска с Windows через PowerShell в каталоге JavaBackend
(генерирует новый временный ключ, для постоянного использования загрузите сохранённый):

```powershell
$jwtKeyBytes = [byte[]]::new(32)
[Security.Cryptography.RandomNumberGenerator]::Fill($jwtKeyBytes)
$env:JWT_SECRET_BASE64 = [Convert]::ToBase64String($jwtKeyBytes)
.\mvnw.cmd spring-boot:run
```

`JWT_TTL_SECONDS` по умолчанию 900 (15 минут), допустимо 60–86400.
`CORS_ALLOWED_ORIGIN` по умолчанию `http://localhost:5173`; при другом адресе
React задайте точный origin. Cookie-аутентификации нет; CSRF отключена для
этого Bearer-only API. При внешнем доступе используйте HTTPS.

Миграция V3 добавляет nullable `production_order.client_id` с FK и индексом.
Новые заказы получают владельца из проверенного principal. Старые заказы остаются
без владельца и доступны только планировщику. Назначение владельца выполняется
явно после проверки принадлежности; готовой автоматической привязки нет:

```sql
UPDATE production_order
SET client_id = :verified_client_id
WHERE id = :verified_order_id AND client_id IS NULL;
```

Владелец должен быть существующим CLIENT. После привязки всех исторических заказов
отдельной будущей миграцией можно сделать client_id обязательным. Старые даты,
статусы и позиции в этом блоке не мигрируются. Baseline существующей БД описан
в `docs/15-local-development.md`.

## Учётные записи

Регистрация открыта без токена, с выбором CLIENT, PLANNER, SHOP_MANAGER, WORKER или REPORTING. Это согласованный режим учебного проекта: пользователь сам выбирает свои полномочия. Первого планировщика также можно создать через API; SQL не требуется.

## Создание пользователя через API

Доступно без авторизации с любой из пяти ролей. В Postman выберите Authorization → No Auth.

`POST /api/users`, Content-Type application/json:

```json
{
  "fullName": "Иван Иванов",
  "username": "worker1",
  "password": "your-new-password",
  "role": "WORKER"
}
```

Ответ 201:

```json
{"id":7,"fullName":"Иван Иванов","username":"worker1","role":"WORKER","active":true}
```

fullName обязателен, до 200 символов; username обязателен, до 100 символов,
без пробельных символов. Логины регистрозависимы. Пароль — от 8 до 72 символов,
не более 72 байт UTF-8. role обязательна: одно из пяти значений выше.
Пароль передаётся обычной строкой по HTTPS, сервер сам создаёт BCrypt-хеш;
пароль и хеш в ответ не включаются. Новый пользователь активен и может сразу
войти через `/api/auth/login`. Изменение и удаление пользователей не добавлены.

Неверные поля или неизвестная роль — 400; занятый логин — 409 USERNAME_TAKEN. Остальные API требуют JWT и соответствующую роль.

## Контракт для React

`POST /api/auth/login` с Content-Type application/json:

```json
{"username":"client-login","password":"your-password"}
```

Успех 200:

```json
{"accessToken":"<JWT>","tokenType":"Bearer","expiresIn":900}
```

Последующие запросы отправляются с `Authorization: Bearer <JWT>`.
Токен хранить в памяти приложения. Refresh-токена нет: после истечения срока
или перезагрузки страницы войти повторно. Выход на фронте удаляет токен;
это не отзывает копии токена до истечения срока. Блокировка пользователя
или его удаление запрещает доступ сразу; роль перечитывается из БД на каждом запросе.

`GET /api/me` возвращает только `{id, fullName, role}`.
JWT содержит issuer, стабильный ID в sub, iat и exp; роль и хеш пароля в JWT не включаются.
Неверный, просроченный токен и неактивный пользователь → 401.
Действие чужой роли → 403. Чужой заказ, включая отмену, → 404.

## Права существующих API

| API | Доступ |
| --- | --- |
| POST /api/auth/login | Без токена |
| GET /api/me | Любой вошедший активный пользователь |
| POST /api/users | Без токена, регистрация с указанием любой роли |
| GET /api/products и /api/products/{id} | CLIENT, PLANNER, SHOP_MANAGER, WORKER |
| POST /api/orders | CLIENT, владелец из principal |
| GET /api/orders и /api/orders/{id} | CLIENT только свои; PLANNER все, включая старые без владельца |
| DELETE /api/orders/{id} | CLIENT только свой NEW; повтор своей CANCELLED возвращает прежний результат |
| /api/planner/** | PLANNER |
| Остальные URL | Запрещены, кроме внутреннего /error |

Списки клиента фильтруются запросом БД, а не после загрузки чужих заказов.
REPORTING имеет доступ к /api/me; endpoint отчётности пока не реализован.
Планировочные методы по-прежнему зависят от незавершённой интеграции Go.
Отмена после начала планирования закрыта ответом 409 до реализации согласованного
процесса отмены с Go. В этом блоке сохранён действующий DTO создания с dueDate;
убрать дату и добавить согласования нужно отдельным бизнес-блоком.

## Проверки

`./mvnw test` проверяет XML, JWT, login, роли, CORS и бизнес-доступ к заказам.
Для PostgreSQL-тестов задайте JAVA_MIGRATION_TEST_URL, JAVA_MIGRATION_TEST_USER,
JAVA_MIGRATION_TEST_PASSWORD. Используйте исключительно отдельную тестовую БД:
AuthPostgresTest запускает приложение на её public и создаёт тестовые данные;
MigrationTest создаёт и удаляет свои уникальные схемы.
