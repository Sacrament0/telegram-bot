.PHONY:
.SOLENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

# создание образа
build-image:
	docker build -t telegram-bot:v0.1 .
# запуск контейнера
start-container:
# имя telegram-bot, локалюный_порт:внутренний_порт_в_контейнере, файл в котором описаны переменные окружения, 
# чтобы прокидывать их при запуске контейнера, в конце название изображения
	docker run --name telegram-bot -p 80:80 --env-file .env telegram-bot:v0.1