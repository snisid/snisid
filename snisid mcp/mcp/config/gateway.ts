import { env } from './env.js';

export const gatewayConfig = {
  baseUrl: env.API_GATEWAY_BASE_URL,
  timeoutMs: 10_000,
  retries: 2,
  requiredHeaders: ['authorization', 'x-correlation-id', 'x-purpose', 'x-device-id'],
  mtlsRequired: true,
  wafRequired: true
} as const;
