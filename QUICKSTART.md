# 🚀 SentinelPay - Quick Start Guide

Welcome to **SentinelPay**, your AI-powered fraud detection and payment orchestration system!

## ⚡ Getting Started in 3 Steps

### Step 1: Prerequisites
Make sure you have **Docker Desktop** installed and running on your Windows machine.
- Download Docker Desktop: https://www.docker.com/products/docker-desktop

### Step 2: Start the Application
Open PowerShell in the project directory and run:

```powershell
.\start.ps1
```

Or manually with Docker Compose:

```powershell
docker-compose up -d
```

### Step 3: Access the Dashboard
Open your browser and navigate to:
- **Admin Dashboard**: http://localhost:3000

That's it! 🎉

---

## 📱 What You'll See

### Admin Dashboard (http://localhost:3000)

The professional dashboard displays:

1. **Statistics Cards**
   - Total Transactions
   - Blocked Transactions  
   - Total Revenue
   - Success Rate

2. **Revenue Chart**
   - Visual breakdown by provider (Stripe vs PayPal)

3. **Fraud Prevention Insights**
   - Average Risk Score
   - Fraud Prevention Rate
   - Average Transaction Value

4. **Blocked Transactions Table**
   - Transaction details
   - Risk levels
   - Timestamps

---

## 🧪 Testing the System

### Test a Low-Risk Transaction

Open PowerShell and run:

```powershell
curl -X POST http://localhost:8080/api/v1/transactions `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer user_001" `
  -d '{
    \"user_id\": \"user_001\",
    \"amount\": 50.00,
    \"currency\": \"USD\",
    \"description\": \"Test payment\"
  }'
```

**Expected Result**: Transaction approved with low risk score (~15%)

### Test a High-Risk Transaction

```powershell
curl -X POST http://localhost:8080/api/v1/transactions `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer user_002" `
  -d '{
    \"user_id\": \"user_002\",
    \"amount\": 8500.00,
    \"currency\": \"USD\",
    \"description\": \"Large payment\"
  }'
```

**Expected Result**: Transaction flagged/blocked with high risk score (~75%+)

### View Dashboard Statistics

```powershell
curl http://localhost:8080/api/v1/admin/stats `
  -H "X-Admin-Key: admin-key-secret-12345"
```

---

## 🎯 How It Works

### Transaction Flow

```
1. Client sends payment request
   ↓
2. Backend receives and validates
   ↓
3. AI Service analyzes fraud risk
   ↓
4. Smart routing selects provider
   ↓
5. Transaction saved to database
   ↓
6. Event published to message queue
   ↓
7. Response sent to client
   ↓
8. Dashboard updates in real-time
```

### Fraud Detection Logic

- **Low Risk (< 40%)**: Amount < $100, established account
- **Medium Risk (40-70%)**: Amount $100-$5000, moderate history
- **High Risk (> 70%)**: Amount > $5000, new account, suspicious patterns

### Smart Routing

- **Small Transactions (< $1000)**: Routes to Stripe (lower fees)
- **Large Transactions (> $1000)**: Routes to PayPal (better for large amounts)

---

## 🔧 Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| **Admin Dashboard** | http://localhost:3000 | Web UI for monitoring |
| **Backend API** | http://localhost:8080 | Core payment service |
| **AI Service** | http://localhost:8001 | Fraud detection |
| **RabbitMQ UI** | http://localhost:15672 | Message queue (guest/guest) |
| **PostgreSQL** | localhost:5432 | Database |

---

## 📊 Understanding the Dashboard

### Statistics Overview
- **Total Transactions**: All transactions processed
- **Blocked Transactions**: High-risk transactions blocked
- **Total Revenue**: Sum of successful payments
- **Success Rate**: Percentage of approved transactions

### Revenue Chart
- Visual comparison of Stripe vs PayPal revenue
- Helps identify preferred payment provider

### Fraud Prevention Metrics
- **Average Risk Score**: Overall risk level (lower is better)
- **Fraud Prevention Rate**: % of fraudulent transactions caught
- **Average Transaction**: Mean transaction value

### Blocked Transactions Table
Shows recently blocked high-risk transactions with:
- Transaction ID
- User ID
- Amount
- Risk score and level
- Provider
- Timestamp

---

## 🛑 Stopping the Application

```powershell
docker-compose down
```

To stop and remove all data:

```powershell
docker-compose down -v
```

---

## 🔍 Troubleshooting

### Dashboard shows no data?
- Wait 10-15 seconds for services to fully start
- Check if backend is running: `curl http://localhost:8080/health`
- Refresh the browser page

### Can't create transactions?
- Ensure all services are running: `docker-compose ps`
- Check logs: `docker-compose logs backend`

### Port already in use?
- Stop conflicting services
- Or modify ports in `docker-compose.yml`

### Docker errors?
- Restart Docker Desktop
- Run: `docker-compose down` then `docker-compose up -d --build`

---

## 📚 Additional Resources

- **Full Documentation**: See `README.md`
- **API Reference**: See `API_DOCS.md`
- **Testing Guide**: See `TESTING.md`
- **Development Guide**: See `DEVELOPMENT.md`

---

## 🎨 UI Features

### Modern Design
- Dark theme for reduced eye strain
- Gradient accents (green for success)
- Responsive layout for all screen sizes

### Interactive Elements
- Hover effects on cards
- Smooth animations
- Color-coded risk levels:
  - 🟢 Green: Low risk
  - 🟡 Yellow: Medium risk
  - 🔴 Red: High risk

### Real-time Updates
- Dashboard refreshes automatically
- Charts update with new data
- Transaction table shows latest entries

---

## 💡 Pro Tips

1. **Monitor RabbitMQ**: Check message flow at http://localhost:15672
2. **Check Logs**: Use `docker-compose logs -f` to watch all services
3. **Database Access**: Connect to PostgreSQL to see raw data
4. **Performance**: System handles 100+ concurrent requests
5. **Testing**: Use the TESTING.md guide for comprehensive tests

---

## ✅ Success Checklist

- [ ] Docker Desktop is running
- [ ] All services started successfully (`docker-compose ps`)
- [ ] Dashboard loads at http://localhost:3000
- [ ] Can create low-risk transaction
- [ ] Can create high-risk transaction
- [ ] Statistics display correctly
- [ ] Charts render with data
- [ ] Blocked transactions appear in table

---

## 🎓 What You've Built

✅ **Microservices Architecture**
- Go backend (high performance)
- Python AI service (machine learning)
- Next.js frontend (modern UI)

✅ **Production-Ready Features**
- Fraud detection with ML
- Smart payment routing
- Real-time analytics
- Message queue integration
- Containerized deployment

✅ **Professional UI**
- Beautiful dark theme
- Interactive charts
- Real-time updates
- Responsive design

---

## 🚀 Next Steps

1. **Explore the API**: Check `API_DOCS.md` for all endpoints
2. **Run Tests**: Follow `TESTING.md` for comprehensive testing
3. **Customize**: Modify the code to add your own features
4. **Learn More**: Review `DEVELOPMENT.md` for technical details

---

**Congratulations!** 🎉 You now have a fully functional AI-powered payment orchestration system!

For questions or issues, review the documentation files or check the service logs.

Happy coding! 💻✨
