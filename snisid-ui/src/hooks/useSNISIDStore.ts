import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'

export interface Alert {
  id: string
  type: 'critical' | 'warning' | 'info' | 'success'
  message: string
  source: string
  timestamp: Date
  acknowledged: boolean
}

export interface Entity {
  id: string
  type: 'person' | 'organization' | 'location' | 'event' | 'device'
  name: string
  riskScore: number
  connections: string[]
  lastSeen: Date
  status: 'active' | 'inactive' | 'watchlist' | 'detained'
}

export interface SurveillanceFeed {
  id: string
  name: string
  location: string
  status: 'online' | 'offline' | 'maintenance'
  thumbnail: string
  alerts: number
}

interface SNISIDState {
  // Authentication
  isAuthenticated: boolean
  user: {
    id: string
    name: string
    role: 'admin' | 'analyst' | 'operator' | 'viewer'
    clearance: 'top-secret' | 'secret' | 'confidential' | 'unclassified'
    agency: 'snisid' | 'fbi' | 'dea' | 'nsa' | 'chinese-partner'
  } | null
  sessionToken: string | null
  
  // Real-time data
  alerts: Alert[]
  entities: Entity[]
  surveillanceFeeds: SurveillanceFeed[]
  activeThreats: number
  totalEntities: number
  
  // System status
  systemStatus: {
    database: 'online' | 'degraded' | 'offline'
    network: 'online' | 'degraded' | 'offline'
    encryption: 'active' | 'compromised'
    auditLog: 'recording' | 'paused'
  }
  
  // Actions
  login: (credentials: { username: string; password: string; mfaCode: string }) => Promise<boolean>
  logout: () => void
  addAlert: (alert: Omit<Alert, 'id' | 'timestamp' | 'acknowledged'>) => void
  acknowledgeAlert: (alertId: string) => void
  updateEntity: (entity: Partial<Entity> & { id: string }) => void
  setSystemStatus: (status: Partial<SNISIDState['systemStatus']>) => void
}

export const useSNISIDStore = create<SNISIDState>()(
  subscribeWithSelector((set, get) => ({
    // Initial state
    isAuthenticated: false,
    user: null,
    sessionToken: null,
    
    alerts: [],
    entities: [],
    surveillanceFeeds: [],
    activeThreats: 0,
    totalEntities: 0,
    
    systemStatus: {
      database: 'online',
      network: 'online',
      encryption: 'active',
      auditLog: 'recording',
    },
    
    // Actions
    login: async ({ username, password, mfaCode }) => {
      // Simulation - In production, this would call the backend API
      if (username && password && mfaCode) {
        set({
          isAuthenticated: true,
          user: {
            id: 'usr_' + Date.now(),
            name: username,
            role: 'admin',
            clearance: 'top-secret',
            agency: 'snisid',
          },
          sessionToken: 'tok_' + Math.random().toString(36).substring(2),
        })
        return true
      }
      return false
    },
    
    logout: () => {
      set({
        isAuthenticated: false,
        user: null,
        sessionToken: null,
      })
    },
    
    addAlert: (alertData) => {
      const newAlert: Alert = {
        ...alertData,
        id: 'alert_' + Date.now(),
        timestamp: new Date(),
        acknowledged: false,
      }
      set((state) => ({
        alerts: [newAlert, ...state.alerts].slice(0, 100), // Keep last 100 alerts
        activeThreats: alertData.type === 'critical' 
          ? state.activeThreats + 1 
          : state.activeThreats,
      }))
    },
    
    acknowledgeAlert: (alertId) => {
      set((state) => ({
        alerts: state.alerts.map((alert) =>
          alert.id === alertId ? { ...alert, acknowledged: true } : alert
        ),
        activeThreats: state.alerts.find((a) => a.id === alertId)?.type === 'critical'
          ? Math.max(0, state.activeThreats - 1)
          : state.activeThreats,
      }))
    },
    
    updateEntity: (entityData) => {
      set((state) => ({
        entities: state.entities.map((entity) =>
          entity.id === entityData.id ? { ...entity, ...entityData } : entity
        ),
      }))
    },
    
    setSystemStatus: (statusUpdate) => {
      set((state) => ({
        systemStatus: { ...state.systemStatus, ...statusUpdate },
      }))
    },
  }))
)

// Auto-generate sample data for demo
setTimeout(() => {
  const store = useSNISIDStore.getState()
  
  // Add sample surveillance feeds
  const feeds: SurveillanceFeed[] = [
    { id: 'cam_001', name: 'Port-au-Prince Centre', location: 'Downtown PAP', status: 'online', thumbnail: '', alerts: 2 },
    { id: 'cam_002', name: 'Aéroport International', location: 'PAP Airport', status: 'online', thumbnail: '', alerts: 0 },
    { id: 'cam_003', name: 'Frontière Dominicaine', location: 'Ouanaminthe', status: 'online', thumbnail: '', alerts: 1 },
    { id: 'cam_004', name: 'Port de Commerce', location: 'PAP Seaport', status: 'maintenance', thumbnail: '', alerts: 0 },
  ]
  
  // Add sample entities
  const entities: Entity[] = [
    { id: 'ent_001', type: 'person', name: 'Jean Pierre', riskScore: 85, connections: ['org_001', 'loc_001'], lastSeen: new Date(), status: 'watchlist' },
    { id: 'ent_002', type: 'organization', name: 'Gang Alpha', riskScore: 95, connections: ['per_001', 'per_003'], lastSeen: new Date(), status: 'active' },
    { id: 'ent_003', type: 'location', name: 'Zone Rouge #12', riskScore: 78, connections: ['org_002'], lastSeen: new Date(), status: 'active' },
  ]
  
  useSNISIDStore.setState({
    surveillanceFeeds: feeds,
    entities: entities,
    totalEntities: entities.length,
    activeThreats: 3,
  })
}, 1000)
