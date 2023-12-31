# ozon-shortener
## Сервис для сокращения ссылок
Докер образ: https://hub.docker.com/r/alexnet1337/dockerhub

Для запуска сервера необходимо запустить `main.go` в папке `cmd/server`. Варианты запуска
1. `go run main.go` - запустит сервис с `PostgreSQL`
2. `go run main.go -r` - запуск с флагом `-m` поднимет `Redis`
3. `go run main.go -m` - запуск с in-memory базой

Добавлен `docker-compose`, чтобы можно было поднять у себя сервис.

После запуска станут доступны два эндпоинта по `localhost:8082`
1. `localhost:8082/shorten` - POST запрос, в теле которого необходимо отправить ссылку для сокращения. Например, POST запрос с телом
```json 
{
    "URL": "https://some-service-webdev-100/serviceweb/documents/info/341813/view/241627"
}
``` 
сделает запись в БД и отдаст сгенерированную короткую ссылку.

2. `localhost:8082/` - GET запрос, которому нужно предоставить короткую ссылку, чтобы получить оригинальную. Например, GET запрос `localhost:8082/uaekGG4vgWp` осуществит redirect на оригинальный адрес.
