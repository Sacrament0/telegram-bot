# Первый этап сборки. Выбор сборки. AS builder оначает, что данный этап будет использоваться как сборщик приложения
FROM golang:1.15-alpine3.12 AS builder
# копирование файлов из текущей директории, в ту, которая обязательно называется как go module
COPY . /github.com/Sacrament0/telegram-bot/
# объявляем эту директорию рабочей директорией, чтобы все последующие команды выполнялись в ней
WORKDIR /github.com/Sacrament0/telegram-bot/
# скачивание всех зависимостей
RUN go mod download
# компиляция бинарных файлов в директорию bin с названием файла bot, далее указываем путь к main.go
RUN go build -o ./bin/bot cmd/bot/main.go

# второй этап сборки
FROM alpine:latest
# объявляем рабочую папку
WORKDIR /root/
# копируем созданный бинарный файл в рабочую дирректорию. --from=0 означает копирование из предыдущего этапа сборки
# точка в конце - в текущую директорию
COPY --from=0 /github.com/Sacrament0/telegram-bot/bin/bot .
# копируем директорию config в папку config, т.к. приложение не запустится без конфига
COPY --from=0 /github.com/Sacrament0/telegram-bot/configs configs/
# пробрасываем порт наружу образа
EXPOSE 80
#  запуск команды в консоли. При запуске контейнера будет запускаться команда
CMD ["./bot"]
