#сборка
FROM golang:1.25.6 AS builder

WORKDIR /app

#модульные файлы и зависимости
COPY go.mod go.sum ./
RUN go mod download

#Копируем  код
COPY . .

#статический бинарник для Linux
#linux
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

#образ 
FROM ubuntu:latest

WORKDIR /app

# бинарник из сборки
COPY --from=builder /app/todo-app .

# фронтенд
COPY --from=builder /app/web ./web

# Переменные окружения
# переопределение через --env-file
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=

# Порт
EXPOSE 7540

# Директория для базы данных 
VOLUME ["/data"]

CMD ["./todo-app"]