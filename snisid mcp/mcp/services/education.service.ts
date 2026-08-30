import { env } from '../config/env.js';
import type { SecurityContext } from '../types/security.types.js';
import { GovernmentApiClient } from './baseClient.js';

export class EducationService extends GovernmentApiClient {
  constructor() { super({ serviceName: 'EducationService', baseURL: env.EDUCATION_API_BASE_URL }); }

  async studentVerification<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/education/student/verify', payload, ctx);
  }

  async diplomaVerification<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/education/diploma/verify', payload, ctx);
  }

  async institutionLookup<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/education/institution/lookup', payload, ctx);
  }

  async academicHistory<T = unknown>(payload: unknown, ctx: SecurityContext): Promise<T> {
    return this.securePost<T>('/education/academic-history', payload, ctx);
  }

}

export const educationService = new EducationService();
