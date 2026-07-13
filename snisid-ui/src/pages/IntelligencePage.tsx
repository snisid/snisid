import { motion } from 'framer-motion'
import { FileText, AlertTriangle, TrendingUp, Clock, CheckCircle } from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'

export default function IntelligencePage() {
  const { alerts, acknowledgeAlert } = useSNISIDStore()

  return (
    <div className="min-h-screen bg-snisid-black relative overflow-hidden">
      <div className="grid-bg" />
      
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // XKEYSCORE INTELLIGENCE CENTER
      </div>

      <header className="fixed top-10 left-0 right-0 z-40 glass-panel mx-4 mt-2 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <FileText className="w-6 h-6 text-snisid-warning" />
          <h1 className="text-lg font-bold">Centre de Renseignement</h1>
        </div>
        <button className="cyber-button text-xs">
          <TrendingUp className="w-4 h-4 mr-2 inline" />
          Rapport
        </button>
      </header>

      <main className="pt-32 pb-8 px-4 relative z-10">
        {/* Alert Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="data-card">
            <div className="flex items-center gap-3">
              <AlertTriangle className="w-8 h-8 text-red-500" />
              <div>
                <p className="text-2xl font-bold">{alerts.filter(a => a.type === 'critical').length}</p>
                <p className="text-xs text-gray-500">Critiques</p>
              </div>
            </div>
          </div>
          <div className="data-card">
            <div className="flex items-center gap-3">
              <AlertTriangle className="w-8 h-8 text-yellow-500" />
              <div>
                <p className="text-2xl font-bold">{alerts.filter(a => a.type === 'warning').length}</p>
                <p className="text-xs text-gray-500">Avertissements</p>
              </div>
            </div>
          </div>
          <div className="data-card">
            <div className="flex items-center gap-3">
              <Clock className="w-8 h-8 text-blue-500" />
              <div>
                <p className="text-2xl font-bold">{alerts.filter(a => !a.acknowledged).length}</p>
                <p className="text-xs text-gray-500">Non acquittées</p>
              </div>
            </div>
          </div>
          <div className="data-card">
            <div className="flex items-center gap-3">
              <CheckCircle className="w-8 h-8 text-green-500" />
              <div>
                <p className="text-2xl font-bold">{alerts.filter(a => a.acknowledged).length}</p>
                <p className="text-xs text-gray-500">Acquittées</p>
              </div>
            </div>
          </div>
        </div>

        {/* Alerts Timeline */}
        <div className="glass-panel p-6">
          <h2 className="text-lg font-bold mb-4 flex items-center gap-2">
            <Clock className="w-5 h-5" />
            Chronologie des Alertes
          </h2>
          <div className="space-y-4 max-h-[600px] overflow-y-auto scrollbar-thin">
            {alerts.length === 0 ? (
              <p className="text-center text-gray-500 py-8">Aucune alerte dans le système</p>
            ) : (
              alerts.map((alert, index) => (
                <motion.div
                  key={alert.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                  className={`relative pl-6 pb-4 border-l-2 ${
                    alert.type === 'critical' ? 'border-red-500' :
                    alert.type === 'warning' ? 'border-yellow-500' :
                    alert.type === 'success' ? 'border-green-500' : 'border-blue-500'
                  }`}
                >
                  {/* Timeline dot */}
                  <div className={`absolute left-[-5px] top-0 w-2.5 h-2.5 rounded-full ${
                    alert.type === 'critical' ? 'bg-red-500' :
                    alert.type === 'warning' ? 'bg-yellow-500' :
                    alert.type === 'success' ? 'bg-green-500' : 'bg-blue-500'
                  } ${!alert.acknowledged ? 'animate-pulse' : ''}`} />

                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`text-xs font-mono px-2 py-0.5 rounded ${
                          alert.type === 'critical' ? 'bg-red-950/50 text-red-400' :
                          alert.type === 'warning' ? 'bg-yellow-950/50 text-yellow-400' :
                          alert.type === 'success' ? 'bg-green-950/50 text-green-400' : 'bg-blue-950/50 text-blue-400'
                        }`}>
                          {alert.type.toUpperCase()}
                        </span>
                        <span className="text-xs text-gray-500 font-mono">{alert.source}</span>
                      </div>
                      <p className="text-sm font-medium">{alert.message}</p>
                      <p className="text-xs text-gray-500 mt-1">
                        {new Date(alert.timestamp).toLocaleString('fr-HT', { timeZone: 'America/Port-au-Prince' })}
                      </p>
                    </div>
                    {!alert.acknowledged && (
                      <button
                        onClick={() => acknowledgeAlert(alert.id)}
                        className="cyber-button text-xs py-1 px-3"
                      >
                        Acquitter
                      </button>
                    )}
                    {alert.acknowledged && (
                      <span className="text-xs text-green-500 flex items-center gap-1">
                        <CheckCircle className="w-3 h-3" />
                        Acquitté
                      </span>
                    )}
                  </div>
                </motion.div>
              ))
            )}
          </div>
        </div>
      </main>

      <div className="classification-top-secret fixed bottom-0 left-0 right-0 z-40">
        TOP SECRET // NOFORN // XKEYSCORE INTELLIGENCE CENTER
      </div>
    </div>
  )
}
