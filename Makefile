SHELL := /bin/bash

.PHONY: run test emulator-start emulator-stop

run:
	@echo "Starting Go server..."
	cd server && go run .

test:
	cd server && go test ./...

emulator-start:
	@echo "Starting Firestore emulator..."
	docker compose up --build -d firestore

emulator-stop:
	@echo "Stopping Firestore emulator..."
	docker compose down
