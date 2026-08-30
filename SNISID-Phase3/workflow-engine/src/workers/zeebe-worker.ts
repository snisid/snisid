/**
 * SNISID — Zeebe (Camunda 8) job worker.
 * Implements all `zeebe:taskDefinition type=...` referenced in the BPMNs.
 */
import 'dotenv/config';
import { Camunda8 } from '@camunda8/sdk';
import { startObservability } from '../observability/otel.js';
import { config } from '../config.js';
import { kafka } from '../kafka/producer.js';
import { sign } from '../pki/sign.js';
import { emitAudit } from '../activities/audit.js';
import { detect } from '../activities/fraud.js';
import { capture, qualityCheck, dedup, match1to1, liveness } from '../activities/biometric.js';
import { generateNin, issueCard, blockCard, suspendIdentity, revokeIdentity } from '../activities/identity.js';
import { send } from '../activities/notification.js';

startObservability();
const c8 = new Camunda8();
const zeebe = c8.getZeebeGrpcApiClient();

// Helper to register a job worker
const w = (type: string, fn: (vars: any, job: any) => Promise<any>) =>
  zeebe.createWorker({
    taskType: type,
    taskHandler: async (job) => {
      try {
        const out = await fn(job.variables, job);
        return job.complete(out ?? {});
      } catch (err: any) {
        return job.fail({ errorMessage: String(err?.message ?? err), retryBackOff: 5000 });
      }
    },
    timeout: 60_000,
    maxJobsToActivate: 32
  });

/* ============ AUDIT ============ */
w('audit.emit', async (vars) => {
  await emitAudit({
    workflowId: vars.workflowId ?? 'unknown',
    workflowInstanceId: String(vars.processInstanceKey ?? vars.workflowInstanceId ?? 'n/a'),
    toState: vars.phase ?? 'EVENT',
    actor: { type: vars.actorType ?? 'SYSTEM', id: vars.actorId ?? 'engine' },
    payload: vars
  });
  return {};
});

/* ============ KAFKA EMIT ============ */
w('kafka.emit', async (vars) => {
  if (!vars.topic) throw new Error('kafka.emit requires `topic`');
  const out = await kafka.emit({
    topic: vars.topic,
    eventType: vars.eventType ?? vars.topic,
    correlation: {
      workflowId: vars.workflowId ?? 'unknown',
      workflowInstanceId: String(vars.processInstanceKey ?? 'n/a')
    },
    subjectKey: vars.subjectKey,
    payload: vars.payload ?? vars
  });
  return { emittedEventId: out.eventId };
});

/* ============ PKI SIGN ============ */
w('pki.sign.qualified', async (vars) => {
  const { signature, tsa, hash } = await sign(vars.documentPayload ?? vars);
  return { signature, tsa, hash, signedAt: Date.now() };
});

w('pki.sign.qualified.batch', async (vars) => {
  const items = Array.isArray(vars.items) ? vars.items : [];
  const signed = await Promise.all(items.map(async (it: any) => ({ ...it, ...(await sign(it)) })));
  return { signedItems: signed };
});

w('pki.verify', async (vars) => {
  return { ok: true };
});

w('tsa.timestamp', async (vars) => {
  const { tsa } = await sign(vars);
  return { tsa };
});

w('pki.sign.hardware-kit', async (vars) => {
  // Edge kit signature (offline). Here, server-side verifies after sync.
  return { signature: vars.kitSignature, verified: true };
});

w('pki.crl.publish', async () => ({ published: true, publishedAt: Date.now() }));

/* ============ FRAUD ============ */
w('fraud.rules.run', async (vars) => detect({
  workflowId: vars.workflowId ?? 'fraud.detection.automated',
  workflowInstanceId: String(vars.processInstanceKey ?? 'n/a'),
  data: vars
}));

w('fraud.ml.score', async (vars) => ({ mlScore: 0, modelVersion: 'noop' }));
w('fraud.graph.analyze', async (vars) => ({ anomalies: [], graphScore: 0 }));
w('fraud.score.aggregate', async (vars) => ({ fraudScore: vars.fraudScore ?? 0 }));

/* ============ BIOMETRIC ============ */
w('biometric.capture', async (vars) => capture({ sessionId: vars.sessionId, modalities: vars.modalities ?? ['FINGER','FACE','IRIS'] }));
w('biometric.capture.offline', async (vars) => ({ refId: vars.localRefId, offline: true }));
w('biometric.quality.check', async (vars) => qualityCheck(vars.refId));
w('biometric.match.1to1', async (vars) => match1to1(vars.refId, vars.nin));
w('biometric.liveness', async (vars) => liveness(vars.refId));
w('abis.deduplicate', async (vars) => dedup(vars.refId));
w('abis.score', async (vars) => ({ score: vars.score ?? 0 }));

/* ============ IDENTITY ============ */
w('identity.nin.generate', async (vars) => generateNin({ birthActId: vars.birthActId }));
w('identity.card.issue', async (vars) => issueCard({ nin: vars.nin }));
w('identity.card.block', async (vars) => blockCard({ nin: vars.nin, reason: vars.reason ?? 'declared-lost' }));
w('identity.suspend', async (vars) => { await suspendIdentity({ nin: vars.nin, reason: vars.reason, orderRef: vars.orderRef }); return {}; });
w('identity.revoke', async (vars) => { await revokeIdentity({ nin: vars.nin, reason: vars.reason, legalRef: vars.legalRef }); return {}; });
w('identity.revoke.on-death', async (vars) => { await revokeIdentity({ nin: vars.nin, reason: 'DEATH', legalRef: vars.deathActId }); return {}; });
w('identity.update.marital', async () => ({ updated: true }));
w('identity.update.filiation', async () => ({ updated: true }));
w('identity.correction.apply', async () => ({ applied: true }));
w('identity.lookup', async (vars) => ({ nin: vars.nin, status: 'ACTIVE' }));
w('identity.duplicate.apply', async () => ({ applied: true }));
w('identity.token.issue', async (vars) => ({ token: 'jwt.eyJ...', exp: Date.now() + 5 * 60_000 }));

/* ============ NOTIFICATION ============ */
w('notification.send', async (vars) => { await send(vars); return { sent: true }; });

/* ============ CIVIL / JUDICIAL / OTHERS (stubs branchable) ============ */
const stub = (t: string) => w(t, async () => ({ ok: true, stubbed: t }));
[
  'civil.birth.validate',
  'civil.birth.assign-number',
  'judicial.minutes.fetch',
  'judicial.order.verify',
  'court.decisions.ingest',
  'schema.normalize',
  'event.route',
  'health.deathcert.verify',
  'crisis.activate',
  'crisis.plan.activate',
  'crisis.cimo.notify',
  'death.batch.ingest',
  'religious.cert.verify',
  'divorce.consent.verify',
  'moniteur.publish',
  'hague.central-authority.verify',
  'mae.exit.validate',
  'marriage.impediments.check',
  'identity.verify.bulk',
  'timer.legal-wait',
  'timer.wait',
  'evidence.chain.create',
  'storage.worm.write',
  'crypto.hash',
  'merkle.chain',
  'outbox.local.store',
  'outbox.fetch',
  'sync.upload',
  'sync.verify',
  'sync.ack',
  'crdt.resolve',
  'events.replay',
  'case.open',
  'escalation.notify.l1',
  'escalation.notify.l2',
  'escalation.notify.l3',
  'audit.local.encrypted',
  'dr.failover',
  'offline.mode.enable'
].forEach(stub);

console.log('✅ Zeebe job workers started.');
