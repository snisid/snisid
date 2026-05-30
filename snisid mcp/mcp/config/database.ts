export const databasePolicy = {
  directDbAccessFromMcpTools: false,
  allowedAccessPattern: 'MCP -> API Gateway -> Microservice -> Database',
  encryptionAtRest: 'AES-256-GCM via national KMS/HSM',
  backupPolicy: 'immutable encrypted backups with sovereign retention policy'
} as const;
