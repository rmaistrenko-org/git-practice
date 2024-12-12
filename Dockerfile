# Базовий образ
FROM golang:1.23

# Робоча директорія
WORKDIR /app

# Копіювання файлів проекту
COPY . .

# Копіювання скрипта wait-for-it
COPY wait-for-it.sh /usr/local/bin/wait-for-it

# Дозвіл на виконання скрипта
RUN chmod +x /usr/local/bin/wait-for-it

# Завантаження залежностей
RUN go mod tidy

# Збірка додатка
RUN go build -o main ./cmd/api/main.go

# Відкриття порту
EXPOSE 8000

# Очікуємо готовність бази даних і запускаємо додаток
CMD ["/usr/local/bin/wait-for-it", "db:3306", "--timeout=30", "--", "./main"]
