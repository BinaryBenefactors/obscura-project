FROM golang:1.24-alpine

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod .
COPY go.sum .

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Создаем папку для загрузок
RUN mkdir -p uploads

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

EXPOSE 8080

# Запускаем приложение
CMD ["./main"]