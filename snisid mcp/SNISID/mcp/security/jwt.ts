import jwt from 'jsonwebtoken';
import { env } from '../config/env.js';
import type { AuthenticatedPrincipal } from '../types/security.types.js';

export interface SnisidJwtClaims extends AuthenticatedPrincipal {
  iss: string;
  aud: string | string[];
  sub: string;
  iat?: number;
  exp?: number;
}

export function verifyAccessToken(token: string): AuthenticatedPrincipal {
  const decoded = jwt.verify(token, env.JWT_SECRET, {
    issuer: env.JWT_ISSUER,
    audience: env.JWT_AUDIENCE,
    clockTolerance: 30
  }) as SnisidJwtClaims;

  return {
    subject: decoded.sub ?? decoded.subject,
    ministry: decoded.ministry,
    roles: decoded.roles,
    permissions: decoded.permissions ?? [],
    clearance: decoded.clearance,
    mfa: Boolean(decoded.mfa),
    sessionId: decoded.sessionId
  };
}

export function signServiceToken(principal: AuthenticatedPrincipal): string {
  return jwt.sign(
    { ...principal, sub: principal.subject },
    env.JWT_SECRET,
    { issuer: env.JWT_ISSUER, audience: env.JWT_AUDIENCE, expiresIn: env.JWT_EXPIRES_IN }
  );
}

export function signMfaToken(subject: string, sessionId?: string): string {
  return jwt.sign({ sub: subject, typ: 'mfa', sessionId }, env.JWT_SECRET, {
    issuer: env.JWT_ISSUER,
    audience: env.JWT_AUDIENCE,
    expiresIn: '5m'
  });
}
