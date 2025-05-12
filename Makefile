
.PHONY: build run

MIGRATIONS_DIR=./migrations

fmt:
	go fmt ./...

lint: fmt
	staticcheck

vet: fmt
	go vet ./...

build:
	go build -o build/xcrawler ./cmd/xcrawler/main.go

run: build
	./build/xcrawler


# Database Migration Commands
gup:
	goose --dir $(MIGRATIONS_DIR) up
gdown:
	goose --dir $(MIGRATIONS_DIR) down
gstatus:
	goose --dir $(MIGRATIONS_DIR) status
gcreate-%:
	goose --dir $(MIGRATIONS_DIR) create ${@:gcreate-%=%} sql