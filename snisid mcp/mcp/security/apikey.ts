import bcrypt from 'bcrypt';
import { createHmac } from 'node:crypto';
import { env } from '../config/env.js';

export function fingerprintApiKey(apiKey: string): string {
  return createHmac('sha256', env.API_KEY_PEPPER).update(apiKey).digest('hex');
}

export async function hashApiKey(apiKey: string): Promise<string> {
  return bcrypt.hash(fingerprintApiKey(apiKey), 12);
}

export async function verifyApiKey(apiKey: string, hash: string): Promise<boolean> {
  return bcrypt.compare(fingerprintApiKey(apiKey), hash);
}

export function rotationDue(createdAt: Date, rotationDays = 30): boolean {
  return Date.now() - createdAt.getTime() > rotationDays * 24 * 60 * 60 * 1000;
}
