import { securityConfig } from '../config/security.js';

interface SessionRecord {
  subject: string;
  deviceId: string;
  createdAt: number;
  lastSeen: number;
}

const sessions = new Map<string, SessionRecord>();

export function bindSession(sessionId: string, subject: string, deviceId: string): void {
  const existing = sessions.get(sessionId);
  const now = Date.now();
  if (existing && (existing.subject !== subject || existing.deviceId !== deviceId)) {
    throw new Error('SESSION_ISOLATION_VIOLATION');
  }
  sessions.set(sessionId, existing ? { ...existing, lastSeen: now } : { subject, deviceId, createdAt: now, lastSeen: now });
}

export function validateSession(sessionId: string | undefined, subject: string, deviceId: string): void {
  if (!sessionId) return;
  const session = sessions.get(sessionId);
  if (!session) return bindSession(sessionId, subject, deviceId);
  if (session.subject !== subject || session.deviceId !== deviceId) throw new Error('SESSION_ISOLATION_VIOLATION');
  if (Date.now() - session.createdAt > securityConfig.maxSessionAgeMs) throw new Error('SESSION_EXPIRED');
  session.lastSeen = Date.now();
}
