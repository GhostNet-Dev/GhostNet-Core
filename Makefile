BINARY_NAME	:= GhostNet
APP_BIN_PATH := $(CURDIR)/bin/$(BINARY_NAME)
APP_SRC_PATH := $(CURDIR)/cmd/server

ifneq ($(OS),Windows_NT)
UNAME_S := $(shell uname -s)
endif

run:
	go run $(APP_SRC_PATH)/main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: build


build:
ifeq ($(OS),Windows_NT)
	GOARCH=amd64 GOOS=window go build -o $(APP_BIN_PATH)-windows ${APP_SRC_PATH}/main.go
else
    ifeq ($(UNAME_S),Linux)
		GOARCH=amd64 GOOS=linux go build -o $(APP_BIN_PATH)-linux ${APP_SRC_PATH}/main.go
    endif
    ifeq ($(UNAME_S),Darwin)
		GOARCH=amd64 GOOS=darwin go build -o $(APP_BIN_PATH)-darwin ${APP_SRC_PATH}/main.go
    endif
endif

runs:
	./${APP_BIN_PATH}

build_and_run: build run

clean:
	go clean
ifeq ($(OS),Windows_NT)
	rm ${APP_BIN_PATH}-windows
else
    ifeq ($(UNAME_S),Linux)
		rm ${APP_BIN_PATH}-linux
    endif
    ifeq ($(UNAME_S),Darwin)
		rm ${APP_BIN_PATH}-darwin
    endif
endif

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
	ROOT :+$(dir $(abspath $(lastword $(ROOT))))
	DIR = $(ROOT)/dir
	COPYDIR = $(ROOT)/copydir
	install : 
		if [! -d $(DIR)]; then mkdir $(NEWDIR); fi
		cp -r $(DIR)/ $(COPYDIR)/

protoc:
	protoc --go_out=./pkg/proto/types --go_opt=paths=source_relative -I ./pkg/proto/types block.proto
	protoc --go_out=./pkg/proto/types --go_opt=paths=source_relative -I ./pkg/proto/types transaction.proto
