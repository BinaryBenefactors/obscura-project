FROM golang:1.24-alpine

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]