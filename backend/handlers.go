package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	db     *DatabaseConnection
	broker *MessageBroker
	cache  *CacheService
}

func NewTransactionHandler(db *DatabaseConnection, broker *MessageBroker, cache *CacheService) *TransactionHandler {
	return &TransactionHandler{db: db, broker: broker, cache: cache}
}

func (th *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	txnID := uuid.New().String()
	now := time.Now()

	// Call AI service for fraud detection
	riskScore, fraudDetected := CallFraudDetectionService(req)

	// Determine provider based on routing logic
	provider := SelectPaymentProvider(req.Amount)

	// Prepare transaction data
	txn := Transaction{
		ID:            txnID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "processing",
		RiskScore:     riskScore,
		FraudDetected: fraudDetected,
		Provider:      provider,
		Description:   req.Description,
		Metadata:      req.Metadata,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save to database
	query := `
		INSERT INTO transactions (id, user_id, amount, currency, status, risk_score, fraud_detected, provider, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	if _, err := th.db.ExecuteQuery(query, txn.ID, txn.UserID, txn.Amount, txn.Currency, txn.Status, txn.RiskScore, txn.FraudDetected, txn.Provider, txn.Description, txn.CreatedAt, txn.UpdatedAt); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to create transaction",
		})
		return
	}

	// Publish to queue for async processing
	th.broker.PublishEvent("payment.created", txn)

	// Cache the transaction
	th.cache.Set(txnID, txn, 5*time.Minute)

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    txn,
		Message: "Transaction created successfully",
	})
}

func (th *TransactionHandler) GetTransaction(c *gin.Context) {
	txnID := c.Param("id")

	// Try cache first
	if cached, ok := th.cache.Get(txnID); ok {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Data:    cached,
		})
		return
	}

	query := `SELECT id, user_id, amount, currency, status, risk_score, fraud_detected, provider, provider_txn_id, description, created_at, updated_at FROM transactions WHERE id = $1`
	row := th.db.QueryRow(query, txnID)

	var txn Transaction
	err := row.Scan(&txn.ID, &txn.UserID, &txn.Amount, &txn.Currency, &txn.Status, &txn.RiskScore, &txn.FraudDetected, &txn.Provider, &txn.ProviderTxnID, &txn.Description, &txn.CreatedAt, &txn.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Transaction not found",
		})
		return
	}

	th.cache.Set(txnID, txn, 5*time.Minute)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    txn,
	})
}

func (th *TransactionHandler) ListTransactions(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	query := `SELECT id, user_id, amount, currency, status, risk_score, fraud_detected, provider, provider_txn_id, description, created_at, updated_at FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := th.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch transactions",
		})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.Amount, &txn.Currency, &txn.Status, &txn.RiskScore, &txn.FraudDetected, &txn.Provider, &txn.ProviderTxnID, &txn.Description, &txn.CreatedAt, &txn.UpdatedAt); err != nil {
			continue
		}
		transactions = append(transactions, txn)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    transactions,
	})
}

func (th *TransactionHandler) GetTransactionStatus(c *gin.Context) {
	txnID := c.Param("id")

	query := `SELECT status FROM transactions WHERE id = $1`
	row := th.db.QueryRow(query, txnID)

	var status string
	if err := row.Scan(&status); err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{"status": status},
	})
}

type PaymentProviderHandler struct {
	db *DatabaseConnection
}

func NewPaymentProviderHandler(db *DatabaseConnection) *PaymentProviderHandler {
	return &PaymentProviderHandler{db: db}
}

func (ph *PaymentProviderHandler) ListProviders(c *gin.Context) {
	query := `SELECT id, name, fee, status, last_checked FROM payment_providers ORDER BY name`
	rows, err := ph.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch providers",
		})
		return
	}
	defer rows.Close()

	var providers []PaymentProvider
	for rows.Next() {
		var p PaymentProvider
		if err := rows.Scan(&p.ID, &p.Name, &p.Fee, &p.Status, &p.LastChecked); err != nil {
			continue
		}
		providers = append(providers, p)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    providers,
	})
}

func (ph *PaymentProviderHandler) GetProviderRates(c *gin.Context) {
	rates := gin.H{
		"stripe": gin.H{"fee": 0.029, "fixed": 0.30, "status": "active"},
		"paypal": gin.H{"fee": 0.034, "fixed": 0.30, "status": "active"},
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    rates,
	})
}

