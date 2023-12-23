.PHONY: test build lint gen all

# Define the default goal. When you run "make" without argument, it will run the "all" target.
default: all

# Capture additional arguments which can optionally be passed in.
ARGS ?=

# Test the code.
test:
	go test -v ./...

# Lint the code.
lint:
	golangci-lint run

# Run any code generation on this step.
gen:
	go run testmode/gen_mock/main.go > testmode/devices.json

# Build the project: run the linter and then build.
build: lint
	go build -v ./...
	@echo "Build: Successful"

# Run all steps: build and then run the application.
# To forward args optionally set ARGS like: make ARGS="--flag1=value1 --flag2=value2"
all: build
	./tips $(ARGS)
