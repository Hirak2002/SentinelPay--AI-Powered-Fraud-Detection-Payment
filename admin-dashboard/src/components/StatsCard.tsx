import React from 'react'

interface StatsCardProps {
  title: string
  value: string | number
  icon: React.ReactNode
  trend?: string
}

export default function StatsCard({ title, value, icon, trend }: StatsCardProps) {
  return (
    <div className="stat-card">
      <div className="flex items-start justify-between mb-4">
        <div className="p-3 bg-primary rounded-lg">
          {icon}
        </div>
        {trend && (
          <div className={`text-sm font-medium ${trend.startsWith('+') ? 'text-accent' : 'text-danger'}`}>
            {trend}
          </div>
        )}
      </div>
      <h3 className="text-gray-400 text-sm font-medium mb-1">{title}</h3>
      <p className="text-3xl font-bold">{value}</p>
    </div>
  )
}
