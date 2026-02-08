SHELL := /bin/bash

.PHONY: run test test-unit test-integration test-coverage emulator-start emulator-stop

run:
	@echo "Starting Go server..."
	cd server && go run .

test: test-unit

test-unit:
	cd server && go test ./...

test-integration:
	cd server && FIRESTORE_EMULATOR_HOST=localhost:8081 go test -tags integration -v -count=1 ./integration/

test-coverage:
	cd server && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

emulator-start:
	@echo "Starting Firestore emulator..."
	docker compose up --build -d firestore

emulator-stop:
	@echo "Stopping Firestore emulator..."
	docker compose down
