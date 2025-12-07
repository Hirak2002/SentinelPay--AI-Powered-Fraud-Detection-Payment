#!/bin/bash

# SentinelPay Quick Start Script

echo "🚀 Starting SentinelPay..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "❌ Docker is not running. Please start Docker and try again."
  exit 1
fi

# Build and start all services
echo "📦 Building Docker images..."
docker-compose build

echo "🔄 Starting services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check service health
echo "🔍 Checking service health..."

# Check PostgreSQL
if docker-compose exec -T postgres pg_isready -U sentinelpay > /dev/null 2>&1; then
  echo "✅ PostgreSQL is ready"
else
  echo "⚠️  PostgreSQL is not ready yet"
fi

# Check Backend
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
  echo "✅ Backend service is ready"
else
  echo "⚠️  Backend service is not ready yet"
fi

# Check AI Service
if curl -f http://localhost:8001/health > /dev/null 2>&1; then
  echo "✅ AI service is ready"
else
  echo "⚠️  AI service is not ready yet"
fi

echo ""
echo "🎉 SentinelPay is running!"
echo ""
echo "📊 Admin Dashboard: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🤖 AI Service: http://localhost:8001"
echo "📬 RabbitMQ Management: http://localhost:15672 (guest/guest)"
echo ""
echo "To stop all services: docker-compose down"
echo "To view logs: docker-compose logs -f"
