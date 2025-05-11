#!/bin/sh
set -e

# Ensure bin directory exists
mkdir -p bin

echo "Building Go binary from ./cmd/main.go..."
go build -o bin/server/subitocart ./cmd/main.go
echo "Build completed. Binary located at ./bin/server"
