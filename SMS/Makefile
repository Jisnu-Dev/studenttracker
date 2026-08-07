include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_PATH=./migrations 

create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

test:
	go test "-coverprofile=coverage.out" ./...

html:
	go tool cover "-html=coverage.out"
