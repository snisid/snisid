/**
 * SNISID — Temporal workflow (code-first) — Birth simple (long-running orchestration).
 * Mirrors the BPMN civil-registry.birth.simple.v1.0.0 with stronger compensation guarantees.
 */
import { proxyActivities, defineSignal, setHandler, sleep, condition } from '@temporalio/workflow';
import type * as biometric from '../activities/biometric.js';
import type * as fraud from '../activities/fraud.js';
import type * as identity from '../activities/identity.js';
import type * as audit from '../activities/audit.js';
import type * as notification from '../activities/notification.js';

const acts = proxyActivities<
  typeof biometric & typeof fraud & typeof identity & typeof audit & typeof notification
>({
  startToCloseTimeout: '1 minute',
  retry: {
    maximumAttempts: 5,
    initialInterval: '2s',
    maximumInterval: '1m',
    backoffCoefficient: 2,
    nonRetryableErrorTypes: ['ValidationError', 'FraudCriticalError']
  }
});

export const supervisorApprovalSignal = defineSignal<[boolean]>('supervisorApproval');

export interface BirthInput {
  declaration: {
    firstName: string;
    lastName: string;
    birthDate: string;
    birthCommune: string;
    birthDepartment: string;
    parents: { father?: any; mother?: any };
  };
  declarant: { nin?: string; phone?: string; email?: string };
  agentId: string;
  channel: 'ONLINE' | 'KIOSK' | 'FIELD' | 'MOBILE';
  workflowInstanceId: string;
}

export async function birthSimpleWorkflow(input: BirthInput): Promise<{ nin: string; cardSerial: string }> {
  const wid = 'civil-registry.birth.simple';

  await acts.emitAudit({
    workflowId: wid, workflowInstanceId: input.workflowInstanceId,
    toState: 'STARTED',
    actor: { type: 'USER', id: input.agentId },
    payload: { channel: input.channel }
  });

  // Fraud detection
  const fr = await acts.detect({
    workflowId: wid,
    workflowInstanceId: input.workflowInstanceId,
    data: input.declaration as any
  });

  if (fr.recommended === 'REJECT' || fr.recommended === 'INVESTIGATE') {
    await acts.emitAudit({
      workflowId: wid, workflowInstanceId: input.workflowInstanceId,
      fromState: 'STARTED', toState: 'FRAUD_HALT',
      actor: { type: 'SYSTEM', id: 'fraud-engine' },
      payload: fr as any
    });
    throw new Error(`FraudCriticalError: score=${fr.fraudScore}`);
  }

  // Wait for supervisor approval (4-eyes)
  let approved: boolean | null = null;
  setHandler(supervisorApprovalSignal, (a) => { approved = a; });
  const got = await condition(() => approved !== null, '24 hours');
  if (!got) throw new Error('SLAExceeded: supervisor approval timed out (24h)');
  if (!approved) throw new Error('ValidationError: supervisor rejected');

  await acts.emitAudit({
    workflowId: wid, workflowInstanceId: input.workflowInstanceId,
    fromState: 'STARTED', toState: 'APPROVED',
    actor: { type: 'USER', id: 'supervisor' }, payload: {}
  });

  const { nin } = await acts.generateNin({ birthActId: input.workflowInstanceId });
  const { cardSerial } = await acts.issueCard({ nin });

  await acts.send({
    channels: ['SMS', 'EMAIL', 'INAPP'],
    recipient: { nin: input.declarant.nin, phone: input.declarant.phone, email: input.declarant.email },
    template: 'birth.registered',
    vars: { nin, firstName: input.declaration.firstName }
  });

  await acts.emitAudit({
    workflowId: wid, workflowInstanceId: input.workflowInstanceId,
    fromState: 'APPROVED', toState: 'COMPLETED',
    actor: { type: 'SYSTEM', id: 'engine' }, payload: { nin, cardSerial }
  });

  return { nin, cardSerial };
}
