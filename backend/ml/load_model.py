import urllib.request
import os

def download_file(url, file_path):
    """Скачивает файл с указанного URL и сохраняет по указанному пути"""
    
    # Создаем директорию, если её нет
    directory = os.path.dirname(file_path)
    if directory:
        os.makedirs(directory, exist_ok=True)
    
    try:
        # Создаем запрос с User-Agent заголовком
        request = urllib.request.Request(
            url, 
            headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'}
        )
        
        # Открываем URL и читаем содержимое
        print("Начинаем скачивание...")
        with urllib.request.urlopen(request) as response:
            content = response.read()
        
        # Сохраняем файл
        with open(file_path, 'wb') as file:
            file.write(content)
        
        print(f"Файл успешно скачан: {file_path}")
        print(f"Размер файла: {len(content)} байт")
        
    except Exception as e:
        print(f"Ошибка при скачивании: {e}")

# Использование
url = 'https://github.com/akanametov/yolo-face/releases/download/v0.0.0/yolov11m-face.pt'
file_path = os.path.join('models', 'yolov11m-face.pt')

download_file(url, file_path)
