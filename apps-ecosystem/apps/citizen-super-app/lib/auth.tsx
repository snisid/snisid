import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import * as SecureStore from 'expo-secure-store';
import * as LocalAuthentication from 'expo-local-authentication';
import { Platform } from 'react-native';

import { getApiClient, type ApiClient } from './api';

interface User {
  id: string;
  fullName: string;
  nnu: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<void>;
  loginWithBiometrics: () => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  api: ApiClient;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: true,
  });

  const api = getApiClient();

  const checkAuth = async () => {
    try {
      if (Platform.OS === 'web') {
        const token = localStorage.getItem('auth_token');
        const userRaw = localStorage.getItem('auth_user');
        if (token && userRaw) {
          const user: User = JSON.parse(userRaw);
          setState({ user, token, isAuthenticated: true, isLoading: false });
          return;
        }
      } else {
        const token = await SecureStore.getItemAsync('auth_token');
        const userRaw = await SecureStore.getItemAsync('auth_user');
        if (token && userRaw) {
          const user: User = JSON.parse(userRaw);
          setState({ user, token, isAuthenticated: true, isLoading: false });
          return;
        }
      }
    } catch {
    } finally {
      setState((prev) => ({ ...prev, isLoading: false }));
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const persistAuth = async (token: string, user: User) => {
    const payload = JSON.stringify(user);
    if (Platform.OS === 'web') {
      localStorage.setItem('auth_token', token);
      localStorage.setItem('auth_user', payload);
    } else {
      await SecureStore.setItemAsync('auth_token', token);
      await SecureStore.setItemAsync('auth_user', payload);
    }
  };

  const clearAuth = async () => {
    if (Platform.OS === 'web') {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
    } else {
      await SecureStore.deleteItemAsync('auth_token');
      await SecureStore.deleteItemAsync('auth_user');
    }
  };

  const login = async (username: string, password: string) => {
    setState((prev) => ({ ...prev, isLoading: true }));
    try {
      const response = await api.login({ username, password });
      const user: User = {
        id: response.user.id,
        fullName: response.user.fullName,
        nnu: response.user.nnu,
      };
      await persistAuth(response.token, user);
      setState({ user, token: response.token, isAuthenticated: true, isLoading: false });
    } catch (error) {
      setState((prev) => ({ ...prev, isLoading: false }));
      throw error;
    }
  };

  const loginWithBiometrics = async () => {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    if (!hasHardware) {
      throw new Error('Biometric authentication not available on this device');
    }
    const enrolled = await LocalAuthentication.isEnrolledAsync();
    if (!enrolled) {
      throw new Error('No biometric credentials enrolled');
    }
    const result = await LocalAuthentication.authenticateAsync({
      promptMessage: 'Authenticate to SNISID',
      fallbackLabel: 'Use password instead',
    });
    if (!result.success) {
      throw new Error('Biometric authentication failed');
    }
    if (Platform.OS === 'web') {
      const token = localStorage.getItem('auth_token');
      const userRaw = localStorage.getItem('auth_user');
      if (!token || !userRaw) throw new Error('No stored credentials for biometric login');
      setState({
        user: JSON.parse(userRaw),
        token,
        isAuthenticated: true,
        isLoading: false,
      });
      return;
    }
    const token = await SecureStore.getItemAsync('auth_token');
    const userRaw = await SecureStore.getItemAsync('auth_user');
    if (!token || !userRaw) throw new Error('No stored credentials for biometric login');
    setState({
      user: JSON.parse(userRaw),
      token,
      isAuthenticated: true,
      isLoading: false,
    });
  };

  const logout = async () => {
    try {
      await api.logout();
    } catch {
    } finally {
      await clearAuth();
      setState({ user: null, token: null, isAuthenticated: false, isLoading: false });
    }
  };

  return (
    <AuthContext.Provider
      value={{ ...state, login, loginWithBiometrics, logout, checkAuth, api }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return ctx;
}
