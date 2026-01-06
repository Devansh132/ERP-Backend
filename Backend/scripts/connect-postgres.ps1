# Connect to PostgreSQL using full path
# This script finds PostgreSQL and connects without needing PATH

Write-Host "Finding PostgreSQL installation..." -ForegroundColor Yellow

$pgVersions = @("17", "16", "15", "14", "13", "12")
$psqlPath = $null

foreach ($version in $pgVersions) {
    $testPath = "C:\Program Files\PostgreSQL\$version\bin\psql.exe"
    if (Test-Path $testPath) {
        $psqlPath = $testPath
        Write-Host "✅ Found PostgreSQL $version" -ForegroundColor Green
        break
    }
}

if (-not $psqlPath) {
    Write-Host "❌ PostgreSQL not found" -ForegroundColor Red
    Write-Host "Please install PostgreSQL or provide the path manually" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "Connecting to PostgreSQL..." -ForegroundColor Yellow
Write-Host "Using: $psqlPath" -ForegroundColor Gray
Write-Host ""

# Connect to PostgreSQL
& $psqlPath -U postgres -h localhost



