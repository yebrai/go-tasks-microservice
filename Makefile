# Makefile for Tasks Microservice

# Variables
DOCKER_COMPOSE = docker-compose
BACKEND_IMAGE = taskservice-backend
FRONTEND_IMAGE = taskservice-frontend
APP_NAME := taskservice
MAIN_PATH := cmd/api/main.go

# Colors for output
YELLOW = \033[33m
GREEN = \033[32m
RED = \033[31m
NC = \033[0m # No Color

.PHONY: help dev infra infra-down build-backend build-frontend build-all run-backend run-frontend run-all clean test lint fmt

# Default target
help: ## Show this help message
	@echo "$(YELLOW)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development
dev: ## Start development environment (all services)
	@echo "$(YELLOW)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) up --build

dev-detached: ## Start development environment in detached mode
	@echo "$(YELLOW)Starting development environment (detached)...$(NC)"
	$(DOCKER_COMPOSE) up --build -d

# Infrastructure
infra: ## Start only infrastructure services (MongoDB, RabbitMQ)
	@echo "$(YELLOW)Starting infrastructure services...$(NC)"
	$(DOCKER_COMPOSE) up mongo rabbitmq -d

infra-down: ## Stop infrastructure services
	@echo "$(YELLOW)Stopping infrastructure services...$(NC)"
	$(DOCKER_COMPOSE) down --volumes

# Build
build-backend: ## Build backend Docker image
	@echo "$(YELLOW)Building backend image...$(NC)"
	docker build -f Dockerfile.backend -t $(BACKEND_IMAGE) .

build-frontend: ## Build frontend Docker image
	@echo "$(YELLOW)Building frontend image...$(NC)"
	docker build -f Dockerfile.frontend -t $(FRONTEND_IMAGE) .

build-all: build-backend build-frontend ## Build all Docker images
	@echo "$(GREEN)All images built successfully!$(NC)"

# Run services individually
run-backend: ## Run only backend service
	@echo "$(YELLOW)Starting backend service...$(NC)"
	$(DOCKER_COMPOSE) up taskservice -d

run-frontend: ## Run only frontend service
	@echo "$(YELLOW)Starting frontend service...$(NC)"
	$(DOCKER_COMPOSE) up frontend -d

run-all: ## Run all services
	@echo "$(YELLOW)Starting all services...$(NC)"
	$(DOCKER_COMPOSE) up -d

# Frontend development
frontend-dev: ## Start frontend in development mode
	@echo "$(YELLOW)Starting frontend development server...$(NC)"
	cd web && npm install && npm run dev

frontend-build: ## Build frontend for production
	@echo "$(YELLOW)Building frontend for production...$(NC)"
	cd web && npm install && npm run build

# Go commands
run: ## Run Go backend locally
	@echo "$(YELLOW)Running Go backend locally...$(NC)"
	go run $(MAIN_PATH)

go-build: ## Build Go binary
	@echo "$(YELLOW)Building Go binary...$(NC)"
	go build -o bin/taskservice $(MAIN_PATH)

test: ## Run Go tests
	@echo "$(YELLOW)Running Go tests...$(NC)"
	go test ./... -v

test-coverage: ## Run Go tests with coverage
	@echo "$(YELLOW)Running Go tests with coverage...$(NC)"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Code quality
lint: ## Run linter
	@echo "$(YELLOW)Running linters...$(NC)"
	golangci-lint run

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting Go code...$(NC)"
	go fmt ./...

# Cleanup
clean: ## Clean up containers, images, and volumes
	@echo "$(YELLOW)Cleaning up...$(NC)"
	$(DOCKER_COMPOSE) down -v --rmi all
	docker system prune -f

# Logs
logs: ## Show logs for all services
	$(DOCKER_COMPOSE) logs -f

logs-backend: ## Show logs for backend service
	$(DOCKER_COMPOSE) logs -f taskservice

logs-frontend: ## Show logs for frontend service
	$(DOCKER_COMPOSE) logs -f frontend

# Database
db-shell: ## Connect to MongoDB shell
	@echo "$(YELLOW)Connecting to MongoDB shell...$(NC)"
	docker exec -it mongo mongosh -u mongoroot -p secret --authenticationDatabase admin taskdb

# RabbitMQ
rabbitmq-mgmt: ## Open RabbitMQ Management UI
	@echo "$(YELLOW)RabbitMQ Management UI: http://localhost:15672$(NC)"
	@echo "$(YELLOW)Username: guest, Password: guest$(NC)"

# Health checks
health: ## Check health of all services
	@echo "$(YELLOW)Checking service health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Backend unhealthy$(NC)"
	@curl -f http://localhost:3000/ || echo "$(RED)Frontend unhealthy$(NC)"

# GCP deployment
gcp-build: ## Build and deploy to GCP using Cloud Build
	@echo "$(YELLOW)Deploying to GCP...$(NC)"
	gcloud builds submit --config=cloudbuild.yaml

# Quick start
quickstart: ## Quick start for new developers
	@echo "$(GREEN)🚀 Quick Start Guide:$(NC)"
	@echo "1. Start infrastructure: $(YELLOW)make infra$(NC)"
	@echo "2. Run backend: $(YELLOW)make run$(NC)"
	@echo "3. Run frontend: $(YELLOW)make frontend-dev$(NC)"
	@echo "4. Or run everything: $(YELLOW)make dev$(NC)"
	@echo ""
	@echo "$(GREEN)📊 Useful URLs:$(NC)"
	@echo "- Frontend: http://localhost:3000"
	@echo "- Backend API: http://localhost:8080"
	@echo "- RabbitMQ Management: http://localhost:15672"
	@echo "- MongoDB: mongodb://localhost:27017"