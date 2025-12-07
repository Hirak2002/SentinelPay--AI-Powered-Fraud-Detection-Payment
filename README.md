# SentinelPay

Modern payment orchestration platform with intelligent fraud detection and dynamic provider routing.

---

## Overview

SentinelPay acts as an intelligent middleware layer between client applications and payment service providers. By leveraging machine learning algorithms and sophisticated routing logic, it minimizes fraud exposure while optimizing transaction costs across multiple payment gateways.

The platform evaluates each transaction through a multi-factor risk assessment system before dynamically selecting the most appropriate payment provider based on cost efficiency, reliability metrics, and transaction characteristics.

## Core Capabilities

**Fraud Detection & Prevention**  
Real-time risk scoring using Isolation Forest anomaly detection. Each transaction undergoes analysis across multiple dimensions including amount patterns, user behavior, and historical data to identify potentially fraudulent activity before processing.

**Dynamic Provider Selection**  
Intelligent routing algorithm that considers transaction size, provider fees, availability status, and historical success rates to optimize both cost and reliability for each payment.

**Scalable Architecture**  
Microservices design built for high throughput and horizontal scaling. Independent service components communicate through message queues, enabling reliable processing of thousands of concurrent transactions.

**Operational Insights**  
Comprehensive analytics dashboard providing real-time visibility into transaction volumes, fraud prevention metrics, revenue distribution, and system performance indicators.

## System Architecture

The platform consists of four primary components working in concert:

**Backend Service** (Go)  
Handles all API requests, transaction orchestration, and business logic. Implements connection pooling, caching strategies, and concurrent request processing for optimal performance.

**AI Engine** (Python)  
Runs the fraud detection model, performs risk calculations, and provides prediction endpoints. Pre-trained on historical transaction patterns with continuous model improvement capabilities.

**Admin Interface** (Next.js)  
React-based dashboard for monitoring transactions, reviewing fraud alerts, analyzing revenue metrics, and managing system configuration.

**Data Layer** (PostgreSQL + RabbitMQ)  
Persistent storage for all transaction records with message queue for asynchronous event processing and inter-service communication.

## Requirements

- Docker Desktop (recommended deployment method)
- Go 1.21 or later (local development)
- Python 3.11 or later (local development)
- Node.js 20 or later (local development)
- PostgreSQL 16 (if running without Docker)

## Quick Start

The fastest way to get SentinelPay running is with Docker:

```bash
git clone <your-repository-url>
cd SenitelPay
docker-compose up -d
```

This starts all services and dependencies. The system will be available at:

- Dashboard: http://localhost:3000
- API: http://localhost:8080
- AI Service: http://localhost:8001
- Queue Management: http://localhost:15672

For Windows users, run `.\start.ps1` which includes health checks and status reporting.

## Development Setup

To run services individually for development:

**Backend**
```bash
cd backend
cp .env.example .env
go mod download
go run .
```

**AI Service**
```bash
cd ai-service
cp .env.example .env
pip install -r requirements.txt
python main.py
```

**Admin Dashboard**
```bash
cd admin-dashboard
cp .env.local.example .env.local
npm install
npm run dev
```

Database schema initializes automatically with Docker. For manual setup, run:
```bash
psql -U sentinelpay -d sentinelpay -f db/schema.sql
```

## API Reference

### Transaction Processing

**Create Transaction**
```http
POST /api/v1/transactions
Authorization: Bearer {user_id}
Content-Type: application/json

{
  "user_id": "user_001",
  "amount": 150.00,
  "currency": "USD",
  "description": "Product purchase"
}
```

Returns transaction details including fraud risk score and selected payment provider.

**Retrieve Transaction**
```http
GET /api/v1/transactions/:id
Authorization: Bearer {user_id}
```

**List Transactions**
```http
GET /api/v1/transactions?limit=20&offset=0
Authorization: Bearer {user_id}
```

### Admin Operations

**System Statistics**
```http
GET /api/v1/admin/stats
X-Admin-Key: {admin_key}
```

Returns comprehensive metrics including transaction volumes, fraud prevention rates, and revenue totals.

