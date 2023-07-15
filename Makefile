#!/usr/bin/make -ef

#export GO = go1.19.11
export GO = foooooo
NAME = myspace
MODULE_NAME = myspace

all: build client

build: tidy
	go build -o $(NAME)

tidy: go.mod
	go mod tidy

serve: build
	./$(NAME)

go.mod:
	go mod init $(MODULE_NAME)

client:
	make -C client

.PHONY: build client serve tidy
