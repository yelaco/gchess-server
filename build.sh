#!/bin/bash

# Variables
BINARY_NAME="go-chess-server"
PKG="./cmd/game-server/server.go"

# Functions
build() {
    echo "Building the project..."
    go build -o $BINARY_NAME $PKG
}

test() {
    echo "Running tests..."
    go test -v "./test/"
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

debug() {
    echo "Debugging the application..."
    go run $PKG
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
    debug)
        debug
        ;;
    *)
        echo "Usage: $0 {build|test|clean|lint|run|debug}"
        exit 1
        ;;
esac
