import { useQuery } from '@tanstack/react-query'
import { BookmarkCheck, AlertTriangle, Search } from 'lucide-react'
import { useState } from 'react'
import { api } from '../services/api'

const WATCHLIST_TYPES = ['TERRORISM', 'ORGANIZED_CRIME', 'NARCOTICS', 'FINANCIAL', 'WANTED', 'IMMIGRATION']
const MOCK_ENTRIES = [
  { id: 'WL-001', niu: 'HT99999999', name: 'JEAN DUPONT', type: 'TERRORISM', risk: 'CRITICAL', added: '2026-01-15', status: 'ACTIVE' },
  { id: 'WL-002', niu: 'HT88888888', name: 'PIERRE PAUL', type: 'ORGANIZED_CRIME', risk: 'HIGH', added: '2026-03-20', status: 'ACTIVE' },
  { id: 'WL-003', niu: 'HT77777777', name: 'MARIE JEAN', type: 'NARCOTICS', risk: 'HIGH', added: '2026-04-10', status: 'ACTIVE' },
]

export default function Watchlist() {
  const [filter, setFilter] = useState('ALL')
  const [search, setSearch] = useState('')

  const entries = MOCK_ENTRIES.filter((e) => {
    if (filter !== 'ALL' && e.type !== filter) return false
    if (search && !e.name.toLowerCase().includes(search.toLowerCase())) return false
    return true
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Watchlist Nationale</h1>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Rechercher..."
            className="bg-gray-800 border border-gray-700 rounded-lg pl-9 pr-4 py-1.5 text-sm w-64"
          />
        </div>
      </div>

      <div className="flex gap-2">
        {['ALL', ...WATCHLIST_TYPES].map((t) => (
          <button
            key={t}
            onClick={() => setFilter(t)}
            className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
              filter === t
                ? 'bg-primary-600/20 text-primary-400 border border-primary-500/30'
                : 'bg-gray-800 text-gray-400 border border-gray-700 hover:border-gray-600'
            }`}
          >
            {t === 'ALL' ? 'Tous' : t.replace('_', ' ')}
          </button>
        ))}
      </div>

      <div className="card">
        <table className="w-full">
          <thead>
            <tr className="text-xs text-gray-500 uppercase border-b border-gray-800">
              <th className="text-left py-3 px-2">ID</th>
              <th className="text-left py-3 px-2">NIU</th>
              <th className="text-left py-3 px-2">Nom</th>
              <th className="text-left py-3 px-2">Type</th>
              <th className="text-left py-3 px-2">Risque</th>
              <th className="text-left py-3 px-2">Ajouté</th>
              <th className="text-left py-3 px-2">Statut</th>
            </tr>
          </thead>
          <tbody>
            {entries.map((e) => (
              <tr key={e.id} className="border-b border-gray-800/50 hover:bg-gray-800/30 cursor-pointer">
                <td className="py-3 px-2 text-sm">{e.id}</td>
                <td className="py-3 px-2 text-sm font-mono">{e.niu}</td>
                <td className="py-3 px-2 text-sm font-medium">{e.name}</td>
                <td className="py-3 px-2">
                  <span className="text-xs bg-gray-800 px-2 py-0.5 rounded">{e.type.replace('_', ' ')}</span>
                </td>
                <td className="py-3 px-2">
                  <span className={`badge-${e.risk.toLowerCase()}`}>{e.risk}</span>
                </td>
                <td className="py-3 px-2 text-sm text-gray-400">{e.added}</td>
                <td className="py-3 px-2">
                  <span className="text-xs text-green-400 bg-green-500/10 px-2 py-0.5 rounded">{e.status}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
