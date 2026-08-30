import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { pnhService } from '../../services/pnh.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerGangAffiliationTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'police.gangAffiliation',
    description: 'Retrieve legally authorized gang-affiliation intelligence assessment.',
    permission: PERMISSIONS.POLICE_THREAT,
    inputShape: {
      nationalId: nationalIdSchema,
      evidenceThreshold: z.number().min(0).max(1).default(0.8),
    },
    handler: async (input, ctx) => {
      return pnhService.gangAffiliation(input, ctx);
    }
  });
}