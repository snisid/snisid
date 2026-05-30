/**
 * SNISID — Audit activity.
 * Emits an immutable, signed audit event for every workflow transition.
 * Used by Temporal activities AND Zeebe service tasks.
 */
import { kafka } from '../kafka/producer.js';
import pino from 'pino';
const log = pino({ name: 'audit' });

export interface AuditInput {
  workflowId: string;
  workflowVersion?: string;
  workflowInstanceId: string;
  fromState?: string;
  toState: string;
  actor: {
    type: 'USER' | 'SYSTEM' | 'SCHEDULER';
    id: string;
    groups?: string[];
  };
  payload: Record<string, unknown>;
}

/** Append-only, signed, Merkle-chained audit entry */
export async function emitAudit(input: AuditInput) {
  const evt = await kafka.emit({
    topic: 'audit.workflow.transition.v1',
    eventType: 'audit.workflow.transition.v1',
    correlation: {
      workflowId: input.workflowId,
      workflowInstanceId: input.workflowInstanceId
    },
    subjectKey: input.workflowInstanceId,
    payload: {
      workflowVersion: input.workflowVersion ?? 'unknown',
      fromState: input.fromState ?? null,
      toState: input.toState,
      actor: input.actor,
      payload: input.payload
    }
  });
  log.info({ eventId: evt.eventId, state: input.toState }, 'audit emitted');
  return evt;
}
