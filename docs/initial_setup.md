# Первоначальная настройка Dev Container

Инструкция перенесена из [DocsOld/initial_setup.md](../DocsOld/initial_setup.md) и сверена с конфигурацией репозитория 21.11.2026. Контейнер в рамках обновления документации не запускался.

## Назначение и требования

Dev Container создаёт общую среду разработки на Debian Trixie. На компьютере нужны Git, Docker Desktop (на Linux — Docker Engine с Compose Plugin), Visual Studio Code и расширение Microsoft Dev Containers. Отдельная установка Java, Go, Node.js и Python на хосте для этого способа работы не требуется.

Официальные страницы установки: [Docker Desktop](https://www.docker.com/products/docker-desktop/), [VS Code](https://code.visualstudio.com/), [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers). Для установки расширения откройте Extensions в VS Code (`Ctrl+Shift+X`).

Запустите Docker и проверьте в терминале хоста:

```bash
git --version
docker --version
docker compose version
docker info
```

Первые три команды показывают версии, последняя должна успешно связаться с Docker Engine. При первой настройке нужен доступ к реестрам контейнеров и зависимостей.

## Получение проекта

Если локальной копии ещё нет:

```bash
git clone https://github.com/domninivan/project__K.git
cd project__K
code .
```

Если проект уже скачан, откройте его существующий корень в VS Code: каталог с `docker-compose.yml` и `.devcontainer/`. Не выполняйте повторное клонирование поверх рабочей копии. Ветку выбирайте по своей задаче; текущая документация подготовлена на `feature/docs`.

## Создание контейнера

1. Откройте палитру команд (`Ctrl+Shift+P`).
2. Выполните `Dev Containers: Reopen in Container`.
3. Дождитесь завершения сборки и команды `postCreateCommand`; ход настройки доступен в журнале Dev Containers.
4. Откройте новый терминал VS Code внутри контейнера. Рабочий каталог — `/workspaces/project__K`, пользователь — `vscode`.

Конфигурация: [.devcontainer/devcontainer.json](../.devcontainer/devcontainer.json) и [docker-compose.yml](../docker-compose.yml). Сервис `app` содержит среду разработки, `postgres` — PostgreSQL 17. Команда `sleep infinity` удерживает `app` запущенным; приложения запускаются отдельно.

Features устанавливают Java 21 и Maven, Go и golangci-lint, Node.js/npm, Python 3.12 и Docker CLI. Для Go указан `latest`, для Node — `lts`: фактические версии следует проверить на совместимость с [go.mod](../GoBackend/go.mod) и [package.json](../Frontend/package.json).

## Автоматическая подготовка

После создания контейнера выполняется [.devcontainer/post-create.sh](../.devcontainer/post-create.sh):

1. Настраивает Git safe.directory и права на каталоги кэшей и зависимостей, делает shell-скрипты в `Scripts/` исполняемыми.
2. Устанавливает PostgreSQL Client, `migrate`, Air и Bruno CLI, если они отсутствуют.
3. Копирует `GoBackend/.env.example` в `GoBackend/.env`.
4. Вызывает `Scripts/install-dependencies.sh`: загружает зависимости Go, npm и Maven.
5. Проверяет PostgreSQL, создаёт отдельную БД `go_db`, если её ещё нет, и применяет Go-миграции через `Scripts/migrate-go-helper.sh up`.
6. Выводит версии инструментов и сообщение `Development environment is ready!` при успешном завершении.

Скрипт при повторном выполнении **перезаписывает `GoBackend/.env`** значениями из примера. Перед пересборкой сохраните свои локальные настройки вне репозитория. Переменные из Compose имеют приоритет над значениями `.env` при загрузке Go-конфигурации.

Для повторной установки только зависимостей, без повторного запуска всей настройки:

```bash
bash Scripts/install-dependencies.sh
```

## Базы данных

| Параметр | Внутри Dev Container | С хоста |
| --- | --- | --- |
| Адрес PostgreSQL | `postgres:5432` | `localhost:5433` |
| БД Java | `postgres` | `postgres` |
| БД Go | `go_db` | `go_db` |
| Пользователь и пароль учебного Compose | `postgres` / `postgres` | `postgres` / `postgres` |

Java получает `DB_URL`, `DB_USERNAME`, `DB_PASSWORD`; Go — `GO_DB_HOST`, `GO_DB_PORT`, `GO_DB_USER`, `GO_DB_PASSWORD`, `GO_DB_NAME` из Compose. Это демонстрационные настройки локальной среды.

При первой инициализации пустого тома PostgreSQL выполняет `database/01-schema.sql` и `database/02-seed.sql` для Java. Go-миграции применяет post-create отдельно. Редактирование Java SQL-файлов не изменяет уже созданную БД автоматически.

Проверка из терминала контейнера (psql запросит пароль):

```bash
pg_isready -h postgres -p 5432 -U postgres
psql -h postgres -U postgres -d postgres -c '\dt'
psql -h postgres -U postgres -d go_db -c '\dt'
bash Scripts/migrate-go-helper.sh version
```

## Запуск приложений

В терминале Dev Container из корня проекта:

```bash
bash Scripts/app-start.sh
```

Скрипт запускает Go через Air, Java через Maven Wrapper и Frontend через Vite. Требуется исполняемый Java Wrapper; если запуск сообщает об отсутствии права, выполните внутри контейнера:

```bash
chmod +x JavaBackend/production-java-backend/mvnw
```

Не закрывайте терминал во время работы. `Ctrl+C` останавливает запущенные скриптом приложения; PostgreSQL продолжает работать отдельно. При завершении одного приложения скрипт останавливает остальные.

| Компонент | Порт | Минимальная проверка |
| --- | --- | --- |
| Frontend | 5173 | Открыть интерфейс изделий и заказов |
| Java | 8080 | `GET /api/orders/hello` |
| Go | 8081 | `GET /api/v1/live`, `GET /api/v1/ready` |

В Compose опубликован только порт PostgreSQL. Если приложения недоступны с хоста, откройте вкладку **Ports** в VS Code и добавьте проброс 5173, 8080 и 8081. Используйте выданные VS Code локальные адреса: при занятом порте хоста они могут отличаться. Vite внутри контейнера использует фиксированный порт 5173 (`strictPort: true`).

Проверки из терминала контейнера:

```bash
curl --fail http://localhost:8080/api/orders/hello
curl --fail http://localhost:8081/api/v1/live
curl --fail http://localhost:8081/api/v1/ready
```

Готовность окружения и health-ответы не подтверждают сквозное планирование. React пока использует локальные данные, контракт Java–Go расходится; см. [архитектуру](03-architecture.md). Полная [программа проверок](../tests/README.md) выполняется отдельно.

## Повторное открытие и пересборка

Для обычной работы снова откройте проект через `Dev Containers: Reopen in Container`. При изменении конфигурации используйте `Dev Containers: Rebuild Container`; для сборки без кэша — `Dev Containers: Rebuild Container Without Cache`.

Пересборка не равнозначна очистке именованных томов: в них сохраняются данные PostgreSQL, npm-зависимости, Go/Maven-кэши и Java target. Сборка без кэша также не означает очистку этих томов. Повторный post-create применяет новые Go-миграции и перезаписывает локальный `.env`, как описано выше.

## Если настройка не завершилась

- Docker недоступен: проверьте Docker Engine командой `docker info` на хосте.
- Не загрузился образ или зависимость: найдите первую ошибку в журнале Dev Containers, проверьте доступ к соответствующему реестру; после устранения причины повторите настройку.
- PostgreSQL не готов: из корня репозитория на хосте выполните `docker compose ps` и `docker compose logs postgres`. Проверьте, свободен ли порт 5433.
- Нет `go_db` или Go-таблиц: проверьте завершение post-create и его шаги создания БД/миграций; отдельный `migrate ... up` не создаёт саму базу.
- Не запускается компонент: проверьте его журнал, доступность `air`, `go`, `npm`, `java` и исполняемость Java Wrapper. Предупреждение об ожидании порта не является успешной проверкой запуска.

[Повседневная разработка и запуск с хоста](15-local-development.md) · [Инструменты](../tools/README.md) · [Оглавление](README.md).
