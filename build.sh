#!/bin/bash

# Detect OS
OS_TYPE=$(uname -s)
BINARY_NAME="revyu"

# Load .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
    echo "Loaded .env file"
else
    echo "Warning: .env file not found"
    echo "Building without embedded API key"
    echo "You'll need to set OPENAI_API_KEY environment variable or create .env file"
fi

# Build with or without embedded API key
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Building $BINARY_NAME without embedded API key..."
    go build -o $BINARY_NAME
else
    echo "Building $BINARY_NAME with embedded OpenAI API key..."
    go build -ldflags "-X main.buildTimeAPIKey=$OPENAI_API_KEY" -o $BINARY_NAME
fi

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Binary created: ./$BINARY_NAME"
    
    # Code sign only on macOS
    if [ "$OS_TYPE" = "Darwin" ]; then
        echo "Detected macOS - signing binary..."
        codesign --force --options runtime --timestamp -s - ./$BINARY_NAME 2>/dev/null
        
        if [ $? -eq 0 ]; then
            echo "Code signing successful!"
        else
            echo "Warning: Code signing failed, but binary was created"
        fi
    fi
    
    echo ""
    echo "To install globally, run:"
    if [ "$OS_TYPE" = "Darwin" ] || [ "$OS_TYPE" = "Linux" ]; then
        echo "  sudo cp $BINARY_NAME /usr/local/bin/"
    else
        echo "  make install"
    fi
else
    echo "Build failed"
    exit 1
fi
