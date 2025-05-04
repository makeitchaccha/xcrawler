
.PHONY: build run

MIGRATIONS_DIR=./migrations

build:
	go build -o build/xcrawler ./cmd/xcrawler/main.go

run: build
	./build/xcrawler


# Database Migration Commands
db-up:
	goose --dir $(MIGRATIONS_DIR) up
db-down:
	goose --dir $(MIGRATIONS_DIR) down
db-status:
	goose --dir $(MIGRATIONS_DIR) status
db-create:
	goose --dir $(MIGRATIONS_DIR) create $(name) sql