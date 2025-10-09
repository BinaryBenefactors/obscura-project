#!/usr/bin/env python3
"""
Демо-скрипт для тестирования обработки видео с CSRT трекингом.
YOLO запускается каждый 10-й кадр, остальные кадры используют трекинг.
"""

import os
import sys
import time
from pathlib import Path

# Добавляем путь к модулям приложения
sys.path.insert(0, str(Path(__file__).parent / "app"))

from app.ml.tools.object_detector import MLObjectDetector


def demo_tracking_video():
    """Демо обработки видео с трекингом"""
    
    # Путь к видео для обработки
    video_path = r"C:\OpenCamera\20241220_100510.MP4"
    
    # Проверяем существование файла
    if not os.path.exists(video_path):
        print(f"❌ Видео файл не найден: {video_path}")
        return
    
    print("=" * 80)
    print("🎬 ДЕМО: Обработка видео с CSRT трекингом")
    print("=" * 80)
    print(f"📁 Входное видео: {video_path}")
    print(f"📊 Алгоритм: YOLO каждый 10-й кадр + CSRT трекинг на остальных")
    print()
    
    # Создаем детектор
    print("🔧 Инициализация моделей...")
    detector = MLObjectDetector(
        face_model_path="models/yolov11m-face.pt",
        general_model_path="models/yolo11m.pt",
        confidence_threshold=0.5
    )
    detector.initialize()
    print("✅ Модели загружены успешно")
    print()
    
    # Параметры обработки
    object_types = ["face"]  # Детектируем лица
    intensity = 25  # Интенсивность размытия
    blur_type = "pixelate"  # Тип размытия: gaussian, pixelate, blackout
    detection_interval = 10  # Запускать YOLO каждый 10-й кадр
    
    print("⚙️  Параметры обработки:")
    print(f"   • Объекты для детекции: {', '.join(object_types)}")
    print(f"   • Тип размытия: {blur_type}")
    print(f"   • Интенсивность: {intensity}")
    print(f"   • Интервал детекции: каждый {detection_interval}-й кадр")
    print()
    
    # Обработка с трекингом
    print("🚀 Начало обработки с CSRT трекингом...")
    start_time = time.time()
    
    try:
        output_tracking = detector.process_video_with_tracking(
            video_path=video_path,
            object_types=object_types,
            intensity=intensity,
            blur_type=blur_type,
            detection_interval=detection_interval
        )
        
        tracking_time = time.time() - start_time
        print(f"✅ Обработка завершена за {tracking_time:.2f} секунд")
        print(f"📁 Результат: {output_tracking}")
        print()
        
    except Exception as e:
        print(f"❌ Ошибка при обработке с трекингом: {e}")
        import traceback
        traceback.print_exc()
        return
    
    # Для сравнения - обработка без трекинга (только YOLO на каждом кадре)
    print("🔄 Для сравнения: обработка БЕЗ трекинга (YOLO на каждом кадре)...")
    start_time_no_tracking = time.time()
    
    try:
        # Копируем видео для обработки без трекинга
        import shutil
        video_dir = os.path.dirname(video_path)
        video_name = os.path.basename(video_path)
        name_without_ext = os.path.splitext(video_name)[0]
        ext = os.path.splitext(video_name)[1]
        
        temp_video = os.path.join(video_dir, f"{name_without_ext}_no_tracking{ext}")
        shutil.copy(video_path, temp_video)
        
        output_no_tracking = detector.process_video(
            video_path=temp_video,
            object_types=object_types,
            intensity=intensity,
            blur_type=blur_type,
        )
        
        no_tracking_time = time.time() - start_time_no_tracking
        print(f"✅ Обработка без трекинга завершена за {no_tracking_time:.2f} секунд")
        print(f"📁 Результат: {output_no_tracking}")
        print()
        
        # Удаляем временную копию
        if os.path.exists(temp_video):
            os.remove(temp_video)
        
    except Exception as e:
        print(f"⚠️  Не удалось обработать без трекинга: {e}")
        no_tracking_time = None
    
    # Статистика
    print("=" * 80)
    print("📊 СТАТИСТИКА")
    print("=" * 80)
    print(f"⏱️  С трекингом (YOLO каждый {detection_interval}-й кадр): {tracking_time:.2f} сек")
    
    if no_tracking_time:
        print(f"⏱️  Без трекинга (YOLO на каждом кадре): {no_tracking_time:.2f} сек")
        speedup = no_tracking_time / tracking_time
        print(f"🚀 Ускорение: {speedup:.2f}x")
        time_saved = no_tracking_time - tracking_time
        print(f"💾 Сэкономлено времени: {time_saved:.2f} сек ({time_saved/60:.2f} мин)")
    
    print()
    print("=" * 80)
    print("✨ Демо завершено!")
    print("=" * 80)


if __name__ == "__main__":
    demo_tracking_video()

