# SentinelPay - Development Notes

## Project Structure

```
SenitelPay/
├── backend/               # Go backend service
│   ├── main.go           # Application entry point
│   ├── database.go       # Database connection management
│   ├── models.go         # Data models
│   ├── handlers.go       # HTTP request handlers
│   ├── middleware.go     # Authentication & CORS
│   ├── broker.go         # RabbitMQ integration
│   ├── cache.go          # In-memory caching
│   ├── providers.go      # Payment provider logic
│   ├── go.mod            # Go dependencies
│   ├── Dockerfile        # Docker configuration
│   └── .env.example      # Environment variables template
│
├── ai-service/           # Python AI service
│   ├── main.py           # FastAPI application
│   ├── requirements.txt  # Python dependencies
│   ├── Dockerfile        # Docker configuration
│   └── .env.example      # Environment variables template
│
├── admin-dashboard/      # Next.js admin UI
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx      # Root layout
│   │   │   ├── page.tsx        # Dashboard page
│   │   │   └── globals.css     # Global styles
│   │   ├── components/
│   │   │   ├── StatsCard.tsx   # Statistics card
│   │   │   ├── TransactionChart.tsx  # Revenue chart
│   │   │   └── TransactionTable.tsx  # Transaction table
│   │   └── lib/
│   │       └── api.ts          # API client
│   ├── package.json      # Node dependencies
│   ├── tsconfig.json     # TypeScript config
│   ├── tailwind.config.ts # Tailwind CSS config
│   ├── Dockerfile        # Docker configuration
│   └── .env.local.example # Environment variables
│
├── db/
│   └── schema.sql        # Database schema
│
├── docker-compose.yml    # Multi-container orchestration
├── README.md             # Project documentation
├── TESTING.md            # Testing guide
├── API_DOCS.md           # API documentation
├── start.sh              # Linux/Mac startup script
├── start.ps1             # Windows startup script
└── .gitignore            # Git ignore rules

```

## Key Features Implemented

### 1. Backend (Go)
- ✅ RESTful API with Gin framework
- ✅ PostgreSQL integration with connection pooling
- ✅ Transaction management with ACID compliance
- ✅ In-memory caching for performance
- ✅ RabbitMQ message broker integration
- ✅ CORS and authentication middleware
- ✅ Smart payment routing logic
- ✅ Webhook handlers for Stripe & PayPal
- ✅ Admin dashboard API endpoints

### 2. AI Service (Python)
- ✅ FastAPI framework for high performance
- ✅ Isolation Forest anomaly detection model
- ✅ Pre-trained with synthetic data
- ✅ Single and batch prediction endpoints
- ✅ Risk scoring and fraud detection
- ✅ Detailed risk factor analysis
- ✅ Background task processing
- ✅ Model information endpoint

### 3. Admin Dashboard (Next.js)
- ✅ Professional dark-themed UI
- ✅ Real-time statistics dashboard
- ✅ Interactive revenue charts (Recharts)
- ✅ Transaction table with filtering
- ✅ Responsive design (Tailwind CSS)
- ✅ Loading states and error handling
- ✅ TypeScript for type safety
- ✅ Component-based architecture

### 4. Infrastructure
- ✅ Docker containers for all services
- ✅ Docker Compose orchestration
- ✅ PostgreSQL database with schema
- ✅ RabbitMQ message broker
- ✅ Health checks for all services
- ✅ Volume persistence for data
- ✅ Network isolation
- ✅ Environment variable management

## Smart Routing Logic

The system implements intelligent payment provider selection:

1. **Amount-based Routing**
   - Transactions > $1000 → PayPal
   - Transactions ≤ $1000 → Stripe

2. **Fee Optimization**
   - Stripe: 2.9% + $0.30
   - PayPal: 3.4% + $0.30

3. **Provider Availability**
   - Checks provider status
   - Fallback to alternative provider

## Fraud Detection Algorithm

The AI service uses Isolation Forest for anomaly detection:

1. **Features Analyzed**
   - Transaction amount
   - Previous transaction count
   - Account age in days

2. **Risk Scoring**
   - Score range: 0.0 - 1.0
   - Low risk: < 0.4
   - Medium risk: 0.4 - 0.7
   - High risk: > 0.7

3. **Actions**
   - APPROVE: Risk < 0.7
   - MANUAL_REVIEW: Risk 0.7 - 0.85
   - BLOCK: Risk > 0.85

4. **Risk Factors Detected**
   - High transaction amounts (> $5000)
   - New accounts (< 30 days)
   - Limited transaction history (< 5 transactions)
   - Anomalous patterns

## Database Schema

### Tables
1. **users** - User accounts
2. **transactions** - Payment transactions
3. **fraud_logs** - Fraud detection events
4. **payment_providers** - Provider configurations

### Indexes
- Transaction user_id, status, created_at
- Fraud logs transaction_id, user_id

## Security Considerations

### Implemented
- ✅ Bearer token authentication
- ✅ Admin key authentication
- ✅ CORS protection
- ✅ SQL injection prevention (parameterized queries)
- ✅ Input validation
- ✅ Connection pooling limits

### For Production
- 🔲 JWT token validation
- 🔲 Rate limiting
- 🔲 API key rotation
- 🔲 HTTPS/TLS encryption
- 🔲 Database encryption at rest
- 🔲 Secrets management (Vault)
- 🔲 WAF integration
- 🔲 DDoS protection

## Performance Optimizations

1. **Backend**
   - Goroutines for concurrent processing
   - Connection pooling (25 max, 5 idle)
   - In-memory caching with TTL
   - Efficient database queries
   - Indexed database columns

2. **AI Service**
   - Pre-trained model (no training overhead)
   - Batch prediction support
   - Background task processing
   - NumPy vectorization

3. **Frontend**
   - Server-side rendering (Next.js)
   - Component code splitting
   - Lazy loading
   - Optimized images

## Monitoring & Observability

### Available
- Health check endpoints
- Docker logs
- RabbitMQ management UI
- Database query logs

### Recommended for Production
- Prometheus metrics
- Grafana dashboards
- ELK stack logging
- APM (Application Performance Monitoring)
- Distributed tracing (Jaeger)
- Error tracking (Sentry)

## Scaling Considerations

### Horizontal Scaling
- Backend: Multiple instances behind load balancer
- AI Service: Multiple instances with load balancing
- Database: Read replicas for queries

### Vertical Scaling
- Increase container resources
- Database connection pool tuning
- Cache size optimization

## Development Workflow

1. **Local Development**
   ```bash
   # Backend
   cd backend
   go run .
   
   # AI Service
   cd ai-service
   python main.py
   
   # Dashboard
   cd admin-dashboard
   npm run dev
   ```

2. **Docker Development**
   ```bash
   docker-compose up -d --build
   ```

3. **Testing**
   - See TESTING.md for comprehensive guide
   - Health checks for all services
   - API endpoint testing
   - Load testing with Apache Bench

## Known Limitations

1. **Payment Integration**
   - Stripe/PayPal integrations are stubbed
   - Need real API keys for production
   - Webhook signature verification needed

2. **Authentication**
   - Simple bearer token (not production-ready)
   - No user registration/login flow
   - Admin key hardcoded

3. **Fraud Detection**
   - Model trained on synthetic data
   - Needs real historical data for accuracy
   - Feature set is minimal

4. **Data Persistence**
   - No backup/restore mechanism
   - No data replication

## Future Enhancements

1. **Features**
   - User management system
   - Email notifications
   - Transaction search/filtering
   - Export to CSV/PDF
   - Multi-currency support
   - Refund processing
   - Chargeback handling

2. **AI Improvements**
   - Neural network models
   - More fraud indicators
   - Real-time learning
   - Explainable AI

3. **Infrastructure**
   - Kubernetes deployment
   - CI/CD pipeline
   - Auto-scaling
   - Multi-region support

## Quick Reference

### Ports
- 3000: Admin Dashboard
- 8080: Backend API
- 8001: AI Service
- 5432: PostgreSQL
- 5672: RabbitMQ
- 15672: RabbitMQ Management

### Default Credentials
- PostgreSQL: sentinelpay / sentinelpay123
- RabbitMQ: guest / guest
- Admin Key: admin-key-secret-12345

### Important URLs
- Dashboard: http://localhost:3000
- API Docs: See API_DOCS.md
- Health: http://localhost:8080/health
- AI Health: http://localhost:8001/health
- RabbitMQ: http://localhost:15672
