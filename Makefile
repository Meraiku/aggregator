include .env

.ONESHELL:


build:
	@go build -o ./.bin/gator ./cmd/gator/

run:build
	@./.bin/gator


docker:
	@docker compose up -d --build

up:
	@go install github.com/pressly/goose/v3/cmd/goose@latest;
	@cd ./sql/migrations;
	@goose postgres ${POSTGRES_DSN} up

down:
	@go install github.com/pressly/goose/v3/cmd/goose@latest;
	@cd ./sql/migrations;
	@goose postgres ${POSTGRES_DSN} down
