#!make
include .env
export $(shell sed 's/=.*//' .env)

build:
	@go build -o bin/rssagg cmd/rssagg/main.go

run: build
	@./bin/rssagg

goose-up:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) up

goose-down:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) down

sqlc-gen:
	@sqlc generate

test:
	@go test -v ./...