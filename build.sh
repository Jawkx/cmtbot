#!/bin/bash

# Variables
PROJECT_NAME="cmtbot"
VERSION="0.0.1"
OUTPUT_DIR="./dist"
EXECUTABLE_NAME="$PROJECT_NAME"

# Ensure the output directory exists
mkdir -p "$OUTPUT_DIR"

# Download dependencies
echo "Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
  echo "Error: Failed to download dependencies."
  exit 1
fi

# Build the application with version information
echo "Building the application..."
LDFLAGS="-X main.version=$VERSION"

go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/$EXECUTABLE_NAME"
if [ $? -ne 0 ]; then
  echo "Error: Failed to build the application."
  exit 1
fi

echo "Packaging complete. Package: $OUTPUT_DIR/$PROJECT_NAME-$VERSION"

echo "Installing the application..."
go install -ldflags "$LDFLAGS"
if [ $? -ne 0 ]; then
    echo "Error: Failed to install the application."
    exit 1
fi

echo "Installation complete."
