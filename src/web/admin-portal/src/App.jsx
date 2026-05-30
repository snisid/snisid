import { useState } from 'react'
import './index.css'

function App() {
  const [niu, setNiu] = useState('')
  const [citizen, setCitizen] = useState(null)
  const [loading, setLoading] = useState(false)

  const handleSearch = async (e) => {
    e.preventDefault()
    if (!niu) return

    setLoading(true)
    
    // Simulating the API call to our Go Citizen Registry Service
    // In production, this would be: fetch(`/v1/registry/citizens/${niu}`)
    setTimeout(() => {
      setCitizen({
        niu: niu,
        fullName: "Jean Doe",
        dateOfBirth: "1990-01-01",
        placeOfBirth: "Port-au-Prince",
        status: "ACTIVE",
        lastUpdated: "2026-05-23T12:00:00Z"
      })
      setLoading(false)
    }, 800)
  }

  return (
    <div className="dashboard-container">
      <header className="header">
        <h1>SNISID Admin Portal</h1>
        <p>Sovereign Citizen Registry & Identity Management</p>
      </header>

      <main>
        <div className="glass-card">
          <form className="search-section" onSubmit={handleSearch}>
            <input 
              type="text" 
              className="search-input" 
              placeholder="Enter National Identification Number (NIU)..." 
              value={niu}
              onChange={(e) => setNiu(e.target.value)}
            />
            <button type="submit" className="btn-primary">
              {loading ? 'Searching...' : 'Lookup Citizen'}
            </button>
          </form>

          {citizen && (
            <div className="citizen-data" style={{ animation: 'fadeIn 0.5s ease-out' }}>
              <div className="data-group">
                <span className="data-label">NIU</span>
                <span className="data-value">{citizen.niu}</span>
              </div>
              <div className="data-group">
                <span className="data-label">Full Name</span>
                <span className="data-value">{citizen.fullName}</span>
              </div>
              <div className="data-group">
                <span className="data-label">Date of Birth</span>
                <span className="data-value">{citizen.dateOfBirth}</span>
              </div>
              <div className="data-group">
                <span className="data-label">Place of Birth</span>
                <span className="data-value">{citizen.placeOfBirth}</span>
              </div>
              <div className="data-group">
                <span className="data-label">Status</span>
                <span className={`status-badge ${citizen.status === 'ACTIVE' ? '' : 'pending'}`}>
                  {citizen.status}
                </span>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}

export default App
