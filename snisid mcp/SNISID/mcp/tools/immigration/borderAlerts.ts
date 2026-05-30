import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { immigrationService } from '../../services/immigration.service.js';
import { nationalIdSchema, passportNumberSchema } from '../../validators/identity.validator.js';

export function registerBorderAlertsTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'immigration.borderAlerts',
    description: 'Query border alerts.',
    permission: PERMISSIONS.IMMIGRATION_READ,
    inputShape: {
      nationalId: nationalIdSchema.optional(),
      passportNumber: passportNumberSchema.optional(),
    },
    handler: async (input, ctx) => {
      return immigrationService.borderAlerts(input, ctx);
    }
  });
}