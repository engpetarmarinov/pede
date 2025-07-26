PEDE_SRC = ./cmd/pede/main.go
PEDE_BIN = pede
IN ?= examples/hello.pede
OUT ?= $(basename $(notdir $(IN)))
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

.PHONY: all build run clean

all: build

build:
	@echo "Building pede binary..."
	OS=$(OS) ARCH=$(ARCH) go build -ldflags="-s -w" -o $(PEDE_BIN) $(PEDE_SRC)

run: build
	@echo "Compiling $(IN) to $(OUT)..."
	./$(PEDE_BIN) build -o $(OUT) --os=$(OS) --arch=$(ARCH) $(IN)
	./$(OUT)

clean:
	@echo "Cleaning..."
	rm -f pede *.ll hello arithmetics
