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

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: valkey-get-all-keys
valkey-get-all-keys:
	@docker exec -it valkey-cache valkey-cli -a valkey_password KEYS "*"

.PHONY: install-test-tool-npm
install-test-tool-npm:
	sudo apt install nodejs npm
	sudo npm install -g autocannon


.PHONY: test
test: 
	@go test -v ./...
	
CONNECTIONS ?= 100
DURATION ?= 2
ENDPOINT ?= http://localhost:8000/v1/users/122
TOKEN ?= 
# Load test target
.PHONY: benchmark
benchmark:
	@if [ -z "$(TOKEN)" ]; then \
		echo "Error: TOKEN is required. Use 'make benchmark TOKEN=your_jwt_token'"; \
		exit 1; \
	fi
	npx autocannon $(ENDPOINT) --connections $(CONNECTIONS) --duration $(DURATION) -H "Authorization: Bearer $(TOKEN)"


