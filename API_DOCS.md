# API Documentation

## Base URLs
- Backend API: `http://localhost:8080/api/v1`
- AI Service: `http://localhost:8001`

## Authentication

### User Authentication
Include in request headers:
```
Authorization: Bearer {user_id}
```

### Admin Authentication
Include in request headers:
```
X-Admin-Key: admin-key-secret-12345
```

---

## Transaction Endpoints

### Create Transaction
**POST** `/transactions`

Creates a new payment transaction with fraud detection.

**Headers:**
- `Content-Type: application/json`
- `Authorization: Bearer {user_id}`

**Request Body:**
```json
{
  "user_id": "user_001",
  "amount": 150.00,
  "currency": "USD",
  "description": "Product purchase",
  "metadata": {
    "order_id": "ORDER-123",
    "product_id": "PROD-456"
  }
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "txn_abc123",
    "user_id": "user_001",
    "amount": 150.00,
    "currency": "USD",
    "status": "processing",
    "risk_score": 0.23,
    "fraud_detected": false,
    "provider": "stripe",
    "provider_txn_id": "",
    "description": "Product purchase",
    "metadata": {
      "order_id": "ORDER-123",
      "product_id": "PROD-456"
    },
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "message": "Transaction created successfully"
}
```

---

### Get Transaction
**GET** `/transactions/:id`

Retrieves details of a specific transaction.

**Headers:**
- `Authorization: Bearer {user_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "txn_abc123",
    "user_id": "user_001",
    "amount": 150.00,
    "currency": "USD",
    "status": "completed",
    "risk_score": 0.23,
    "fraud_detected": false,
    "provider": "stripe",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### List Transactions
**GET** `/transactions?limit=20&offset=0`

Lists transactions with pagination.

**Query Parameters:**
- `limit` (optional): Number of results (1-100, default: 20)
- `offset` (optional): Pagination offset (default: 0)

**Headers:**
- `Authorization: Bearer {user_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "txn_abc123",
      "user_id": "user_001",
      "amount": 150.00,
      "currency": "USD",
      "status": "completed",
      "risk_score": 0.23,
      "fraud_detected": false,
      "provider": "stripe",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### Get Transaction Status
**GET** `/transactions/:id/status`

Gets current status of a transaction.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "status": "completed"
  }
}
```

---

## Admin Endpoints

### Get Dashboard Statistics
**GET** `/admin/stats`

Retrieves overall dashboard statistics.

**Headers:**
- `X-Admin-Key: admin-key-secret-12345`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total_transactions": 1247,
    "blocked_transactions": 89,
    "total_revenue": 125840.50,
    "average_risk_score": 0.23,
    "fraud_prevention_rate": 7.14,
    "transaction_success_rate": 92.86
  }
}
```

---

### Get Blocked Transactions
**GET** `/admin/transactions/blocked`

Lists all blocked transactions.

**Headers:**
- `X-Admin-Key: admin-key-secret-12345`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "txn_xyz789",
      "user_id": "user_002",
      "amount": 8500.00,
      "currency": "USD",
      "status": "blocked",
      "risk_score": 0.89,
      "fraud_detected": true,
      "provider": "stripe",
      "created_at": "2024-01-01T11:30:00Z"
    }
  ]
}
```

---

### Get Revenue Metrics
**GET** `/admin/revenue`

Retrieves revenue breakdown by provider.

**Headers:**
- `X-Admin-Key: admin-key-secret-12345`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "stripe_revenue": 68420.30,
    "paypal_revenue": 57420.20,
    "total_revenue": 125840.50,
    "transaction_count": 1158,
    "average_per_txn": 108.68
  }
}
```

---

### Get Fraud Logs
**GET** `/admin/fraud-logs`

Retrieves fraud detection logs.

