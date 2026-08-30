import { appendFile, mkdir, readFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { randomUUID } from 'node:crypto';
import { SNISID } from '../config/constants.js';
import type { AuditEvent } from '../types/security.types.js';
import { redact } from '../utils/logger.js';
import { sha256 } from '../utils/crypto.js';

async function lastHash(path: string): Promise<string> {
  try {
    const data = await readFile(path, 'utf8');
    const lines = data.trim().split('\n').filter(Boolean);
    if (!lines.length) return 'GENESIS';
    const last = JSON.parse(lines.at(-1) ?? '{}') as AuditEvent;
    return last.hash ?? 'GENESIS';
  } catch {
    return 'GENESIS';
  }
}

export async function writeAuditEvent(event: Omit<AuditEvent, 'id' | 'timestamp' | 'previousHash' | 'hash'>): Promise<AuditEvent> {
  const path = join(process.cwd(), SNISID.auditLogPath);
  await mkdir(dirname(path), { recursive: true });
  const previousHash = await lastHash(path);
  const base: AuditEvent = {
    id: randomUUID(),
    timestamp: new Date().toISOString(),
    ...event,
    metadata: redact(event.metadata) as Record<string, unknown>,
    previousHash
  };
  const hash = sha256(JSON.stringify(base));
  const complete = { ...base, hash };
  await appendFile(path, `${JSON.stringify(complete)}\n`, { encoding: 'utf8', mode: 0o600 });
  return complete;
}
