import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Shield, Lock, Key, Eye, EyeOff, AlertTriangle } from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'

export default function LoginPage() {
  const navigate = useNavigate()
  const login = useSNISIDStore((state) => state.login)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [mfaCode, setMfaCode] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const success = await login({ username, password, mfaCode })
      if (success) {
        navigate('/dashboard')
      } else {
        setError('Échec de l\'authentification. Vérifiez vos identifiants.')
      }
    } catch (err) {
      setError('Erreur système. Contactez l\'administrateur.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center relative overflow-hidden">
      {/* Background Effects */}
      <div className="grid-bg" />
      <div className="absolute inset-0 bg-radial-gradient opacity-30" />
      
      {/* Classification Banner */}
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // SNISID CLASSIFIED SYSTEM
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="glass-panel p-8 w-full max-w-md mx-4 z-10"
      >
        {/* Header */}
        <div className="text-center mb-8">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
            className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gradient-to-br from-snisid-primary to-snisid-secondary mb-4 shadow-lg shadow-snisid-primary/30"
          >
            <Shield className="w-10 h-10 text-white" />
          </motion.div>
          <h1 className="text-2xl font-bold mb-2 text-shadow-glow">
            SNISID
          </h1>
          <p className="text-sm text-gray-400 font-mono">
            Système National d'Identification Sécurisée
          </p>
          <p className="text-xs text-gray-500 mt-1">
            République d'Haïti - Sécurité Nationale
          </p>
        </div>

        {/* Warning Banner */}
        <div className="alert-warning mb-6 flex items-start gap-2">
          <AlertTriangle className="w-4 h-4 mt-0.5 flex-shrink-0" />
          <p className="text-xs">
            ACCÈS RESTREINT - Système classifié. Toute tentative d'accès non autorisé sera poursuivie selon les lois haïtiennes et internationales.
          </p>
        </div>

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Username */}
          <div>
            <label className="block text-xs font-mono text-gray-400 mb-1">
              IDENTIFIANT
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full bg-snisid-dark/50 border border-white/10 rounded-md py-2 pl-10 pr-4 text-sm focus:outline-none focus:border-snisid-primary/50 focus:ring-1 focus:ring-snisid-primary/50 transition-colors font-mono"
                placeholder="Entrez votre identifiant"
                required
              />
            </div>
          </div>

          {/* Password */}
          <div>
            <label className="block text-xs font-mono text-gray-400 mb-1">
              MOT DE PASSE
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-snisid-dark/50 border border-white/10 rounded-md py-2 pl-10 pr-10 text-sm focus:outline-none focus:border-snisid-primary/50 focus:ring-1 focus:ring-snisid-primary/50 transition-colors font-mono"
                placeholder="Entrez votre mot de passe"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300"
              >
                {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
          </div>

          {/* MFA Code */}
          <div>
            <label className="block text-xs font-mono text-gray-400 mb-1">
              CODE MFA (6 chiffres)
            </label>
            <div className="relative">
              <Key className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
              <input
                type="text"
                value={mfaCode}
                onChange={(e) => setMfaCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                className="w-full bg-snisid-dark/50 border border-white/10 rounded-md py-2 pl-10 pr-4 text-sm focus:outline-none focus:border-snisid-primary/50 focus:ring-1 focus:ring-snisid-primary/50 transition-colors font-mono tracking-widest"
                placeholder="000000"
                maxLength={6}
                required
              />
            </div>
          </div>

          {/* Error Message */}
          {error && (
            <div className="alert-critical flex items-center gap-2">
              <AlertTriangle className="w-4 h-4" />
              <span className="text-xs">{error}</span>
            </div>
          )}

          {/* Submit Button */}
          <button
            type="submit"
            disabled={loading}
            className="cyber-button w-full py-3 font-mono text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <span className="flex items-center justify-center gap-2">
                <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                AUTHENTIFICATION...
              </span>
            ) : (
              'ACCÉDER AU SYSTÈME'
            )}
          </button>
        </form>

        {/* Footer */}
        <div className="mt-6 pt-6 border-t border-white/10">
          <p className="text-xs text-gray-500 text-center font-mono">
            Version: 1.0.0-secure • Build: {new Date().toISOString().split('T')[0]}
          </p>
          <p className="text-xs text-gray-600 text-center mt-1">
            © 2024 Gouvernement Haïtien - Tous droits réservés
          </p>
        </div>
      </motion.div>

      {/* Bottom Classification Banner */}
      <div className="classification-top-secret fixed bottom-0 left-0 right-0">
        TOP SECRET // NOFORN // SNISID CLASSIFIED SYSTEM
      </div>
    </div>
  )
}
