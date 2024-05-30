.PHONY: build, help, all

all: help

BUILD_PATH=target
APP_NAME=niter
TARGET_FILE=cmd/wasm/main.go

build:
	@echo "Building wasm"
	@GOOS=js GOARCH=wasm go build -o ${BUILD_PATH}/${APP_NAME}.wasm ${TARGET_FILE}

build-dev: build
	cp ${BUILD_PATH}/${APP_NAME}.wasm ../ui/public/wasm/${APP_NAME}.wasm

run:
	@echo "Running niter-node"
	@go run cmd/cli/main.go

help:
	@echo "NITER-CORE"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build		Build wasm"
	@echo "  build-dev	Build wasm and copy to ui"
	@echo "  help		Display this help message"
