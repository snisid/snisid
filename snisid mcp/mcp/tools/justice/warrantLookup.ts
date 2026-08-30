import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { mjspService } from '../../services/mjsp.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerWarrantLookupTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'justice.warrantLookup',
    description: 'Lookup active warrants via MJSP API.',
    permission: PERMISSIONS.JUSTICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      warrantId: z.string().min(4).max(64).optional(),
    },
    handler: async (input, ctx) => {
      return mjspService.warrantLookup(input, ctx);
    }
  });
}