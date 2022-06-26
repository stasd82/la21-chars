SHELL := /bin/bash

tidy:
	go mod tidy
	go mod vendor

VERSION := 1.0

all: chars-api

chars-api:
	docker build \
		-f infra/docker/dockerfile.chars-api \
		-t chars-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
