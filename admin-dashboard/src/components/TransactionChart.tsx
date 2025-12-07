'use client'

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { TrendingUp } from 'lucide-react'

interface TransactionChartProps {
  revenueMetrics: any
}

export default function TransactionChart({ revenueMetrics }: TransactionChartProps) {
  const data = [
    {
      name: 'Stripe',
      revenue: revenueMetrics?.stripe_revenue || 0,
    },
    {
      name: 'PayPal',
      revenue: revenueMetrics?.paypal_revenue || 0,
    },
  ]

  return (
    <div className="card">
      <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
        <TrendingUp className="h-6 w-6 text-accent" />
        Revenue by Provider
      </h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="name" stroke="#9ca3af" />
          <YAxis stroke="#9ca3af" />
          <Tooltip
            contentStyle={{
              backgroundColor: '#1e293b',
              border: '1px solid #374151',
              borderRadius: '8px',
              color: '#fff'
            }}
          />
          <Legend />
          <Bar dataKey="revenue" fill="#10b981" radius={[8, 8, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
