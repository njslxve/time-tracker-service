.PHONY: run

run:
	@go run ./cmd/time-tracker/main.go

.PHONY: up
up:
	@cd ./deploy && docker compose up -d

.PHONY: down
down:
	@cd ./deploy && docker compose down

.PHONY: migrate
migrate:
	@go run ./cmd/migration/main.go

.PHONY: swagger
swagger:
	@swag init -g ./cmd/time-tracker/main.go