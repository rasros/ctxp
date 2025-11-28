# Makefile for ctxp

GO           ?= go
BINARY_NAME  ?= ctxp
CMD_DIR      := ./cmd/ctxp
BIN_DIR      := ./bin

.PHONY: all build install fmt test run clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

install:
	$(GO) install $(CMD_DIR)

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

run: build
	$(BIN_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BIN_DIR)

