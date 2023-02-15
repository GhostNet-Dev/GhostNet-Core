BINARY_CONTAINER_NAME	:= ghostnet
BINARY_DEAMON_NAME	:= ghostd
BINARY_CLI_NAME	:= ghost

APP_BIN_CONTAINER_PATH := $(CURDIR)/bin/$(BINARY_CONTAINER_NAME)
APP_BIN_DEAMON_PATH := $(CURDIR)/bin/$(BINARY_DEAMON_NAME)
APP_BIN_CLI_PATH := $(CURDIR)/bin/$(BINARY_CLI_NAME)

APP_CLI_SRC_PATH := $(CURDIR)/cmd/cli
APP_DEAMON_SRC_PATH := $(CURDIR)/cmd/ghostd
APP_CONTAINER_SRC_PATH := $(CURDIR)/cmd/container

ifneq ($(OS),Windows_NT)
UNAME_S := $(shell uname -s)
endif

run:
	go run $(APP_CONTAINER_SRC_PATH)/main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: cli ghostd container

ghostd:
ifeq ($(OS),Windows_NT)
	set GOARCH=amd64& set GOOS=windows& go build -o $(APP_BIN_DEAMON_PATH).exe ${APP_DEAMON_SRC_PATH}/main.go
else
    ifeq ($(UNAME_S),Linux)
	GOARCH=amd64 GOOS=linux go build -o $(APP_BIN_DEAMON_PATH) ${APP_DEAMON_SRC_PATH}/main.go
    endif
endif


container:
ifeq ($(OS),Windows_NT)
	set GOARCH=amd64& set GOOS=windows& go build -o $(APP_BIN_CONTAINER_PATH).exe ${APP_CONTAINER_SRC_PATH}/main.go
else
    ifeq ($(UNAME_S),Linux)
	GOARCH=amd64 GOOS=linux go build -o $(APP_BIN_CONTAINER_PATH) ${APP_CONTAINER_SRC_PATH}/main.go
    endif
endif

cli:
ifeq ($(UNAME_S),Linux)
	GOARCH=amd64 GOOS=linux go build -o $(APP_BIN_CLI_PATH) ${APP_CLI_SRC_PATH}/main.go
endif

runs:
	./${APP_BIN_CONTAINER_PATH}

build_and_run: build run

clean:
	go clean
ifeq ($(OS),Windows_NT)
	rm ${APP_BIN_CONTAINER_PATH}-windows
else
    ifeq ($(UNAME_S),Linux)
		rm ${APP_BIN_CONTAINER_PATH}-linux
    endif
    ifeq ($(UNAME_S),Darwin)
		rm ${APP_BIN_CONTAINER_PATH}-darwin
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
	protoc --go_out=./pkg/proto/rpc --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/proto/rpc --go-grpc_opt=paths=source_relative -I ./pkg/proto/ptypes -I ./pkg/proto/rpc controls.proto
	protoc --go_out=./pkg/proto/ptypes --go_opt=paths=source_relative -I ./pkg/proto/ptypes block.proto
	protoc --go_out=./pkg/proto/ptypes --go_opt=paths=source_relative -I ./pkg/proto/ptypes transaction.proto
	protoc --go_out=./pkg/proto/ptypes --go_opt=paths=source_relative -I ./pkg/proto/ptypes user.proto
	protoc --go_out=./pkg/proto/packets --go_opt=paths=source_relative -I ./pkg/proto/ptypes -I ./pkg/proto/packets masternet_packet.proto
	protoc --go_out=./pkg/proto/packets --go_opt=paths=source_relative -I ./pkg/proto/ptypes -I ./pkg/proto/packets block_packet.proto
	protoc --go_out=./pkg/proto/packets --go_opt=paths=source_relative -I ./pkg/proto/ptypes -I ./pkg/proto/packets common_packet.proto
	protoc --go_out=./pkg/proto/packets --go_opt=paths=source_relative -I ./pkg/proto/ptypes -I ./pkg/proto/packets file_packet.proto
