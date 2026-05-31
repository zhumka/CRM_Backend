.PHONY: run build tidy test swagger fmt vet docker-up docker-down docker-logs

# Локальный запуск (нужен запущенный PostgreSQL и .env)
run:
	go run ./cmd/api

build:
	go build -o bin/crm ./cmd/api

tidy:
	go mod tidy

test:
	go test ./... -race -count=1

fmt:
	gofmt -w .

vet:
	go vet ./...

# Генерация Swagger-документации (нужен swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g cmd/api/main.go -o docs --parseInternal --parseDepth 2

# Полный запуск в Docker (API + PostgreSQL)
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api
