.PHONY: setup build run-api run-bot run-collector run-emulator run-all stop clean help test-api test-bot docker-up docker-down docker-logs docker-rebuild

# Настройки
SHELL := /bin/bash
GOFLAGS := -ldflags="-s -w"

# Создание структуры
setup:
	@echo "Preparing environment..."
	@mkdir -p bin
	@[ -f .env ] || cp .env.example .env || echo "Warning: .env.example not found"

# Сборка всех компонентов
build: setup
	@echo "Building components..."
	@go build $(GOFLAGS) -o bin/api cmd/smart-house-api/main.go
	@go build $(GOFLAGS) -o bin/bot cmd/smart-house-bot/main.go
	@go build $(GOFLAGS) -o bin/collector cmd/smart-house-collector/main.go
	@go build $(GOFLAGS) -o bin/emulator cmd/smart-house-emulators/main.go
	@echo "Build complete"

# Запуск компонентов
run-api: build
	@echo "Starting API..."
	@./bin/api

run-bot: build
	@echo "Starting Telegram bot..."
	@./bin/bot

run-collector: build
	@echo "Starting collector..."
	@./bin/collector

run-emulator: build
	@echo "Starting emulators..."
	@./bin/emulator

# Запуск всех компонентов (в фоне)
run-all: build
	@echo "Starting all components..."
	@mkdir -p logs
	@./bin/emulator > logs/emulator.log 2>&1 & echo $$! > .pid.emulator
	@sleep 2
	@./bin/collector > logs/collector.log 2>&1 & echo $$! > .pid.collector
	@sleep 2
	@./bin/api > logs/api.log 2>&1 & echo $$! > .pid.api
	@sleep 2
	@./bin/bot > logs/bot.log 2>&1 & echo $$! > .pid.bot
	@echo "All components started"

# Остановка всех компонентов
stop:
	@echo "Stopping all components..."
	@-[ -f .pid.emulator ] && kill $$(cat .pid.emulator) 2>/dev/null || true
	@-[ -f .pid.collector ] && kill $$(cat .pid.collector) 2>/dev/null || true
	@-[ -f .pid.api ] && kill $$(cat .pid.api) 2>/dev/null || true
	@-[ -f .pid.bot ] && kill $$(cat .pid.bot) 2>/dev/null || true
	@rm -f .pid.*
	@echo "All components stopped"

# Очистка
clean:
	@echo "Cleaning up..."
	@rm -rf bin/*
	@rm -f .pid.*
	@echo "Clean complete"

# Помощь
help:
	@echo "Smart House Commands"
	@echo "-------------------"
	@echo "make setup      - Prepare environment"
	@echo "make build      - Build all components"
	@echo "make run-api    - Run API server"
	@echo "make run-bot    - Run Telegram bot"
	@echo "make run-collector - Run data collector"
	@echo "make run-emulator - Run sensor emulators"
	@echo "make run-all    - Run all components"
	@echo "make stop       - Stop all components"
	@echo "make clean      - Clean build artifacts"
	@echo "make help       - Show this help"

# Тестирование компонентов
test-api:
	@echo "Testing API..."
	@curl -s http://localhost:8080/api/v1/sensors | jq . || echo "Failed: jq not installed or API not running"

test-bot:
	@echo "Testing bot connection..."
	@curl -s "https://api.telegram.org/bot$$(grep TELEGRAM_BOT_TOKEN .env | cut -d= -f2)/getMe" | jq . || echo "Failed: jq not installed or invalid token"

# Docker Compose команды
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-rebuild:
	docker-compose down
	docker-compose build
	docker-compose up -d