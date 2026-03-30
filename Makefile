# --------------------------------
# Configuration
# --------------------------------
APP_NAME := credit-card-api
BINARY_NAME := credit-card-api
DOCKER_IMAGE := credit-card-api
DB_CONTAINER := postgres
DB_PORT := 5432

# --------------------------------
# Targets
# --------------------------------

# Build Go Application Binary locally
build:
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) .

start:
	@echo "Starting Go REST API application container as well as database..."
	docker-compose up $(APP_NAME) -d --build

stop:
	@echo "Stopping Go REST API application container..."
	docker-compose down $(APP_NAME)

start_postgres:
	@echo "Starting Postgres DB container..."
	docker-compose up $(DB_CONTAINER) -d

stop_postgres:
	@echo "Stopping Postgres DB container..."
	docker-compose down $(DB_CONTAINER)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)




