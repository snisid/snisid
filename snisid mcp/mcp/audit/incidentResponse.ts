import { recordSecurityEvent } from './securityEvents.js';

export async function openIncident(params: {
  title: string;
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  correlationId?: string;
  metadata?: Record<string, unknown>;
}) {
  return recordSecurityEvent({
    action: 'incident.open',
    resource: 'security.incident',
    outcome: 'ERROR',
    severity: params.severity,
    ...(params.correlationId ? { correlationId: params.correlationId } : {}),
    metadata: { title: params.title, ...(params.metadata ?? {}) }
  });
}
