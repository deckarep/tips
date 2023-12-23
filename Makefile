.PHONY: test build lint all

# Define the default goal. When you run "make" without argument, it will run the "all" target.
default: all

# Capture additional arguments
ARGS ?=

# Test the code.
test:
	go test ./...

# Lint the code.
lint:
	golangci-lint run

# Build the project: run the linter and then build.
build: lint
	go build
	@echo "Build: Successful"

# Run all steps: build and then run the application.
# To forward args optionally set ARGS like: make ARGS="--flag1=value1 --flag2=value2"
all: build
	./tips $(ARGS)

