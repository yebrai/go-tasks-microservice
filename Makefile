.PHONY: test run fmt infra infra-down

# Variables
APP_NAME := taskservice
MAIN_PATH := cmd/api/main.go

test:
	go test ./...

run:
	go run $(MAIN_PATH)

fmt:
	go fmt ./...

# Inicia la infraestructura necesaria para desarrollo local
infra:
	docker-compose up -d

# Detiene y elimina la infraestructura de desarrollo
infra-down:
	docker-compose down --volumes