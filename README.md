# Your Project Name

## Introduction

This is a Makefile for building and running components of your project. The Makefile commands allow you to build the API
server, create Docker images for the API, CSV store, and SQLite store, run the components as Docker containers, and
generate API documentation.

## Prerequisites

Before using the Makefile commands, make sure you have the following tools and dependencies installed:

- Go: Make sure you have Go installed to build the API server and run tests.
- Docker: Install Docker to build and run Docker images.
- Docker Compose (optional): If you plan to use Docker Compose for managing multiple containers, install Docker Compose.

## Environment Variables

The Makefile uses the `.env` file to set environment variables. Make sure to include the required environment variables
in the `.env` file to customize the behavior of the Makefile commands.

## Makefile Commands

### `make run`

Precentiles the API server in normal mode using `go run`.

### `make build`

Build the API server binary.

### `make docker-build`

Build the API server as a Docker image.

### `make docker-build-csv-store`

Build the CSV store as a Docker image.

### `make docker-build-sqlite-store`

Build the SQLite store as a Docker image.

### `make docker-run`

Precentiles the API server as a Docker container.

### `make docker-run-csv-store`

Precentiles the CSV store as a Docker container.

### `make docker-run-sqlite-store`

Precentiles the SQLite store as a Docker container.

### `make docker-start`

Build Docker API image and run it.

### `make docker-stop`

Stop the Docker API container.

### `make docker-remove`

Remove the Docker images and containers.

### `make docker-compose-up`

Build Docker images for API, CSV store, and SQLite store, and run them using Docker Compose.

### `make docker-compose-stop`

Stop the Docker Compose services.

### `make api-docs`

Generate OpenAPI3 specification and save it as `docs/openapi.json`.

### `make test`

Precentiles tests for the project.

### `make coverage`

Measure code coverage for the tests.

### `make coverage-html`

Generate an HTML report of the code coverage.

---

Please note that the `.env` file is used to set environment variables for the Docker image builds and runs. Make sure to
customize the Makefile and `.env` file according to your project's specific requirements.

Remember to replace `Your Project Name` with the actual name of your project. This README file provides an overview of
the available Makefile commands to help users build, run, and manage the components of your project.
