# Obscura API

## Запуск
```bash
docker-compose up --build
```

## Роуты

| Method | Path | Auth | Описание |
|--------|------|------|----------|
| GET | `/` | ❌ | API info |
| GET | `/health` | ❌ | Health check |
| POST | `/api/register` | ❌ | Регистрация |
| POST | `/api/login` | ❌ | Авторизация |
| GET | `/api/user/profile` | ✅ | Профиль |
| PUT | `/api/user/profile/update` | ✅ | Обновить профиль |
| GET | `/api/user/stats` | ✅ | Статистика |
| POST | `/api/upload` | 🔶 | Загрузить файл |
| GET | `/api/files` | ✅ | Список файлов |
| GET | `/api/files/{id}` | 🔶 | Скачать файл |
| DELETE | `/api/files/{id}` | ✅ | Удалить файл |

## Авторизация
```bash
# Получить токен
curl -X POST http://localhost:8080/api/login \
  -d '{"email":"user@example.com","password":"password123"}'

# Использовать
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/user/profile
```

## Ограничения

### Анонимы
- 3 загрузки в день
- Файлы удаляются через 24ч
- Заголовки: `X-RateLimit-Remaining`

### Файлы
- JPG, PNG, GIF, WebP, BMP, TIFF, MP4, AVI, MOV, WebM, MKV, WMV, FLV
- Максимум 50 MB

### Валидация
- email: валидный формат
- password: минимум 6 символов
- name: минимум 2 символа

## Загрузка файлов
```bash
curl -X POST http://localhost:8080/api/upload -F "file=@photo.jpg"
```
