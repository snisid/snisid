import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { immigrationService } from '../../services/immigration.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerTravelHistoryTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'immigration.travelHistory',
    description: 'Retrieve minimized travel history.',
    permission: PERMISSIONS.IMMIGRATION_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      fromDate: z.string().date().optional(),
      toDate: z.string().date().optional(),
    },
    handler: async (input, ctx) => {
      return immigrationService.travelHistory(input, ctx);
    }
  });
}