type AdminHandler struct {
	db *DatabaseConnection
}

func NewAdminHandler(db *DatabaseConnection) *AdminHandler {
	return &AdminHandler{db: db}
}

func (ah *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats := &DashboardStats{}

	// Get total transactions
	query := `SELECT COUNT(*) FROM transactions`
	ah.db.QueryRow(query).Scan(&stats.TotalTransactions)

	// Get blocked transactions
	query = `SELECT COUNT(*) FROM transactions WHERE status = 'blocked'`
	ah.db.QueryRow(query).Scan(&stats.BlockedTransactions)

	// Get total revenue
	query = `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE status = 'completed'`
	ah.db.QueryRow(query).Scan(&stats.TotalRevenue)

	// Get average risk score
	query = `SELECT COALESCE(AVG(risk_score), 0) FROM transactions`
	ah.db.QueryRow(query).Scan(&stats.AverageRiskScore)

	if stats.TotalTransactions > 0 {
		stats.FraudPreventionRate = (float64(stats.BlockedTransactions) / float64(stats.TotalTransactions)) * 100
		successCount := stats.TotalTransactions - stats.BlockedTransactions
		stats.TransactionSuccessRate = (float64(successCount) / float64(stats.TotalTransactions)) * 100
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

func (ah *AdminHandler) GetBlockedTransactions(c *gin.Context) {
	query := `SELECT id, user_id, amount, currency, status, risk_score, fraud_detected, provider, created_at FROM transactions WHERE status = 'blocked' ORDER BY created_at DESC LIMIT 50`
	rows, err := ah.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch blocked transactions",
		})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.Amount, &txn.Currency, &txn.Status, &txn.RiskScore, &txn.FraudDetected, &txn.Provider, &txn.CreatedAt); err != nil {
			continue
		}
		transactions = append(transactions, txn)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    transactions,
	})
}

func (ah *AdminHandler) GetRevenueMetrics(c *gin.Context) {
	metrics := &RevenueMetrics{}

	query := `SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM transactions WHERE provider = 'stripe' AND status = 'completed'`
	ah.db.QueryRow(query).Scan(&metrics.StripeRevenue, &metrics.TransactionCount)

	query = `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE provider = 'paypal' AND status = 'completed'`
	ah.db.QueryRow(query).Scan(&metrics.PayPalRevenue)

	metrics.TotalRevenue = metrics.StripeRevenue + metrics.PayPalRevenue
	if metrics.TransactionCount > 0 {
		metrics.AveragePerTxn = metrics.TotalRevenue / float64(metrics.TransactionCount)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    metrics,
	})
}

func (ah *AdminHandler) GetFraudLogs(c *gin.Context) {
	query := `SELECT id, transaction_id, user_id, risk_score, action, detected_at FROM fraud_logs ORDER BY detected_at DESC LIMIT 100`
	rows, err := ah.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch fraud logs",
		})
		return
	}
	defer rows.Close()

	var logs []FraudLog
	for rows.Next() {
		var log FraudLog
		if err := rows.Scan(&log.ID, &log.TransactionID, &log.UserID, &log.RiskScore, &log.Action, &log.DetectedAt); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    logs,
	})
}

type WebhookHandler struct {
	db     *DatabaseConnection
	broker *MessageBroker
}

func NewWebhookHandler(db *DatabaseConnection, broker *MessageBroker) *WebhookHandler {
	return &WebhookHandler{db: db, broker: broker}
}

func (wh *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	txnID := fmt.Sprintf("%v", payload["metadata"].(map[string]interface{})["txn_id"])
	status := fmt.Sprintf("%v", payload["status"])

	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	wh.db.ExecuteQuery(query, status, time.Now(), txnID)

	c.JSON(http.StatusOK, APIResponse{Success: true})
}

func (wh *WebhookHandler) HandlePayPalWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	txnID := fmt.Sprintf("%v", payload["txn_id"])
	status := fmt.Sprintf("%v", payload["status"])

	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	wh.db.ExecuteQuery(query, status, time.Now(), txnID)

	c.JSON(http.StatusOK, APIResponse{Success: true})
}
