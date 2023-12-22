#!/usr/bin/env bash

# brew install golangci-lint
golangci-lint run && go build && echo "Build: Successful"