# PowerShell build script for Windows

$BinaryName = "revyu.exe"
$EnvFile = ".env"

Write-Host "Building Revyu for Windows..." -ForegroundColor Cyan

# Load .env file if it exists
if (Test-Path $EnvFile) {
    Get-Content $EnvFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "Loaded .env file" -ForegroundColor Green
} else {
    Write-Host "Warning: .env file not found" -ForegroundColor Yellow
    Write-Host "Building without embedded API key" -ForegroundColor Yellow
    Write-Host "You'll need to set OPENAI_API_KEY environment variable or create .env file" -ForegroundColor Yellow
}

# Build with or without embedded API key
$ApiKey = $env:OPENAI_API_KEY

if ([string]::IsNullOrEmpty($ApiKey)) {
    Write-Host "Building $BinaryName without embedded API key..." -ForegroundColor Yellow
    go build -o $BinaryName
} else {
    Write-Host "Building $BinaryName with embedded OpenAI API key..." -ForegroundColor Green
    go build -ldflags "-X main.buildTimeAPIKey=$ApiKey" -o $BinaryName
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Binary created: .\$BinaryName" -ForegroundColor Green
    Write-Host ""
    Write-Host "To install globally, you can:" -ForegroundColor Cyan
    Write-Host "  1. Run: make install" -ForegroundColor White
    Write-Host "  2. Or manually copy to a directory in your PATH" -ForegroundColor White
    Write-Host "     Example: Copy-Item $BinaryName `$env:USERPROFILE\bin\" -ForegroundColor White
} else {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}

