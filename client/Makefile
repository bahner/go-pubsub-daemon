#!/usr/bin/make -ef

NAME = go-pubsub-daemon-client
PREFIX ?= /usr/local

$(NAME): tidy
	go build -o $(NAME)

tidy: go.mod
	go mod tidy

clean:
	rm -f client $(NAME)

install: $(NAME)
	sudo mv $(NAME) $(PREFIX)/bin

.PHONY: build clean tidy
