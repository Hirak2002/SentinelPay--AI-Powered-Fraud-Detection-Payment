# SentinelPay - System Verification Script

Write-Host "╔════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║     SentinelPay System Verification      ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

$allGood = $true

# Check Docker
Write-Host "🔍 Checking Docker..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    Write-Host "✅ Docker is installed" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not installed" -ForegroundColor Red
    $allGood = $false
}

# Check if Docker is running
Write-Host "🔍 Checking if Docker is running..." -ForegroundColor Yellow
try {
    docker ps | Out-Null
    Write-Host "✅ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker Desktop" -ForegroundColor Red
    $allGood = $false
}

Write-Host ""
Write-Host "🔍 Checking project structure..." -ForegroundColor Yellow

# Check directories
$directories = @(
    "backend",
    "ai-service",
    "admin-dashboard",
    "db"
)

foreach ($dir in $directories) {
    if (Test-Path $dir) {
        Write-Host "✅ $dir directory exists" -ForegroundColor Green
    } else {
        Write-Host "❌ $dir directory missing" -ForegroundColor Red
        $allGood = $false
    }
}

# Check key files
Write-Host ""
Write-Host "🔍 Checking key files..." -ForegroundColor Yellow

$files = @(
    "docker-compose.yml",
    "backend\main.go",
    "backend\Dockerfile",
    "ai-service\main.py",
    "ai-service\Dockerfile",
    "admin-dashboard\package.json",
    "admin-dashboard\Dockerfile",
    "db\schema.sql"
)

foreach ($file in $files) {
    if (Test-Path $file) {
        Write-Host "✅ $file exists" -ForegroundColor Green
    } else {
        Write-Host "❌ $file missing" -ForegroundColor Red
        $allGood = $false
    }
}

Write-Host ""
if ($allGood) {
    Write-Host "╔════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║   ✅ All checks passed!                   ║" -ForegroundColor Green
    Write-Host "║   System is ready to start                ║" -ForegroundColor Green
    Write-Host "╚════════════════════════════════════════════╝" -ForegroundColor Green
    Write-Host ""
    Write-Host "To start SentinelPay, run:" -ForegroundColor Cyan
    Write-Host "  .\start.ps1" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Or manually with:" -ForegroundColor Cyan
    Write-Host "  docker-compose up -d" -ForegroundColor Yellow
} else {
    Write-Host "╔════════════════════════════════════════════╗" -ForegroundColor Red
    Write-Host "║   ⚠️  Some checks failed                  ║" -ForegroundColor Red
    Write-Host "║   Please fix the issues above             ║" -ForegroundColor Red
    Write-Host "╚════════════════════════════════════════════╝" -ForegroundColor Red
}

Write-Host ""
