
# HTTP rest api
Первый проект на golang

Сервис простой авторизации

Используемые технологии: 
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Gorilla (веб фреймворк)
- golang-migrate/migrate (для миграций БД)

# Usage
[![forthebadge](https://forthebadge.com/images/featured/featured-built-with-love.svg)](https://forthebadge.com)

**Скопируйте проект**
```bash
  git clone https://github.com/artemiyKew/json-rpc-lamoda.git
```

**Перейдите в каталог проекта**
```bash
  cd http-rest-api
```

**Запустите сервер**
```bash
  make compose
```

## Examples
- [Регистрация](#регистрация)
- [Аутентификация](#аутентификация)
- [Получение данных о пользователе](#получение-данных-о-пользователе)

## Регистрация
Регистрация пользователя: 

```bash
curl -X POST \
  http://localhost:1234/sign-up \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "test@test.com",
    "password": "password"
    }'
```
Пример ответа: 
```json
{
    "id":2,
    "email":"test@test.com"
}
```

## Аутентификация
Аутентификация пользователя:
```bash
curl -X POST \
  http://localhost:1234/sign-in \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "test@test.com",
    "password": "password"
    }'
```
Пример ответа: 
```jwt
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE1MzcyMjIsInN1YiI6Mn0.HkLT91ZrAboXhJGuW4wSR1IkEaA6ezyBInlVmli9izA

```

## Получение данных о пользователе
Получение данных о пользователе:

```bash
curl -X GET \
  http://localhost:1234/private/whoami \
  -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE1MzcyMjIsInN1YiI6Mn0.HkLT91ZrAboXhJGuW4wSR1IkEaA6ezyBInlVmli9izA'
```
Пример ответа: 
```json
{
    "id":2,
    "email":"test@test.com"
}
```



