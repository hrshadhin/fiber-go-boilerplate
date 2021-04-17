.PHONY: clean test security build run

APP_NAME = fiber-go-boilerplate
BUILD_DIR = ./build
MIGRATIONS_FOLDER = ./platform/migrations
DB_NAME = fiber_go_api
DB_USER = dev
DB_PASS = dev
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@localhost/$(DB_NAME)?sslmode=disable

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out

security:
	gosec -quiet ./...

test: security
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

swag:
	swag init

build: swag clean
	CGO_ENABLED=0  go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

seed:
	PGPASSWORD=$(DB_PASS) psql -h localhost -p 5432 -U$(DB_USER) -d $(DB_NAME) -a -f platform/seeds/001_seed_user_table.sql
	PGPASSWORD=$(DB_PASS) psql -h localhost -p 5432 -U$(DB_USER) -d $(DB_NAME) -a -f platform/seeds/002_seed_book_table.sql

docker.run: docker.setup docker.postgres docker.fiber migrate.up
	@echo "\n===========FGB==========="
	@echo "App is running...\nVisit: http://localhost:5000 OR http://localhost:5000/swagger/"

docker.setup:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network
	docker volume create fibergb-pgdata

docker.fiber.build: swag
	docker build -t fibergb:latest .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name fibergb-api \
		--network dev-network \
		-p 5000:5000 \
		fibergb

docker.postgres:
	docker run --rm -d \
		--name fibergb-postgres \
		--network dev-network \
		-e POSTGRES_USER=dev \
		-e POSTGRES_PASSWORD=dev \
		-e POSTGRES_DB=fiber_go_api \
		-v fibergb-pgdata:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.stop: docker.stop.fiber docker.stop.postgres

docker.stop.fiber:
	docker stop fibergb-api || true

docker.stop.postgres:
	docker stop fibergb-postgres || true

docker.dev:
	docker-compose up
