# Test Database Connection Script
# This script helps you test and configure your database connection

param(
    [string]$DBUser = "postgres",
    [string]$DBPassword = "",
    [string]$DBHost = "localhost",
    [string]$DBPort = "5432",
    [string]$DBName = "school_erp"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Database Connection Test" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if psql is available
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    Write-Host "⚠️  psql not found in PATH" -ForegroundColor Yellow
    Write-Host "Trying default PostgreSQL installation path..." -ForegroundColor Yellow
    
    $defaultPaths = @(
        "C:\Program Files\PostgreSQL\17\bin\psql.exe",
        "C:\Program Files\PostgreSQL\16\bin\psql.exe",
        "C:\Program Files\PostgreSQL\15\bin\psql.exe",
        "C:\Program Files\PostgreSQL\14\bin\psql.exe"
    )
    
    $found = $false
    foreach ($path in $defaultPaths) {
        if (Test-Path $path) {
            $env:PATH += ";$(Split-Path $path)"
            $found = $true
            Write-Host "✅ Found PostgreSQL at: $path" -ForegroundColor Green
            break
        }
    }
    
    if (-not $found) {
        Write-Host "❌ PostgreSQL not found. Please install PostgreSQL or add it to PATH." -ForegroundColor Red
        exit 1
    }
}

Write-Host "Testing connection to PostgreSQL..." -ForegroundColor Yellow
Write-Host "User: $DBUser" -ForegroundColor White
Write-Host "Host: $DBHost" -ForegroundColor White
Write-Host "Port: $DBPort" -ForegroundColor White
Write-Host ""

# Set password environment variable
if ($DBPassword) {
    $env:PGPASSWORD = $DBPassword
}

# Test connection
Write-Host "Attempting to connect..." -ForegroundColor Yellow
$testQuery = "SELECT version();"
$result = & psql -U $DBUser -h $DBHost -p $DBPort -d postgres -c $testQuery 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Connection successful!" -ForegroundColor Green
    Write-Host ""
    
    # Check if database exists
    Write-Host "Checking if database '$DBName' exists..." -ForegroundColor Yellow
    $dbCheck = & psql -U $DBUser -h $DBHost -p $DBPort -d postgres -t -c "SELECT 1 FROM pg_database WHERE datname = '$DBName';" 2>&1
    
    if ($dbCheck -match "1") {
        Write-Host "✅ Database '$DBName' already exists" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Database '$DBName' does not exist" -ForegroundColor Yellow
        Write-Host "Creating database..." -ForegroundColor Yellow
        
        $createDb = & psql -U $DBUser -h $DBHost -p $DBPort -d postgres -c "CREATE DATABASE $DBName;" 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Database '$DBName' created successfully!" -ForegroundColor Green
        } else {
            Write-Host "❌ Failed to create database: $createDb" -ForegroundColor Red
        }
    }
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "✅ Database is ready!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Update your .env file with:" -ForegroundColor Yellow
    Write-Host "DB_PASSWORD=$DBPassword" -ForegroundColor White
    Write-Host ""
    
} else {
    Write-Host "❌ Connection failed!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Error details:" -ForegroundColor Yellow
    Write-Host $result -ForegroundColor Red
    Write-Host ""
    Write-Host "Troubleshooting:" -ForegroundColor Yellow
    Write-Host "1. Check if PostgreSQL is running" -ForegroundColor White
    Write-Host "2. Verify username and password" -ForegroundColor White
    Write-Host "3. Check if port $DBPort is correct" -ForegroundColor White
    Write-Host "4. Try connecting manually: psql -U $DBUser -h $DBHost" -ForegroundColor White
    Write-Host ""
    Write-Host "To reset password:" -ForegroundColor Yellow
    Write-Host "  psql -U postgres" -ForegroundColor White
    Write-Host "  ALTER USER postgres WITH PASSWORD 'new_password';" -ForegroundColor White
    Write-Host ""
}

# Clean up
if ($env:PGPASSWORD) {
    Remove-Item Env:\PGPASSWORD
}



