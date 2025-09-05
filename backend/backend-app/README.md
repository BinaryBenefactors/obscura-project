# 🧪 Obscura API Testing Guide

## 📋 **Простая схема работы:**

```
1. Фронтенд отправляет файл → POST /api/upload
2. Бэкенд возвращает fileId и статус "processing"  
3. Фронтенд делает polling → GET /api/files/{fileId} каждые 2-3 сек
4. Когда статус "completed" → показывает кнопку скачивания
5. Скачивание → GET /api/files/{fileId}?type=processed
```

## 🚀 **Пошаговое тестирование:**

### **1. Регистрация/Авторизация**

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123", 
    "name": "Test User"
  }'
```

**Авторизация:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```
**Ответ:** Сохраните `token` для дальнейших запросов!

### **2. Загрузка файла**

**Авторизованный пользователь (полные параметры):**
```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "file=@/path/to/image.jpg" \
  -F "blur_type=gaussian" \
  -F "intensity=7" \
  -F "object_types=face,person"
```

**Авторизованный пользователь (минимальные параметры):**
```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "file=@/path/to/image.jpg"
```

**Анонимный пользователь:**
```bash  
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/image.jpg" \
  -F "intensity=5" \
  -F "object_types=face"
```

**Ответ:**
```json
{
  "message": "File uploaded successfully and processing started",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "original_name": "image.jpg",
    "status": "processing", 
    "uploaded_at": "2024-01-15T09:00:00Z"
  }
}
```

### **3. Проверка статуса обработки (Polling)**

**Для авторизованного пользователя:**
```bash
curl -X GET http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Статус "processing":**
```json
{
  "message": "File info retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing",
    "original_name": "image.jpg",
    "uploaded_at": "2024-01-15T09:00:00Z"
  }
}
```

**Статус "completed":**
```json
{
  "message": "File info retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000", 
    "status": "completed",
    "original_name": "image.jpg",
    "processed_name": "550e8400-e29b-41d4-a716-446655440000_processed.jpg",
    "file_size": 1048576,
    "processed_size": 987654,
    "uploaded_at": "2024-01-15T09:00:00Z",
    "processed_at": "2024-01-15T09:02:30Z"
  }
}
```

**Статус "failed":**
```json
{
  "message": "File info retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "failed",
    "original_name": "image.jpg",
    "error_message": "No objects detected for blurring",
    "uploaded_at": "2024-01-15T09:00:00Z"
  }
}
```

### **4. Скачивание файлов**

**Скачать оригинал:**
```bash
curl -X GET "http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000?type=original" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -o original_image.jpg
```

**Скачать обработанный (только если status=completed):**
```bash
curl -X GET "http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000?type=processed" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -o processed_image.jpg
```

**Анонимный пользователь (скачивание обработанного без токена):**
```bash
curl -X GET "http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000?type=processed" \
  -o processed_image.jpg
```

### **5. Список файлов пользователя**

```bash
curl -X GET http://localhost:8080/api/files \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Ответ:**
```json
{
  "message": "Files retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "original_name": "photo1.jpg",
      "status": "completed",
      "uploaded_at": "2024-01-15T09:00:00Z",
      "processed_at": "2024-01-15T09:02:30Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "original_name": "photo2.jpg", 
      "status": "processing",
      "uploaded_at": "2024-01-15T09:05:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "original_name": "photo3.jpg", 
      "status": "failed",
      "error_message": "Unsupported file format",
      "uploaded_at": "2024-01-15T09:07:00Z"
    }
  ]
}
```

### **6. Удаление файла**

```bash
curl -X DELETE http://localhost:8080/api/files/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Ответ:**
```json
{
  "message": "File deleted successfully"
}
```

## ⚙️ **Параметры обработки (упрощенные):**

```bash
-F "blur_type=gaussian"           # gaussian, motion, pixelate (default: gaussian)
-F "intensity=7"                  # 1-10 (default: 5)
-F "object_types=face,person"     # face,person,car,plate,text,logo
```

**Примеры комбинаций:**
```bash
# Размытие лиц средней интенсивности
-F "blur_type=gaussian" -F "intensity=5" -F "object_types=face"

