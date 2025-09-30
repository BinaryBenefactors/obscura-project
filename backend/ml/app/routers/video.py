import asyncio
import os
import time

from fastapi import APIRouter, UploadFile, File, Form

from app.schemas.uploadfile import (
    ProcessRequest,
    ProcessResponse,
    SuccessResponse,
    ErrorResponse,
)
from app.ml.tools.object_detector import MLObjectDetector
from app.ml.ml_executor import get_ml_executor
from app.tools.generate_name_file import generate_name_file


router = APIRouter()

UPLOAD_FOLDER = "uploads"
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

detector = MLObjectDetector()
try:
    detector.initialize()
except Exception:
    pass


@router.post("/process", response_model=ProcessResponse)
async def process_file(request: ProcessRequest) -> ProcessResponse:
    """Обработка файла через очередь воркеров (синхронно для совместимости)"""
    start = time.time()
    executor = get_ml_executor()
    
    # Добавляем файл в очередь
    success = executor.add_to_queue(request.file_path, request.options)
    if not success:
        return ErrorResponse(success=False, error_message="File already in processing")
    
    # ЖДЕМ завершения обработки (для совместимости с Go backend)
    timeout = 300  # 5 минут
    check_interval = 1  # проверяем каждую секунду
    elapsed = 0
    
    while elapsed < timeout:
        status = executor.get_status(request.file_path)
        if status and status["status"] == "completed":
            processing_time_ms = int((time.time() - start) * 1000)
            processed_size = os.path.getsize(status["result"]) if status["result"] else 0
            return SuccessResponse(
                success=True,
                processed_path=status["result"],
                processed_size=processed_size,
                processing_time_ms=processing_time_ms,
            )
        elif status and status["status"] == "error":
            return ErrorResponse(success=False, error_message=status.get("error", "Processing failed"))
            
        await asyncio.sleep(check_interval)
        elapsed += check_interval
    
    # Таймаут
    return ErrorResponse(success=False, error_message="Processing timeout")


@router.post("/uploadfile", response_model=ProcessResponse)
async def upload_file(
    file: UploadFile = File(...),
    blur_amount: int = Form(..., ge=1, le=10),
    blur_type: str = Form(...),
    object_types: str = Form(...),
) -> ProcessResponse:
    start = time.time()
    try:
        contents = await file.read()
        file_ext = file.filename.split(".")[-1].lower()
        filename = file.filename
        file_path = os.path.join(UPLOAD_FOLDER, filename)
        with open(file_path, "wb") as f:
            f.write(contents)

        blur_map = {
            "gaus": "gaussian",
            "gaussian": "gaussian",
            "motion": "motion",
            "pixelization": "pixelate",
            "pixelate": "pixelate",
        }
        mapped_blur = blur_map.get(blur_type.lower())
        if not mapped_blur:
            return ErrorResponse(
                success=False,
                error_message="Unsupported blur type",
            )

        # Создаем объект Options
        from app.schemas.uploadfile import Options
        object_types_list = [obj.strip() for obj in object_types.split(",") if obj.strip()]
        options = Options(
            blur_type=mapped_blur,
            intensity=blur_amount,
            object_types=object_types_list
        )

        processed_path = detector.process_file(
            file_path,
            options.object_types,
            options.intensity,
            options.blur_type,
        )
        processed_size = os.path.getsize(processed_path)
        processing_time_ms = int((time.time() - start) * 1000)
        return SuccessResponse(
            success=True,
            processed_path=processed_path,
            processed_size=processed_size,
            processing_time_ms=processing_time_ms,
        )
    except Exception as e:
        return ErrorResponse(success=False, error_message=str(e))



