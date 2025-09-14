import os

from ultralytics import YOLO
from typing import List

class Model:
    """Модуль для работы с YOLO моделью и обработки результатов детекции"""
    
    def __init__(self, model_path: str = "models/yolov11m-face.pt", confidence_threshold: float = 0.7):
        """
        Инициализация модуля
        
        Args:
            model_path: Путь к файлу весов модели (.pt)
            confidence_threshold: Порог уверенности для фильтрации детекций
        """
        self.model_path = model_path
        self.confidence_threshold = confidence_threshold
        self.model = None
        self.class_names = None
        
    def load_model(self):
        """Загрузка модели YOLO"""
        try:
            os.environ['YOLO_VERBOSE'] = 'False'  # Глобальное отключение вывода
            self.model = YOLO(self.model_path, verbose=False)
            self.class_names = self.model.names
            print(f"Модель успешно загружена из {self.model_path}")
        except Exception as e:
            print(f"Ошибка загрузки модели: {e}")
            raise
    
    def predict(self, image_source) -> List:
        """
        Выполнение предсказания на изображении
        
        Args:
            image_source: Путь к изображению, numpy array или URL
            
        Returns:
            Список результатов детекции
        """
        if self.model is None:
            raise ValueError("Модель не загружена. Вызовите load_model() сначала")
            
        results = self.model(image_source, conf=self.confidence_threshold, verbose=False)
        return results
    
    def extract_boxes(self, results) -> List[dict]:
        """
        Извлечение информации о боксах из результатов
        
        Args:
            results: Результаты предсказания модели
            
        Returns:
            Список словарей с информацией о боксах
        """
        boxes_info = []
        
        for result in results:
            if result.boxes is not None:
                for box in result.boxes:
                    # Получаем координаты бокса
                    coords = box.xyxy[0].tolist()
                    coords = [int(x) for x in coords]
                    
                    box_info = {
                        'coordinates': coords,  # [x_min, y_min, x_max, y_max]
                        'confidence': round(box.conf.item(), 3),
                        'class_id': int(box.cls.item()),
                        'class_name': self.class_names[int(box.cls.item())]
                    }
                    boxes_info.append(box_info)
        
        return boxes_info