# Сильная пикселизация людей и машин
-F "blur_type=pixelate" -F "intensity=9" -F "object_types=person,car"

# Motion blur номерных знаков
-F "blur_type=motion" -F "intensity=6" -F "object_types=plate"

# Минимальные параметры (используются defaults)
# blur_type=gaussian, intensity=5, object_types=все доступные
```

## 💻 **Frontend JavaScript Examples:**

### **Авторизованный пользователь (с polling):**
```javascript
// 1. Загрузка файла с параметрами
const formData = new FormData();
formData.append('file', fileInput.files[0]);
formData.append('blur_type', 'gaussian');
formData.append('intensity', '7');
formData.append('object_types', 'face,person');

const uploadResponse = await fetch('/api/upload', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: formData
});

const { data: fileInfo } = await uploadResponse.json();
const fileId = fileInfo.id;

// 2. Показать статус "Загружено, обрабатывается..."
document.getElementById('status').textContent = '⏳ Обрабатывается...';
document.getElementById('progress').style.display = 'block';

// 3. Polling статуса каждые 2.5 секунды
const pollStatus = async () => {
  try {
    const response = await fetch(`/api/files/${fileId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const { data } = await response.json();
    
    console.log('Status:', data.status);
    
    if (data.status === 'completed') {
      // Обработка завершена успешно
      document.getElementById('status').textContent = '✅ Обработка завершена!';
      document.getElementById('progress').style.display = 'none';
      showDownloadButtons(fileId);
      
    } else if (data.status === 'failed') {
      // Ошибка обработки
      document.getElementById('status').textContent = `❌ Ошибка: ${data.error_message}`;
      document.getElementById('progress').style.display = 'none';
      
    } else if (data.status === 'processing') {
      // Продолжить polling
      document.getElementById('status').textContent = '⏳ Обрабатывается...';
      setTimeout(pollStatus, 2500);
    }
    
  } catch (error) {
    console.error('Polling error:', error);
    document.getElementById('status').textContent = '❌ Ошибка получения статуса';
  }
};

// Запуск polling
pollStatus();

// 4. Функции скачивания
function showDownloadButtons(fileId) {
  const buttonsDiv = document.getElementById('downloadButtons');
  buttonsDiv.innerHTML = `
    <button onclick="downloadOriginal('${fileId}')">⬇️ Скачать оригинал</button>
    <button onclick="downloadProcessed('${fileId}')">⬇️ Скачать обработанный</button>
  `;
  buttonsDiv.style.display = 'block';
}

function downloadOriginal(fileId) {
  const link = document.createElement('a');
  link.href = `/api/files/${fileId}?type=original`;
  link.target = '_blank';
  link.click();
}

function downloadProcessed(fileId) {
  const link = document.createElement('a');
  link.href = `/api/files/${fileId}?type=processed`;
  link.target = '_blank';
  link.click();
}
```

### **Анонимный пользователь (упрощенный):**
```javascript
// 1. Загрузка файла без авторизации
const formData = new FormData();
formData.append('file', fileInput.files[0]);
formData.append('intensity', '6');
formData.append('object_types', 'face');

const uploadResponse = await fetch('/api/upload', {
  method: 'POST',
  body: formData // БЕЗ Authorization header
});

const { data: fileInfo } = await uploadResponse.json();
const fileId = fileInfo.id;

// 2. Показать прогресс (эмуляция времени обработки)
document.getElementById('status').textContent = '⏳ Обрабатывается...';

let progress = 0;
const progressInterval = setInterval(() => {
  progress += 20;
  document.getElementById('progressBar').style.width = progress + '%';
  
  if (progress >= 100) {
    clearInterval(progressInterval);
    
    // 3. Показать готовность к скачиванию
    document.getElementById('status').textContent = '✅ Готово!';
    document.getElementById('downloadProcessed').onclick = () => {
      window.open(`/api/files/${fileId}?type=processed`);
    };
    document.getElementById('downloadProcessed').style.display = 'block';
  }
}, 900); // ~4.5 секунд общее время
```

### **Получение списка файлов пользователя:**
```javascript
const getFilesList = async () => {
  try {
    const response = await fetch('/api/files', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const { data: files } = await response.json();
    
    const filesList = document.getElementById('filesList');
    filesList.innerHTML = '';
    
    files.forEach(file => {
      const fileDiv = document.createElement('div');
      fileDiv.className = 'file-item';
      
      let statusIcon = '';
      let actions = '';
      
      switch(file.status) {
        case 'completed':
          statusIcon = '✅';
          actions = `
            <button onclick="downloadOriginal('${file.id}')">Оригинал</button>
            <button onclick="downloadProcessed('${file.id}')">Обработанный</button>
            <button onclick="deleteFile('${file.id}')">Удалить</button>
          `;
          break;
        case 'processing':
          statusIcon = '⏳';
          actions = '<span>Обрабатывается...</span>';
          break;
        case 'failed':
          statusIcon = '❌';
          actions = `<span>Ошибка: ${file.error_message}</span>`;
          break;
        default:
          statusIcon = '📁';
      }
      
      fileDiv.innerHTML = `
        <div class="file-info">
          ${statusIcon} <strong>${file.original_name}</strong>
          <small>${new Date(file.uploaded_at).toLocaleString()}</small>
        </div>
        <div class="file-actions">${actions}</div>
      `;
      
      filesList.appendChild(fileDiv);
    });
    
  } catch (error) {
    console.error('Error loading files:', error);
  }
};

// Функция удаления файла
const deleteFile = async (fileId) => {
  if (!confirm('Удалить файл?')) return;
  
  try {
    await fetch(`/api/files/${fileId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    // Обновить список
    getFilesList();
  } catch (error) {
    alert('Ошибка удаления файла');
  }
};
```

## 🔍 **Системные проверки:**

**Health check:**
```bash
curl http://localhost:8080/health
```

**Swagger UI:**
```
http://localhost:8080/swagger/
```

**Статистика пользователя:**
```bash
curl -X GET http://localhost:8080/api/user/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Admin статистика:**
```bash
curl http://localhost:8080/api/admin/stats
```

## 🧪 **Тестирование Rate Limiting (для анонимных):**

```bash
# Отправить 4 файла подряд (лимит 3/24ч)
for i in {1..4}; do
  echo "Upload $i:"
  curl -X POST http://localhost:8080/api/upload \
    -F "file=@/path/to/image.jpg" \
    -F "intensity=$i"
  echo -e "\n---"
done
```

**Ожидаемый результат:** первые 3 загрузки успешны, 4-я вернет ошибку 429.

---

## ✅ **Чеклист тестирования:**

- [ ] **Регистрация/Логин** → получение JWT токена
- [ ] **Загрузка файла** → возврат fileId и статус "processing" 
- [ ] **Polling статуса** → статус меняется на "completed" через ~3-5 сек
- [ ] **Скачивание** → можно загрузить оригинал и обработанный файл
- [ ] **Список файлов** → показывает историю с корректными статусами
- [ ] **Удаление файла** → файл пропадает из списка и с диска
- [ ] **Анонимная загрузка** → работает без регистрации
- [ ] **Rate Limiting** → блокирует 4-й файл анонимного пользователя
- [ ] **Параметры обработки** → blur_type, intensity, object_types корректно обрабатываются
- [ ] **Health check** → система здорова
- [ ] **Swagger** → документация доступна

## 🎯 **Примеры curl для быстрого тестирования:**

```bash
# Быстрый тест с минимальными параметрами
curl -X POST http://localhost:8080/api/upload -F "file=@test.jpg"

# Полный тест с авторизацией
TOKEN="your-jwt-token"
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test.jpg" \
  -F "blur_type=pixelate" \
  -F "intensity=8" \
  -F "object_types=face,person"
```

**Всё готово к тестированию! 🚀**