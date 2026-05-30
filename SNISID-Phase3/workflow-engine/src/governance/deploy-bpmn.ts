/**
 * SNISID — BPMN deployer with signature verification.
 *  - Verifies every .bpmn under BPMN/ has been signed by the WGO
 *  - Refuses unsigned BPMN deployment
 *  - Pushes to Zeebe with versionTag
 *  - Emits Kafka `governance.bpmn.deployed.v1`
 */
import 'dotenv/config';
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';
import { createHash } from 'node:crypto';
import { Camunda8 } from '@camunda8/sdk';
import { sign, verify } from '../pki/sign.js';
import { kafka } from '../kafka/producer.js';

const ROOT = process.env.BPMN_ROOT ?? 'BPMN';
const SIG_DIR = process.env.BPMN_SIG_DIR ?? '.bpmn-signatures';

function walk(dir: string): string[] {
  const out: string[] = [];
  for (const e of readdirSync(dir)) {
    const p = join(dir, e);
    if (statSync(p).isDirectory()) out.push(...walk(p));
    else if (p.endsWith('.bpmn')) out.push(p);
  }
  return out;
}

async function main() {
  const c8 = new Camunda8();
  const zeebe = c8.getZeebeGrpcApiClient();
  const files = walk(ROOT);
  console.log(`🔍 Found ${files.length} BPMN files.`);

  for (const file of files) {
    const buf = readFileSync(file);
    const hash = createHash('sha384').update(buf).digest('hex');

    // signature lookup
    const sigFile = join(SIG_DIR, relative(ROOT, file) + '.sig');
    let signature: string | null = null;
    try { signature = readFileSync(sigFile, 'utf8'); } catch { /* missing */ }

    if (!signature) {
      console.error(`❌ ${file} — UNSIGNED. Refused by governance.`);
      process.exitCode = 2;
      continue;
    }

    const ok = await verify({ hash }, signature);
    if (!ok) {
      console.error(`❌ ${file} — INVALID SIGNATURE. Refused.`);
      process.exitCode = 2;
      continue;
    }

    const res = await zeebe.deployResource({ name: file, process: buf });
    console.log(`✅ Deployed ${file} → key=${res.deployments?.[0]?.process?.processDefinitionKey}`);

    await kafka.emit({
      topic: 'governance.bpmn.deployed.v1',
      eventType: 'governance.bpmn.deployed.v1',
      correlation: { workflowId: 'governance.bpmn.deploy', workflowInstanceId: hash },
      payload: { file: relative(ROOT, file), hash, deployedAt: Date.now() }
    });
  }
}

main().catch((e) => { console.error(e); process.exit(1); });
