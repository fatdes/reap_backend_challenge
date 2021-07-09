.PHONY: test pre-commit generate install local local-clean local-automated-test

.DEFAULT_GOAL := test

test:
	go test -v -race -cover ./...

pre-commit:
	go mod tidy
	go mod vendor
	go vet
	go fmt ./...

generate:
	go generate ./...

install:
	go install

local:
	docker compose -f .local/docker-compose.yaml up --build test

local-clean:
	docker compose -f .local/docker-compose.yaml down
	docker compose -f .local/docker-compose.yaml rm

local-automated-test:
	docker compose -f .local/docker-compose.yaml build automated_test
	docker compose -f .local/docker-compose.yaml run automated_test
