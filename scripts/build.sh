#!/bin/sh
set -e

echo "Building Go binary from ./cmd/main.go..."
go build -o bin/server ./cmd/main.go
echo "Build completed. Binary located at ./bin/server"
