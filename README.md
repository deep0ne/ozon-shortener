# ozon-shortener
## Сервис для сокращения ссылок
Докер образ: https://hub.docker.com/r/alexnet1337/dockerhub

Для запуска сервера необходимо запустить `main.go` в папке `cmd/server`. Варианты запуска
1. `go run main.go` - запустит сервис с `PostgreSQL` в docker контейнере
2. `go run main.go -m` - запуск с флагом `-m` поднимет `Redis` в docker контейнере

Запуск контейнеров осуществляется командой из `Makefile`.
После запуска станут доступны два эндпоинта по `localhost:8082`
1. `localhost:8082/shorten` - POST запрос, в теле которого необходимо отправить ссылку для сокращения. Например, POST запрос с телом
`json
{
    "URL": "https://desud-webdev-382/landocsweb/documents/341813/view/241627"
}` 
сделает запись в БД и отдаст сгенерированную короткую ссылку.
2. `localhost:8082/` - GET запрос, которому нужно предоставить короткую ссылку, чтобы получить оригинальную. Например, GET запрос `localhost:8082/uaekGG4vgWp` вернет соответствующую этой короткой ссылке полную ссылку.
