import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { oniService } from '../../services/oni.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerIdentityRiskTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.identityRisk',
    description: 'Compute identity risk from authoritative APIs.',
    permission: PERMISSIONS.IDENTITY_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      signals: z.array(z.string().max(64)).max(20).optional(),
    },
    handler: async (input, ctx) => {
      return oniService.identityRisk(input, ctx);
    }
  });
}