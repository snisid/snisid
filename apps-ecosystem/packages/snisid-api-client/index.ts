import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export interface TokenProvider {
  getToken: () => Promise<string | null>;
  setToken: (token: string | null) => void;
  onUnauthorized: () => void;
}

export interface Identity {
  id: string;
  nnu: string;
  fullName: string;
  dateOfBirth: string;
  nationality: string;
  gender: string;
  photoUrl?: string;
  status: 'active' | 'expired' | 'suspended' | 'pending';
  expirationDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface Document {
  id: string;
  type: 'birth_certificate' | 'national_id' | 'passport' | 'drivers_license';
  name: string;
  documentNumber: string;
  issueDate: string;
  expirationDate: string;
  status: 'valid' | 'expired' | 'pending';
  issuer: string;
}

export interface Alert {
  id: string;
  type: string;
  severity: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW';
  title: string;
  description: string;
  identityId: string;
  identityName: string;
  status: 'open' | 'acknowledged' | 'resolved';
  assignedTo?: string;
  createdAt: string;
}

export interface LoginResponse {
  token: string;
  user: { id: string; fullName: string; nnu: string; role?: string };
}

export interface SearchParams {
  q?: string;
  status?: string;
  agency?: string;
  page?: number;
  limit?: number;
}

export interface NotificationItem {
  id: string;
  type: 'identity_update' | 'verification' | 'document_expiry' | 'security_alert' | 'system';
  title: string;
  description: string;
  timestamp: string;
  read: boolean;
}

export interface AccessLog {
  id: string;
  application: string;
  action: string;
  timestamp: string;
  ipAddress: string;
  location?: string;
}

export interface ActiveSession {
  id: string;
  deviceName: string;
  deviceType: string;
  lastActive: string;
  loggedInAt: string;
  current: boolean;
}

export interface FAQItem {
  id: string;
  question: string;
  answer: string;
  category: string;
}

export interface SupportTicket {
  id: string;
  subject: string;
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  priority: 'low' | 'medium' | 'high';
  createdAt: string;
  updatedAt: string;
  lastMessage?: string;
}

export interface EnrollmentData {
  personalInfo: Record<string, unknown>;
  documents: Record<string, unknown>;
  biometrics: Record<string, unknown>;
}

interface RetryConfig {
  maxRetries: number;
  retryDelay: number;
}

export class ApiClient {
  private client: AxiosInstance;
  private tokenProvider: TokenProvider | null = null;
  private retryConfig: RetryConfig = { maxRetries: 2, retryDelay: 1000 };

  constructor(baseUrl: string, config?: Partial<RetryConfig>) {
    if (config) {
      this.retryConfig = { ...this.retryConfig, ...config };
    }

    this.client = axios.create({
      baseURL: baseUrl,
      timeout: 15000,
      headers: { 'Content-Type': 'application/json' },
    });

    this.client.interceptors.request.use(async (config) => {
      const token = await this.tokenProvider?.getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          this.tokenProvider?.onUnauthorized();
        }
        const config = error.config;
        if (!config || config._retryCount >= this.retryConfig.maxRetries) {
          return Promise.reject(error);
        }
        config._retryCount = (config._retryCount || 0) + 1;
        await new Promise((resolve) => setTimeout(resolve, this.retryConfig.retryDelay));
        return this.client.request(config);
      }
    );
  }

  setTokenProvider(provider: TokenProvider) {
    this.tokenProvider = provider;
  }

  private async request<T>(config: AxiosRequestConfig): Promise<T> {
    return this.client.request<T>(config).then((r) => r.data);
  }

  getIdentity() {
    return this.request<{ identity: Identity }>({ url: '/api/v1/identity/me' });
  }

  getDocuments() {
    return this.request<{ documents: Document[] }>({ url: '/api/v1/documents' });
  }

  verifyIdentity(data: { qrData: string }) {
    return this.request<{ verified: boolean; message: string; identity?: Identity }>({
      url: '/api/v1/identity/verify',
      method: 'POST',
      data,
    });
  }

  login(credentials: { username: string; password: string }) {
    return this.request<LoginResponse>({
      url: '/api/v1/auth/login',
      method: 'POST',
      data: credentials,
    });
  }

