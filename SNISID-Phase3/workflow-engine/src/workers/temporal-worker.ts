/**
 * SNISID — Temporal worker. Hosts code-first workflows.
 */
import 'dotenv/config';
import { Worker, NativeConnection } from '@temporalio/worker';
import { startObservability } from '../observability/otel.js';
import { config } from '../config.js';
import * as biometric from '../activities/biometric.js';
import * as fraud from '../activities/fraud.js';
import * as identity from '../activities/identity.js';
import * as audit from '../activities/audit.js';
import * as notification from '../activities/notification.js';

async function main() {
  startObservability();

  const connection = await NativeConnection.connect({
    address: config.TEMPORAL_ADDRESS,
    tls: {}
  });

  const worker = await Worker.create({
    connection,
    namespace: config.TEMPORAL_NAMESPACE,
    taskQueue: config.TEMPORAL_TASK_QUEUE,
    workflowsPath: new URL('../workflows/', import.meta.url).pathname,
    activities: { ...biometric, ...fraud, ...identity, ...audit, ...notification }
  });

  console.log(`✅ Temporal worker started on queue ${config.TEMPORAL_TASK_QUEUE}`);
  await worker.run();
}

main().catch((e) => { console.error(e); process.exit(1); });
