import { useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Shield, Activity, Users, MapPin, AlertTriangle, 
  TrendingUp, Database, Lock, Wifi, FileText,
  Search, Eye, Target, Globe, Server
} from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'
import { Link } from 'react-router-dom'

export default function DashboardPage() {
  const { 
    user, alerts, entities, surveillanceFeeds, 
    activeThreats, totalEntities, systemStatus, logout 
  } = useSNISIDStore()

  return (
    <div className="min-h-screen bg-snisid-black relative overflow-hidden">
      {/* Background */}
      <div className="grid-bg" />
      
      {/* Top Classification Banner */}
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // SNISID COMMAND CENTER
      </div>

      {/* Header */}
      <header className="fixed top-10 left-0 right-0 z-40 glass-panel mx-4 mt-2 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Shield className="w-6 h-6 text-snisid-primary" />
          <h1 className="text-lg font-bold text-shadow-glow">SNISID</h1>
          <nav className="hidden md:flex items-center gap-1 ml-8">
            <Link to="/dashboard" className="px-4 py-2 rounded-md bg-snisid-primary/20 text-snisid-primary text-sm font-medium">
              Dashboard
            </Link>
            <Link to="/surveillance" className="px-4 py-2 rounded-md hover:bg-white/5 text-gray-400 hover:text-white text-sm font-medium transition-colors">
              Surveillance
            </Link>
            <Link to="/entities" className="px-4 py-2 rounded-md hover:bg-white/5 text-gray-400 hover:text-white text-sm font-medium transition-colors">
              Entités
            </Link>
            <Link to="/intelligence" className="px-4 py-2 rounded-md hover:bg-white/5 text-gray-400 hover:text-white text-sm font-medium transition-colors">
              Renseignement
            </Link>
            <Link to="/settings" className="px-4 py-2 rounded-md hover:bg-white/5 text-gray-400 hover:text-white text-sm font-medium transition-colors">
              Paramètres
            </Link>
          </nav>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right hidden sm:block">
            <p className="text-xs font-mono text-gray-400">{user?.name}</p>
            <p className="text-xs font-mono text-snisid-primary uppercase">{user?.clearance}</p>
          </div>
          <button onClick={logout} className="danger-button text-xs py-2">
            DÉCONNEXION
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="pt-32 pb-8 px-4 relative z-10">
        {/* System Status Bar */}
        <div className="glass-panel p-4 mb-6 flex items-center justify-between flex-wrap gap-4">
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2">
              <Database className={`w-4 h-4 ${systemStatus.database === 'online' ? 'text-green-500' : 'text-red-500'}`} />
              <span className="text-xs font-mono">BDD: {systemStatus.database.toUpperCase()}</span>
            </div>
            <div className="flex items-center gap-2">
              <Wifi className={`w-4 h-4 ${systemStatus.network === 'online' ? 'text-green-500' : 'text-red-500'}`} />
              <span className="text-xs font-mono">RÉSEAU: {systemStatus.network.toUpperCase()}</span>
            </div>
            <div className="flex items-center gap-2">
              <Lock className={`w-4 h-4 ${systemStatus.encryption === 'active' ? 'text-green-500' : 'text-red-500'}`} />
              <span className="text-xs font-mono">CRYPTO: {systemStatus.encryption.toUpperCase()}</span>
            </div>
            <div className="flex items-center gap-2">
              <FileText className={`w-4 h-4 ${systemStatus.auditLog === 'recording' ? 'text-green-500' : 'text-yellow-500'}`} />
              <span className="text-xs font-mono">AUDIT: {systemStatus.auditLog.toUpperCase()}</span>
            </div>
          </div>
          <div className="text-xs font-mono text-gray-500">
            {new Date().toLocaleString('fr-HT', { timeZone: 'America/Port-au-Prince' })}
          </div>
        </div>

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="data-card"
          >
            <div className="flex items-center justify-between mb-2">
              <AlertTriangle className="w-5 h-5 text-snisid-danger" />
              <span className="text-xs font-mono text-gray-500">TEMPS RÉEL</span>
            </div>
            <p className="text-3xl font-bold text-snisid-danger">{activeThreats}</p>
            <p className="text-xs text-gray-400 mt-1">Menaces Actives</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="data-card"
          >
            <div className="flex items-center justify-between mb-2">
              <Users className="w-5 h-5 text-snisid-info" />
              <span className="text-xs font-mono text-gray-500">BASE NATIONALE</span>
            </div>
            <p className="text-3xl font-bold text-snisid-info">{totalEntities}</p>
            <p className="text-xs text-gray-400 mt-1">Entités Enregistrées</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="data-card"
          >
            <div className="flex items-center justify-between mb-2">
              <Eye className="w-5 h-5 text-snisid-success" />
              <span className="text-xs font-mono text-gray-500">SURVEILLANCE</span>
            </div>
            <p className="text-3xl font-bold text-snisid-success">{surveillanceFeeds.filter(f => f.status === 'online').length}/{surveillanceFeeds.length}</p>
            <p className="text-xs text-gray-400 mt-1">Caméras en Ligne</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
            className="data-card"
          >
            <div className="flex items-center justify-between mb-2">
              <Activity className="w-5 h-5 text-snisid-warning" />
              <span className="text-xs font-mono text-gray-500">DERNIÈRES 24H</span>
            </div>
            <p className="text-3xl font-bold text-snisid-warning">{alerts.length}</p>
            <p className="text-xs text-gray-400 mt-1">Alertes Générées</p>
          </motion.div>
        </div>

        {/* Alerts & Feeds Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Recent Alerts */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.5 }}
            className="glass-panel p-6"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-snisid-danger" />
                Alertes Récentes
              </h2>
              <Link to="/intelligence" className="text-xs text-snisid-primary hover:underline">
                Voir tout →
              </Link>
            </div>
            <div className="space-y-3 max-h-96 overflow-y-auto scrollbar-thin">
              {alerts.length === 0 ? (
                <p className="text-sm text-gray-500 text-center py-8">Aucune alerte récente</p>
              ) : (
                alerts.slice(0, 10).map((alert) => (
                  <div key={alert.id} className={`p-3 rounded-md border-l-4 ${
                    alert.type === 'critical' ? 'bg-red-950/30 border-red-500' :
                    alert.type === 'warning' ? 'bg-yellow-950/30 border-yellow-500' :
                    alert.type === 'success' ? 'bg-green-950/30 border-green-500' :
                    'bg-blue-950/30 border-blue-500'
                  }`}>
                    <div className="flex items-start justify-between gap-2">
                      <div className="flex-1">
                        <p className="text-sm font-medium">{alert.message}</p>
                        <p className="text-xs text-gray-500 mt-1 font-mono">
                          {alert.source} • {new Date(alert.timestamp).toLocaleTimeString()}
                        </p>
                      </div>
                      {!alert.acknowledged && (
                        <span className="w-2 h-2 rounded-full bg-snisid-danger animate-pulse" />
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </motion.div>

          {/* Surveillance Feeds */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.6 }}
            className="glass-panel p-6"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold flex items-center gap-2">
                <Eye className="w-5 h-5 text-snisid-success" />
                Flux de Surveillance
              </h2>
              <Link to="/surveillance" className="text-xs text-snisid-primary hover:underline">
                Voir tout →
              </Link>
            </div>
            <div className="grid grid-cols-2 gap-3">
              {surveillanceFeeds.map((feed) => (
                <div key={feed.id} className="data-card p-3">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <div className={`w-2 h-2 rounded-full ${
                        feed.status === 'online' ? 'bg-green-500 animate-pulse' :
                        feed.status === 'maintenance' ? 'bg-yellow-500' : 'bg-red-500'
                      }`} />
                      <span className="text-xs font-medium truncate">{feed.name}</span>
                    </div>
                    {feed.alerts > 0 && (
                      <span className="text-xs bg-snisid-danger text-white px-2 py-0.5 rounded-full">
                        {feed.alerts}
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-gray-500">{feed.location}</p>
                  <p className="text-xs font-mono text-gray-600 mt-1 uppercase">{feed.status}</p>
                </div>
              ))}
            </div>
          </motion.div>
        </div>

        {/* Quick Actions */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7 }}
          className="glass-panel p-6 mt-6"
        >
          <h2 className="text-lg font-bold mb-4 flex items-center gap-2">
            <Target className="w-5 h-5 text-snisid-primary" />
            Actions Rapides
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
            <Link to="/entities" className="cyber-button text-xs py-3 flex items-center justify-center gap-2">
              <Search className="w-4 h-4" />
              Recherche
            </Link>
            <Link to="/surveillance" className="cyber-button text-xs py-3 flex items-center justify-center gap-2">
              <Globe className="w-4 h-4" />
              Carte Live
            </Link>
            <Link to="/intelligence" className="cyber-button text-xs py-3 flex items-center justify-center gap-2">
              <FileText className="w-4 h-4" />
              Rapport
            </Link>
            <Link to="/entities" className="cyber-button text-xs py-3 flex items-center justify-center gap-2">
              <Users className="w-4 h-4" />
              Entités
            </Link>
            <Link to="/settings" className="cyber-button text-xs py-3 flex items-center justify-center gap-2">
              <Server className="w-4 h-4" />
              Système
            </Link>
            <button className="danger-button text-xs py-3 flex items-center justify-center gap-2">
              <AlertTriangle className="w-4 h-4" />
              Urgence
            </button>
          </div>
        </motion.div>
      </main>

      {/* Bottom Classification Banner */}
      <div className="classification-top-secret fixed bottom-0 left-0 right-0 z-40">
        TOP SECRET // NOFORN // SNISID COMMAND CENTER
      </div>
    </div>
  )
}
