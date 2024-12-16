# Test task BackDev

## Используемые зависимости:

- Go 1.23.4
- PostgreSQL 16.3
- Docker
- UUID
- Air (hot reloading) `https://github.com/air-verse/air`
- Chi (routing) `https://github.com/go-chi/chi`
- JWT (access token) `https://github.com/golang-jwt/jwt`
- Godotenv (.env variables) `https://github.com/joho/godotenv`
- Bcrypt (refresh token hash/compare) `https://pkg.go.dev/golang.org/x/crypto/bcrypt`

## Самый простой способ запустить проект (если нет docker, то нужно настроить подключение к базе данных в .env):

```bash
docker-compose up --build -d
```

## Для упрощения работы с проектом используется Makefile (если нет make, то можно просто запустить команду из Makefile в терминале)

### Работа с миграциями осуществляется с помощью migrate (`https://github.com/golang-migrate/migrate`) примеры команд:

```bash
make migrate-create <name>
make migrate-up
make migrate-down
```

### Исходя из постановки задачи, регистрация пользователей не предусмотрено, поэтому для наполнения базы данных используется seed:

```bash
make seed
```

### Оповещение о смене IP реализовано моковым методом mockSendEmail, который выводит в консоль сообщение.

## API

### По маршруту /tokens необходимо указать в query параметре user_id, который является id пользователя, которому нужно выдать токены. Access token выдается в теле ответа, refresh token устанавливается в cookie.

### refresh token имеет флаг HttpOnly, поэтому его нельзя будет изменить из клиента.

```bash
curl -X GET http://localhost:8080/api/auth/tokens?user_id=<user_id>
```

## Для защиты маршрута на /refresh эндпоинт используется middleware, который проверяет наличие access token в заголовке Authorization, если токен не валидный или не предоставлен, то возвращается ошибка 401 Unauthorized, так же если токен истек, то возвращается ошибка 401 Unauthorized.
