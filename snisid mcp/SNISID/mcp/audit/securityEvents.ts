import { writeAuditEvent } from './auditLogger.js';

export async function recordSecurityEvent(input: {
  actor?: string;
  action: string;
  resource: string;
  purpose?: string;
  correlationId?: string;
  outcome: 'ALLOW' | 'DENY' | 'ERROR';
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  metadata?: Record<string, unknown>;
}) {
  return writeAuditEvent({
    actor: input.actor ?? 'anonymous',
    action: input.action,
    resource: input.resource,
    purpose: input.purpose ?? 'UNSPECIFIED',
    correlationId: input.correlationId ?? 'UNSPECIFIED',
    outcome: input.outcome,
    severity: input.severity,
    metadata: input.metadata
  });
}
