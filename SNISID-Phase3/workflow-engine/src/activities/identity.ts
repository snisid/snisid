/**
 * SNISID — Identity activities (NIN generation, card issuance, suspend, revoke).
 */
import axios from 'axios';
import { kafka } from '../kafka/producer.js';

const IDSVC = process.env.IDENTITY_GATEWAY ?? 'https://identity.snisid.ht';

export async function generateNin(input: { birthActId: string }): Promise<{ nin: string }> {
  const r = await axios.post(`${IDSVC}/nin/generate`, input, { timeout: 5_000 });
  return r.data;
}

export async function issueCard(input: { nin: string }): Promise<{ cardSerial: string }> {
  const r = await axios.post(`${IDSVC}/card/issue`, input, { timeout: 30_000 });
  return r.data;
}

export async function blockCard(input: { nin: string; reason: string }) {
  await axios.post(`${IDSVC}/card/block`, input, { timeout: 5_000 });
}

export async function suspendIdentity(input: { nin: string; reason: string; orderRef?: string }) {
  await axios.post(`${IDSVC}/suspend`, input, { timeout: 5_000 });
  await kafka.emit({
    topic: 'identity.suspended.v1',
    eventType: 'identity.suspended.v1',
    correlation: { workflowId: 'identity.suspension.judicial', workflowInstanceId: input.orderRef ?? 'n/a' },
    subjectKey: input.nin,
    payload: input
  });
}

export async function revokeIdentity(input: { nin: string; reason: string; legalRef?: string }) {
  await axios.post(`${IDSVC}/revoke`, input, { timeout: 5_000 });
  await kafka.emit({
    topic: 'identity.revoked.v1',
    eventType: 'identity.revoked.v1',
    correlation: { workflowId: 'identity.revocation.administrative', workflowInstanceId: input.legalRef ?? 'n/a' },
    subjectKey: input.nin,
    payload: input
  });
}
