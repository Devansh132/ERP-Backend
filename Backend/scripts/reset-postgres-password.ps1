# Reset PostgreSQL Password Script
# This script helps you reset the postgres user password

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PostgreSQL Password Reset Helper" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Add PostgreSQL to PATH
$pgBin = "C:\Program Files\PostgreSQL\17\bin"
if (Test-Path $pgBin) {
    $env:Path += ";$pgBin"
}

$pgData = "C:\Program Files\PostgreSQL\17\data"
$pgHbaPath = Join-Path $pgData "pg_hba.conf"

Write-Host "This script will help you reset the PostgreSQL password." -ForegroundColor Yellow
Write-Host ""
Write-Host "Method 1: Try Windows Authentication first" -ForegroundColor Cyan
Write-Host ""

# Try Windows authentication
Write-Host "Attempting Windows authentication..." -ForegroundColor Yellow
$winAuthTest = & psql -U $env:USERNAME -h localhost -d postgres -c "SELECT 1;" 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Windows authentication works!" -ForegroundColor Green
    Write-Host ""
    $newPassword = Read-Host "Enter new password for 'postgres' user" -AsSecureString
    $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($newPassword))
    
    $resetCmd = "ALTER USER postgres WITH PASSWORD '$passwordPlain';"
    $result = & psql -U $env:USERNAME -h localhost -d postgres -c $resetCmd 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Password reset successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Update your .env file with:" -ForegroundColor Yellow
        Write-Host "DB_PASSWORD=$passwordPlain" -ForegroundColor White
        exit 0
    } else {
        Write-Host "❌ Failed to reset password: $result" -ForegroundColor Red
    }
} else {
    Write-Host "⚠️  Windows authentication not available" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Method 2: Use pg_hba.conf trust method" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  This will temporarily allow connections without password" -ForegroundColor Yellow
Write-Host "⚠️  We'll restore security after resetting password" -ForegroundColor Yellow
Write-Host ""

$continue = Read-Host "Continue with Method 2? (y/n)"
if ($continue -ne "y" -and $continue -ne "Y") {
    Write-Host "Cancelled." -ForegroundColor Yellow
    exit 0
}

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "⚠️  This method requires Administrator privileges" -ForegroundColor Yellow
    Write-Host "Please run PowerShell as Administrator and try again" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Alternative: Use pgAdmin GUI to reset password" -ForegroundColor Cyan
    exit 1
}

# Backup pg_hba.conf
Write-Host "Backing up pg_hba.conf..." -ForegroundColor Yellow
if (Test-Path $pgHbaPath) {
    Copy-Item $pgHbaPath "$pgHbaPath.backup" -Force
    Write-Host "✅ Backup created" -ForegroundColor Green
} else {
    Write-Host "❌ pg_hba.conf not found at: $pgHbaPath" -ForegroundColor Red
    exit 1
}

# Read and modify pg_hba.conf
Write-Host "Modifying pg_hba.conf..." -ForegroundColor Yellow
$content = Get-Content $pgHbaPath

$modified = $false
$newContent = $content | ForEach-Object {
    if ($_ -match "^\s*host\s+all\s+all\s+(127\.0\.0\.1/32|::1/128)\s+scram-sha-256") {
        $modified = $true
        $_ -replace "scram-sha-256", "trust"
    } else {
        $_
    }
}

if ($modified) {
    $newContent | Set-Content $pgHbaPath
    Write-Host "✅ pg_hba.conf modified (trust mode enabled)" -ForegroundColor Green
} else {
    Write-Host "⚠️  Could not find expected lines in pg_hba.conf" -ForegroundColor Yellow
    Write-Host "You may need to edit it manually" -ForegroundColor Yellow
}

# Restart PostgreSQL
Write-Host "Restarting PostgreSQL service..." -ForegroundColor Yellow
Restart-Service postgresql-x64-17 -ErrorAction SilentlyContinue
Start-Sleep -Seconds 3

# Try to connect
Write-Host "Attempting to connect without password..." -ForegroundColor Yellow
$testConn = & psql -U postgres -h localhost -c "SELECT 1;" 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Connected successfully!" -ForegroundColor Green
    Write-Host ""
    
    $newPassword = Read-Host "Enter new password for 'postgres' user" -AsSecureString
    $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($newPassword))
    
    $resetCmd = "ALTER USER postgres WITH PASSWORD '$passwordPlain';"
    $result = & psql -U postgres -h localhost -c $resetCmd 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Password reset successfully!" -ForegroundColor Green
        
        # Restore pg_hba.conf
        Write-Host "Restoring pg_hba.conf security..." -ForegroundColor Yellow
        $restoreContent = $newContent | ForEach-Object {
            $_ -replace "trust", "scram-sha-256"
        }
        $restoreContent | Set-Content $pgHbaPath
        Restart-Service postgresql-x64-17 -ErrorAction SilentlyContinue
        
        Write-Host "✅ Security restored" -ForegroundColor Green
        Write-Host ""
        Write-Host "Update your .env file with:" -ForegroundColor Yellow
        Write-Host "DB_PASSWORD=$passwordPlain" -ForegroundColor White
    } else {
        Write-Host "❌ Failed to reset password: $result" -ForegroundColor Red
        Write-Host "Restoring pg_hba.conf from backup..." -ForegroundColor Yellow
        Copy-Item "$pgHbaPath.backup" $pgHbaPath -Force
        Restart-Service postgresql-x64-17 -ErrorAction SilentlyContinue
    }
} else {
    Write-Host "❌ Could not connect: $testConn" -ForegroundColor Red
    Write-Host "Restoring pg_hba.conf from backup..." -ForegroundColor Yellow
    Copy-Item "$pgHbaPath.backup" $pgHbaPath -Force
    Restart-Service postgresql-x64-17 -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan



