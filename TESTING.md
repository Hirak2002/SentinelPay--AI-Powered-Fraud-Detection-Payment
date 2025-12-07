# SentinelPay Testing Guide

## 🧪 Testing the Complete System

### Prerequisites
Ensure all services are running:
```bash
docker-compose up -d
```

### 1. Health Checks

**Backend Service:**
```bash
curl http://localhost:8080/health
```
Expected: `{"status":"healthy"}`

**AI Service:**
```bash
curl http://localhost:8001/health
```
Expected: `{"status":"healthy","service":"sentinelpay-ai","model_trained":true}`

### 2. Create a Low-Risk Transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_001" \
  -d '{
    "user_id": "user_001",
    "amount": 50.00,
    "currency": "USD",
    "description": "Low risk test payment"
  }'
```

Expected: Transaction created with low risk score (~0.15)

### 3. Create a High-Risk Transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user_002" \
  -d '{
    "user_id": "user_002",
    "amount": 8500.00,
    "currency": "USD",
    "description": "High risk test payment"
  }'
```

Expected: Transaction created with high risk score (~0.75+) and fraud_detected=true

### 4. List Transactions

```bash
curl http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer user_001"
```

### 5. Get Dashboard Statistics

```bash
curl http://localhost:8080/api/v1/admin/stats \
  -H "X-Admin-Key: admin-key-secret-12345"
```

### 6. Get Blocked Transactions

```bash
curl http://localhost:8080/api/v1/admin/transactions/blocked \
  -H "X-Admin-Key: admin-key-secret-12345"
```

### 7. Get Revenue Metrics

```bash
curl http://localhost:8080/api/v1/admin/revenue \
  -H "X-Admin-Key: admin-key-secret-12345"
```

### 8. Test AI Service Directly

**Single Prediction:**
```bash
curl -X POST http://localhost:8001/predict \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test_001" \
  -d '{
    "user_id": "user_001",
    "amount": 5000.00,
    "currency": "USD",
    "previous_transaction_count": 3,
    "account_age_days": 10
  }'
```

**Batch Prediction:**
```bash
curl -X POST http://localhost:8001/batch-predict \
  -H "Content-Type: application/json" \
  -d '[
    {
      "user_id": "user_001",
      "amount": 100.00,
      "currency": "USD",
      "previous_transaction_count": 20,
      "account_age_days": 180
    },
    {
      "user_id": "user_002",
      "amount": 9000.00,
      "currency": "USD",
      "previous_transaction_count": 1,
      "account_age_days": 3
    }
  ]'
```

### 9. Access Admin Dashboard

Open browser: http://localhost:3000

You should see:
- Dashboard with statistics
- Transaction charts
- Blocked transactions table
- Revenue metrics
- Real-time updates

### 10. Check Service Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f ai-service
docker-compose logs -f admin-dashboard
```

### 11. Database Verification

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U sentinelpay -d sentinelpay

# Check tables
\dt

# View transactions
SELECT * FROM transactions LIMIT 10;

# View users
SELECT * FROM users;

# View fraud logs
SELECT * FROM fraud_logs;
```

### 12. RabbitMQ Management

Open browser: http://localhost:15672
- Username: guest
- Password: guest

Check queues and message rates.

## 🎯 Expected Results

### Low-Risk Transaction (Amount < $100)
- Risk Score: 0.10 - 0.20
- Fraud Detected: false
- Status: processing
- Provider: stripe (lower fees)

### Medium-Risk Transaction ($100 - $5000)
- Risk Score: 0.30 - 0.50
- Fraud Detected: false
- Status: processing
- Provider: stripe or paypal

### High-Risk Transaction (> $5000)
- Risk Score: 0.70 - 0.95
- Fraud Detected: true
- Status: blocked or manual_review
- Provider: paypal (for large amounts)

## 🔍 Performance Testing

### Load Test (using Apache Bench)
```bash
# Install Apache Bench
# Windows: Download from Apache website
# Linux: apt-get install apache2-utils

# Run 1000 requests with 10 concurrent connections
ab -n 1000 -c 10 -p payload.json -T application/json \
  -H "Authorization: Bearer user_001" \
  http://localhost:8080/api/v1/transactions
```

Create `payload.json`:
```json
{
  "user_id": "user_001",
  "amount": 150.00,
  "currency": "USD",
  "description": "Load test payment"
}
```

## ✅ Success Criteria

- [ ] All health checks return healthy status
- [ ] Low-risk transactions are processed successfully
- [ ] High-risk transactions are blocked appropriately
- [ ] Admin dashboard displays accurate statistics
- [ ] Charts render correctly with data
- [ ] Database contains transaction records
- [ ] RabbitMQ shows message activity
- [ ] No error logs in any service
- [ ] Response times < 500ms for most requests
- [ ] System handles 100+ concurrent requests

## 🐛 Common Issues

### Issue: Backend can't connect to database
**Solution:** Ensure PostgreSQL is running and credentials match

### Issue: AI service returns errors
**Solution:** Check Python dependencies are installed correctly

### Issue: Dashboard shows no data
**Solution:** Verify API_URL in dashboard .env.local

### Issue: RabbitMQ connection failed
**Solution:** RabbitMQ may take longer to start, wait 30s and retry

## 📊 Monitoring

Check service status:
```bash
docker-compose ps
```

View resource usage:
```bash
docker stats
```

Stop all services:
```bash
docker-compose down
```

Clean up and restart:
```bash
docker-compose down -v
docker-compose up -d --build
```
