import os
import json
import logging
from datetime import datetime
from typing import Optional, Dict, Any
from dotenv import load_dotenv

import numpy as np
import pandas as pd
from fastapi import FastAPI, HTTPException, BackgroundTasks, Header
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler

load_dotenv()
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="SentinelPay AI Service",
    description="AI-powered fraud detection service",
    version="1.0.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class TransactionData(BaseModel):
    user_id: str
    amount: float
    currency: str
    merchant_category: Optional[str] = None
    device_fingerprint: Optional[str] = None
    ip_country: Optional[str] = None
    previous_transaction_count: int = 0
    account_age_days: int = 0

class FraudDetectionResponse(BaseModel):
    transaction_id: str
    risk_score: float
    is_fraud: bool
    risk_factors: list
    confidence: float
    recommendation: str

class FraudDetectionModel:
    def __init__(self):
        self.model = IsolationForest(
            contamination=0.05,
            random_state=42,
            n_estimators=100
        )
        self.scaler = StandardScaler()
        self.is_trained = False
        self._train_model()

    def _train_model(self):
        """Train the model with synthetic historical data"""
        # Generate synthetic training data
        np.random.seed(42)
        n_samples = 1000
        
        amounts = np.random.exponential(scale=500, size=n_samples)
        amounts = np.clip(amounts, 1, 10000)
        
        transaction_counts = np.random.poisson(lam=20, size=n_samples)
        account_ages = np.random.exponential(scale=200, size=n_samples)
        
        # Add some anomalies
        anomaly_indices = np.random.choice(n_samples, size=50, replace=False)
        amounts[anomaly_indices] = np.random.uniform(8000, 10000, len(anomaly_indices))
        transaction_counts[anomaly_indices] = np.random.poisson(lam=2, size=len(anomaly_indices))
        
        features = np.column_stack([amounts, transaction_counts, account_ages])
        features_scaled = self.scaler.fit_transform(features)
        
        self.model.fit(features_scaled)
        self.is_trained = True
        logger.info("Fraud detection model trained successfully")

    def predict(self, data: TransactionData) -> tuple:
        """Predict fraud probability and risk score"""
        if not self.is_trained:
            return 0.3, False, []

        features = np.array([[
            data.amount,
            data.previous_transaction_count,
            data.account_age_days
        ]])
        
        features_scaled = self.scaler.transform(features)
        prediction = self.model.predict(features_scaled)[0]
        anomaly_score = self.model.score_samples(features_scaled)[0]
        
        risk_score = 1 / (1 + np.exp(-anomaly_score))
        is_fraud = prediction == -1
        
        risk_factors = []
        if data.amount > 5000:
            risk_factors.append("High transaction amount")
        if data.previous_transaction_count < 5:
            risk_factors.append("New account with limited history")
        if data.account_age_days < 30:
            risk_factors.append("Very recent account creation")
        if is_fraud:
            risk_factors.append("Anomalous transaction pattern detected")
        
        confidence = float(np.abs(anomaly_score))
        
        return float(risk_score), bool(is_fraud), risk_factors, confidence

# Initialize fraud detection model
fraud_detector = FraudDetectionModel()

@app.on_event("startup")
async def startup_event():
    logger.info("SentinelPay AI Service started")

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("SentinelPay AI Service stopped")

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "sentinelpay-ai",
        "model_trained": fraud_detector.is_trained
    }

@app.post("/predict", response_model=FraudDetectionResponse)
async def predict_fraud(
    data: TransactionData,
    x_request_id: Optional[str] = Header(None),
    background_tasks: BackgroundTasks = BackgroundTasks()
):
    """Predict fraud probability for a transaction"""
    try:
        if not x_request_id:
            x_request_id = f"txn_{datetime.now().timestamp()}"

        risk_score, is_fraud, risk_factors, confidence = fraud_detector.predict(data)
        
        if is_fraud:
            recommendation = "BLOCK"
        elif risk_score > 0.7:
            recommendation = "MANUAL_REVIEW"
        else:
            recommendation = "APPROVE"
        
        response = FraudDetectionResponse(
            transaction_id=x_request_id,
            risk_score=risk_score,
            is_fraud=is_fraud,
            risk_factors=risk_factors,
            confidence=confidence,
            recommendation=recommendation
        )
        
        # Log prediction asynchronously
        background_tasks.add_task(log_prediction, x_request_id, response)
        
        return response
    
    except Exception as e:
        logger.error(f"Error in fraud prediction: {str(e)}")
        raise HTTPException(status_code=500, detail="Fraud prediction failed")

@app.post("/batch-predict")
async def batch_predict(transactions: list[TransactionData]):
    """Batch prediction for multiple transactions"""
    try:
        results = []
        for txn in transactions:
            risk_score, is_fraud, risk_factors, confidence = fraud_detector.predict(txn)
            
            if is_fraud:
                recommendation = "BLOCK"
            elif risk_score > 0.7:
                recommendation = "MANUAL_REVIEW"
            else:
                recommendation = "APPROVE"
            
            results.append({
                "user_id": txn.user_id,
                "risk_score": float(risk_score),
                "is_fraud": bool(is_fraud),
                "risk_factors": risk_factors,
                "recommendation": recommendation
            })
        
        return {
            "status": "success",
            "count": len(results),
            "predictions": results
        }
    
    except Exception as e:
        logger.error(f"Error in batch prediction: {str(e)}")
        raise HTTPException(status_code=500, detail="Batch prediction failed")

@app.get("/model-info")
async def get_model_info():
    """Get information about the current model"""
    return {
        "model_type": "Isolation Forest",
        "model_trained": fraud_detector.is_trained,
        "contamination": 0.05,
        "n_estimators": 100,
        "features": ["amount", "transaction_count", "account_age_days"]
    }

async def log_prediction(transaction_id: str, response: FraudDetectionResponse):
    """Log prediction results for audit trail"""
    logger.info(f"Prediction logged - TxnID: {transaction_id}, RiskScore: {response.risk_score}, Recommendation: {response.recommendation}")

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", 8001))
    uvicorn.run(app, host="0.0.0.0", port=port)
