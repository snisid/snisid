import type { Permission } from '../config/permissions.js';
import { SENSITIVE_PERMISSIONS } from '../config/permissions.js';
import { securityConfig } from '../config/security.js';
import type { SecurityContext } from '../types/security.types.js';

export function enforceZeroTrust(ctx: SecurityContext, permission: Permission): void {
  if (securityConfig.requirePurpose && ctx.purpose.trim().length < 6) throw new Error('PURPOSE_REQUIRED');
  if (securityConfig.requireCorrelationId && ctx.correlationId.trim().length < 8) throw new Error('CORRELATION_ID_REQUIRED');
  if (ctx.riskScore > securityConfig.maxRiskScore) throw new Error('RISK_SCORE_TOO_HIGH');
  if (SENSITIVE_PERMISSIONS.includes(permission) && !ctx.principal.mfa) throw new Error('MFA_REQUIRED');
}
