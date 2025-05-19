#!/bin/bash

echo "Building Public IP Updater for DigitalOcean DNS..."

# Check if build directory exists, create if not
if [ ! -d "build" ]; then
    mkdir build
fi

# Build the executable
if [[ "$OSTYPE" == "darwin"* ]]; then
    # Mac OS
    echo "Building for macOS..."
    go build -o build/pub_ip_updater main.go
else
    # Linux
    echo "Building for Linux..."
    go build -o build/pub_ip_updater main.go
fi

# Check if build was successful
if [ $? -eq 0 ]; then
    echo -e "\nBuild successful! Executable created at build/pub_ip_updater"
    # Make the file executable
    chmod +x build/pub_ip_updater
else
    echo -e "\nBuild failed with error code $?"
fi
