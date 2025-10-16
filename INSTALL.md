# Installation Guide

This guide covers installation on macOS, Linux, and Windows.

## Prerequisites

- Go 1.19 or later
- Git
- OpenAI API key

## Quick Start

### Option 1: Using Make (Recommended for macOS/Linux)

```bash
# Clone the repository
git clone https://github.com/yourusername/revyu.git
cd revyu

# Create .env file with your API key (optional but recommended)
echo "OPENAI_API_KEY=your-api-key-here" > .env

# Build and install
make install
```

### Option 2: Using Go Install

```bash
go install github.com/yourusername/revyu@latest

# Set environment variable
export OPENAI_API_KEY=your-api-key-here
```

## Platform-Specific Installation

### macOS

#### Method 1: Make (Recommended)
```bash
# Build with embedded API key
make install

# The binary will be installed to /usr/local/bin
# Code signing will be applied automatically
```

#### Method 2: Build Script
```bash
# Build
./build.sh

# Install manually
sudo cp revyu /usr/local/bin/
```

### Linux

#### Method 1: Make (Recommended)
```bash
# Build and install
make install

# The binary will be installed to /usr/local/bin
```

#### Method 2: Build Script
```bash
# Make script executable
chmod +x build.sh

# Build
./build.sh

# Install manually
sudo cp revyu /usr/local/bin/
```

#### Method 3: Build from source
```bash
# Build
go build -o revyu

# Install
sudo cp revyu /usr/local/bin/
sudo chmod +x /usr/local/bin/revyu
```

### Windows

#### Method 1: PowerShell Build Script (Recommended)
```powershell
# Build
.\build.ps1

# Install (creates directory if needed)
make install

# Or manually copy to a directory in your PATH
New-Item -ItemType Directory -Force -Path $env:USERPROFILE\bin
Copy-Item revyu.exe $env:USERPROFILE\bin\

# Add to PATH if not already there
$env:PATH += ";$env:USERPROFILE\bin"
```

#### Method 2: Make (requires Make for Windows)
```bash
# Using Git Bash, WSL, or Make for Windows
make install
```

#### Method 3: Build from source
```bash
# Build
go build -o revyu.exe

# Copy to a directory in your PATH
copy revyu.exe %USERPROFILE%\bin\
```

## Configuration

### API Key Setup

You have three options for providing your OpenAI API key:

#### 1. Build-time Embedding (Most Secure)
Create a `.env` file and build with the API key embedded:

```bash
echo "OPENAI_API_KEY=your-api-key-here" > .env
make build
```

#### 2. Environment Variable
Set the environment variable in your shell:

**macOS/Linux:**
```bash
export OPENAI_API_KEY=your-api-key-here
```

Add to `~/.bashrc`, `~/.zshrc`, or `~/.profile` for persistence.

**Windows (PowerShell):**
```powershell
$env:OPENAI_API_KEY="your-api-key-here"
```

Add to your PowerShell profile for persistence:
```powershell
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', 'your-api-key-here', 'User')
```

**Windows (CMD):**
```cmd
setx OPENAI_API_KEY "your-api-key-here"
```

#### 3. .env File (Development)
Create a `.env` file in the project directory:

```bash
OPENAI_API_KEY=your-api-key-here
```

## Verification

After installation, verify it works:

```bash
# Check version and help
revyu

# Test on a file
echo "test" > test.txt
git init
git add test.txt
git commit -m "test"
echo "modified" > test.txt
revyu test.txt
```

## Updating

### Using Make
```bash
cd revyu
git pull
make reinstall
```

### Using Go Install
```bash
go install github.com/yourusername/revyu@latest
```

## Uninstalling

### Using Make
```bash
make uninstall
```

### Manual Removal

**macOS/Linux:**
```bash
sudo rm /usr/local/bin/revyu
```

**Windows:**
```powershell
Remove-Item $env:USERPROFILE\bin\revyu.exe
```

## Troubleshooting

### macOS: "revyu cannot be opened because the developer cannot be verified"
The build process includes code signing, but if you still see this:

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/revyu

# Or re-sign manually
codesign --force --options runtime --timestamp -s - /usr/local/bin/revyu
```

### Linux: Permission denied
Make sure the binary is executable:

```bash
chmod +x /usr/local/bin/revyu
```

### Windows: 'revyu' is not recognized
Make sure the installation directory is in your PATH:

```powershell
# Check PATH
$env:PATH

# Add to PATH (adjust path as needed)
[System.Environment]::SetEnvironmentVariable('PATH', $env:PATH + ";$env:USERPROFILE\bin", 'User')
```

### API Key Not Found
If you see "OPENAI_API_KEY not found":

1. Check if the environment variable is set: `echo $OPENAI_API_KEY` (Unix) or `echo %OPENAI_API_KEY%` (Windows)
2. Create a `.env` file in the current directory
3. Rebuild with embedded key using `make build`

## Build Options

### Development Build (No API Key)
```bash
go build -o revyu
```

### Production Build (With Embedded API Key)
```bash
# Using Make
make build

# Using build script
./build.sh           # macOS/Linux
.\build.ps1          # Windows

# Manual
go build -ldflags "-X main.buildTimeAPIKey=$OPENAI_API_KEY" -o revyu
```

### Cross-Platform Builds

Build for different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o revyu-linux

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o revyu-macos-intel

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o revyu-macos-arm

# Windows
GOOS=windows GOARCH=amd64 go build -o revyu.exe
```

## Support

For issues or questions:
- GitHub Issues: https://github.com/yourusername/revyu/issues
- Documentation: https://github.com/yourusername/revyu

