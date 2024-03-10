#!make
include .env
export $(shell sed 's/=.*//' .env)

build:
	@go build -C ./cmd/rssagg -o ../../bin/rssagg

run:
	@./bin/rssagg

dev:
	@make build && make run

goose-up:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) up && make sqlc-gen

goose-down:
	@goose -dir=./internal/sql/schema postgres $(DB_URL) down && make sqlc-gen

sqlc-gen:
	@sqlc generate

test:
	@go test -v ./...