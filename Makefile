.PHONY: test build lint gen all clean update docker run format stats

PKG=github.com/deckarep/tips
VERSION := `git fetch --tags && git tag | sort -V | tail -1`
#VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS=-ldflags "-X=github.com/deckarep/tips/pkg.AppVersion=$(VERSION)"
COVER=--cover --coverprofile=cover.out

# Define the default goal. When you run "make" without argument, it will run the "all" target.
default: all

# Capture additional arguments which can optionally be passed in.
ARGS ?=

stats:
	git diff --stat HEAD~$(ARGS) HEAD

update:
	go get -u
	go mod tidy

clean:
	rm -f *db.bolt

format:
	go fmt ./...

# Test the code.
test:
	go test -v ./... --race $(COVER) $(PKG)
	go tool cover -html=cover.out
	go tool cover -func cover.out | grep statements

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
	go build $(LDFLAGS)
	@echo "Build: Successful"

# Run all steps: build and then run the application.
# To forward args optionally set ARGS like: make ARGS="--flag1=value1 --flag2=value2"
all: build
	./tips $(ARGS)
