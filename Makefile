.PHONY: help tidy run watch seeder migrate both up down reset build-docker docker-migrate docker-seeder docker-both

# Panggil dengan: make up ENV=prod
ENV ?= local
APP_NAME ?= rextra-backend
COMPOSE_FILE = docker-compose.$(ENV).yml
PROJECT_NAME = $(APP_NAME)-$(ENV)

ifeq ($(ENV),local)
	ENV_FILE = .env.dev
else
	ENV_FILE = .env.$(ENV)
endif

ifeq ($(ENV),local)
	DOCKER_RUN_CMD = go run main.go
else
	DOCKER_RUN_CMD = /app/main
endif

# Target Go Local
tidy:
	go mod tidy

run:
	go run main.go

watch:
	go run main.go --watch

seeder:
	go run main.go --seeder

reset-db:
	go run main.go --reset-db

migrate:
	go run main.go --migrate

both:
	go run main.go --migrate --seeder

# Docker targets
up:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

down:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down

reset:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down -v

build-docker:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d --build

# Docker exec targets
docker-migrate:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --migrate

docker-seeder:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --seeder

docker-both:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --migrate --seeder

deploy-prod:
	docker compose -f docker-compose.prod.yml up -d --build app

# Help
help:
	@echo "Usage: make [target] [ENV=local|dev|prod]"
	@echo "  (ENV default ke 'local' jika tidak diset)"
	@echo ""
	@echo "Targets Lokal:"
	@echo "  tidy        Tidy dependencies"
	@echo "  run         Run the application"
	@echo "  migrate     Run database migrations locally"
	@echo "  seeder      Seed the database locally"
	@echo "  watch       Run program with auto loading"
	@echo ""
	@echo "Targets Docker (Gunakan 'ENV=prod' untuk produksi):"
	@echo "  up          Start docker container"
	@echo "  down        Stop docker container"
	@echo "  reset       Stop and remove docker container and volumes"
	@echo "  build-docker Build docker container"
	@echo "  docker-migrate Run migrations inside Docker container"
	@echo "  docker-seeder  Run seeder inside Docker container"
	@echo "  docker-both    Run migrate and seeder inside Docker container"
	@echo ""
	@echo "Targets Deployment:"
	@echo "  deploy-prod Deploy to production"
