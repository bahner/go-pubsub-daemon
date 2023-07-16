#!/usr/bin/make -ef

NAME = client
MODULE_NAME = myspace/$(NAME)

build: tidy
	go build -o $(NAME)

tidy: go.mod
	go mod tidy

go.mod:
	go mod init $(MODULE_NAME)

clean:
	rm -f client $(NAME)

.PHONY: build clean tidy
