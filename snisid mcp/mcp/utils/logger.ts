import winston from 'winston';
import { env } from '../config/env.js';

export const logger = winston.createLogger({
  level: env.LOG_LEVEL,
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: false }),
    winston.format.json()
  ),
  defaultMeta: { service: 'snisid-mcp' },
  transports: [new winston.transports.Console({ stderrLevels: ['error', 'warn', 'info', 'debug'] })]
});

export function redact(value: unknown): unknown {
  if (Array.isArray(value)) return value.map(redact);
  if (value && typeof value === 'object') {
    const out: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(value)) {
      if (/token|secret|password|key|biometric|face|fingerprint|iris|image/i.test(k)) out[k] = '[REDACTED]';
      else out[k] = redact(v);
    }
    return out;
  }
  return value;
}
