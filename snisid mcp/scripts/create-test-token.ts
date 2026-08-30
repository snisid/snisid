import 'dotenv/config';

import { signServiceToken, signMfaToken } from '../mcp/security/jwt.js';
import { ROLES } from '../mcp/config/roles.js';
import { PERMISSIONS, type Permission } from '../mcp/config/permissions.js';
import type { AuthenticatedPrincipal } from '../mcp/types/security.types.js';

const sessionId = process.env['SNISID_DEV_SESSION_ID'] ?? 'sess-test-0001';

const principal: AuthenticatedPrincipal = {
  subject: process.env['SNISID_DEV_SUBJECT'] ?? 'dev-admin-001',
  ministry: process.env['SNISID_DEV_MINISTRY'] ?? 'SNISID',
  roles: [ROLES.SNISID_ADMIN],
  permissions: Object.values(PERMISSIONS) as Permission[],
  clearance: 'TOP_SECRET',
  mfa: true,
  sessionId
};

const accessToken = signServiceToken(principal);
const mfaToken = signMfaToken(principal.subject, sessionId);

console.log(
  JSON.stringify(
    {
      accessToken,
      mfaToken,
      sessionId,
      deviceId: 'device-test-001',
      purpose: 'verification legale de test',
      correlationId: 'corr-test-0001'
    },
    null,
    2
  )
);
