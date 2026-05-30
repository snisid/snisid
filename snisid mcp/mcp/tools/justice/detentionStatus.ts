import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { mjspService } from '../../services/mjsp.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerDetentionStatusTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'justice.detentionStatus',
    description: 'Check detention status through MJSP API.',
    permission: PERMISSIONS.JUSTICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
    },
    handler: async (input, ctx) => {
      return mjspService.detentionStatus(input, ctx);
    }
  });
}