**Blocked Transactions**
```http
GET /api/v1/admin/transactions/blocked
X-Admin-Key: {admin_key}
```

**Revenue Analytics**
```http
GET /api/v1/admin/revenue
X-Admin-Key: {admin_key}
```

### Fraud Detection

**Risk Prediction**
```http
POST http://localhost:8001/predict
Content-Type: application/json

{
  "user_id": "user_001",
  "amount": 500.00,
  "currency": "USD",
  "previous_transaction_count": 15,
  "account_age_days": 180
}
```

Returns risk score, fraud likelihood, identified risk factors, and recommended action.

Complete API documentation available in `API_DOCS.md`.

## Testing

Verify system health:
```bash
curl http://localhost:8080/health
curl http://localhost:8001/health
```

Create a test transaction:
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_001" \
  -d '{
    "user_id": "user_001",
    "amount": 150.00,
    "currency": "USD",
    "description": "Test payment"
  }'
```

Detailed testing procedures and scenarios documented in `TESTING.md`.

## Technology Stack

**Backend Infrastructure**
- Go 1.21 with Gin web framework
- PostgreSQL 16 for transaction storage
- RabbitMQ 3.12 for message queuing
- In-memory caching layer

**AI/ML Components**
- Python 3.11 with FastAPI
- Scikit-learn for model training
- Isolation Forest algorithm for anomaly detection
- NumPy and Pandas for data processing

**Frontend**
- Next.js 14 with React 18
- TypeScript for type safety
- Tailwind CSS for styling
- Recharts for data visualization

**Infrastructure**
- Docker and Docker Compose
- Multi-stage builds for optimization
- Health checks and graceful shutdowns

## Project Structure

```
SenitelPay/
├── backend/              Go service handling core business logic
├── ai-service/           Python ML service for fraud detection
├── admin-dashboard/      Next.js admin interface
├── db/                   Database schemas and migrations
├── docker-compose.yml    Service orchestration configuration
└── docs/                 Additional documentation
```

## Configuration

Environment variables are managed through `.env` files in each service directory. Templates are provided as `.env.example` files.

**Backend Configuration**
- Database credentials
- RabbitMQ connection string
- AI service endpoint
- Server port and environment mode

**AI Service Configuration**
- Model parameters
- Prediction thresholds
- Server port

**Dashboard Configuration**
- Backend API endpoint URL

Never commit actual `.env` files to version control.

## Security Considerations

Current implementation includes:
- Bearer token authentication for API requests
- Admin key verification for sensitive endpoints
- SQL injection prevention through parameterized queries
- CORS configuration for cross-origin requests
- Input validation on all endpoints

For production deployment, additionally implement:
- JWT-based authentication with token expiration
- Rate limiting per client/IP
- HTTPS/TLS encryption for all communications
- Secrets management system (e.g., HashiCorp Vault)
- Database encryption at rest
- API key rotation policies
- Web Application Firewall (WAF)

## Performance

The system is designed to handle high transaction volumes:

- Concurrent request processing through Go goroutines
- Database connection pooling (configurable, default 25 max connections)
- In-memory caching with TTL for frequently accessed data
- Asynchronous event processing via message queue
- Indexed database columns for optimized queries

Observed performance metrics:
- API response time: < 200ms (p95)
- Fraud detection latency: < 50ms
- Throughput: 100+ transactions/second per instance

Horizontal scaling supported for all services except PostgreSQL (use read replicas for scaling reads).

## Monitoring

Service health endpoints:
- Backend: `GET /health`
- AI Service: `GET /health`

Log aggregation recommended for production:
- Structured JSON logging enabled in all services
- Compatible with ELK stack, Splunk, or CloudWatch
- Correlation IDs for request tracing

## Contributing

When contributing to this project:

1. Follow existing code style and conventions
2. Write tests for new functionality
3. Update documentation for API changes
4. Ensure Docker build succeeds
5. Verify all services start without errors

## License

This project is available under the MIT License.

## Documentation

- `API_DOCS.md` - Complete API endpoint reference
- `TESTING.md` - Testing guide and scenarios
- `DEVELOPMENT.md` - Architecture and technical details
- `QUICKSTART.md` - Simplified getting started guide
