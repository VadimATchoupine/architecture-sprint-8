# Этап 1: Сборка приложения
FROM golang:1.23-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git

# Создание рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

# Этап 2: Финальный образ
FROM alpine:3.18

# Установка необходимых зависимостей
RUN apk add --no-cache bash

# Создание директории для приложения
WORKDIR /root/

# Копирование скомпилированного приложения из первого этапа
COPY --from=builder /app/api .

# Команда запуска
CMD ["./api"]