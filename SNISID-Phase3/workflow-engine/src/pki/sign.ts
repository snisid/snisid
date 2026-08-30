/**
 * SNISID — PKI signing client.
 * Calls the National PKI signing service and an RFC 3161 TSA.
 * Returns base64-encoded detached signature + timestamp token.
 */
import axios from 'axios';
import { createHash } from 'node:crypto';
import { config } from '../config.js';

export interface SignResult {
  signature: string; // base64 PKCS#7 detached signature
  tsa: string;       // base64 RFC 3161 timestamp token
  hash: string;      // sha384 hex
}

function canonicalize(obj: unknown): string {
  // Stable JSON: sorted keys at every level
  const sorter = (x: any): any => {
    if (Array.isArray(x)) return x.map(sorter);
    if (x && typeof x === 'object') {
      return Object.keys(x).sort().reduce((acc: any, k) => {
        acc[k] = sorter(x[k]);
        return acc;
      }, {});
    }
    return x;
  };
  return JSON.stringify(sorter(obj));
}

export async function sign(payload: unknown): Promise<SignResult> {
  const canonical = canonicalize(payload);
  const hash = createHash('sha384').update(canonical).digest('hex');

  // 1. PKI sign
  const sigRes = await axios.post(config.PKI_SIGN_ENDPOINT, {
    algo: 'rsa-pss-sha384',
    hash,
    purpose: 'workflow-event'
  }, { timeout: 5000 });
  const signature = sigRes.data.signature as string;

  // 2. RFC 3161 timestamp
  const tsaRes = await axios.post(config.PKI_TSA_ENDPOINT, {
    hash,
    algo: 'sha384'
  }, { timeout: 5000 });
  const tsa = tsaRes.data.token as string;

  return { signature, tsa, hash };
}

export async function verify(payload: unknown, signature: string): Promise<boolean> {
  const canonical = canonicalize(payload);
  const hash = createHash('sha384').update(canonical).digest('hex');
  const res = await axios.post(`${config.PKI_SIGN_ENDPOINT}/verify`, {
    hash, signature
  }, { timeout: 5000 });
  return res.data.valid === true;
}
