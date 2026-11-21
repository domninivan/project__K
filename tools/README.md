# Инструменты разработки

`tools/` содержит справочники [labels.json](labels.json) и [milestones.json](milestones.json), унаследованные от шаблона. Они не создают метки/вехи в GitVerse и не являются официальным форматом импорта. Согласованный календарь ведётся в [плане работ](../docs/12-work-plan.md); значения JSON не подтверждают сроки или готовность продукта.

Рабочие скрипты расположены в [Scripts](../Scripts/), а не в этом каталоге. Для них нужен Bash/Linux (Dev Container или совместимая среда), не нативный PowerShell.

| Скрипт / инструмент | Назначение и условия |
| --- | --- |
| `bash Scripts/install-dependencies.sh` | Загружает Go, npm и Maven зависимости. Требует доступных Go, Node/npm, Java и Maven Wrapper. |
| `bash Scripts/app-start.sh` | Запускает Go через `air`, Java через Maven Wrapper и Vite. Нужны зависимости, БД, `air`, Bash и `setsid`. Ctrl+C останавливает группы процессов. |
| `bash Scripts/migrate-go-helper.sh version` / `up` | Проверяет версию / применяет миграции Go; нужен `migrate`, GO_DB_* или `GoBackend/.env`. |
| `bash Scripts/run-sql.sh путь.sql` | Выполняет SQL через psql; параметры DB_HOST/PORT/USERNAME/PASSWORD/NAME. Значения по умолчанию рассчитаны на Compose-сеть. |
| [Bruno__K](../Bruno__K/) | Коллекции ручных HTTP-запросов. Локальные адреса необходимо сверять с текущими маршрутами. |
| [.devcontainer](../.devcontainer/) | Среда разработки Debian, Java, Go, Node, Python; использует корневой Compose. |

`down`, `down-all`, `force` у помощника миграций изменяют состояние БД; применяйте только осознанно на своей тестовой базе. SQL и миграции выполняет владелец соответствующего сервиса. Данные Compose и миграций демонстрационные; секреты в примеры не добавляются.

[Первоначальная настройка](../docs/initial_setup.md) · [Запуск](../docs/15-local-development.md) · [Проверки](../tests/README.md) · [Правила работы](../CONTRIBUTING.md).
