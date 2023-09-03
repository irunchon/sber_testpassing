include .env
export


.PHONY: all
all:
	go run cmd/sber_testpassing/main.go

.PHONY: run
build:
	go run cmd/sber_testpassing/main.go

.PHONY: test
test:
	go test ./internal/...

.PHONY: test-coverage
test-coverage:
	go test ./internal/... -coverprofile=coverage.out  && go tool cover -html=coverage.out
