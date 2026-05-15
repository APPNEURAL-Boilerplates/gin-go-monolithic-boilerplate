APP_NAME := gin-monolithic-boilerplate

.PHONY: dev test test-cover build run tidy fmt vet check docker-up docker-down

dev:
	go run ./cmd/api

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

build:
	mkdir -p bin
	go build -trimpath -o bin/api ./cmd/api

run: build
	./bin/api

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

check: fmt tidy vet test

docker-up:
	docker compose up --build

docker-down:
	docker compose down --remove-orphans
