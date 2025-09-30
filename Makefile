# --------------------------------
# Configuration
# --------------------------------
APP_NAME := credit-card-api
BINARY_NAME := credit-card-api
DOCKER_IMAGE := credit-card-api
MONGO_CONTAINER := mongodb
MONGO_PORT := 27017

# --------------------------------
# Targets
# --------------------------------

# Build Go Application Binary locally
build:
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) .

build_docker_image:
	@echo "Building Docker image..."
	docker build -f infrastructure/Dockerfile -t $(DOCKER_IMAGE) .

start: build_docker_image
	@echo "Running Go app container..."
	docker-compose up $(APP_NAME) -d

stop:
	@echo "Stopping Go app container..."
	docker-compose down $(APP_NAME)

start_mongodb:
	@echo "Stopping Go app container..."
	docker-compose up $(MONGO_CONTAINER) -d

stop_mongodb:
	@echo "Stopping Go app container..."
	docker-compose down $(MONGO_CONTAINER)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)




