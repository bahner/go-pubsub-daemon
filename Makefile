#!/usr/bin/make -ef

GO ?= go

NAME = go-pubsub-daemon
CLIENT_NAME = $(NAME)-client
MODULE_NAME = github.com/bahner/go-pubsub-daemon
PREFIX ?= /usr/local

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: build image client

build: tidy $(NAME)

$(NAME): tidy
	$(GO) build -o $(NAME)

tidy: go.mod
	$(GO) mod tidy

serve: build
	./$(NAME)

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

go.mod:
	$(GO) mod init $(MODULE_NAME)

install: build client
	sudo install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	sudo install -Dm755 client/$(CLIENT_NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)-client

client:
	make -C client

clean:
	rm -f $(NAME)
	make -C client clean

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

.PHONY: build client serve tidy
