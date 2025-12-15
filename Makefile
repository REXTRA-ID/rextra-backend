OS := $(shell uname -s 2>/dev/null || echo Windows)

dep: 
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

# Parsing argument --dev or --prod
ENV ?= dev
COMPOSE_FILE = docker-compose.$(ENV).yml
APP_CONTAINER = rextra-backend-$(ENV)
DB_CONTAINER  = rextra-backend-db-$(ENV)

# Docker targets
up:
	docker compose -f $(COMPOSE_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) down

reset:
	docker compose -f $(COMPOSE_FILE) down -v

build-docker:
	docker compose -f $(COMPOSE_FILE) up -d --build

DOCKER_CONTAINER = $(shell docker ps --filter "name=rextra-backend-$(ENV)" --format "{{.Names}}")

docker-migrate:
	docker exec -it $(DOCKER_CONTAINER) /bin/sh -c "go run main.go --migrate"

docker-seeder:
	docker exec -it $(DOCKER_CONTAINER) /bin/sh -c "go run main.go --seed"

docker-both:
	docker exec -it $(DOCKER_CONTAINER) /bin/sh -c "go run main.go --migrate --seed"

# Help
help:
	@echo "Usage: make [target] [--dev|--prod]"
	@echo "Targets:"
	@echo "  tidy        Tidy dependencies"
	@echo "  run         Run the application"
	@echo "  migrate     Run database migrations locally"
	@echo "  seeder      Seed the database locally"
	@echo "  watch       Run program with auto loading"
	@echo "  up          Start docker container (--dev|--prod)"
	@echo "  down        Stop docker container (--dev|--prod)"
	@echo "  reset       Stop and remove docker container and volumes (--dev|--prod)"
	@echo "  build-docker Build docker container (--dev|--prod)"
	@echo "  docker-migrate Run migrations inside Docker container"
	@echo "  docker-seeder  Run seeder inside Docker container"
	@echo "  docker-both    Run migrate and seeder inside Docker container"
