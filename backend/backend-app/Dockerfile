FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod .
COPY go.sum .

# Скачиваем зависимости
RUN go mod download

# Устанавливаем swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем исходный код
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Генерируем swagger документацию
RUN swag init -g cmd/main.go -o docs/

# Создаем папку для загрузок
RUN mkdir -p uploads

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем бинарник и папку docs
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs/
COPY --from=builder /app/uploads ./uploads/

EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
