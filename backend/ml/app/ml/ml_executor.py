
from app.ml.tools.object_detector import MLObjectDetector
from app.schemas.uploadfile import Options
import os
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
    Класс для параллельной обработки файлов в очереди.
    Поддерживает настраиваемое количество воркеров.
    """
    
    def __init__(self, model_path: str = "models/yolov11m-face.pt", confidence_threshold: float = 0.5):
        """
        Инициализация процессора файлов.
        
        Args:
            model_path: Путь к модели YOLO
            confidence_threshold: Порог уверенности для детекции
        """
        self.model_path = model_path
        self.confidence_threshold = confidence_threshold
        self.num_workers = self._get_worker_count()
        
        self._file_queue = queue.Queue()
        self._file_statuses: Dict[str, Dict[str, Any]] = {}
        self._status_lock = threading.Lock()
        self._workers = []
        self._running = False
        
    def _get_worker_count(self) -> int:
        """Получение количества воркеров из переменной окружения"""
        return int(os.getenv('ML_WORKERS', '3'))
        
    def start(self):
        """Запуск всех воркеров"""
        if self._running:
            return
            
        self._running = True
        self._start_workers()
        
    def stop(self):
        """Остановка всех воркеров"""
        self._running = False
        self._stop_workers()
        
    def _start_workers(self):
        """Создание и запуск воркеров"""
        for i in range(self.num_workers):
            worker = threading.Thread(
                target=self._worker,
                name=f"MLWorker-{i+1}",
                daemon=True
            )
            worker.start()
            self._workers.append(worker)
            
    def _stop_workers(self):
        """Остановка и ожидание завершения воркеров"""
        for worker in self._workers:
            if worker.is_alive():
                worker.join(timeout=5)
        self._workers.clear()
            
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
        # Каждый воркер создает свой экземпляр детектора
        detector = self._create_detector()
        worker_name = threading.current_thread().name
        
        while self._running:
            try:
                filename, options = self._file_queue.get(timeout=1)
                self._process_file(detector, filename, options, worker_name)
                self._file_queue.task_done()
                
            except queue.Empty:
                continue
            except Exception as e:
                print(f"Ошибка в воркере {worker_name}: {e}")
                
    def _create_detector(self) -> MLObjectDetector:
        """Создание экземпляра детектора для воркера"""
        detector = MLObjectDetector(
            face_model_path=self.model_path,
            confidence_threshold=self.confidence_threshold
        )
        detector.initialize()
        return detector
        
    def _process_file(self, detector: MLObjectDetector, filename: str, options: 'Options', worker_name: str):
        """Обработка одного файла"""
        self._update_status(filename, FileStatus.PROCESSING, start_time=time.time())
        
        try:
            result = detector.process_file(
                filename,
                options.object_types,
                options.intensity,
                options.blur_type,
            )
            result = result.replace('\\', '/')
            self._update_status(filename, FileStatus.COMPLETED, result=result, end_time=time.time())
            
        except Exception as e:
            self._update_status(filename, FileStatus.ERROR, error=str(e), end_time=time.time())
            
    def _update_status(self, filename: str, status: FileStatus, **kwargs):
        """Потокобезопасное обновление статуса файла"""
        with self._status_lock:
            if filename in self._file_statuses:
                self._file_statuses[filename]["status"] = status
                for key, value in kwargs.items():
                    self._file_statuses[filename][key] = value

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
