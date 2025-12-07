# SentinelPay Quick Start Script for Windows

Write-Host "🚀 Starting SentinelPay..." -ForegroundColor Green

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker is not running. Please start Docker and try again." -ForegroundColor Red
    exit 1
}

# Build and start all services
Write-Host "📦 Building Docker images..." -ForegroundColor Yellow
docker-compose build

Write-Host "🔄 Starting services..." -ForegroundColor Yellow
docker-compose up -d

# Wait for services to be ready
Write-Host "⏳ Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check service health
Write-Host "🔍 Checking service health..." -ForegroundColor Yellow

# Check PostgreSQL
try {
    docker-compose exec -T postgres pg_isready -U sentinelpay | Out-Null
    Write-Host "✅ PostgreSQL is ready" -ForegroundColor Green
} catch {
    Write-Host "⚠️  PostgreSQL is not ready yet" -ForegroundColor Yellow
}

# Check Backend
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing -TimeoutSec 5
    Write-Host "✅ Backend service is ready" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Backend service is not ready yet" -ForegroundColor Yellow
}

# Check AI Service
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8001/health" -UseBasicParsing -TimeoutSec 5
    Write-Host "✅ AI service is ready" -ForegroundColor Green
} catch {
    Write-Host "⚠️  AI service is not ready yet" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🎉 SentinelPay is running!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Admin Dashboard: http://localhost:3000" -ForegroundColor Cyan
Write-Host "🔧 Backend API: http://localhost:8080" -ForegroundColor Cyan
Write-Host "🤖 AI Service: http://localhost:8001" -ForegroundColor Cyan
Write-Host "📬 RabbitMQ Management: http://localhost:15672 (guest/guest)" -ForegroundColor Cyan
Write-Host ""
Write-Host "To stop all services: docker-compose down" -ForegroundColor Gray
Write-Host "To view logs: docker-compose logs -f" -ForegroundColor Gray
