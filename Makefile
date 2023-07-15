#!/usr/bin/make -ef


NAME = myspace-pubsub-daemon

GO_VERSION = 1.19.11
export GO = go$(GO_VERSION)

BUILD_IMAGE = golang:$(GO_VERSION)-alpine
IMAGE = $(NAME):latest

all: build image client

build: tidy
	go build -o $(NAME)

tidy: go.mod
	go mod tidy

serve: build
	./$(NAME)

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

go.mod:
	go mod init $(NAME)

client:
	make -C client

.PHONY: build client serve tidy
