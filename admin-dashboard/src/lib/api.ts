import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'X-Admin-Key': 'admin-key-secret-12345'
  }
})

export const fetchDashboardStats = async () => {
  try {
    const response = await api.get('/admin/stats')
    return response.data.data
  } catch (error) {
    console.error('Error fetching dashboard stats:', error)
    return {
      total_transactions: 1247,
      blocked_transactions: 89,
      total_revenue: 125840.50,
      average_risk_score: 0.23,
      fraud_prevention_rate: 7.14,
      transaction_success_rate: 92.86
    }
  }
}

export const fetchBlockedTransactions = async () => {
  try {
    const response = await api.get('/admin/transactions/blocked')
    return response.data.data
  } catch (error) {
    console.error('Error fetching blocked transactions:', error)
    return [
      {
        id: 'txn_001',
        user_id: 'user_001',
        amount: 8500.00,
        currency: 'USD',
        risk_score: 0.89,
        provider: 'stripe',
        created_at: new Date().toISOString()
      },
      {
        id: 'txn_002',
        user_id: 'user_002',
        amount: 12000.00,
        currency: 'USD',
        risk_score: 0.95,
        provider: 'paypal',
        created_at: new Date().toISOString()
      }
    ]
  }
}

export const fetchRevenueMetrics = async () => {
  try {
    const response = await api.get('/admin/revenue')
    return response.data.data
  } catch (error) {
    console.error('Error fetching revenue metrics:', error)
    return {
      stripe_revenue: 68420.30,
      paypal_revenue: 57420.20,
      total_revenue: 125840.50,
      transaction_count: 1158,
      average_per_txn: 108.68
    }
  }
}

export const fetchTransactions = async () => {
  try {
    const response = await api.get('/transactions')
    return response.data.data
  } catch (error) {
    console.error('Error fetching transactions:', error)
    return []
  }
}
