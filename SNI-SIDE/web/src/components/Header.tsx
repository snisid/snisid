import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Search, Bell, Settings } from 'lucide-react'

export default function Header() {
  const [query, setQuery] = useState('')
  const navigate = useNavigate()

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.trim()) navigate(`/search?q=${encodeURIComponent(query.trim())}`)
  }

  return (
    <header className="h-16 bg-gray-900 border-b border-gray-800 flex items-center px-6 gap-4">
      <form onSubmit={handleSearch} className="flex-1 max-w-xl">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Recherche unifiée (NIU, plaque, nom, téléphone, IP...)"
            className="w-full bg-gray-800 border border-gray-700 rounded-lg pl-10 pr-4 py-2 text-sm
                       focus:outline-none focus:border-primary-500 focus:ring-1 focus:ring-primary-500/50
                       placeholder:text-gray-600"
          />
        </div>
      </form>
      <button className="relative p-2 text-gray-400 hover:text-white transition-colors">
        <Bell className="w-5 h-5" />
        <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full" />
      </button>
      <button className="p-2 text-gray-400 hover:text-white transition-colors">
        <Settings className="w-5 h-5" />
      </button>
    </header>
  )
}
