import { create } from 'zustand';

interface User {
  id: string;
  name: string;
  role: string;
  agency: string;
}

interface Alert {
  id: string;
  type: 'fraud' | 'identity' | 'system';
  message: string;
  time: string;
  severity: 'high' | 'medium' | 'low';
}

interface AppState {
  user: User | null;
  alerts: Alert[];
  setUser: (user: User | null) => void;
  addAlert: (alert: Alert) => void;
  clearAlerts: () => void;
}

export const useStore = create<AppState>((set) => ({
  user: {
    id: '1',
    name: 'Jean-Luc Martí',
    role: 'Super Admin',
    agency: 'National Intelligence'
  },
  alerts: [
    { id: '1', type: 'fraud', message: 'Potential Deepfake detected in Port-au-Prince', time: '2 mins ago', severity: 'high' },
    { id: '2', type: 'identity', message: 'Bulk identity creation attempt blocked', time: '5 mins ago', severity: 'medium' },
  ],
  setUser: (user) => set({ user }),
  addAlert: (alert) => set((state) => ({ alerts: [alert, ...state.alerts].slice(0, 50) })),
  clearAlerts: () => set({ alerts: [] }),
}));
