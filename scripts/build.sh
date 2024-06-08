#!/bin/bash

# Variables
BINARY_NAME="robinhood_chess"
PKG="../..."

# Functions
build() {
    echo "Building the project..."
    go build -o $BINARY_NAME $PKG
}

test() {
    echo "Running tests..."
    go test -v $PKG
}

clean() {
    echo "Cleaning up..."
    rm -f $BINARY_NAME
}

lint() {
    echo "Running linter..."
    golangci-lint run
}

run() {
    build
    echo "Running the application..."
    ./$BINARY_NAME
}

# Main script
case $1 in
    build)
        build
        ;;
    test)
        test
        ;;
    clean)
        clean
        ;;
    lint)
        lint
        ;;
    run)
        run
        ;;
    *)
        echo "Usage: $0 {build|test|clean|lint|run}"
        exit 1
        ;;
esac
