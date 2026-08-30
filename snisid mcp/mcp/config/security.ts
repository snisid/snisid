export const securityConfig = {
  denyByDefault: true,
  requirePurpose: true,
  requireCorrelationId: true,
  requireDeviceTrust: true,
  requireMfaForSensitivePermissions: true,
  tokenClockToleranceSeconds: 30,
  maxSessionAgeMs: 15 * 60 * 1000,
  maxRiskScore: 70,
  apiKeyRotationDays: 30,
  allowedPromptInstructionSources: ['SNISID_SYSTEM', 'GOV_POLICY', 'AUTHORIZED_OPERATOR']
} as const;
