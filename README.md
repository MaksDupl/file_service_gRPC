📌 Описание проекта

Проект: gRPC-сервис для хранения и управления файлами (изображениями, документами и т.д.).

🔹 Основной функционал:
✅ Загрузка файлов на сервер
✅ Просмотр списка файлов (имя, дата создания, размер)
✅ Скачивание файлов
✅ Ограничение конкурентных подключений




🚀 Быстрый старт

1. Установка зависимостей

bash
# Установите protoc
brew install protobuf  # macOS
sudo apt install protobuf-compiler  # Linux

# Установите Go-плагины
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

2. Запуск сервера

bash
go run server/main.go
Сервер запустится на localhost:50051.

3. Запуск клиента

bash
# Загрузить файл
go run client/main.go upload ./test.jpg

# Получить список файлов
go run client/main.go list

# Скачать файл
go run client/main.go download test.jpg




🔧 Настройка

1. Конфигурация сервера

Порт: 50051 (можно изменить в server/main.go)
Директория для файлов: ./uploads (меняется в storage/file_storage.go)

2. Ограничения

Макс. одновременных загрузок/скачиваний: 10
Макс. запросов списка файлов: 100




📌 Контакты

Автор: Максим Д.
Почта: 9956576054@mail.ru
Telegram: @Deitsev
