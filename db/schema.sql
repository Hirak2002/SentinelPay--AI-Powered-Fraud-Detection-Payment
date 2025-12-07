-- SentinelPay Database Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    transaction_count INT DEFAULT 0,
    total_spent DECIMAL(12, 2) DEFAULT 0.00,
    is_blacklisted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    risk_score DECIMAL(5, 4) DEFAULT 0.0000,
    fraud_detected BOOLEAN DEFAULT FALSE,
    provider VARCHAR(50),
    provider_txn_id VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Fraud logs table
CREATE TABLE IF NOT EXISTS fraud_logs (
    id VARCHAR(255) PRIMARY KEY,
    transaction_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    risk_score DECIMAL(5, 4) NOT NULL,
    action VARCHAR(50) NOT NULL,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Payment providers table
CREATE TABLE IF NOT EXISTS payment_providers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    fee DECIMAL(5, 4) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_fraud_logs_transaction_id ON fraud_logs(transaction_id);
CREATE INDEX idx_fraud_logs_user_id ON fraud_logs(user_id);

-- Insert default payment providers
INSERT INTO payment_providers (id, name, fee, status) VALUES
    ('stripe_001', 'Stripe', 0.0290, 'active'),
    ('paypal_001', 'PayPal', 0.0340, 'active')
ON CONFLICT DO NOTHING;

-- Insert sample users
INSERT INTO users (id, email, name, transaction_count, total_spent) VALUES
    ('user_001', 'john.doe@example.com', 'John Doe', 25, 5420.50),
    ('user_002', 'jane.smith@example.com', 'Jane Smith', 18, 3280.75),
    ('user_003', 'bob.wilson@example.com', 'Bob Wilson', 42, 12450.00)
ON CONFLICT DO NOTHING;
