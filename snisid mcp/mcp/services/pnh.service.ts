import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class PnhService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'PnhService', baseURL: env.PNH_API_BASE_URL }); }

  async wantedPerson<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/police/wanted-person', payload, ctx);
  }

  async incidentLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/police/incidents', payload, ctx);
  }

  async gangAffiliation<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/police/gang-affiliation', payload, ctx);
  }

  async weaponPermit<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/police/weapon-permit', payload, ctx);
  }

  async threatMonitoring<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/police/threat-monitoring', payload, ctx);
  }

}

export const pnhService = new PnhService();
