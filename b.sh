#!/bin/bash

echo "Building Public IP Updater for DigitalOcean DNS..."

# Check if build directory exists, create if not
if [ ! -d "build" ]; then
    mkdir build
fi

# Check if pre-commit hook is installed
if [ ! -f ".git/hooks/pre-commit" ]; then
    echo "Installing pre-commit hook for security..."
    cp pre-commit.hook .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
fi

# Build the executable with custom name
if [[ "$OSTYPE" == "darwin"* ]]; then
    # Mac OS
    echo "Building for macOS..."
    go build -o build/ip_updater main.go
else
    # Linux
    echo "Building for Linux..."
    go build -o build/ip_updater main.go
fi

# Check if build was successful
if [ $? -eq 0 ]; then
    echo -e "\nBuild successful! Executable created at build/ip_updater"
    # Make the file executable
    chmod +x build/ip_updater
    
    echo -e "\n=== SECURITY REMINDER ==="
    echo "- Never commit your .env file with API tokens"
    echo "- The pre-commit hook helps prevent accidental token exposure"
    echo "- Remember to configure your .env file before running the application"
    echo "=========================="
else
    echo -e "\nBuild failed with error code $?"
fi
