# Banner-service

## Описание

Сервис  для взаимодействия с баннерами, фичами и тегами, взаимодействие по HTTP.

## Запуск

Для запуска используйте следующие команды:

1. `make run` (для локального запуска на вашей машине - должен быть запущен Docker и установлен make)
2. `make container` (для запуска приложения через docker-compose)
3. `go run ./cmd/banners/main.go` для локального запуска (требует запущенного postgres)

## Примеры запросов