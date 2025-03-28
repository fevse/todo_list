BIN := "./bin/todolist"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# build:
# 	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/todolist

# run: build
# 	$(BIN) 

run:
	docker-compose up

down:
	docker-compose down
	docker rmi todo_list-app

test:
	go test -race -v ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run test lint