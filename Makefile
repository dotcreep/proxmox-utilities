.PHONY: dev
dev:
	@go run ./cmd/server/

.PHONY: build
build:
	@go build -o ./bin/main ./cmd/main.go

.PHONY: run
run:
	@go build -o ./bin/main ./cmd/main.go
	@go run ./bin/main