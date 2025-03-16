.PHONY: build

VERSION := 0.1.0
BINARY_NAME=peykar


build:
	go build \
		-trimpath \
		-mod=vendor \
		-gcflags "-trimpath $(PWD)" \
		-asmflags "-trimpath $(PWD)" \
		-o ./build/peykar \


install: build
	install -m 755 ./build/peykar $(GOPATH)/bin/