interface Transaction {
  id: string
  user_id: string
  amount: number
  currency: string
  risk_score: number
  provider: string
  created_at: string
}

interface TransactionTableProps {
  transactions: Transaction[]
}

export default function TransactionTable({ transactions }: TransactionTableProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const getRiskLevel = (score: number) => {
    if (score > 0.7) return { level: 'High', class: 'badge-danger' }
    if (score > 0.4) return { level: 'Medium', class: 'badge-warning' }
    return { level: 'Low', class: 'badge-success' }
  }

  if (transactions.length === 0) {
    return (
      <div className="text-center py-8 text-gray-400">
        No blocked transactions found
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-gray-700">
            <th className="text-left py-3 px-4 text-gray-400 font-medium">Transaction ID</th>
            <th className="text-left py-3 px-4 text-gray-400 font-medium">User ID</th>
            <th className="text-left py-3 px-4 text-gray-400 font-medium">Amount</th>
            <th className="text-left py-3 px-4 text-gray-400 font-medium">Risk Score</th>
            <th className="text-left py-3 px-4 text-gray-400 font-medium">Provider</th>
            <th className="text-left py-3 px-4 text-gray-400 font-medium">Date</th>
          </tr>
        </thead>
        <tbody>
          {transactions.map((txn) => {
            const riskInfo = getRiskLevel(txn.risk_score)
            return (
              <tr key={txn.id} className="border-b border-gray-800 hover:bg-primary transition-colors">
                <td className="py-3 px-4 font-mono text-sm">{txn.id}</td>
                <td className="py-3 px-4 font-mono text-sm">{txn.user_id}</td>
                <td className="py-3 px-4 font-semibold">
                  ${txn.amount.toLocaleString()} {txn.currency}
                </td>
                <td className="py-3 px-4">
                  <span className={riskInfo.class}>
                    {riskInfo.level} ({(txn.risk_score * 100).toFixed(0)}%)
                  </span>
                </td>
                <td className="py-3 px-4">
                  <span className="badge-info capitalize">{txn.provider}</span>
                </td>
                <td className="py-3 px-4 text-sm text-gray-400">
                  {formatDate(txn.created_at)}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
