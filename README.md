# MerchShop

Учебный проект на Go. REST API сервиса для внутреннего магазина мерча.  
Реализованы авторизация по JWT, переводы монет между пользователями, покупка товаров, учёт заказов и истории транзакций.

## Стек
- Go 1.24
- Gin
- PostgreSQL (pgx)
- Docker / docker-compose
- bcrypt / JWT
- testify / gomock

## Основные задачи
1. Настроить слоистую архитектуру (handlers, repo, model, middleware).
2. Сделать корректную работу с транзакциями в PostgreSQL.
3. Реализовать JWT-аутентификацию и middleware проверки токенов.
4. Покрыть код юнит-тестами.
5. Настроить миграции, Makefile и docker-compose для локального запуска.

## Проблемы, которые пришлось решать
- Конфликты транзакций при параллельных переводах (решено через `Serializable` + retry).
- Сложности с моками интерфейсов (добавил `go:generate mockgen`).
- JWT-токен с проверкой `iss`, `aud` и временем жизни.
- Корректная обработка ошибок в хендлерах (возврат JSON без паник).
- Тестирование через `httptest` и `gomock`.

## Тесты
- Покрытие кода по проекту: ~65–70%.
- Есть юнит-тесты для хендлеров (`auth`, `buy`, `send`, `info`).
- Проверяются все ветки: валидные/невалидные запросы, недостающие поля, ошибки репозитория.
- Интеграционные тесты планируются отдельно.

```bash
	github.com/6ermvH/MerchShop/internal/db		coverage: 0.0% of statements
ok  	github.com/6ermvH/MerchShop/internal/hasher	(cached)	coverage: 100.0% of statements
ok  	github.com/6ermvH/MerchShop/internal/http/handlers	(cached)	coverage: 92.7% of statements
ok  	github.com/6ermvH/MerchShop/internal/http/middleware	(cached)	coverage: 52.0% of statements
	github.com/6ermvH/MerchShop/internal/jwtutil		coverage: 0.0% of statements
	github.com/6ermvH/MerchShop/internal/logx		coverage: 0.0% of statements
?   	github.com/6ermvH/MerchShop/internal/model	[no test files]
	github.com/6ermvH/MerchShop/internal/repo		coverage: 0.0% of statements
```

## Запуск
```bash
make up
