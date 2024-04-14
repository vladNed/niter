.PHONY: build, help, all

all: help

BUILD_PATH=target
APP_NAME=niter
TARGET_FILE=cmd/wasm/main.go

build:
	@echo "Building wasm"
	@GOOS=js GOARCH=wasm go build -o ${BUILD_PATH}/${APP_NAME}.wasm ${TARGET_FILE}

help:
	@echo "NITER-CORE"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build		Build wasm"
	@echo "  help		Display this help message"
