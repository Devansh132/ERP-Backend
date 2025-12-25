# Quick Database Setup Script
# This script sets up everything you need

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Quick Database Setup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Add PostgreSQL to PATH for this session
$pgBin = "C:\Program Files\PostgreSQL\17\bin"
if (Test-Path $pgBin) {
    $env:Path += ";$pgBin"
    Write-Host "✅ Added PostgreSQL to PATH" -ForegroundColor Green
} else {
    Write-Host "❌ PostgreSQL not found at $pgBin" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Now you can use 'psql' command!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Test connection: psql -U postgres -h localhost" -ForegroundColor White
Write-Host "2. Create database: CREATE DATABASE school_erp;" -ForegroundColor White
Write-Host "3. Exit: \q" -ForegroundColor White
Write-Host "4. Update .env file with your password" -ForegroundColor White
Write-Host ""

# Ask if user wants to create database now
$createDb = Read-Host "Do you want to create the database now? (y/n)"
if ($createDb -eq "y" -or $createDb -eq "Y") {
    Write-Host ""
    Write-Host "Creating database 'school_erp'..." -ForegroundColor Yellow
    
    # Try to create database (will prompt for password)
    $result = & psql -U postgres -h localhost -c "CREATE DATABASE school_erp;" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database 'school_erp' created successfully!" -ForegroundColor Green
    } else {
        if ($result -like "*already exists*") {
            Write-Host "ℹ️  Database 'school_erp' already exists" -ForegroundColor Yellow
        } else {
            Write-Host "⚠️  Could not create database automatically" -ForegroundColor Yellow
            Write-Host "Please create it manually:" -ForegroundColor White
            Write-Host "  psql -U postgres -h localhost" -ForegroundColor Green
            Write-Host "  CREATE DATABASE school_erp;" -ForegroundColor Green
            Write-Host "  \q" -ForegroundColor Green
        }
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Setup complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Remember: PATH change is only for this PowerShell session." -ForegroundColor Yellow
Write-Host "To make it permanent, run: .\scripts\add-postgres-to-path.ps1 (as Admin)" -ForegroundColor Yellow
Write-Host ""

