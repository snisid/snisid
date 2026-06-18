import { useState } from 'react'
import { Bell, AlertTriangle, Filter, CheckCircle, Clock } from 'lucide-react'

const MOCK_ALERTS = [
  { id: '1', type: 'ALPR_WANTED_HIT', severity: 'CRITICAL', title: 'Véhicule recherché détecté',
    description: 'AB-123-CD repéré à Delmas, Route de l\'Aéroport. Propriétaire lié au Gang 400 Mawozo.',
    timestamp: Date.now() - 120_000, source: 'sniside-alpr-ingest', confidence: 0.98 },
  { id: '2', type: 'BORDER_WANTED', severity: 'CRITICAL', title: 'Personne recherchée à la frontière',
    description: 'JEAN DUPONT (HT12345678) détecté au poste-frontière de Malpasse. Mandat d\'arrêt international.',
    timestamp: Date.now() - 300_000, source: 'sniside-border-api', confidence: 0.95 },
  { id: '3', type: 'DNA_MATCH', severity: 'HIGH', title: 'Correspondance ADN positive',
    description: 'Profil ADN de scène de crime #SCI-2024-0891 correspond à un profil CODIS existant.',
    timestamp: Date.now() - 900_000, source: 'sniside-codis-api', confidence: 0.999 },
  { id: '4', type: 'AML_ALERT', severity: 'HIGH', title: 'Transaction blanchiment suspect',
    description: 'Transfert de $50,000 USD depuis compte lié à une société-écran vers bénéficiaire en double NIU.',
    timestamp: Date.now() - 1_320_000, source: 'sniside-aml-engine', confidence: 0.87 },
  { id: '5', type: 'WATCHLIST_MATCH', severity: 'MEDIUM', title: 'Correspondance watchlist',
    description: 'MARIE JEAN (HT87654321) correspond à une entrée watchlist catégorie NARCOTICS.',
    timestamp: Date.now() - 3_600_000, source: 'sniside-watchlist-api', confidence: 0.72 },
  { id: '6', type: 'CYBER_INCIDENT', severity: 'CRITICAL', title: 'Incident cyber — Campagne C2',
    description: 'IP 185.220.101.42 identifié comme serveur C2 actif ciblant les infrastructures gouvernementales.',
    timestamp: Date.now() - 4_200_000, source: 'sniside-cyber-api', confidence: 0.96 },
  { id: '7', type: 'BEHAVIORAL_ANOMALY', severity: 'MEDIUM', title: 'Anomalie comportementale',
    description: 'Agent PNH-4423 accède à 150+ dossiers sensibles en 1 heure. Pattern inhabituel.',
    timestamp: Date.now() - 7_200_000, source: 'sniside-ai-fusion', confidence: 0.81 },
  { id: '8', type: 'GRAPH_INSIGHT', severity: 'LOW', title: 'Nouveau réseau détecté',
    description: 'Connexion identifiée entre 5 personnes via des transactions financières croisées.',
    timestamp: Date.now() - 14_400_000, source: 'sniside-graphrag-engine', confidence: 0.65 },
]

const SEVERITIES = ['ALL', 'CRITICAL', 'HIGH', 'MEDIUM', 'LOW']

export default function Alerts() {
  const [filter, setFilter] = useState('ALL')
  const alerts = filter === 'ALL' ? MOCK_ALERTS : MOCK_ALERTS.filter((a) => a.severity === filter)

  const severityBadge = (s: string) => `badge-${s.toLowerCase()}`
  const timeAgo = (ts: number) => {
    const diff = Date.now() - ts
    const mins = Math.floor(diff / 60000)
    if (mins < 60) return `${mins} min`
    const hrs = Math.floor(mins / 60)
    return `${hrs}h`
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Centre d'Alertes</h1>
        <div className="flex gap-2">
          {SEVERITIES.map((s) => (
            <button
              key={s}
              onClick={() => setFilter(s)}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                filter === s
                  ? 'bg-primary-600/20 text-primary-400 border border-primary-500/30'
                  : 'bg-gray-800 text-gray-400 border border-gray-700 hover:border-gray-600'
              }`}
            >
              {s === 'ALL' ? 'Toutes' : s}
            </button>
          ))}
        </div>
      </div>

      <div className="space-y-3">
        {alerts.map((alert) => (
          <div key={alert.id} className="card hover:border-gray-600 transition-colors">
            <div className="flex items-start gap-4">
              <div className={`p-2 rounded-lg ${
                alert.severity === 'CRITICAL' ? 'bg-red-500/10' :
                alert.severity === 'HIGH' ? 'bg-orange-500/10' :
                'bg-gray-800'
              }`}>
                {alert.severity === 'CRITICAL' ? <AlertTriangle className="w-5 h-5 text-red-400" /> :
                 <Bell className="w-5 h-5 text-gray-400" />}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className={severityBadge(alert.severity)}>{alert.severity}</span>
                  <span className="text-xs text-gray-500">{alert.type}</span>
                  <span className="text-xs text-gray-600">{timeAgo(alert.timestamp)}</span>
                </div>
                <p className="font-medium">{alert.title}</p>
                <p className="text-sm text-gray-400 mt-1">{alert.description}</p>
                <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                  <span>Source: {alert.source}</span>
                  <span>Confiance: {(alert.confidence * 100).toFixed(0)}%</span>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
