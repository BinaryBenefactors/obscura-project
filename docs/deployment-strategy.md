# Развёртывание проекта Obscura

В этом документе описано, как развернуть проект Obscura в различных окружениях.

## Локальное развертывание

Для локального развертывания используется Docker Compose.

### Требования
- Docker
- Docker Compose

### Запуск

```bash
docker-compose up -d
```

Эта команда запустит все сервисы проекта:
- Бэкенд приложение (Golang)
- Фронтенд приложение (React в Nginx)
- Nginx как reverse proxy
- PostgreSQL базу данных

### Остановка

```bash
docker-compose down
```

### Просмотр логов

```bash
docker-compose logs -f
```

## Production развертывание

Production развертывание осуществляется через GitHub Actions при пуше в ветку `master`.

### Что происходит при развертывании
1. Развертывание на сервере через SSH

### Требуемые секреты GitHub Actions
- `SSH_HOST` - адрес сервера
- `SSH_USER` - пользователь для подключения
- `SSH_PRIVATE_KEY` - приватный SSH ключ

### Ручное развертывание

Если необходимо развернуть вручную на сервере:

```bash
git pull origin master
docker-compose down
docker-compose up -d --build
```

## Конфигурация

### Переменные окружения

Основные переменные окружения для бэкенда:
- `DB_HOST` - хост базы данных
- `DB_USER` - пользователь базы данных
- `DB_PASSWORD` - пароль базы данных
- `DB_NAME` - имя базы данных
- `DB_PORT` - порт базы данных
- `JWT_SECRET` - секрет для JWT токенов

Эти переменные задаются в `docker-compose.yml` файле.