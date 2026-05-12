# Файлы для итогового задания

Директория pkg/db содержит код для работы с базой данных SQLite: инициализацию, структуру задачи и все операции с таблицей scheduler.

Директория pkg/api содержит HTTP-обработчики всех API-эндпоинтов, аутентификации и функцию вычисления следующей даты повторения NextDate.

Директория pkg/server содержит функцию запуска веб-сервера: инициализацию БД, регистрацию обработчиков и запуск на указанном порту.

Директория web содержит файлы фронтенда: HTML, CSS и JavaScript.

Директория tests содержит тесты для проверки API и логики вычисления дат. Настройки тестов находятся в файле tests/settings.go.

Файл scheduler.db — база данных SQLite, создаётся автоматически при первом запуске в текущей директории.

:0

# Выполненные задания со звёздочкой

Поиск задач — обработчик /api/tasks?search= поддерживает поиск по подстроке в заголовке и комментарии, а также по конкретной дате в формате 02012006

Аутентификация — защита API через JWT-токен, пароль задаётся переменной окружения TODO_PASSWORD

Docker — многоступенчатый Dockerfile, база данных монтируется с хоста

;p

# Запускать только из корневой директории проекта:
```bash
cd todo-app
go run main.go
```

Запуск без пароля:go run .
go run .   Затем, откройте браузер по адресу: http://localhost:7540

Запуск с паролем для Linux
TODO_PASSWORD=mysecret go run .

Для Windows PowerShell
$env:TODO_PASSWORD="parolb"
go run .

Запуск с указанием порта и пути к БД
Linux
TODO_PORT=8080 TODO_DBFILE=/tmp/scheduler.db go run .

Windows PowerShell
$env:TODO_PORT="8080"
$env:TODO_DBFILE="C:\data\scheduler.db"
go run .

Откройте браузер по адресу: http://localhost:7540 (или другой порт, если переопределили TODO_PORT).


# тесты  запуск

Настройка tests/settings.go
govar Port = 7540                 
var DBFile = "../scheduler.db"  
var FullNextDate = true         
var Search = true               
var Token = ``           // JWT-токен (нужен только если задан TODO_PASSWORD)


#  Аутентификации

Запуск без аутентификации: go test ./tests -v
Запуск с аутентификацией: 
'1' Запустите сервер с паролем: $env:TODO_PASSWORD="mysecret"
go run .  
'2' Получите JWT-токен и сохраните в файл: $r = Invoke-RestMethod -Uri http://localhost:7540/api/signin -Method POST -ContentType "application/json" -Body '{"password":"mysecret"}'
$r.token | Out-File token.txt 
'3' Вставьте содержимое token.txt в tests/settings.go: var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiNTNjOGEwZWViZmUyNTBjZTY1M2EwZTIwYjMwZjg1MzE1NDMzNzE4MDM3NGE0NDkzM2FjY2NlYzNhZDI0YjBmNSJ9.Sg9e35P9Qw7zRFLsULWkxMubS50CBWmYi8ueWNiJSuY` 
'4' Запустите тесты с той же переменной окружения: $env:TODO_PASSWORD="mysecret"
go test ./tests -v


# Сборка и запуск через Docker
Сборка образа
В корне проекта (рядом с Dockerfile):
bashdocker build -t todo-app .

Запуск контейнера
Linux
docker run -d \
  -p 7540:7540 \
  -v /home/user/todo-app:/data \
  -e TODO_PASSWORD=mysecret \
  --name todo-app \
  todo-app

У меня шиндоувс -_-
PowerShell
docker run -d `
  -p 7540:7540 `
  -v C:\Users\user\Documents\todo-app:/data `
  -e TODO_PASSWORD=mysecret `
  --name todo-app `
  todo-app

