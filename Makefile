.PHONY: build

VERSION := 0.1.0
BINARY_NAME=peykar


build:
	GOOS=linux go build -o ./build/$(BINARY_NAME)

install:
	go install