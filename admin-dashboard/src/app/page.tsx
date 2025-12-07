'use client'

import { useState, useEffect } from 'react'
import { DollarSign, TrendingUp, Shield, Activity, AlertTriangle, CheckCircle } from 'lucide-react'
import StatsCard from '@/components/StatsCard'
import TransactionChart from '@/components/TransactionChart'
import TransactionTable from '@/components/TransactionTable'
import { fetchDashboardStats, fetchBlockedTransactions, fetchRevenueMetrics } from '@/lib/api'

export default function Dashboard() {
  const [stats, setStats] = useState<any>(null)
  const [blockedTransactions, setBlockedTransactions] = useState<any[]>([])
  const [revenueMetrics, setRevenueMetrics] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboardData()
  }, [])

  const loadDashboardData = async () => {
    try {
      const [statsData, blockedData, revenueData] = await Promise.all([
        fetchDashboardStats(),
        fetchBlockedTransactions(),
        fetchRevenueMetrics()
      ])
      
      setStats(statsData)
      setBlockedTransactions(blockedData)
      setRevenueMetrics(revenueData)
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-accent"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        <header className="mb-12">
          <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-accent to-green-400 bg-clip-text text-transparent">
            SentinelPay Dashboard
          </h1>
          <p className="text-gray-400">AI-Powered Fraud Detection & Payment Orchestration</p>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatsCard
            title="Total Transactions"
            value={stats?.total_transactions?.toLocaleString() || '0'}
            icon={<Activity className="h-8 w-8 text-accent" />}
            trend="+12.5%"
          />
          <StatsCard
            title="Blocked Transactions"
            value={stats?.blocked_transactions?.toLocaleString() || '0'}
            icon={<AlertTriangle className="h-8 w-8 text-danger" />}
            trend={`${stats?.fraud_prevention_rate?.toFixed(1)}%`}
          />
          <StatsCard
            title="Total Revenue"
            value={`$${revenueMetrics?.total_revenue?.toLocaleString() || '0'}`}
            icon={<DollarSign className="h-8 w-8 text-accent" />}
            trend="+8.2%"
          />
          <StatsCard
            title="Success Rate"
            value={`${stats?.transaction_success_rate?.toFixed(1)}%`}
            icon={<CheckCircle className="h-8 w-8 text-accent" />}
            trend="+2.1%"
          />
        </div>

        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <TransactionChart revenueMetrics={revenueMetrics} />
          <div className="card">
            <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <Shield className="h-6 w-6 text-accent" />
              Fraud Prevention Insights
            </h3>
            <div className="space-y-4">
              <div className="flex justify-between items-center p-4 bg-primary rounded-lg">
                <span className="text-gray-400">Average Risk Score</span>
                <span className="text-2xl font-bold text-warning">
                  {(stats?.average_risk_score * 100)?.toFixed(1)}%
                </span>
              </div>
              <div className="flex justify-between items-center p-4 bg-primary rounded-lg">
                <span className="text-gray-400">Fraud Prevention Rate</span>
                <span className="text-2xl font-bold text-accent">
                  {stats?.fraud_prevention_rate?.toFixed(1)}%
                </span>
              </div>
              <div className="flex justify-between items-center p-4 bg-primary rounded-lg">
                <span className="text-gray-400">Average Transaction</span>
                <span className="text-2xl font-bold text-blue-400">
                  ${revenueMetrics?.average_per_txn?.toFixed(2) || '0'}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Blocked Transactions Table */}
        <div className="card">
          <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <AlertTriangle className="h-6 w-6 text-danger" />
            Recent Blocked Transactions
          </h3>
          <TransactionTable transactions={blockedTransactions} />
        </div>
      </div>
    </div>
  )
}
