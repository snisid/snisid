import { create } from 'zustand';

// Use Zustand for application state (non-PII).
// No persisting to localStorage to comply with security constraints.

type SecurityState = {
  roles: string[];
  permissions: string[];
  isAuthenticated: boolean;
  agencyId: string | null;
  setAuthData: (roles: string[], permissions: string[], agencyId: string | null) => void;
  clearAuthData: () => void;
  hasRole: (role: string) => boolean;
  hasPermission: (permission: string) => boolean;
};

export const useSecurityStore = create<SecurityState>((set, get) => ({
  roles: [],
  permissions: [],
  isAuthenticated: false,
  agencyId: null,
  
  setAuthData: (roles, permissions, agencyId) => set({
    roles,
    permissions,
    agencyId,
    isAuthenticated: true
  }),
  
  clearAuthData: () => set({
    roles: [],
    permissions: [],
    agencyId: null,
    isAuthenticated: false
  }),
  
  hasRole: (role: string) => {
    const { roles } = get();
    return roles.includes(role) || roles.includes('SUPER_ADMIN');
  },
  
  hasPermission: (permission: string) => {
    const { permissions, roles } = get();
    return permissions.includes(permission) || roles.includes('SUPER_ADMIN');
  }
}));
