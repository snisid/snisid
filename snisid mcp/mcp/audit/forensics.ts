import { readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { SNISID } from '../config/constants.js';
import { sha256 } from '../utils/crypto.js';
import type { AuditEvent } from '../types/security.types.js';

export async function verifyAuditChain(): Promise<{ valid: boolean; brokenAt?: string }> {
  const path = join(process.cwd(), SNISID.auditLogPath);
  const data = await readFile(path, 'utf8').catch(() => '');
  let previous = 'GENESIS';
  for (const line of data.split('\n').filter(Boolean)) {
    const event = JSON.parse(line) as AuditEvent;
    const { hash, ...withoutHash } = event;
    if (event.previousHash !== previous) return { valid: false, brokenAt: event.id };
    if (sha256(JSON.stringify(withoutHash)) !== hash) return { valid: false, brokenAt: event.id };
    previous = hash ?? previous;
  }
  return { valid: true };
}
