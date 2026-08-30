import { describe, expect, it, beforeAll } from 'vitest';

beforeAll(() => {
  process.env['JWT_SECRET'] = 'x'.repeat(64);
  process.env['ENCRYPTION_KEY_B64'] = Buffer.alloc(32, 1).toString('base64');
});

describe('security baseline', () => {
  it('redacts sensitive keys', async () => {
    const { redact } = await import('../utils/logger.js');
    expect(redact({ token: 'abc', nested: { apiKey: 'def', ok: true } })).toEqual({ token: '[REDACTED]', nested: { apiKey: '[REDACTED]', ok: true } });
  });
});
