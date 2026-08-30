import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class BiometricService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'BiometricService', baseURL: env.BIOMETRIC_API_BASE_URL }); }

  async biometricMatch<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/biometric/match', payload, ctx);
  }

  async faceVerification<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/biometric/face/verify', payload, ctx);
  }

}

export const biometricService = new BiometricService();
