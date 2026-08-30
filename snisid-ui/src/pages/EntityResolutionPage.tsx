import { motion } from 'framer-motion'
import { Users, Link, Search, Filter, Shield, AlertTriangle } from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'

export default function EntityResolutionPage() {
  const { entities } = useSNISIDStore()

  return (
    <div className="min-h-screen bg-snisid-black relative overflow-hidden">
      <div className="grid-bg" />
      
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // PALANTIR GOTHAM ENTITY RESOLUTION
      </div>

      <header className="fixed top-10 left-0 right-0 z-40 glass-panel mx-4 mt-2 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Users className="w-6 h-6 text-snisid-info" />
          <h1 className="text-lg font-bold">Résolution d'Entités</h1>
        </div>
        <div className="flex items-center gap-2">
          <button className="cyber-button text-xs">
            <Search className="w-4 h-4 mr-2 inline" />
            Recherche
          </button>
          <button className="cyber-button text-xs">
            <Filter className="w-4 h-4 mr-2 inline" />
            Filtres
          </button>
        </div>
      </header>

      <main className="pt-32 pb-8 px-4 relative z-10">
        {/* Stats Bar */}
        <div className="glass-panel p-4 mb-6 flex items-center justify-between flex-wrap gap-4">
          <div className="flex items-center gap-6">
            <div>
              <p className="text-2xl font-bold text-snisid-info">{entities.length}</p>
              <p className="text-xs text-gray-500">Total Entités</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-snisid-danger">{entities.filter(e => e.riskScore > 80).length}</p>
              <p className="text-xs text-gray-500">Haut Risque</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-snisid-warning">{entities.filter(e => e.status === 'watchlist').length}</p>
              <p className="text-xs text-gray-500">Surveillance</p>
            </div>
          </div>
        </div>

        {/* Entity Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {entities.map((entity, index) => (
            <motion.div
              key={entity.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="data-card"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                    entity.type === 'person' ? 'bg-blue-950/50 border border-blue-500/50' :
                    entity.type === 'organization' ? 'bg-purple-950/50 border border-purple-500/50' :
                    entity.type === 'location' ? 'bg-green-950/50 border border-green-500/50' :
                    'bg-gray-950/50 border border-gray-500/50'
                  }`}>
                    <Users className="w-5 h-5 text-white" />
                  </div>
                  <div>
                    <h3 className="font-medium">{entity.name}</h3>
                    <p className="text-xs text-gray-500 uppercase">{entity.type}</p>
                  </div>
                </div>
                <div className={`px-2 py-1 rounded text-xs font-mono ${
                  entity.riskScore > 80 ? 'bg-red-950/50 text-red-400' :
                  entity.riskScore > 50 ? 'bg-yellow-950/50 text-yellow-400' :
                  'bg-green-950/50 text-green-400'
                }`}>
                  RISK: {entity.riskScore}
                </div>
              </div>

              <div className="space-y-2 text-xs">
                <div className="flex items-center justify-between">
                  <span className="text-gray-500">Statut:</span>
                  <span className={`font-mono ${
                    entity.status === 'active' ? 'text-green-400' :
                    entity.status === 'watchlist' ? 'text-yellow-400' :
                    entity.status === 'detained' ? 'text-red-400' : 'text-gray-400'
                  }`}>
                    {entity.status.toUpperCase()}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-500">Connexions:</span>
                  <span className="font-mono flex items-center gap-1">
                    <Link className="w-3 h-3" />
                    {entity.connections.length}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-500">Dernière activité:</span>
                  <span className="font-mono">
                    {new Date(entity.lastSeen).toLocaleDateString()}
                  </span>
                </div>
              </div>

              <div className="mt-4 pt-4 border-t border-white/10 flex items-center gap-2">
                <button className="cyber-button text-xs flex-1 py-2">
                  <Search className="w-3 h-3 mr-1 inline" />
                  Analyser
                </button>
                {entity.riskScore > 70 && (
                  <button className="danger-button text-xs px-3 py-2">
                    <AlertTriangle className="w-3 h-3" />
                  </button>
                )}
              </div>
            </motion.div>
          ))}
        </div>
      </main>

      <div className="classification-top-secret fixed bottom-0 left-0 right-0 z-40">
        TOP SECRET // NOFORN // PALANTIR GOTHAM ENTITY RESOLUTION
      </div>
    </div>
  )
}
