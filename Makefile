## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
## run: run the cmd/api application
.PHONY: run
run:
	go run ./cmd/api
.PHONY: docs
docs:
	swag init -g cmd/api/main.go --output docs/greenlight
## migrations name=$1: create a new database migration
.PHONY: migration
migration:
	@echo "Creating migration files for ${name}"
	migrate create -seq -ext=.sql -dir=./migrations ${name}
## db/migrations/up: apply all up database migrations
.PHONY: up
up: confirm
	@echo "Running up migrations..."
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up