import { useQuery } from '@tanstack/react-query'
import {
  Activity, AlertTriangle, Users, Car, Database, GitBranch,
  Shield, TrendingUp
} from 'lucide-react'
import { api } from '../services/api'

const stats = [
  { key: 'total_alerts_24h', label: 'Alertes 24h', icon: AlertTriangle, color: 'text-red-400' },
  { key: 'critical_alerts', label: 'Critiques', icon: Activity, color: 'text-red-500' },
  { key: 'persons_wanted', label: 'Personnes recherchées', icon: Users, color: 'text-orange-400' },
  { key: 'vehicles_wanted', label: 'Véhicules recherchés', icon: Car, color: 'text-yellow-400' },
  { key: 'events_per_second', label: 'Événements/s', icon: TrendingUp, color: 'text-green-400' },
  { key: 'graph_nodes', label: 'Nœuds graphe', icon: GitBranch, color: 'text-blue-400' },
  { key: 'graph_relationships', label: 'Relations graphe', icon: Shield, color: 'text-purple-400' },
  { key: 'active_cases', label: 'Cas actifs', icon: Database, color: 'text-cyan-400' },
]

function StatCard({ label, value, icon: Icon, color }: any) {
  return (
    <div className="card flex items gap-4">
      <div className={`p-3 rounded-lg bg-gray-800 ${color}`}>
        <Icon className="w-6 h-6" />
      </div>
      <div>
        <p className="stat-value">{value ?? '—'}</p>
        <p className="stat-label">{label}</p>
      </div>
    </div>
  )
}

export default function Dashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: () => api.getDashboard().catch(() => null),
  })

  const recentAlerts = [
    { type: 'ALPR_WANTED_HIT', severity: 'CRITICAL', desc: 'Véhicule recherché détecté — AB-123-CD à Delmas', time: '2 min' },
    { type: 'BORDER_WANTED', severity: 'CRITICAL', desc: 'Personne recherchée à la frontière — HT12345678', time: '5 min' },
    { type: 'DNA_MATCH', severity: 'HIGH', desc: 'Correspondance ADN — Cas #2024-0891', time: '15 min' },
    { type: 'AML_ALERT', severity: 'HIGH', desc: 'Transaction suspecte — $50,000 USD', time: '22 min' },
    { type: 'WATCHLIST_MATCH', severity: 'MEDIUM', desc: 'Correspondance watchlist — Jean Dupont', time: '1h' },
  ]

  const severityClass = (s: string) => {
    if (s === 'CRITICAL') return 'badge-critical'
    if (s === 'HIGH') return 'badge-high'
    return 'badge-medium'
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Tableau de Bord National</h1>
      <div className="grid grid-cols-4 gap-4">
        {stats.map((s) => (
          <StatCard key={s.key} {...s} value={data?.[s.key]} />
        ))}
      </div>
      <div className="grid grid-cols-3 gap-6">
        <div className="card col-span-2">
          <h2 className="text-lg font-semibold mb-4">Alertes Récentes</h2>
          <div className="space-y-3">
            {recentAlerts.map((a, i) => (
              <div key={i} className="flex items-start gap-3 pb-3 border-b border-gray-800 last:border-0">
                <span className={severityClass(a.severity)}>{a.severity}</span>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-gray-200 truncate">{a.desc}</p>
                  <p className="text-xs text-gray-500">{a.type}</p>
                </div>
                <span className="text-xs text-gray-500 whitespace-nowrap">{a.time}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Santé du Système</h2>
          <div className="space-y-3">
            {[
              { name: 'NCID', status: 'healthy' },
              { name: 'Biometrics', status: 'healthy' },
              { name: 'CODIS', status: 'healthy' },
              { name: 'ALPR', status: 'healthy' },
              { name: 'Neo4j', status: 'healthy' },
              { name: 'Kafka', status: 'degraded' },
            ].map((s) => (
              <div key={s.name} className="flex items-center justify-between py-1">
                <span className="text-sm text-gray-300">{s.name}</span>
                <span className={`text-xs px-2 py-0.5 rounded ${
                  s.status === 'healthy' ? 'bg-green-500/10 text-green-400' : 'bg-yellow-500/10 text-yellow-400'
                }`}>
                  {s.status === 'healthy' ? 'OK' : 'Dégradé'}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
