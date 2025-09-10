from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles

from app.routers import video

# Создание экземпляра FastAPI приложения
app = FastAPI(
    title="My API",
    description="API с поддержкой CORS",
    version="1.0.0"
)

# Настройка CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Монтирование папки uploads для раздачи медиафайлов
app.mount("/uploads", StaticFiles(directory="uploads"), name="uploads")

# Основной маршрут
@app.get("/")
async def read_root():
    return {"message": "Hello World"}

# Подключаем роутеры
app.include_router(video.router, prefix="/api", tags=["api"])

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
