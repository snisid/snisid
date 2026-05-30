import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class ImmigrationService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'ImmigrationService', baseURL: env.IMMIGRATION_API_BASE_URL }); }

  async borderAlerts<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/immigration/border-alerts', payload, ctx);
  }

  async travelHistory<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/immigration/travel-history', payload, ctx);
  }

  async visaLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/immigration/visa', payload, ctx);
  }

  async entryExit<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/immigration/entry-exit', payload, ctx);
  }

  async watchlistScan<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/immigration/watchlist-scan', payload, ctx);
  }

}

export const immigrationService = new ImmigrationService();
