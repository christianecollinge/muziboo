.PHONY: run tidy

# Run the API server
run:
	go run ./cmd/api

# Tidy up dependencies
tidy:
	go mod tidy

# Build binary
build:
	go build -o bin/api ./cmd/api
