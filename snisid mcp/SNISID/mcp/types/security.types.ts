import type { Permission } from '../config/permissions.js';
import type { Role } from '../config/roles.js';

export interface AuthenticatedPrincipal {
  subject: string;
  ministry: string;
  roles: Role[];
  permissions: Permission[];
  clearance: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET';
  mfa: boolean;
  sessionId?: string;
}

export interface SecurityContext {
  principal: AuthenticatedPrincipal;
  correlationId: string;
  purpose: string;
  deviceId: string;
  sourceIp?: string;
  userAgent?: string;
  riskScore: number;
}

export interface ToolAuthInput {
  accessToken: string;
  apiKey?: string;
  mfaToken?: string;
  deviceId: string;
  purpose: string;
  correlationId: string;
  sessionId?: string;
}

export interface AuditEvent {
  id: string;
  timestamp: string;
  actor: string;
  action: string;
  resource: string;
  purpose: string;
  correlationId: string;
  outcome: 'ALLOW' | 'DENY' | 'ERROR';
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  metadata?: Record<string, unknown>;
  previousHash?: string;
  hash?: string;
}
