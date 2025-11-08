.PHONY: help tidy run watch seeder migrate both up down reset build-docker docker-migrate docker-seeder docker-both

# Panggil dengan: make up ENV=prod
ENV ?= dev
APP_NAME ?= rextra-backend
COMPOSE_FILE = docker-compose.$(ENV).yml
PROJECT_NAME = $(APP_NAME)-$(ENV)
ENV_FILE = .env.$(ENV)

ifeq ($(ENV),prod)
	DOCKER_RUN_CMD = /app/main
else
	DOCKER_RUN_CMD = go run main.go
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
	docker compose -p $(PROJECT_NAME) -f $(COMFİLE) --env-file $(ENV_FILE) down -v

build-docker:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d --build

# Docker exec targets
docker-migrate:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --migrate

docker-seeder:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --seeder

docker-both:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) --env-file $(ENV_FILE) exec app $(DOCKER_RUN_CMD) --migrate --seeder

# Help
help:
	@echo "Usage: make [target] [ENV=prod]"
	@echo "  (ENV default ke 'dev' jika tidak diset)"
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
