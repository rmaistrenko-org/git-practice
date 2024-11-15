# Базовий образ для збірки
FROM golang:1.23 AS builder

# Робоча директорія
WORKDIR /app

# Кешуємо залежності для швидшої збірки при повторних змінах
COPY go.mod go.sum ./
RUN go mod download

# Копіюємо решту файлів проекту
COPY . .

# Збірка додатка
RUN go build -o main ./cmd/api/main.go

# Базовий образ для запуску
FROM alpine:latest

# Встановлюємо необхідні інструменти
RUN apk --no-cache add ca-certificates bash

# Копіюємо зібраний додаток з першого образу
COPY --from=builder /app/main /main

# Копіюємо скрипт wait-for-it
COPY wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

# Відкриваємо порт
EXPOSE 8000

# Запуск додатка з очікуванням бази даних
CMD ["wait-for-it", "db:3306", "--timeout=30", "--", "/main"]
