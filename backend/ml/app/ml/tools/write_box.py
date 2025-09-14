from typing import List
import cv2
import numpy as np


class BoxProcessor:
    """Класс для обработки и визуализации боксов с помощью OpenCV"""

    @staticmethod
    def draw_boxes(
        image: np.ndarray,
        boxes_info: List[dict],
        intensity: int = 5,
        blur_type: str = "gaussian",
    ) -> np.ndarray:
        """Применение размытия внутри бокса на изображении.

        Args:
            image: Изображение в формате numpy array
            boxes_info: Информация о боксах
            intensity: Степень размытия от 1 до 10
            blur_type: Тип размытия: "gaussian", "motion" или "pixelate"

        Returns:
            Изображение с размытыми боксами
        """
        result_image = image.copy()

        for box in boxes_info:
            x_min, y_min, x_max, y_max = box["coordinates"]

            roi = result_image[y_min:y_max, x_min:x_max]

            if blur_type == "pixelate":
                h, w = roi.shape[:2]
                factor = max(1, intensity * 5)
                down_w = max(1, w // factor)
                down_h = max(1, h // factor)
                temp = cv2.resize(roi, (down_w, down_h), interpolation=cv2.INTER_LINEAR)
                blurred_roi = cv2.resize(temp, (w, h), interpolation=cv2.INTER_NEAREST)
            elif blur_type == "motion":
                k = intensity * 2 + 1
                kernel = np.zeros((k, k))
                kernel[int((k - 1) / 2), :] = 1.0 / k
                blurred_roi = cv2.filter2D(roi, -1, kernel)
            else:
                k = intensity * 2 + 1
                blurred_roi = cv2.GaussianBlur(roi, (k, k), 0)

            result_image[y_min:y_max, x_min:x_max] = blurred_roi

        return result_image
    
    @staticmethod
    def filter_boxes_by_area(boxes_info: List[dict], min_area: int = 500) -> List[dict]:
        """
        Фильтрация боксов по минимальной площади
        
        Args:
            boxes_info: Информация о боксах
            min_area: Минимальная площадь бокса
            
        Returns:
            Отфильтрованный список боксов
        """
        filtered_boxes = []
        
        for box in boxes_info:
            x_min, y_min, x_max, y_max = box['coordinates']
            area = (x_max - x_min) * (y_max - y_min)
            
            if area >= min_area:
                box['area'] = area
                filtered_boxes.append(box)
        
        return filtered_boxes
    
    @staticmethod
    def filter_boxes_by_confidence(boxes_info: List[dict], min_confidence: float = 0.7) -> List[dict]:
        """
        Фильтрация боксов по уверенности
        
        Args:
            boxes_info: Информация о боксах
            min_confidence: Минимальная уверенность
            
        Returns:
            Отфильтрованный список боксов
        """
        return [box for box in boxes_info if box['confidence'] >= min_confidence]