**Headers:**
- `X-Admin-Key: admin-key-secret-12345`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "fraud_log_001",
      "transaction_id": "txn_xyz789",
      "user_id": "user_002",
      "risk_score": 0.89,
      "action": "blocked",
      "detected_at": "2024-01-01T11:30:00Z"
    }
  ]
}
```

---

## AI Service Endpoints

### Predict Fraud
**POST** `/predict`

Predicts fraud probability for a transaction.

**Headers:**
- `Content-Type: application/json`
- `X-Request-ID: {transaction_id}` (optional)

**Request Body:**
```json
{
  "user_id": "user_001",
  "amount": 5000.00,
  "currency": "USD",
  "merchant_category": "electronics",
  "device_fingerprint": "abc123",
  "ip_country": "US",
  "previous_transaction_count": 15,
  "account_age_days": 180
}
```

**Response:** `200 OK`
```json
{
  "transaction_id": "txn_abc123",
  "risk_score": 0.45,
  "is_fraud": false,
  "risk_factors": [
    "High transaction amount"
  ],
  "confidence": 0.82,
  "recommendation": "APPROVE"
}
```

**Recommendations:**
- `APPROVE`: Low risk, proceed with transaction
- `MANUAL_REVIEW`: Medium risk, requires manual review
- `BLOCK`: High risk, block transaction

---

### Batch Predict
**POST** `/batch-predict`

Batch fraud prediction for multiple transactions.

**Request Body:**
```json
[
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
]
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "count": 2,
  "predictions": [
    {
      "user_id": "user_001",
      "risk_score": 0.15,
      "is_fraud": false,
      "risk_factors": [],
      "recommendation": "APPROVE"
    },
    {
      "user_id": "user_002",
      "risk_score": 0.92,
      "is_fraud": true,
      "risk_factors": [
        "High transaction amount",
        "New account with limited history",
        "Very recent account creation",
        "Anomalous transaction pattern detected"
      ],
      "recommendation": "BLOCK"
    }
  ]
}
```

---

### Get Model Info
**GET** `/model-info`

Returns information about the fraud detection model.

**Response:** `200 OK`
```json
{
  "model_type": "Isolation Forest",
  "model_trained": true,
  "contamination": 0.05,
  "n_estimators": 100,
  "features": [
    "amount",
    "transaction_count",
    "account_age_days"
  ]
}
```

---

## Payment Provider Endpoints

### List Providers
**GET** `/providers`

Lists all available payment providers.

**Headers:**
- `Authorization: Bearer {user_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "stripe_001",
      "name": "Stripe",
      "fee": 0.0290,
      "status": "active",
      "last_checked": "2024-01-01T12:00:00Z"
    },
    {
      "id": "paypal_001",
      "name": "PayPal",
      "fee": 0.0340,
      "status": "active",
      "last_checked": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### Get Provider Rates
**GET** `/providers/rates`

Gets current provider rates and status.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "stripe": {
      "fee": 0.029,
      "fixed": 0.30,
      "status": "active"
    },
    "paypal": {
      "fee": 0.034,
      "fixed": 0.30,
      "status": "active"
    }
  }
}
```

---

## Webhook Endpoints

### Stripe Webhook
**POST** `/webhooks/stripe`

Handles Stripe webhook events.

**Request Body:** Stripe webhook payload

---

### PayPal Webhook
**POST** `/webhooks/paypal`

Handles PayPal webhook events.

**Request Body:** PayPal webhook payload

---

## Health Check

### Backend Health
**GET** `/health`

**Response:** `200 OK`
```json
{
  "status": "healthy"
}
```

### AI Service Health
**GET** `/health` (AI Service)

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "service": "sentinelpay-ai",
  "model_trained": true
}
```

---

## Error Responses

All endpoints may return error responses:

**400 Bad Request**
```json
{
  "success": false,
  "error": "Invalid request format"
}
```

**401 Unauthorized**
```json
{
  "success": false,
  "error": "Missing authorization header"
}
```

**403 Forbidden**
```json
{
  "success": false,
  "error": "Unauthorized access"
}
```

**404 Not Found**
```json
{
  "success": false,
  "error": "Transaction not found"
}
```

**500 Internal Server Error**
```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## Rate Limiting

Currently no rate limiting is implemented. In production, consider:
- 100 requests/minute for authenticated users
- 1000 requests/minute for admin endpoints
- 10 requests/minute for unauthenticated endpoints
