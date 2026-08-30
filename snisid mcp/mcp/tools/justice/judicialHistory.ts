import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { mjspService } from '../../services/mjsp.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerJudicialHistoryTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'justice.judicialHistory',
    description: 'Retrieve judicial history under strict RBAC/MFA.',
    permission: PERMISSIONS.JUSTICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      fromDate: z.string().date().optional(),
      toDate: z.string().date().optional(),
    },
    handler: async (input, ctx) => {
      return mjspService.judicialHistory(input, ctx);
    }
  });
}