import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class DgiService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'DgiService', baseURL: env.DGI_API_BASE_URL }); }

  async verifyNif<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/tax/nif/verify', payload, ctx);
  }

  async taxCompliance<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/tax/compliance', payload, ctx);
  }

  async businessRegistry<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/tax/business-registry', payload, ctx);
  }

  async financialRisk<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/tax/financial-risk', payload, ctx);
  }

}

export const dgiService = new DgiService();
