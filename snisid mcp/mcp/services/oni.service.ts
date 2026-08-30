import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class OniService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'OniService', baseURL: env.ONI_API_BASE_URL }); }

  async verifyIdentity<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/identity/verify', payload, ctx);
  }

  async citizenProfile<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/identity/profile', payload, ctx);
  }

  async birthCertificate<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/civil/birth-certificate', payload, ctx);
  }

  async nationalityCheck<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/identity/nationality', payload, ctx);
  }

  async identityRisk<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/identity/risk', payload, ctx);
  }

}

export const oniService = new OniService();
