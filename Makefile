.PHONY: all build-server build-client

all: build-server build-client

build-server:
	go build -v ./cmd/server/

build-client:
	go build -v ./cmd/client/

.DEFAULT_GOAL := all
