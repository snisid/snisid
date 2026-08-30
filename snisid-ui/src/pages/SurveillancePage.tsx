import { motion } from 'framer-motion'
import { MapPin, Eye, Wifi, AlertTriangle, Maximize2, Settings } from 'lucide-react'
import { useSNISIDStore } from '../hooks/useSNISIDStore'

export default function SurveillancePage() {
  const { surveillanceFeeds } = useSNISIDStore()

  return (
    <div className="min-h-screen bg-snisid-black relative overflow-hidden">
      <div className="grid-bg" />
      
      <div className="classification-top-secret">
        TOP SECRET // NOFORN // SKYTET SURVEILLANCE SYSTEM
      </div>

      <header className="fixed top-10 left-0 right-0 z-40 glass-panel mx-4 mt-2 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Eye className="w-6 h-6 text-snisid-success" />
          <h1 className="text-lg font-bold">Surveillance - Skynet</h1>
        </div>
        <button className="cyber-button text-xs">
          <Settings className="w-4 h-4 mr-2 inline" />
          Configuration
        </button>
      </header>

      <main className="pt-32 pb-8 px-4 relative z-10">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {surveillanceFeeds.map((feed, index) => (
            <motion.div
              key={feed.id}
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: index * 0.1 }}
              className="glass-panel overflow-hidden"
            >
              {/* Camera Feed Placeholder */}
              <div className="relative aspect-video bg-snisid-dark flex items-center justify-center">
                <div className="absolute inset-0 radar-scan" />
                <div className="text-center">
                  <MapPin className="w-12 h-12 text-gray-600 mx-auto mb-2" />
                  <p className="text-xs text-gray-500 font-mono">{feed.location}</p>
                </div>
                
                {/* Status Indicator */}
                <div className="absolute top-2 left-2 flex items-center gap-2">
                  <div className={`w-2 h-2 rounded-full ${
                    feed.status === 'online' ? 'bg-green-500 animate-pulse' :
                    feed.status === 'maintenance' ? 'bg-yellow-500' : 'bg-red-500'
                  }`} />
                  <span className="text-xs font-mono uppercase bg-black/50 px-2 py-0.5 rounded">
                    {feed.status}
                  </span>
                </div>

                {/* Alerts Badge */}
                {feed.alerts > 0 && (
                  <div className="absolute top-2 right-2 bg-snisid-danger text-white text-xs px-2 py-1 rounded-full flex items-center gap-1 animate-pulse">
                    <AlertTriangle className="w-3 h-3" />
                    {feed.alerts}
                  </div>
                )}

                {/* REC Indicator */}
                {feed.status === 'online' && (
                  <div className="absolute bottom-2 right-2 text-red-500 font-mono text-xs flex items-center gap-1">
                    <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse" />
                    REC
                  </div>
                )}

                {/* Fullscreen Button */}
                <button className="absolute bottom-2 left-2 p-1 bg-black/50 hover:bg-black/70 rounded transition-colors">
                  <Maximize2 className="w-4 h-4 text-gray-400" />
                </button>
              </div>

              {/* Camera Info */}
              <div className="p-3">
                <h3 className="text-sm font-medium mb-1">{feed.name}</h3>
                <div className="flex items-center justify-between text-xs text-gray-500 font-mono">
                  <span>ID: {feed.id}</span>
                  <Wifi className={`w-3 h-3 ${feed.status === 'online' ? 'text-green-500' : 'text-red-500'}`} />
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </main>

      <div className="classification-top-secret fixed bottom-0 left-0 right-0 z-40">
        TOP SECRET // NOFORN // SKYTET SURVEILLANCE SYSTEM
      </div>
    </div>
  )
}