  logout() {
    return this.request<{ success: boolean }>({
      url: '/api/v1/auth/logout',
      method: 'POST',
    });
  }

  searchIdentities(params?: SearchParams) {
    return this.request<{ identities: Identity[]; total: number; page: number }>({
      url: '/api/v1/admin/identities',
      params,
    });
  }

  getAdminIdentity(id: string) {
    return this.request<{
      identity: Identity;
      documents: Document[];
      biometrics: unknown;
      enrollmentHistory: unknown[];
    }>({ url: `/api/v1/admin/identities/${id}` });
  }

  getStats() {
    return this.request<{
      totalCitizens: number;
      pendingVerifications: number;
      fraudAlerts: number;
      syncStatus: string;
      recentActivity: { id: string; action: string; timestamp: string; details: string }[];
    }>({ url: '/api/v1/admin/stats' });
  }

  getAlerts(params?: { severity?: string; status?: string; page?: number; limit?: number }) {
    return this.request<{ alerts: Alert[]; total: number; page: number }>({
      url: '/api/v1/admin/alerts',
      params,
    });
  }

  enrollCitizen(data: EnrollmentData) {
    return this.request<{ id: string }>({
      url: '/api/v1/admin/enroll',
      method: 'POST',
      data,
    });
  }

  verifyAdminIdentity(id: string) {
    return this.request<{ verified: boolean; message: string }>({
      url: `/api/v1/admin/identities/${id}/verify`,
      method: 'POST',
    });
  }

  suspendIdentity(id: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/admin/identities/${id}/suspend`,
      method: 'POST',
    });
  }

  flagIdentity(id: string, reason: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/admin/identities/${id}/flag`,
      method: 'POST',
      data: { reason },
    });
  }

  acknowledgeAlert(id: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/admin/alerts/${id}/acknowledge`,
      method: 'POST',
    });
  }

  assignAlert(id: string, investigatorId: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/admin/alerts/${id}/assign`,
      method: 'POST',
      data: { investigatorId },
    });
  }

  getNotifications(params?: { page?: number; limit?: number }) {
    return this.request<{ notifications: NotificationItem[]; total: number; page: number }>({
      url: '/api/v1/notifications',
      params,
    });
  }

  markNotificationRead(id: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/notifications/${id}/read`,
      method: 'POST',
    });
  }

  getAccessLogs(params?: { page?: number; limit?: number }) {
    return this.request<{ logs: AccessLog[]; total: number; page: number }>({
      url: '/api/v1/privacy/access-logs',
      params,
    });
  }

  getActiveSessions() {
    return this.request<{ sessions: ActiveSession[] }>({
      url: '/api/v1/auth/sessions',
    });
  }

  terminateSession(id: string) {
    return this.request<{ success: boolean }>({
      url: `/api/v1/auth/sessions/${id}`,
      method: 'DELETE',
    });
  }

  enrollBiometrics(data: { faceImage: string; livenessConfidence: number }) {
    return this.request<{ success: boolean; enrollmentId: string }>({
      url: '/api/v1/biometrics/enroll',
      method: 'POST',
      data,
    });
  }

  updateIdentity(data: Partial<{ fullName: string; dateOfBirth: string; nationality: string; gender: string }>) {
    return this.request<{ identity: Identity }>({
      url: '/api/v1/identity/update',
      method: 'PUT',
      data,
    });
  }

  getSupportFAQ() {
    return this.request<{ faq: FAQItem[] }>({
      url: '/api/v1/support/faq',
    });
  }

  getSupportTickets(params?: { page?: number; limit?: number }) {
    return this.request<{ tickets: SupportTicket[]; total: number; page: number }>({
      url: '/api/v1/support/tickets',
      params,
    });
  }

  createSupportTicket(data: { subject: string; message: string; category: string }) {
    return this.request<{ ticket: SupportTicket }>({
      url: '/api/v1/support/tickets',
      method: 'POST',
      data,
    });
  }
}

export const createApiClient = (baseUrl: string, config?: Partial<RetryConfig>) =>
  new ApiClient(baseUrl, config);
