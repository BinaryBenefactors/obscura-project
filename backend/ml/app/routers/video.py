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
    start = time.time()
    try:
        processed_path = detector.process_file(
            request.file_path,
            request.options.object_types,
            request.options.intensity,
            request.options.blur_type,
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

