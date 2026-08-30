import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class PassportService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'PassportService', baseURL: env.PASSPORT_API_BASE_URL }); }

  async passportLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/passport/lookup', payload, ctx);
  }

  async passportStatus<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/passport/status', payload, ctx);
  }

}

export const passportService = new PassportService();
