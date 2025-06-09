# Go Worker Pool

Простой пример реализации `worker pool` на языке Go с возможностью **динамического добавления и удаления воркеров**.

## Описание

- Входящие строки (задания) отправляются в канал.
  
- Каждый воркер запускается как отдельная горутина и обрабатывает задания.

- Можно в любой момент добавить или удалить воркера.

- Завершение воркеров осуществляется через `context.Context`.

## Возможности

-  Динамическое добавление воркеров
  
-  Удаление конкретных воркеров
  
-  Очередь заданий
  
-  Безопасное завершение через `Shutdown()`

-  Обработка через `WaitGroup` и `mutex`

## Как запустить

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Mukam21/go-worker-pool.git
   cd go-worker-pool


Пример вывода:

Worker 0 started

Worker 1 started

Worker 1 processing job: Task 1

Worker 0 processing job: Task 0

Worker 0 processing job: Task 2

Worker 1 processing job: Task 3

Worker 0 processing job: Task 4

Worker 0 processing job: Task 6

Worker 2 started

Worker 2 processing job: Task 7

Worker 1 processing job: Task 5

Worker 0 stopped

Shutting down pool...

Worker 2 stopped

Worker 1 stopped

Pool shutdown complete.

