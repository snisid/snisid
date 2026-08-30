import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class MjspService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'MjspService', baseURL: env.MJSP_API_BASE_URL }); }

  async criminalRecord<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/justice/criminal-record', payload, ctx);
  }

  async warrantLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/justice/warrants', payload, ctx);
  }

  async courtCases<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/justice/court-cases', payload, ctx);
  }

  async detentionStatus<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/justice/detention-status', payload, ctx);
  }

  async judicialHistory<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/justice/history', payload, ctx);
  }

}

export const mjspService = new MjspService();
