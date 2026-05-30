import jwt from 'jsonwebtoken';
import { env } from '../config/env.js';

export function verifyMfaToken(token: string | undefined, subject: string, sessionId?: string): boolean {
  if (!token) return false;
  try {
    const decoded = jwt.verify(token, env.JWT_SECRET, {
      issuer: env.JWT_ISSUER,
      audience: env.JWT_AUDIENCE,
      clockTolerance: 30
    }) as { sub?: string; typ?: string; sessionId?: string };
    return decoded.typ === 'mfa' && decoded.sub === subject && (!sessionId || decoded.sessionId === sessionId);
  } catch {
    return false;
  }
}
