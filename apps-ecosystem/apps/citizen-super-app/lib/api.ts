import { createApiClient, type ApiClient } from '@snisid/api-client';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';

const API_BASE_URL =
  process.env.EXPO_PUBLIC_API_URL ?? 'https://api.snisid.gov.za';

let clientInstance: ApiClient | null = null;

export function getApiClient(): ApiClient {
  if (!clientInstance) {
    clientInstance = createApiClient(API_BASE_URL);
    clientInstance.setTokenProvider({
      getToken: async () => {
        if (Platform.OS === 'web') {
          return localStorage.getItem('auth_token');
        }
        return SecureStore.getItemAsync('auth_token');
      },
      setToken: async (token: string | null) => {
        if (Platform.OS === 'web') {
          if (token) {
            localStorage.setItem('auth_token', token);
          } else {
            localStorage.removeItem('auth_token');
          }
          return;
        }
        if (token) {
          await SecureStore.setItemAsync('auth_token', token);
        } else {
          await SecureStore.deleteItemAsync('auth_token');
        }
      },
      onUnauthorized: () => {
        if (Platform.OS === 'web') {
          localStorage.removeItem('auth_token');
        } else {
          SecureStore.deleteItemAsync('auth_token');
        }
      },
    });
  }
  return clientInstance;
}

export {
  type ApiClient,
  type Identity,
  type Document,
  type NotificationItem,
  type AccessLog,
  type ActiveSession,
  type FAQItem,
  type SupportTicket,
} from '@snisid/api-client';
