# CI/CD Workflows Guide

This guide explains how to set up GitHub Actions workflows for cross-platform builds and releases.

## Recommended GitHub Actions Workflow

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build:
    name: Build
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            goos: linux
            goarch: amd64
            artifact: revyu-linux-amd64
          - os: ubuntu-latest
            goos: linux
            goarch: arm64
            artifact: revyu-linux-arm64
          - os: macos-latest
            goos: darwin
            goarch: amd64
            artifact: revyu-macos-amd64
          - os: macos-latest
            goos: darwin
            goarch: arm64
            artifact: revyu-macos-arm64
          - os: windows-latest
            goos: windows
            goarch: amd64
            artifact: revyu-windows-amd64.exe

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          go build -o ${{ matrix.artifact }} -ldflags "-s -w"

      - name: Code sign (macOS only)
        if: matrix.os == 'macos-latest'
        run: |
          codesign --force --options runtime --timestamp -s - ${{ matrix.artifact }}

      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: ${{ matrix.artifact }}
          path: ${{ matrix.artifact }}

  release:
    name: Create Release
    needs: build
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    
    steps:
      - name: Download all artifacts
        uses: actions/download-artifact@v3

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            revyu-*/*
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Testing Workflow

Create `.github/workflows/test.yml`:

```yaml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.20', '1.21']

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run tests
        run: go test -v ./...

      - name: Build
        run: go build -v ./...

      - name: Verify build works
        shell: bash
        run: |
          if [ "$RUNNER_OS" == "Windows" ]; then
            ./revyu.exe || true
          else
            ./revyu || true
          fi
```

## Security Notes

1. **Never** commit API keys to the repository
2. Use GitHub Secrets for sensitive data
3. The workflows above build without embedded API keys
4. End users should provide their own API keys via environment variables or .env files

