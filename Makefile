SHELL := /bin/sh
PACKAGE_NAME := email-api-back
PROJECT_DIR := $(shell pwd)
VERSION ?= 0.0.0-SNAPSHOT
IMAGE_DEV := $(PACKAGE_NAME)-dev
IMAGE_RELEASE := $(REGISTRY)$(PACKAGE_NAME):$(VERSION)

.PHONY: run tidy test fmt docker-build docker-shell

run:
	go run ./cmd/main.go

fmt:
	gofmt -w ./api ./cmd ./healthcheck ./middleware

tidy:
	go mod tidy

test:
	go test ./...

docker-build:
	docker build -t $(IMAGE_DEV) -f Dockerfile.dev $(PROJECT_DIR)

docker-shell:
	docker run -it --rm --env-file=.env -v $(PROJECT_DIR):/app --net host -w /app --entrypoint=/bin/sh $(IMAGE_DEV)


