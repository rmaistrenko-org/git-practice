# Базовий образ
FROM golang:1.23

# Встановлення утиліти migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz -o migrate.tar.gz \
    && tar -xvf migrate.tar.gz \
    && chmod +x migrate \
    && mv migrate /usr/local/bin/

# Робоча директорія
WORKDIR /app

# Копіювання файлів проекту
COPY . .

# Копіювання SQL-файлів міграцій
COPY ./migrations /app/migrations

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

# Очікуємо готовність бази даних, запускаємо міграції і додаток
CMD ["/usr/local/bin/wait-for-it", "db:3306", "--timeout=30", "--", "sh", "-c", "migrate -path /app/migrations -database 'mysql://root:1111@tcp(db:3306)/go_crud_api' up && ./main"]
