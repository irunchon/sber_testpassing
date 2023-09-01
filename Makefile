include .env
export


.PHONY: all
all:
	go run cmd/sber_testpassing/main.go
