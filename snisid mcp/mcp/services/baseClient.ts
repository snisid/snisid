import axios, { type AxiosInstance } from 'axios';
import { randomUUID } from 'node:crypto';
import { env } from '../config/env.js';
import { SNISID } from '../config/constants.js';
import type { SecurityContext } from '../types/security.types.js';
import { logger, redact } from '../utils/logger.js';
import { sleep } from '../utils/helpers.js';

export interface ApiClientOptions {
  serviceName: string;
  baseURL: string;
  timeoutMs?: number;
  retries?: number;
}

export class GovernmentApiClient {
  protected readonly client: AxiosInstance;
  protected readonly retries: number;

  constructor(private readonly options: ApiClientOptions) {
    this.retries = options.retries ?? SNISID.retryCount;
    this.client = axios.create({
      baseURL: options.baseURL,
      timeout: options.timeoutMs ?? SNISID.requestTimeoutMs,
      headers: {
        'x-snisid-service': options.serviceName,
        'x-api-key': env.GOV_API_KEY
      },
      maxBodyLength: SNISID.maxPayloadBytes,
      maxContentLength: SNISID.maxPayloadBytes,
      validateStatus: (status) => status >= 200 && status < 500
    });
  }

  protected async securePost<T>(path: string, payload: unknown, ctx: SecurityContext): Promise<T> {
    const requestId = randomUUID();
    const headers = {
      'x-request-id': requestId,
      'x-correlation-id': ctx.correlationId,
      'x-purpose': ctx.purpose,
      'x-device-id': ctx.deviceId,
      'x-actor-subject': ctx.principal.subject,
      'x-actor-ministry': ctx.principal.ministry
    };

    let lastError: unknown;
    for (let attempt = 0; attempt <= this.retries; attempt++) {
      try {
        logger.info('government_api_request', { service: this.options.serviceName, path, requestId, attempt, payload: redact(payload) });
        const response = await this.client.post<T>(path, payload, { headers });
        if (response.status >= 400) throw new Error(`Government API ${this.options.serviceName} returned ${response.status}`);
        return response.data;
      } catch (error) {
        lastError = error;
        if (attempt < this.retries) await sleep(100 * (attempt + 1));
      }
    }
    logger.error('government_api_failure', { service: this.options.serviceName, path, requestId, error: String(lastError) });
    throw lastError instanceof Error ? lastError : new Error('GOVERNMENT_API_FAILURE');
  }
}
