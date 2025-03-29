include .env

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@, $(MAKECMDGOALS))

.PHONY: up-migrate
up-migrate:
	@migrate -path $(MIGRATIONS_PATH) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: down-migrate
down-migrate:
	@migrate -path $(MIGRATIONS_PATH) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: concurrency-test
concurrency-test:
	/usr/local/go/bin/go run scripts/test_concurrency.go

.PHONY: seed-data
seed-data:
	/usr/local/go/bin/go run cmd/migrate/seed/main.go
