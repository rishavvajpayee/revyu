#!/bin/bash

# Load .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "Error: .env file not found"
    echo "Please create a .env file with your OPENAI_API_KEY"
    exit 1
fi

# Check if API key is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Error: OPENAI_API_KEY not found in .env file"
    exit 1
fi

# Build with embedded API key
echo "Building revyu with embedded OpenAI API key..."
CGO_ENABLED=1 go build -ldflags "-X main.buildTimeAPIKey=$OPENAI_API_KEY" -o revyu

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Binary created: ./revyu"
    
    # Code sign the binary to prevent macOS from killing it
    echo "Signing binary..."
    codesign --options runtime --timestamp -s - ./revyu
    
    if [ $? -eq 0 ]; then
        echo "Code signing successful!"
    else
        echo "Warning: Code signing failed, but binary was created"
    fi
else
    echo "Build failed"
    echo "Please check your API key and try again"
    exit 1
fi
