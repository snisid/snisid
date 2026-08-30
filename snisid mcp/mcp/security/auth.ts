import { z } from 'zod';
import type { Permission } from '../config/permissions.js';
import type { SecurityContext, ToolAuthInput } from '../types/security.types.js';
import { verifyAccessToken } from './jwt.js';
import { requirePermission } from './rbac.js';
import { verifyMfaToken } from './mfa.js';
import { assessDeviceTrust } from './deviceTrust.js';
import { validateSession } from './session.js';
import { enforceZeroTrust } from './zerotrust.js';

export const authContextSchema = z.object({
  accessToken: z.string().min(20),
  apiKey: z.string().optional(),
  mfaToken: z.string().optional(),
  deviceId: z.string().min(8),
  purpose: z.string().min(6).max(512),
  correlationId: z.string().min(8).max(128),
  sessionId: z.string().min(8).optional()
});

export async function authenticateAndAuthorize(auth: ToolAuthInput, permission: Permission): Promise<SecurityContext> {
  const parsed = authContextSchema.parse(auth);
  const principal = verifyAccessToken(parsed.accessToken);
  const mfaOk = principal.mfa || verifyMfaToken(parsed.mfaToken, principal.subject, parsed.sessionId ?? principal.sessionId);
  const withMfa = { ...principal, mfa: mfaOk };
  requirePermission(withMfa, permission);
  validateSession(parsed.sessionId ?? principal.sessionId, principal.subject, parsed.deviceId);
  const device = assessDeviceTrust(parsed.deviceId);
  const ctx: SecurityContext = {
    principal: withMfa,
    correlationId: parsed.correlationId,
    purpose: parsed.purpose,
    deviceId: parsed.deviceId,
    riskScore: device.riskScore
  };
  enforceZeroTrust(ctx, permission);
  return ctx;
}
