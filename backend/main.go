package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Application struct {
	server   *http.Server
	router   *gin.Engine
	db       *DatabaseConnection
	broker   *MessageBroker
	cache    *CacheService
}

func init() {
	godotenv.Load()
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	app := &Application{}
	
	// Initialize database
	var err error
	app.db, err = NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer app.db.Close()

	// Initialize message broker
	app.broker, err = NewMessageBroker()
	if err != nil {
		log.Fatalf("Failed to connect to message broker: %v", err)
	}
	defer app.broker.Close()

	// Initialize cache service
	app.cache = NewCacheService()

	// Setup router
	app.router = gin.New()
	app.setupRoutes()

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.server = &http.Server{
		Addr:         ":" + port,
		Handler:      app.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting SentinelPay Backend on port %s", port)
	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func (app *Application) setupRoutes() {
	// Middleware
	app.router.Use(gin.Logger())
	app.router.Use(CORSMiddleware())
	app.router.Use(AuthenticationMiddleware(app.db))

	// Health check
	app.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Transaction routes
	transactionHandler := NewTransactionHandler(app.db, app.broker, app.cache)
	transactionRoutes := app.router.Group("/api/v1/transactions")
	{
		transactionRoutes.POST("", transactionHandler.CreateTransaction)
		transactionRoutes.GET("/:id", transactionHandler.GetTransaction)
		transactionRoutes.GET("", transactionHandler.ListTransactions)
		transactionRoutes.GET("/:id/status", transactionHandler.GetTransactionStatus)
	}

	// Provider routes
	providerHandler := NewPaymentProviderHandler(app.db)
	providerRoutes := app.router.Group("/api/v1/providers")
	{
		providerRoutes.GET("", providerHandler.ListProviders)
		providerRoutes.GET("/rates", providerHandler.GetProviderRates)
	}

	// Admin routes
	adminHandler := NewAdminHandler(app.db)
	adminRoutes := app.router.Group("/api/v1/admin")
	adminRoutes.Use(AdminMiddleware())
	{
		adminRoutes.GET("/stats", adminHandler.GetDashboardStats)
		adminRoutes.GET("/transactions/blocked", adminHandler.GetBlockedTransactions)
		adminRoutes.GET("/revenue", adminHandler.GetRevenueMetrics)
		adminRoutes.GET("/fraud-logs", adminHandler.GetFraudLogs)
	}

	// Webhook routes
	webhookHandler := NewWebhookHandler(app.db, app.broker)
	webhookRoutes := app.router.Group("/api/v1/webhooks")
	{
		webhookRoutes.POST("/stripe", webhookHandler.HandleStripeWebhook)
		webhookRoutes.POST("/paypal", webhookHandler.HandlePayPalWebhook)
	}
}
