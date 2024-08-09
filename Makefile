include .env

# migration path
MIGRATION_PATH=./src/database/migrations

# database url
DATABASE_URL="$(DB_CONNECTION)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"

.PHONY: dev
dev:
	./bin/air server --port $(APP_PORT)

.PHONY: migration-create
migration-create: 
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_user_table

.PHONE: migration-up
migration-up:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose up

.PHONE: migration-down
migration-down:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose down
