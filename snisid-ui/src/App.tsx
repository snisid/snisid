import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useSNISIDStore } from './hooks/useSNISIDStore'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import SurveillancePage from './pages/SurveillancePage'
import EntityResolutionPage from './pages/EntityResolutionPage'
import IntelligencePage from './pages/IntelligencePage'
import SettingsPage from './pages/SettingsPage'

function App() {
  const isAuthenticated = useSNISIDStore((state) => state.isAuthenticated)

  return (
    <Router>
      <Routes>
        <Route 
          path="/login" 
          element={!isAuthenticated ? <LoginPage /> : <Navigate to="/dashboard" />} 
        />
        <Route 
          path="/dashboard" 
          element={isAuthenticated ? <DashboardPage /> : <Navigate to="/login" />} 
        />
        <Route 
          path="/surveillance" 
          element={isAuthenticated ? <SurveillancePage /> : <Navigate to="/login" />} 
        />
        <Route 
          path="/entities" 
          element={isAuthenticated ? <EntityResolutionPage /> : <Navigate to="/login" />} 
        />
        <Route 
          path="/intelligence" 
          element={isAuthenticated ? <IntelligencePage /> : <Navigate to="/login" />} 
        />
        <Route 
          path="/settings" 
          element={isAuthenticated ? <SettingsPage /> : <Navigate to="/login" />} 
        />
        <Route 
          path="/" 
          element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} />} 
        />
      </Routes>
    </Router>
  )
}

export default App
