import { useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Search as SearchIcon, User, Car, FileText, Globe, AlertTriangle, Fingerprint } from 'lucide-react'
import { api } from '../services/api'

const typeIcons: Record<string, any> = {
  NIU: User, PLATE: Car, VIN: Car, PASSPORT: FileText,
  PHONE: Globe, EMAIL: Globe, IP: Globe, WALLET: Globe,
  DOMAIN: Globe, NAME: User, CASE: FileText, UNKNOWN: SearchIcon,
}

export default function SearchPage() {
  const [params] = useSearchParams()
  const initialQ = params.get('q') || ''
  const [query, setQuery] = useState(initialQ)
  const navigate = useNavigate()

  const { data, isLoading, error } = useQuery({
    queryKey: ['search', query],
    queryFn: () => api.search(query),
    enabled: query.length > 0,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.trim()) navigate(`/search?q=${encodeURIComponent(query.trim())}`)
  }

  const detectedType = data?.detected_type || 'UNKNOWN'
  const Icon = typeIcons[detectedType] || SearchIcon
  const results: any[] = data?.results || []

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Recherche Unifiée</h1>
      <form onSubmit={handleSubmit} className="relative max-w-2xl">
        <SearchIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Rechercher par NIU, plaque, nom, téléphone, IP, wallet..."
          className="w-full bg-gray-900 border border-gray-700 rounded-xl pl-12 pr-4 py-3.5 text-sm
                     focus:outline-none focus:border-primary-500 focus:ring-1 focus:ring-primary-500/50
                     placeholder:text-gray-600"
        />
      </form>

      {query && (
        <div className="flex items-center gap-2 text-sm text-gray-400">
          <Icon className="w-4 h-4" />
          Type détecté : <span className="text-primary-400 font-medium">{detectedType}</span>
          {data?.execution_time_ms && (
            <span className="text-gray-600">— {data.execution_time_ms}ms</span>
          )}
        </div>
      )}

      {isLoading && (
        <div className="flex items-center justify-center py-20">
          <div className="w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin" />
        </div>
      )}

      {error && (
        <div className="card text-red-400 flex items-center gap-2">
          <AlertTriangle className="w-4 h-4" />
          Erreur de recherche
        </div>
      )}

      {results.length > 0 && (
        <div className="space-y-3">
          <p className="text-sm text-gray-400">{results.length} résultat(s)</p>
          {results.map((r: any, i: number) => (
            <div
              key={i}
              className="card cursor-pointer hover:border-gray-600 transition-colors"
              onClick={() => navigate(`/entity/${r.entity_type || 'unknown'}/${r.entity_id}`)}
            >
              <div className="flex items-start justify-between">
                <div>
                  <p className="font-medium">{r.label || r.entity_id}</p>
                  <p className="text-sm text-gray-400">{r.entity_type} · Score: {(r.score * 100).toFixed(0)}%</p>
                  {r.match_type && (
                    <span className="badge-high mt-1 inline-block">{r.match_type}</span>
                  )}
                </div>
                <span className="text-xs text-gray-500">{r.source}</span>
              </div>
              {r.details && Object.keys(r.details).length > 0 && (
                <div className="mt-2 text-xs text-gray-500 line-clamp-2">
                  {JSON.stringify(r.details).slice(0, 200)}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {query && !isLoading && results.length === 0 && !error && (
        <div className="card text-center py-12">
          <SearchIcon className="w-12 h-12 text-gray-700 mx-auto mb-3" />
          <p className="text-gray-500">Aucun résultat pour "{query}"</p>
        </div>
      )}
    </div>
  )
}
