.PHONY: all ready up migrate psql generator eventlistener eventstats down topic_created proto fmt clean build

all: ready clean build

ready: up migrate

up:
	docker compose up -d

migrate:
	docker compose exec database sh -c 'psql -U casino < /db/migrations/00001.create_base.sql'
	docker compose exec database sh -c 'psql -U casino < /db/migrations/00002.events.sql'

psql:
	docker run -it --rm --link db:postgres --net event-processing_default -e POSTGRES_USER=casino -e POSTGRES_PASSWORD=casino -p 5432 postgres psql postgresql://casino:casino@db

# easier to simply call the command here
generator:
	bin/generator

eventlistener:
	bin/eventlistener

eventstats:
	bin/eventstats

down:
	docker compose down -v

topic_created:
	docker container exec -it kafka /bin/kafka-console-consumer --bootstrap-server localhost:9092 --topic event-created -from-beginning

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/ports/event.proto

fmt:
	gofmt -l -w -e ./

clean: fmt
	go clean
	rm -f bin/*

build:
	go build -o bin/generator cmd/generator/main.go
	go build -o bin/eventlistener cmd/eventlistener/main.go
	go build -o bin/eventstats cmd/eventstats/main.go
