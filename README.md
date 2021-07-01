# go-shortener

В данном проекте будет использоваться мультиплексор github.com/gorilla/mux
так как он поддерживает параметры URI что необходимо при реализации REST и
также он наиболее популярен в среде разработчиков.


Add link:
curl -i -X POST -H "Content-Type: application/json" -d '{"resource": "http://url.com/something", "custom_name": "Custom"}' http://127.0.0.1:8084/_api/link/

Get all links:
curl -i -X GET http://127.0.0.1:8084/_api/link/

Get specified link:
curl -i -X GET http://127.0.0.1:8084/_api/link/1
