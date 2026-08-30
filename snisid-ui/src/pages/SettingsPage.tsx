import { motion } from 'framer-motion'
import { Settings, Shield, Database, Lock, Wifi, Server, Key, Bell, User } from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'

export default function SettingsPage() {
  const { user, systemStatus, setSystemStatus } = useSNISIDStore()

  return (
    <div className="min-h-screen bg-snisid-black relative overflow-hidden">
      <div className="grid-bg" />
      
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // SYSTEM CONFIGURATION
      </div>

      <header className="fixed top-10 left-0 right-0 z-40 glass-panel mx-4 mt-2 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Settings className="w-6 h-6 text-snisid-primary" />
          <h1 className="text-lg font-bold">Paramètres Système</h1>
        </div>
      </header>

      <main className="pt-32 pb-8 px-4 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* User Profile */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="glass-panel p-6"
          >
            <h2 className="text-lg font-bold mb-4 flex items-center gap-2">
              <User className="w-5 h-5" />
              Profil Utilisateur
            </h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <span className="text-sm text-gray-400">Nom</span>
                <span className="font-mono">{user?.name}</span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <span className="text-sm text-gray-400">Rôle</span>
                <span className="font-mono uppercase">{user?.role}</span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <span className="text-sm text-gray-400">Habilitation</span>
                <span className={`font-mono uppercase ${
                  user?.clearance === 'top-secret' ? 'text-red-400' :
                  user?.clearance === 'secret' ? 'text-orange-400' :
                  user?.clearance === 'confidential' ? 'text-blue-400' : 'text-green-400'
                }`}>
                  {user?.clearance}
                </span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <span className="text-sm text-gray-400">Agence</span>
                <span className="font-mono uppercase">{user?.agency}</span>
              </div>
            </div>
          </motion.div>

          {/* System Status */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="glass-panel p-6"
          >
            <h2 className="text-lg font-bold mb-4 flex items-center gap-2">
              <Server className="w-5 h-5" />
              État du Système
            </h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2">
                  <Database className={`w-4 h-4 ${systemStatus.database === 'online' ? 'text-green-500' : 'text-red-500'}`} />
                  <span className="text-sm">Base de Données</span>
                </div>
                <span className={`font-mono uppercase ${systemStatus.database === 'online' ? 'text-green-400' : 'text-red-400'}`}>
                  {systemStatus.database}
                </span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2">
                  <Wifi className={`w-4 h-4 ${systemStatus.network === 'online' ? 'text-green-500' : 'text-red-500'}`} />
                  <span className="text-sm">Réseau</span>
                </div>
                <span className={`font-mono uppercase ${systemStatus.network === 'online' ? 'text-green-400' : 'text-red-400'}`}>
                  {systemStatus.network}
                </span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2">
                  <Lock className={`w-4 h-4 ${systemStatus.encryption === 'active' ? 'text-green-500' : 'text-red-500'}`} />
                  <span className="text-sm">Chiffrement</span>
                </div>
                <span className={`font-mono uppercase ${systemStatus.encryption === 'active' ? 'text-green-400' : 'text-red-400'}`}>
                  {systemStatus.encryption}
                </span>
              </div>
              <div className="flex items-center justify-between p-3 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2">
                  <Bell className={`w-4 h-4 ${systemStatus.auditLog === 'recording' ? 'text-green-500' : 'text-yellow-500'}`} />
                  <span className="text-sm">Journal d'Audit</span>
                </div>
                <span className={`font-mono uppercase ${systemStatus.auditLog === 'recording' ? 'text-green-400' : 'text-yellow-400'}`}>
                  {systemStatus.auditLog}
                </span>
              </div>
            </div>
          </motion.div>

          {/* Security Settings */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="glass-panel p-6 lg:col-span-2"
          >
            <h2 className="text-lg font-bold mb-4 flex items-center gap-2">
              <Shield className="w-5 h-5" />
              Paramètres de Sécurité
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2 mb-2">
                  <Key className="w-4 h-4 text-snisid-primary" />
                  <span className="text-sm font-medium">Gestion des Clés</span>
                </div>
                <p className="text-xs text-gray-500 mb-3">Configuration TPM 2.0 et clés cryptographiques</p>
                <button className="cyber-button text-xs w-full py-2">Configurer</button>
              </div>
              <div className="p-4 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2 mb-2">
                  <Lock className="w-4 h-4 text-snisid-primary" />
                  <span className="text-sm font-medium">Chiffrement</span>
                </div>
                <p className="text-xs text-gray-500 mb-3">AES-512, SM4, FIPS 140-2</p>
                <button className="cyber-button text-xs w-full py-2">Configurer</button>
              </div>
              <div className="p-4 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2 mb-2">
                  <Shield className="w-4 h-4 text-snisid-primary" />
                  <span className="text-sm font-medium">STIG NSA</span>
                </div>
                <p className="text-xs text-gray-500 mb-3">Hardening système selon standards NSA</p>
                <button className="cyber-button text-xs w-full py-2">Vérifier</button>
              </div>
              <div className="p-4 bg-snisid-dark/50 rounded-md">
                <div className="flex items-center gap-2 mb-2">
                  <Bell className="w-4 h-4 text-snisid-primary" />
                  <span className="text-sm font-medium">Audit & Logs</span>
                </div>
                <p className="text-xs text-gray-500 mb-3">Configuration auditd et rétention</p>
                <button className="cyber-button text-xs w-full py-2">Configurer</button>
              </div>
            </div>
          </motion.div>

          {/* System Info */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="glass-panel p-6 lg:col-span-2"
          >
            <h2 className="text-lg font-bold mb-4">Informations Système</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
              <div>
                <p className="text-gray-500">Version SNISID</p>
                <p className="font-mono">1.0.0-secure</p>
              </div>
              <div>
                <p className="text-gray-500">Build Date</p>
                <p className="font-mono">{new Date().toISOString().split('T')[0]}</p>
              </div>
              <div>
                <p className="text-gray-500">Kernel</p>
                <p className="font-mono">Linux 6.1-hardened</p>
              </div>
              <div>
                <p className="text-gray-500">Architecture</p>
                <p className="font-mono">x86_64 Sovereign</p>
              </div>
              <div>
                <p className="text-gray-500">Classification</p>
                <p className="font-mono text-red-400">TOP SECRET // NOFORN</p>
              </div>
              <div>
                <p className="text-gray-500">Juridiction</p>
                <p className="font-mono">République d'Haïti</p>
              </div>
            </div>
          </motion.div>
        </div>
      </main>

      <div className="classification-top-secret fixed bottom-0 left-0 right-0 z-40">
        TOP SECRET // NOFORN // SYSTEM CONFIGURATION
      </div>
    </div>
  )
}
