.PHONY: up down migrate run test

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	go run cmd/migrate/main.go up

run:
	go run cmd/bot/main.go

test:
	go test -v ./...
