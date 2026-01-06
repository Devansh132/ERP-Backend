# Add PostgreSQL to PATH (Run as Administrator)
# This script adds PostgreSQL bin directory to system PATH

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PostgreSQL PATH Configuration" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Find PostgreSQL installation
$pgVersions = @("17", "16", "15", "14", "13", "12")
$pgPath = $null

foreach ($version in $pgVersions) {
    $testPath = "C:\Program Files\PostgreSQL\$version\bin"
    if (Test-Path $testPath) {
        $pgPath = $testPath
        Write-Host "✅ Found PostgreSQL $version at: $pgPath" -ForegroundColor Green
        break
    }
}

if (-not $pgPath) {
    Write-Host "❌ PostgreSQL not found in standard location" -ForegroundColor Red
    Write-Host "Please install PostgreSQL or provide the installation path" -ForegroundColor Yellow
    exit 1
}

# Check if already in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -like "*$pgPath*") {
    Write-Host "✅ PostgreSQL is already in PATH" -ForegroundColor Green
    Write-Host ""
    Write-Host "To use it in current session, run:" -ForegroundColor Yellow
    Write-Host "`$env:Path += `";$pgPath`"" -ForegroundColor White
    exit 0
}

# Add to PATH
Write-Host "Adding PostgreSQL to system PATH..." -ForegroundColor Yellow

try {
    $newPath = $currentPath + ";$pgPath"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    Write-Host "✅ PostgreSQL added to PATH successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "⚠️  You need to restart PowerShell for changes to take effect" -ForegroundColor Yellow
    Write-Host "Or run this in your current session:" -ForegroundColor Yellow
    Write-Host "`$env:Path += `";$pgPath`"" -ForegroundColor White
} catch {
    Write-Host "❌ Failed to add to PATH. Please run PowerShell as Administrator" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
}



