#!/usr/bin/make -ef

GO ?= go1.19.11

NAME = client
MODULE_NAME = myspace/$(NAME)

build: tidy
	go build -o $(NAME)

tidy: go.mod
	go mod tidy

go.mod:
	go mod init $(MODULE_NAME)

.PHONY: build tidy
