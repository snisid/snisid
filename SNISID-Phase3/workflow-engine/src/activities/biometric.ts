/**
 * SNISID — Biometric activities (capture, quality, 1:1 match, ABIS dedup).
 * Calls Phase-1 services via gRPC/HTTP.
 */
import axios from 'axios';
import pino from 'pino';
const log = pino({ name: 'biometric' });

const BIO_GW = process.env.BIO_GATEWAY ?? 'https://biometric.snisid.ht';

export async function capture(input: { sessionId: string; modalities: string[] }) {
  const r = await axios.post(`${BIO_GW}/capture`, input, { timeout: 30_000 });
  log.info({ sessionId: input.sessionId }, 'biometric captured');
  return r.data;
}

export async function qualityCheck(refId: string) {
  const r = await axios.get(`${BIO_GW}/quality/${refId}`, { timeout: 5_000 });
  return r.data as { nfiq2: number; faceQuality: number; irisQuality: number; pass: boolean };
}

export async function dedup(refId: string) {
  const r = await axios.post(`${BIO_GW}/abis/dedup`, { refId }, { timeout: 60_000 });
  return r.data as { duplicateFound: boolean; topMatches: Array<{ nin: string; score: number }> };
}

export async function match1to1(refId: string, nin: string) {
  const r = await axios.post(`${BIO_GW}/match/1to1`, { refId, nin }, { timeout: 10_000 });
  return r.data as { match: boolean; score: number };
}

export async function liveness(refId: string) {
  const r = await axios.post(`${BIO_GW}/liveness`, { refId }, { timeout: 5_000 });
  return r.data as { live: boolean; spoofScore: number };
}
