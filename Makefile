-include .env

PROJECT_NAME := $(shell basename "$(PWD)" | tr '[:upper:]' '[:lower:]')

KUBERNETES_NAMESPACE="titanic"
DOCKER_API_IMAGE_NAME := "$(PROJECT_NAME):v1"
DOCKER_STORE_IMAGE_NAME := "$(PROJECT_NAME)-store:v1"

## run: Run the API server alone in normal mode
run:
	CSV_STORE_PATH=${CSV_STORE_PATH} \
	SQLITE_STORE_PATH=${SQLITE_STORE_PATH} \
	API_PORT=${API_PORT} \
	go run -mod=vendor ./cmd/api/main.go

## build: Build the API server binary
build: api-docs
	CGO_ENABLED=1 go build -mod=vendor -o ${PROJECT_NAME} ./cmd/api/main.go

## docker-build: Build the API server as a docker image
docker-build:
	$(info ---> Building Docker Image: ${DOCKER_API_IMAGE_NAME})
	docker build --progress=plain -t ${DOCKER_API_IMAGE_NAME} . \
		--build-arg port=${API_PORT} \
		--build-arg csv_path=${CSV_STORE_PATH} \
		--build-arg sqlite_path=${SQLITE_STORE_PATH}

## docker-build-store: Build the data store as a docker image
docker-build-store:
	$(info ---> Building Docker Image: ${DOCKER_STORE_IMAGE_NAME})
	docker build -t ${DOCKER_STORE_IMAGE_NAME} -f Dockerfile_store .

## docker-run: Run the API server as a docker container
docker-run:
	$(info ---> Running Docker Container: ${DOCKER_API_IMAGE_NAME})
	docker run -p ${API_PORT}:${API_PORT} -it $(DOCKER_API_IMAGE_NAME)

## docker-run-store: Run the data store as a docker container
docker-run-store:
	$(info ---> Running Docker Container: ${DOCKER_STORE_IMAGE_NAME})
	docker run -it $(DOCKER_STORE_IMAGE_NAME)

## docker-start: Builds Docker API image and runs it.
docker-start: docker-build docker-run

## docker-remove: Removes the docker images and containers for API and data store
docker-remove:
	-@docker stop $(DOCKER_API_IMAGE_NAME)
	-@docker stop $(DOCKER_STORE_IMAGE_NAME)
	-@docker rm -f $(DOCKER_API_IMAGE_NAME)
	-@docker rm -f $(DOCKER_STORE_IMAGE_NAME)
	-@docker rmi -f $(DOCKER_API_IMAGE_NAME)
	-@docker rmi -f $(DOCKER_STORE_IMAGE_NAME)

## docker-compose-start: Deploy docker compose containers
docker-compose-start: docker-build docker-build-store
	STORE_IMAGE_NAME=${DOCKER_STORE_IMAGE_NAME} \
    API_IMAGE_NAME=${DOCKER_API_IMAGE_NAME} \
    API_PORT=${API_PORT} \
    docker-compose up -d

## docker-compose-remove: Remove docker compose containers
docker-compose-remove:
	STORE_IMAGE_NAME=${DOCKER_STORE_IMAGE_NAME} \
    API_IMAGE_NAME=${DOCKER_API_IMAGE_NAME} \
    API_PORT=${API_PORT} \
	docker-compose down

## k8s-deploy: Deploy kubernetes resources
k8s-deploy: docker-build docker-build-store
	$(info ---> Deploying Kubernetes Deployment...)
	-kubectl create namespace ${KUBERNETES_NAMESPACE}
	kubectl apply -f deploy/k8s/configmap.yml
	kubectl apply -f deploy/k8s/deployment.yml
	kubectl apply -f deploy/k8s/service.yml


## docker-compose-remove: Remove kubernetes resources
k8s-remove:
	$(info ---> Deleting Kubernetes Deployment...)
	-kubectl delete -f deploy/k8s/deployment.yml
	-kubectl delete -f deploy/k8s/service.yml
	-kubectl delete -f deploy/k8s/configmap.yml
	-kubectl delete namespace ${KUBERNETES_NAMESPACE}

## helm-deploy: Deploy kubernetes resources using helm release
helm-deploy: docker-build docker-build-store
	$(info ---> Deploying Helm Chart Release...)
	-kubectl create namespace ${KUBERNETES_NAMESPACE}
	helm install $(PROJECT_NAME) deploy/helm/ --namespace ${KUBERNETES_NAMESPACE}

## helm-remove: Remove kubernetes resources using helm release
helm-remove:
	$(info ---> Deleting Helm Chart Release...)
	-helm uninstall $(PROJECT_NAME) --namespace ${KUBERNETES_NAMESPACE}
	-kubectl delete namespace ${KUBERNETES_NAMESPACE}

## api-docs: Generate OpenAPI3 Spec
api-docs:
	@go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go
	curl -X POST "https://converter.swagger.io/api/convert" \
		-H "accept: application/json" \
		-H "Content-Type: application/json" \
		-d @docs/swagger.json > docs/openapi.json

## test: Run tests
test:
	go test -v ./...

## coverage: Measures code coverage
coverage:
	go test ./... -v -coverprofile coverage.out -covermode count
	go tool cover -func=coverage.out

## coverage-html: Opens html code coverage
coverage-html: coverage
	go tool cover -html=coverage.out

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |	sed -e 's/^/ /'
	@echo
