import { NavLink } from 'react-router-dom'
import { Shield, Search, Bell, GitBranch, Users, FileText, BookmarkCheck, Activity } from 'lucide-react'

const links = [
  { to: '/', label: 'Dashboard', icon: Activity },
  { to: '/search', label: 'Recherche', icon: Search },
  { to: '/alerts', label: 'Alertes', icon: Bell },
  { to: '/intelligence', label: 'Intelligence', icon: GitBranch },
  { to: '/cases', label: 'Cas', icon: FileText },
  { to: '/watchlist', label: 'Watchlist', icon: BookmarkCheck },
  { to: '/search?type=person', label: 'Personnes', icon: Users },
]

export default function Sidebar() {
  return (
    <aside className="w-64 bg-gray-900 border-r border-gray-800 flex flex-col">
      <div className="p-5 border-b border-gray-800">
        <div className="flex items-center gap-3">
          <Shield className="w-8 h-8 text-primary-500" />
          <div>
            <h1 className="text-lg font-bold">SNI-SIDE</h1>
            <p className="text-xs text-gray-400">Intelligence Dashboard</p>
          </div>
        </div>
      </div>
      <nav className="flex-1 p-3 space-y-1">
        {links.map((l) => (
          <NavLink
            key={l.to}
            to={l.to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-primary-600/20 text-primary-400'
                  : 'text-gray-400 hover:text-white hover:bg-gray-800'
              }`
            }
          >
            <l.icon className="w-5 h-5" />
            {l.label}
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t border-gray-800">
        <div className="flex items-center gap-2 text-xs text-gray-500">
          <div className="w-2 h-2 rounded-full bg-green-500" />
          Système opérationnel
        </div>
        <p className="text-xs text-gray-600 mt-1">v2.1.0</p>
      </div>
    </aside>
  )
}
