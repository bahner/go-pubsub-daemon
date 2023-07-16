#!/usr/bin/make -ef


NAME = myspace-pubsub-daemon
PREFIX ?= /usr/local

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: build image client

build: tidy $(NAME)

$(NAME): tidy
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

install: build
	install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

client:
	make -C client|

clean:
	rm -f $(NAME)
	make -C client clean

dist-clean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

.PHONY: build client serve tidy
