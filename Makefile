SHELL := /bin/bash

.DEFAULT_GOAL := help

COMPOSE ?= docker compose
BACKEND_DIR := backend
FRONTEND_DIR := frontend

.PHONY: help install backend-deps frontend-deps up up-d down logs ps restart dev dev-db dev-api dev-frontend

help:
	@printf "%-18s %s\n" "make up" "Start the full stack with Docker Compose"
	@printf "%-18s %s\n" "make up-d" "Start the full stack in detached mode"
	@printf "%-18s %s\n" "make down" "Stop Docker Compose services"
	@printf "%-18s %s\n" "make logs" "Follow Docker Compose logs"
	@printf "%-18s %s\n" "make ps" "Show Docker Compose service status"
	@printf "%-18s %s\n" "make install" "Install local development dependencies"
	@printf "%-18s %s\n" "make dev" "Start only the database for local development"
	@printf "%-18s %s\n" "make dev-api" "Run the Go API locally"
	@printf "%-18s %s\n" "make dev-frontend" "Run the Next.js frontend locally"

install: backend-deps frontend-deps

backend-deps:
	cd $(BACKEND_DIR) && go mod download

frontend-deps:
	cd $(FRONTEND_DIR) && npm ci

up:
	$(COMPOSE) up --build

up-d:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

restart: down up

dev: dev-db
	@echo "Database is running in Docker."
	@echo "Run 'make dev-api' and 'make dev-frontend' in separate terminals."

dev-db:
	$(COMPOSE) up -d db

dev-api:
	@set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	cd $(BACKEND_DIR) && go run ./cmd/server

dev-frontend:
	cd $(FRONTEND_DIR) && npm run dev
