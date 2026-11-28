# Исходный код

Код размещён в каталогах компонентов; `src/` служит указателем и не дублирует исходники.

| Каталог | Назначение и точка входа |
| --- | --- |
| [Frontend](../Frontend/) | React: изделия, заказы, форма заказа. Вход `src/main.jsx`, маршруты `src/App.jsx`. Данные пока локальные: `src/mocks/mockData.js`, `src/context/OrdersContext.jsx`. |
| [JavaBackend/production-java-backend](../JavaBackend/production-java-backend/) | Java API, модули `product`, `order`, `planning`; вход `ProductionApplication.java`. DTO → сервис → DAO, маппинг MapStruct. |
| [GoBackend](../GoBackend/) | Gin API рабочих центров и состояния; вход `cmd/api/main.go`. Слои router → handler → service → repository; конфигурация и middleware в `internal/`. |
| [Миграции Java](../JavaBackend/production-java-backend/src/main/resources/db/migration/) | Схема и демонстрационные данные Java PostgreSQL; применяются Flyway при запуске JavaBackend. |
| [GoBackend/migrations](../GoBackend/migrations/) | Отдельные миграции и начальные данные Go. |
| [Scripts](../Scripts/) | Установка зависимостей, запуск, SQL и миграции. |
| [Bruno__K](../Bruno__K/) | Коллекции запросов Java, Go и интеграции. |

Целевая цепочка: React → Java → Go; Java формирует файл для 1С. На этой ветке подключение React к Java и планировщик Go ещё не реализованы. Исходников 1С в ветке нет.

[Архитектура](../docs/03-architecture.md) · [Запуск](../docs/15-local-development.md) · [Проверки](../tests/README.md).
