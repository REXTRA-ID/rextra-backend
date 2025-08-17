OS := $(shell uname -s 2>/dev/null || echo Windows)

dep: 
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

# build: 
# 	go build -o main main.go

# run-build: build
# 	./main

# test:
# 	go test -v ./tests

build-docker-dev:
	docker compose -f docker-compose.dev.yml up -d --build

up-dev: 
	docker-compose -f docker-compose.dev.yml up -d

down-dev:
	docker-compose -f docker-compose.dev.yml down

# logs:
# 	docker-compose logs -f

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  tidy        Tidy dependencies"
	@echo "  run         Run the application"
	@echo "  migrate     Run database migrations"
	@echo "  seeder      Seed the database"
	@echo "  watch       Run program with auto loading"