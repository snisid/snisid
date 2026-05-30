import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class IntelligenceService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'IntelligenceService', baseURL: env.INTELLIGENCE_API_BASE_URL }); }

  async fusionAnalysis<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/intelligence/fusion-analysis', payload, ctx);
  }

  async riskScore<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/intelligence/risk-score', payload, ctx);
  }

  async networkAnalysis<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/intelligence/network-analysis', payload, ctx);
  }

  async threatDetection<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/intelligence/threat-detection', payload, ctx);
  }

  async behaviorAnalysis<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/intelligence/behavior-analysis', payload, ctx);
  }

}

export const intelligenceService = new IntelligenceService();
