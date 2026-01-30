# Rextra Backend

Backend service for the Rextra application, built with Go.

## Getting Started

This guide will walk you through setting up and running the project for both development and production environments. The project is fully containerized using Docker and managed with a `Makefile` for simplicity.

### Prerequisites

Make sure you have the following software installed on your machine:

- [Go](https://go.dev/doc/install) (for local development without Docker)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Configuration

The project uses `.env` files for environment variable configuration. You'll need to create two separate files for development and production.

1.  **Create Development Environment File (`.env.dev`)**

    Create a file named `.env.dev` in the root directory and add the following variables. You can change the values as needed.

    ```env
    # Application
    APP_NAME=rextra-backend
    APP_PORT=8081

    # Nginx
    NGINX_PORT=81

    # Database (PostgreSQL)
    DB_USER=postgres
    DB_PASS=password
    DB_NAME=rextra
    ```

2.  **Create Production Environment File (`.env.prod`)**

    Create a file named `.env.prod` for your production server with your production-ready credentials.

    ```env
    # Application
    APP_NAME=rextra-backend

    # Nginx (Port 8888 is used as port 80 might be in use)
    NGINX_PORT=8888

    # Database (PostgreSQL) - Use strong, secret credentials
    DB_USER=your_production_user
    DB_PASS=your_secret_production_password
    DB_NAME=rextra_prod
    ```

## Usage (Development)

All development tasks are managed via `make` commands.

1.  **Start the Environment**

    To build the Docker images and start all services (app, database, nginx) for development, run:

    ```sh
    make up
    ```

    - The application will be accessible via Nginx at `http://localhost:81`.
    - The Go app itself is exposed on `http://localhost:8081`.
    - The PostgreSQL database is exposed on port `5433`.

2.  **Run Database Migrations and Seeders**

    After starting the environment, you can run migrations and seed the database with initial data using:

    ```sh
    # Run migrations only
    make docker-migrate

    # Run seeders only
    make docker-seeder

    # Run both
    make docker-both
    ```

3.  **Stop the Environment**

    To stop all running Docker containers for the development environment, use:

    ```sh
    make down
    ```

    To stop the containers and also remove the database volume, use `make reset`.

## Deployment (Production)

The deployment process is also streamlined with `make`.

1.  **Initial Deployment**

    On your production server, after cloning the repository and creating the `.env.prod` file, run the following command to build and start all production services:

    ```sh
    make up ENV=prod
    ```

2.  **Updating the Application**

    For subsequent updates, you only need to pull the latest code changes and run a single command to rebuild and redeploy just the application container, leaving the database untouched.

    ```sh
    make deploy-prod
    ```

## Available `make` Commands

Here is a list of the most common commands available in the `Makefile`.

| Command               | `ENV`           | Description                                                 |
| --------------------- | --------------- | ----------------------------------------------------------- |
| `make up`             | `dev` (default) | Starts all services using Docker Compose.                   |
| `make up`             | `prod`          | Starts all production services.                             |
| `make down`           | `dev`/`prod`    | Stops all services.                                         |
| `make reset`          | `dev`/`prod`    | Stops services and removes the database volume.             |
| `make build-docker`   | `dev`/`prod`    | Forces a rebuild of all images and restarts services.       |
| `make docker-migrate` | `dev`/`prod`    | Runs database migrations inside the Docker container.       |
| `make docker-seeder`  | `dev`/`prod`    | Runs database seeders inside the Docker container.          |
| `make deploy-prod`    | (prod only)     | Rebuilds and deploys only the `app` service for production. |
| `make run`            | (local)         | Runs the Go application locally without Docker.             |
| `make watch`          | (local)         | Runs the app locally with hot-reloading.                    |

