import { writeAuditEvent } from './auditLogger.js';
import type { SecurityContext } from '../types/security.types.js';

export async function trackActivity(ctx: SecurityContext, action: string, resource: string, metadata?: Record<string, unknown>) {
  return writeAuditEvent({
    actor: ctx.principal.subject,
    action,
    resource,
    purpose: ctx.purpose,
    correlationId: ctx.correlationId,
    outcome: 'ALLOW',
    severity: 'LOW',
    ...(metadata ? { metadata } : {})
  });
}
