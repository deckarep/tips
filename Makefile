.PHONY: test build lint gen all clean update docker run

# Define the default goal. When you run "make" without argument, it will run the "all" target.
default: all

# Capture additional arguments which can optionally be passed in.
ARGS ?=

update:
	go get -u
	go mod tidy

clean:
	rm -f *db.bolt

# Test the code.
test:
	go test -v ./...

# Lint the code.
lint:
	golangci-lint run

# Run any code generation on this step.
gen:
	go run testmode/gen_mock/main.go > testmode/devices.json

docker:
	docker build -t docker-tips .

run:
	docker run -v .:/home/tipsuser docker-tips:latest

# Build the project: run the linter and then build.
build: lint
	go build
	@echo "Build: Successful"

# Run all steps: build and then run the application.
# To forward args optionally set ARGS like: make ARGS="--flag1=value1 --flag2=value2"
all: build
	./tips $(ARGS)
