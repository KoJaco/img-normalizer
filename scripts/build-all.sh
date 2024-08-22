#!/bin/bash

# Create the output directory if it doesn't exist
mkdir -p build

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/img-normalizer-linux-amd64 ./cmd/img-normalizer


echo "Build complete. Binaries are located in the build/ directory."