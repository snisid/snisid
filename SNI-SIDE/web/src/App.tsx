import { Routes, Route, Navigate } from 'react-router-dom'
import Sidebar from './components/Sidebar'
import Header from './components/Header'
import Dashboard from './pages/Dashboard'
import Alerts from './pages/Alerts'
import Search from './pages/Search'
import Intelligence from './pages/Intelligence'
import EntityDetail from './pages/EntityDetail'
import Watchlist from './pages/Watchlist'
import Cases from './pages/Cases'

export default function App() {
  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header />
        <main className="flex-1 overflow-y-auto p-6">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/alerts" element={<Alerts />} />
            <Route path="/search" element={<Search />} />
            <Route path="/intelligence" element={<Intelligence />} />
            <Route path="/entity/:type/:id" element={<EntityDetail />} />
            <Route path="/watchlist" element={<Watchlist />} />
            <Route path="/cases" element={<Cases />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}
