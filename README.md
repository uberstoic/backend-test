# Marketplace REST API

Проект реализует REST-API для условного маркетплейса на Go.

## Функционал

- Регистрация и авторизация пользователей (JWT).
- Создание объявлений.
- Просмотр ленты объявлений с пагинацией, сортировкой и фильтрацией.

## Технологии

- **Язык:** Go
- **Фреймворк:** Gin
- **База данных:** PostgreSQL
- **Контейнеризация:** Docker, Docker Compose

## Запуск проекта

1.  Убедитесь, что у вас установлен Docker и Docker Compose.
2.  Создайте файл `.env` в корне проекта и заполните его? примерно вот так:
    ```
    POSTGRES_USER=user
    POSTGRES_PASSWORD=password
    POSTGRES_DB=marketplace
    JWT_SECRET=your_super_secret_key
    ```
3.  Выполните команду для запуска:
    ```bash
    docker-compose up --build
    ```

API будет доступен по адресу `http://localhost:8080`.

## Краткая справка по эндпоинтам (localhost):

Регистрация нового пользователя.
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{"login": "biba", "password": "boba123"}'
```

Получение JWT токена.
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{"login": "biba", "password": "boba123"}'
```

Создание объявления (требует Authorization: Bearer <token>).
```bash
curl -X POST http://localhost:8080/api/v1/ads \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <token>" \
-d '{"title": "My New Ad", "text": "A great description of my ad.", "image_url": "http://example.com/image.jpg", "price": 150.50}'
```

Получение списка объявлений.
```bash
curl -X GET http://localhost:8080/api/v1/ads \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <token>"
```

## Тесты:
```bash
cd tests && go test -v 
```

## Команды для ендпоинтов через хостинг (Railway)

Регистрация нового пользователя.
```bash
curl -X POST https://backend-test-production-b721.up.railway.app/api/v1/auth/register -H "Content-Type: application/json" -d '{"login": "biba", "password": "boba123"}'
```

Получение JWT токена.
```bash
curl -X POST https://backend-test-production-b721.up.railway.app/api/v1/auth/login -H "Content-Type: application/json" -d '{"login": "biba", "password": "boba123"}'
```

Создание объявления (требует Authorization: Bearer <token>).
```bash
curl -X POST https://backend-test-production-b721.up.railway.app/api/v1/ads -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"title": "My New Ad", "text": "A great description of my ad.", "image_url": "http://example.com/image.jpg", "price": 150.50}'
```

Получение списка объявлений.
```bash
curl -X GET https://backend-test-production-b721.up.railway.app/api/v1/ads -H "Content-Type: application/json" -H "Authorization: Bearer <token>"
```
