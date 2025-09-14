
from app.ml.tools.object_detector import MLObjectDetector
from app.schemas.uploadfile import Options
import queue
import threading
import time
from enum import Enum
from typing import Optional, Dict, Any

class FileStatus(Enum):
    """Статусы обработки файлов"""
    PENDING = "pending"      # В очереди, ожидает обработки
    PROCESSING = "processing" # Обрабатывается в данный момент
    COMPLETED = "completed"   # Успешно обработан
    ERROR = "error"          # Ошибка при обработке

class MLExecutor:
    """
    Класс для последовательной обработки файлов в очереди.
    Гарантирует обработку не более 1 файла одновременно.
    """
    
    def __init__(self, model_path:str = "models/yolov11m-face.pt", confidence_threshold=0.5):
        """
        Инициализация процессора файлов.
        
        Args:
            process_function: Функция для обработки файла, принимает имя файла
        """
        self.detector = MLObjectDetector(model_path, confidence_threshold=confidence_threshold)
        self.detector.initialize()
        self._file_queue = queue.Queue()
        self._file_statuses: Dict[str, Dict[str, Any]] = {}
        self._status_lock = threading.Lock()
        self._worker_thread = None
        self._running = False
        
    def start(self):
        """Запуск обработчика"""
        if self._running:
            return
            
        self._running = True
        self._worker_thread = threading.Thread(target=self._worker, daemon=True)
        self._worker_thread.start()
        
    def stop(self):
        """Остановка обработчика"""
        self._running = False
        if self._worker_thread:
            self._worker_thread.join()
            
    def add_to_queue(self, filename: str, options: 'Options') -> bool:
        """
        Добавить файл в очередь на обработку.

        Args:
            filename: Имя файла для обработки
            options: Объект Options с параметрами обработки

        Returns:
            True если файл добавлен, False если уже в обработке
        """
        with self._status_lock:
            # Проверяем, не обрабатывается ли уже этот файл
            if filename in self._file_statuses:
                status = self._file_statuses[filename]["status"]
                if status in [FileStatus.PENDING, FileStatus.PROCESSING]:
                    return False

            # Добавляем файл в очередь
            self._file_statuses[filename] = {
                "status": FileStatus.PENDING,
                "added_time": time.time(),
                "start_time": None,
                "end_time": None,
                "error": None,
                "result": None
            }
        self._file_queue.put((filename, options))
        return True
        
    def get_status(self, filename: str) -> Optional[Dict[str, Any]]:
        """
        Получить статус обработки файла.
        
        Args:
            filename: Имя файла
            
        Returns:
            Словарь с информацией о статусе или None если файл не найден
        """
        with self._status_lock:
            if filename in self._file_statuses:
                status_info = self._file_statuses[filename].copy()
                status_info["status"] = status_info["status"].value
                return status_info
            return None
            
    def get_all_statuses(self) -> Dict[str, Dict[str, Any]]:
        """
        Получить статусы всех файлов.
        
        Returns:
            Словарь со статусами всех файлов
        """
        with self._status_lock:
            result = {}
            for filename, info in self._file_statuses.items():
                status_info = info.copy()
                status_info["status"] = status_info["status"].value
                result[filename] = status_info
            return result
            
    def clear_completed(self):
        """Очистить записи о завершенных файлах"""
        with self._status_lock:
            to_remove = []
            for filename, info in self._file_statuses.items():
                if info["status"] in [FileStatus.COMPLETED, FileStatus.ERROR]:
                    to_remove.append(filename)
                    
            for filename in to_remove:
                del self._file_statuses[filename]
                
    def _worker(self):
        """Рабочий поток для обработки файлов"""
        while self._running:
            try:
                # Получаем файл из очереди с таймаутом
                filename, options = self._file_queue.get(timeout=1)
                
                # Обновляем статус на "обрабатывается"
                with self._status_lock:
                    if filename in self._file_statuses:
                        self._file_statuses[filename]["status"] = FileStatus.PROCESSING
                        self._file_statuses[filename]["start_time"] = time.time()
                
                try:
                    # Выполняем обработку файла
                    result = self.detector.process_file(
                        filename,
                        options.object_types,
                        options.intensity,
                        options.blur_type,
                    )
                    result = result.replace('\\', '/')
                    # Обновляем статус на "завершено"
                    with self._status_lock:
                        if filename in self._file_statuses:
                            self._file_statuses[filename]["status"] = FileStatus.COMPLETED
                            self._file_statuses[filename]["end_time"] = time.time()
                            self._file_statuses[filename]["result"] = result
                            
                except Exception as e:
                    # Обновляем статус на "ошибка"
                    with self._status_lock:
                        if filename in self._file_statuses:
                            self._file_statuses[filename]["status"] = FileStatus.ERROR
                            self._file_statuses[filename]["end_time"] = time.time()
                            self._file_statuses[filename]["error"] = str(e)
                
                # Отмечаем задачу как выполненную
                self._file_queue.task_done()
                
            except queue.Empty:
                # Таймаут при получении из очереди - продолжаем
                continue
            except Exception as e:
                # Неожиданная ошибка в рабочем потоке
                print(f"Ошибка в рабочем потоке: {e}")

processor = MLExecutor()
processor.start()

def get_ml_executor():
    return processor

# Пример использования
if __name__ == "__main__":

    # Создаем процессор
    processor = MLExecutor()
    processor.start()
    
    # Добавляем файлы в очередь
    files = ["image.jpg", "video.mp4"]
    for file in files:
        options = Options(blur_type="gaussian", intensity=5, object_types=[])
        success = processor.add_to_queue(file, options)
        print(f"Файл {file} {'добавлен' if success else 'уже в обработке'}")
    
    # Проверяем статусы
    time.sleep(1)
    for file in files:
        status = processor.get_status(file)
        if status:
            print(f"Статус {file}: {status['status']}")
    
    # Ждем завершения обработки
    time.sleep(10)
    
    # Проверяем финальные статусы
    print("\nФинальные статусы:")
    all_statuses = processor.get_all_statuses()
    for filename, info in all_statuses.items():
        print(f"{filename}: {info['status']}")
        if info['result']:
            print(f"  Результат: {info['result']}")
        if info['error']:
            print(f"  Ошибка: {info['error']}")
    
    processor.stop()
