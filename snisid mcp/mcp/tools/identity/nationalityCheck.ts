import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { oniService } from '../../services/oni.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerNationalityCheckTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.nationalityCheck',
    description: 'Check nationality status via ONI API.',
    permission: PERMISSIONS.IDENTITY_VERIFY,
    inputShape: {
      nationalId: nationalIdSchema,
    },
    handler: async (input, ctx) => {
      return oniService.nationalityCheck(input, ctx);
    }
  });
}