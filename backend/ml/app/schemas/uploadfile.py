from typing import List, Literal, Union

from pydantic import BaseModel, Field


class Options(BaseModel):
    """Параметры обработки файла."""

    blur_type: Literal["gaussian", "motion", "pixelate"]
    intensity: int = Field(..., ge=1, le=10)
    object_types: List[str] = Field(default_factory=list)


class ProcessRequest(BaseModel):
    """Схема входящего запроса на обработку файла."""

    file_id: str
    file_path: str
    mime_type: str
    options: Options


class SuccessResponse(BaseModel):
    """Ответ успешной обработки."""

    success: Literal[True]
    processed_path: str
    processed_size: int
    processing_time_ms: int


class ErrorResponse(BaseModel):
    """Ответ при возникновении ошибки."""

    success: Literal[False]
    error_message: str


ProcessResponse = Union[SuccessResponse, ErrorResponse]

