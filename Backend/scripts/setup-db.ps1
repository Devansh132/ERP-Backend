# PowerShell script to set up database for School ERP System
# This script helps create the database and run migrations

param(
    [string]$DBDriver = "postgres",
    [string]$DBHost = "localhost",
    [string]$DBPort = "5432",
    [string]$DBUser = "postgres",
    [string]$DBPassword = "postgres",
    [string]$DBName = "school_erp"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "School ERP - Database Setup Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($DBDriver -eq "postgres") {
    Write-Host "Setting up PostgreSQL database..." -ForegroundColor Yellow
    
    # Create database if it doesn't exist
    Write-Host "Creating database '$DBName' if it doesn't exist..." -ForegroundColor Green
    $env:PGPASSWORD = $DBPassword
    psql -U $DBUser -h $DBHost -p $DBPort -c "SELECT 1 FROM pg_database WHERE datname = '$DBName'" | Out-Null
    if ($LASTEXITCODE -ne 0) {
        psql -U $DBUser -h $DBHost -p $DBPort -c "CREATE DATABASE $DBName"
        Write-Host "Database '$DBName' created successfully!" -ForegroundColor Green
    } else {
        Write-Host "Database '$DBName' already exists." -ForegroundColor Yellow
    }
    
    # Run migrations
    Write-Host "Running migrations..." -ForegroundColor Green
    $env:PGPASSWORD = $DBPassword
    psql -U $DBUser -h $DBHost -p $DBPort -d $DBName -f "..\migrations\001_initial_schema.sql"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migrations completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Error running migrations. Please check the SQL file." -ForegroundColor Red
    }
    
} elseif ($DBDriver -eq "mysql") {
    Write-Host "Setting up MySQL database..." -ForegroundColor Yellow
    
    # Create database if it doesn't exist
    Write-Host "Creating database '$DBName' if it doesn't exist..." -ForegroundColor Green
    mysql -u $DBUser -p$DBPassword -h $DBHost -P $DBPort -e "CREATE DATABASE IF NOT EXISTS $DBName"
    
    # Run migrations
    Write-Host "Running migrations..." -ForegroundColor Green
    mysql -u $DBUser -p$DBPassword -h $DBHost -P $DBPort $DBName -e "source migrations/001_initial_schema.sql"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migrations completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Error running migrations. Please check the SQL file." -ForegroundColor Red
    }
} else {
    Write-Host "Unsupported database driver: $DBDriver" -ForegroundColor Red
    Write-Host "Supported drivers: postgres, mysql" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Database setup completed!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Update .env file with your database credentials" -ForegroundColor White
Write-Host "2. Run: go run cmd/server/main.go" -ForegroundColor White
Write-Host "3. Access Swagger at: http://localhost:8080/swagger/index.html" -ForegroundColor White



