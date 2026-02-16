## Описание проекта

**Todo Scheduler** — это веб-приложение для управления задачами. Проект реализован на Go и предоставляет как REST API, так и готовый веб-интерфейс.

**Для чего создан:**
- Управление персональными задачами
- Автоматическое вычисление дат на основе правил повторения

## Задания со звездочкой

- работа с переменными окружения
- обработка правил **m**, **w**
- обработка параметра search
- добавлена аутентификация
- создан докер-образ

## Запуск приложения

### 1. Клонирование репозитория

```bash
git clone https://github.com/Chugerbuller/task-scheduler
cd todo-scheduler
```

### 2. Запуск приложения
```bash
go mod tidy
go run ./cmd/main.go
```
### 3. Откройте браузер
```text
http://localhost:7540/login.html
```

## Тесты

### Настройка 
Откройте файл tests/settings.go и вставьте полученный токен:

```go
package tests

var (
    Port        = "7540"
    DBFile      = "../database/scheduler_test.db"
    FullNextDate = true
    Search      = true
    Token       = "" // Вставьте ваш токен здесь!
)
```
## Запуск тестов

```bash
go test ./tests
```
## Запуск через Docker
```bash
# Сборка образа
docker build -t todo-scheduler .

# Запуск контейнера
 docker run -p 7540:7540 todo-scheduler
```
