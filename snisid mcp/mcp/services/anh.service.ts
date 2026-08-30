import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class AnhService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'AnhService', baseURL: env.ANH_API_BASE_URL }); }

  async archiveLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/archives/lookup', payload, ctx);
  }

  async documentAttestation<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/archives/attestation', payload, ctx);
  }

}

export const anhService = new AnhService();
