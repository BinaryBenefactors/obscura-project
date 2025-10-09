import os
from typing import List, Optional, Tuple, Union

import cv2
import numpy as np
from moviepy.editor import VideoFileClip

from app.ml.tools.model import Model
from app.ml.tools.write_box import BoxProcessor


class MLObjectDetector:
    """Класс, объединяющий несколько моделей для детекции объектов."""

    def __init__(
        self,
        face_model_path: str = "models/yolov11m-face.pt",
        general_model_path: str = "models/yolo11m.pt",
        confidence_threshold: float = 0.5,
    ) -> None:
        self.face_model = Model(face_model_path, confidence_threshold)
        self.general_model = Model(general_model_path, confidence_threshold)
        self.box_processor = BoxProcessor()

    def initialize(self) -> None:
        """Загрузка моделей."""
        self.face_model.load_model()
        self.general_model.load_model()

    def _get_output_filename(self, input_path: str) -> str:
        """Возвращает путь для сохранения результата с суффиксом _processed."""
        dir_name = os.path.dirname(input_path)
        base_name = os.path.basename(input_path)
        name, ext = os.path.splitext(base_name)
        output_name = f"{name}_processed{ext}"
        return os.path.join(dir_name, output_name)

    def _run_models(
        self, image_source: Union[str, np.ndarray], object_types: List[str]
    ) -> List[dict]:
        boxes: List[dict] = []
        run_face = "face" in object_types
        run_general = set(object_types)|set(self.general_model.class_names.values())

        if run_face:
            results = self.face_model.predict(image_source)
            boxes.extend(self.face_model.extract_boxes(results))

        if len(run_general)!=0:
            results = self.general_model.predict(image_source)
            boxes.extend(self.general_model.extract_boxes(results))

        if object_types:
            boxes = [b for b in boxes if b["class_name"] in object_types]

        return boxes

    def detect_objects(
        self,
        image_source: Union[str, np.ndarray],
        object_types: List[str],
        intensity: int,
        blur_type: str,
        min_area: Optional[int] = None,
        min_confidence: Optional[float] = None,
    ) -> Tuple[List[dict], np.ndarray]:
        """Полный цикл детекции объектов для изображений."""

        if isinstance(image_source, str):
            image = cv2.imread(image_source)
        else:
            image = image_source

        boxes_info = self._run_models(image_source, object_types)

        if min_area is not None:
            boxes_info = self.box_processor.filter_boxes_by_area(boxes_info, min_area)
        if min_confidence is not None:
            boxes_info = self.box_processor.filter_boxes_by_confidence(
                boxes_info, min_confidence
            )

        result_image = self.box_processor.draw_boxes(
            image, boxes_info, intensity, blur_type
        )
        return boxes_info, result_image

    def process_image(
        self,
        image_path: str,
        object_types: List[str],
        intensity: int,
        blur_type: str,
    ) -> str:
        boxes_info, result_image = self.detect_objects(
            image_path, object_types, intensity, blur_type
        )
        output_path = self._get_output_filename(image_path)
        cv2.imwrite(output_path, result_image)
        return output_path

    def process_video(
        self,
        video_path: str,
        object_types: List[str],
        intensity: int,
        blur_type: str,
    ) -> str:
        cap = cv2.VideoCapture(video_path)
        if not cap.isOpened():
            raise ValueError(f"Не удалось открыть видео файл: {video_path}")

        fps = int(cap.get(cv2.CAP_PROP_FPS))
        width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))

        temp_output = self._get_output_filename(video_path).replace(".", "_temp.")
        fourcc = cv2.VideoWriter_fourcc(*"mp4v")
        out = cv2.VideoWriter(temp_output, fourcc, fps, (width, height))

        try:
            while True:
                ret, frame = cap.read()
                if not ret:
                    break
                _, processed_frame = self.detect_objects(
                    frame, object_types, intensity, blur_type
                )
                out.write(processed_frame)
        finally:
            cap.release()
            out.release()

        output_path = self._add_audio_to_video(video_path, temp_output)
        if os.path.exists(temp_output):
            os.remove(temp_output)
        return output_path

    def process_video_with_tracking(
        self,
        video_path: str,
        object_types: List[str],
        intensity: int,
        blur_type: str,
        detection_interval: int = 10,
    ) -> str:
        """
        Обработка видео с использованием CSRT трекера.
        YOLO запускается только каждый N-й кадр, остальные - трекинг.
        
        Args:
            video_path: Путь к видео файлу
            object_types: Типы объектов для детекции
            intensity: Интенсивность размытия
            blur_type: Тип размытия
            detection_interval: Интервал кадров для запуска детекции (по умолчанию 10)
            
        Returns:
            Путь к обработанному видео
        """
        cap = cv2.VideoCapture(video_path)
        if not cap.isOpened():
            raise ValueError(f"Не удалось открыть видео файл: {video_path}")

        fps = int(cap.get(cv2.CAP_PROP_FPS))
        width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))

        temp_output = self._get_output_filename(video_path).replace(".", "_temp.")
        fourcc = cv2.VideoWriter_fourcc(*"mp4v")
        out = cv2.VideoWriter(temp_output, fourcc, fps, (width, height))

        # Список трекеров и их боксов
        trackers = []
        frame_count = 0

        try:
            while True:
                ret, frame = cap.read()
                if not ret:
                    break

                frame_count += 1
                boxes_info = []

                # Каждый N-й кадр - запускаем детекцию YOLO
                if frame_count % detection_interval == 1:
                    raw_boxes = self._run_models(frame, object_types)
                    
                    # Создаем новые трекеры для всех обнаруженных объектов
                    trackers = []
                    boxes_info = []
                    for box in raw_boxes:
                        x1, y1, x2, y2 = box["coordinates"]
                        
                        # Валидация координат - ограничиваем по границам кадра
                        x1 = max(0, min(x1, width - 1))
                        y1 = max(0, min(y1, height - 1))
                        x2 = max(0, min(x2, width))
                        y2 = max(0, min(y2, height))
                        
                        # Проверяем что бокс имеет положительные размеры
                        if x2 > x1 and y2 > y1:
                            # Добавляем валидированный бокс для отрисовки
                            boxes_info.append({
                                "coordinates": [x1, y1, x2, y2],
                                "class_name": box["class_name"],
                                "confidence": box["confidence"]
                            })
                            
                            # Инициализируем трекер
                            tracker = cv2.legacy.TrackerCSRT_create()
                            bbox = (x1, y1, x2 - x1, y2 - y1)
                            tracker.init(frame, bbox)
                            trackers.append({
                                "tracker": tracker,
                                "class_name": box["class_name"],
                                "confidence": box["confidence"]
                            })
                else:
                    # На остальных кадрах - обновляем трекеры
                    updated_trackers = []
                    for tracker_info in trackers:
                        success, bbox = tracker_info["tracker"].update(frame)
                        if success:
                            x, y, w, h = [int(v) for v in bbox]
                            # Валидация координат - ограничиваем по границам кадра
                            x1 = max(0, x)
                            y1 = max(0, y)
                            x2 = min(width, x + w)
                            y2 = min(height, y + h)
                            
                            # Проверяем что бокс имеет положительные размеры
                            if x2 > x1 and y2 > y1:
                                boxes_info.append({
                                    "coordinates": [x1, y1, x2, y2],
                                    "class_name": tracker_info["class_name"],
                                    "confidence": tracker_info["confidence"]
                                })
                                updated_trackers.append(tracker_info)
                    trackers = updated_trackers

                # Применяем размытие к найденным областям
                processed_frame = self.box_processor.draw_boxes(
                    frame, boxes_info, intensity, blur_type
                )
                out.write(processed_frame)

        finally:
            cap.release()
            out.release()

        output_path = self._add_audio_to_video(video_path, temp_output)
        if os.path.exists(temp_output):
            os.remove(temp_output)
        return output_path

    def _add_audio_to_video(self, original_video_path: str, processed_video_path: str) -> str:
        output_path = self._get_output_filename(original_video_path)
        try:
            original_clip = VideoFileClip(original_video_path)
            processed_clip = VideoFileClip(processed_video_path)
            if original_clip.audio is not None:
                final_clip = processed_clip.set_audio(original_clip.audio)
            else:
                final_clip = processed_clip
            final_clip.write_videofile(
                output_path,
                codec="libx264",
                audio_codec="aac",
                temp_audiofile="{output_path}-temp-audio.m4a",
                remove_temp=True,
            )
            original_clip.close()
            processed_clip.close()
            final_clip.close()
        except Exception:
            os.rename(processed_video_path, output_path)
        return output_path

    def process_file(
        self,
        file_path: str,
        object_types: List[str],
        intensity: int,
        blur_type: str,
    ) -> str:
        """
        Обработка файла.
        
        Args:
            file_path: Путь к файлу
            object_types: Типы объектов для детекции
            intensity: Интенсивность размытия
            blur_type: Тип размытия
            
        Returns:
            Путь к обработанному файлу
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Файл не найден: {file_path}")

        file_ext = os.path.splitext(file_path)[-1].lower()
        image_extensions = {".jpg", ".jpeg", ".png", ".bmp", ".tiff", ".tif"}
        video_extensions = {".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv"}

        if file_ext in image_extensions:
            return self.process_image(file_path, object_types, intensity, blur_type)
        elif file_ext in video_extensions:
            return self.process_video(file_path, object_types, intensity, blur_type)
        else:
            raise ValueError(f"Неподдерживаемый формат файла: {file_ext}")

