import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { oniService } from '../../services/oni.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerVerifyIdentityTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.verifyIdentity',
    description: 'Verify a citizen identity via ONI API with purpose limitation.',
    permission: PERMISSIONS.IDENTITY_VERIFY,
    inputShape: {
      nationalId: nationalIdSchema,
      consentReference: z.string().min(4).max(128).optional(),
    },
    handler: async (input, ctx) => {
      return oniService.verifyIdentity(input, ctx);
    }
  });
}