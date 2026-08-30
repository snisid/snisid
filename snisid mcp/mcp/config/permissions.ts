export const PERMISSIONS = {
  IDENTITY_VERIFY: 'identity:verify',
  IDENTITY_READ: 'identity:read',
  IDENTITY_BIOMETRIC: 'identity:biometric',
  JUSTICE_READ: 'justice:read',
  POLICE_READ: 'police:read',
  POLICE_THREAT: 'police:threat',
  IMMIGRATION_READ: 'immigration:read',
  EDUCATION_READ: 'education:read',
  TAX_READ: 'tax:read',
  TAX_RISK: 'tax:risk',
  INTELLIGENCE_READ: 'intelligence:read',
  INTELLIGENCE_ANALYZE: 'intelligence:analyze',
  AUDIT_READ: 'audit:read',
  SECURITY_ADMIN: 'security:admin',
  AI_ORCHESTRATE: 'ai:orchestrate'
} as const;

export type Permission = (typeof PERMISSIONS)[keyof typeof PERMISSIONS];

export const SENSITIVE_PERMISSIONS: Permission[] = [
  PERMISSIONS.IDENTITY_BIOMETRIC,
  PERMISSIONS.JUSTICE_READ,
  PERMISSIONS.POLICE_READ,
  PERMISSIONS.POLICE_THREAT,
  PERMISSIONS.TAX_RISK,
  PERMISSIONS.INTELLIGENCE_READ,
  PERMISSIONS.INTELLIGENCE_ANALYZE,
  PERMISSIONS.SECURITY_ADMIN
];
