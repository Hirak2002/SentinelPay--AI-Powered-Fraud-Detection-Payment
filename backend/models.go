package main

import "time"

// Transaction represents a payment transaction
type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"` // pending, processing, completed, failed, blocked
	RiskScore       float64   `json:"risk_score"`
	FraudDetected   bool      `json:"fraud_detected"`
	Provider        string    `json:"provider"` // stripe or paypal
	ProviderTxnID   string    `json:"provider_txn_id"`
	Description     string    `json:"description"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	TransactionCount int      `json:"transaction_count"`
	TotalSpent      float64   `json:"total_spent"`
	IsBlacklisted   bool      `json:"is_blacklisted"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FraudLog represents fraud detection events
type FraudLog struct {
	ID              string    `json:"id"`
	TransactionID   string    `json:"transaction_id"`
	UserID          string    `json:"user_id"`
	RiskFactors     []string  `json:"risk_factors"`
	RiskScore       float64   `json:"risk_score"`
	Action          string    `json:"action"` // approved, blocked, manual_review
	DetectedAt      time.Time `json:"detected_at"`
}

// PaymentProvider represents payment provider configuration
type PaymentProvider struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Fee         float64 `json:"fee"`
	Status      string  `json:"status"`
	LastChecked time.Time `json:"last_checked"`
}

// APIResponse standard response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaymentRequest incoming payment request
type PaymentRequest struct {
	UserID      string                 `json:"user_id" binding:"required"`
	Amount      float64                `json:"amount" binding:"required,gt=0"`
	Currency    string                 `json:"currency" binding:"required"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalTransactions    int64     `json:"total_transactions"`
	BlockedTransactions  int64     `json:"blocked_transactions"`
	TotalRevenue         float64   `json:"total_revenue"`
	AverageRiskScore     float64   `json:"average_risk_score"`
	FraudPreventionRate  float64   `json:"fraud_prevention_rate"`
	TransactionSuccessRate float64  `json:"transaction_success_rate"`
}

// RevenueMetrics represents revenue breakdown
type RevenueMetrics struct {
	StripeRevenue     float64 `json:"stripe_revenue"`
	PayPalRevenue     float64 `json:"paypal_revenue"`
	TotalRevenue      float64 `json:"total_revenue"`
	TransactionCount  int64   `json:"transaction_count"`
	AveragePerTxn     float64 `json:"average_per_txn"`
}
