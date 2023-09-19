# Titanic

## Introduction

---
Titanic API exposes several endpoints over titanic dataset, enabling querying passengers and passengers aggregated data.

## Prerequisites

---
Before using the Makefile commands, make sure you have the following tools and dependencies installed:

- Go: Make sure you have Go installed to build the API server and run tests.
- Docker: Install Docker to build and run Docker images.
- Docker Compose (optional): If you plan to use Docker Compose to deploy the application.
- Kubernetes Cluster: If you plan to use Kubernetes and helm or kubectl to deploy the application

## Environment Variables

---
The application uses environment variables located in `.env`. Make sure to set the environment variables if you are not using the Makefile.
If you are using the Makefile environment variables will be loaded from `.env`.

## Folder Structure

---

Project structure follows unofficial standard however simple one for API layout - https://github.com/golang-standards/project-layout

#### project uses vendor folder while building the binary to not download dependencies in runtime
## OpenAPI Docs / Custom UI

---
Once you run the API you can access the OpenAPI UI in `/api/docs/` or `/api/docs/index.html`.
Also you can access the custom UI built with HTMX under `/ui`.

#### Notice the default host and ports are http://localhost:8089

## Store

---
### CSV store
Dataset used in the API is the Titanic CSV data under folder `/data/csv/titanic.csv`
### SQLite store
Dataset is a copy of `/data/csv/titanic.csv` data located in `/data/sqlite/titanic.db`
created with the following in sqlite terminal:

```
CREATE TABLE passengers (
    id INTEGER PRIMARY KEY,
    survived INTEGER,
    class INTEGER,
    name TEXT,
    sex TEXT,
    age TEXT,
    siblings_spouses INTEGER,
    parents_children INTEGER,
    ticket TEXT,
    fare REAL,
    cabin TEXT,
    embarked TEXT
);
```
and then 
```
.mode csv
.import data/csv/titanic.csv passengers

```

#### NOTICE: If you want to check both implementation you can set the store type in `config.yaml` to `SQLITE/CSV`.

## Tests

---

Not all test files are complete however you can find testing for edge case examples in API handlers test files for reference

## Makefile Commands

---
### `make run`

Run the API server in standalone mode using `go run`.

### `make build`

Build the API server binary using `go build`.

### `make docker-build`

Build the API server as a Docker image.

### `make docker-build-store`

Build the data store as a Docker image.

### `make docker-run`

Run the API server as a Docker container.

### `make docker-run-store`

Run the data store as a Docker container.

### `make docker-start`

Build Docker API image and run it.

### `make docker-remove`

Remove the Docker images and containers.

### `make docker-compose-start`

Build Docker images for API & store and deploy using Docker Compose.

### `make docker-compose-remove`

Stop the Docker Compose services.

### `make k8s-deploy`

Deploy Kubernetes resources.

### `make k8s-remove`

Remove Kubernetes resources.

### `make helm-deploy`

Deploy Kubernetes resources using helm.

### `make helm-remove`

Remove Kubernetes resources using helm.

### `make api-docs`

Generate OpenAPI3 specification and save it as `docs/openapi.json`.

### `make test`

Precentiles tests for the project.

### `make coverage`

Measure code coverage for the tests.

### `make coverage-html`

Generate an HTML report of the code coverage.
