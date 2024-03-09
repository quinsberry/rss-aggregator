#!make
include .env
export $(shell sed 's/=.*//' .env)

build:
	@go build -o bin/rss-aggregator

run:
	./bin/rss-aggregator

dev:
	make build && make run

goose-up:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) up

goose-down:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) down

sqlc-gen:
	@sqlc generate

test:
	@go test -v ./...