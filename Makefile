BUILD_FLAGS := -ldflags "-X 'github.com/michimani/cfkvs/cli.version=v1.0.0' -X 'github.com/michimani/cfkvs/cli.revision=$(shell git rev-parse --short HEAD)'"
BUILD_OUTPUT := ./cfkvs
BUILD_CMD_DIR := ./cmd/cfkvs

# # Default target
# .PHONY: all
# all: build

# Build target
.PHONY: build
build:
	go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT) $(BUILD_CMD_DIR)

# Clean target
.PHONY: clean
clean:
	rm -f $(BUILD_OUTPUT) | true