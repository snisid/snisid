export const SNISID = {
  systemName: 'Système National d’Identification et de Sécurité Intelligente',
  shortName: 'SNISID',
  country: 'HT',
  version: '1.0.0',
  mcpName: 'snisid-sovereign-mcp',
  classificationDefault: 'CONFIDENTIAL',
  maxPayloadBytes: 1_000_000,
  requestTimeoutMs: 10_000,
  retryCount: 2,
  auditLogPath: 'mcp/logs/audit.log',
  securityLogPath: 'mcp/logs/security.log'
} as const;

export const DATA_CLASSIFICATIONS = ['PUBLIC', 'INTERNAL', 'CONFIDENTIAL', 'SECRET', 'TOP_SECRET'] as const;
export type DataClassification = (typeof DATA_CLASSIFICATIONS)[number];
