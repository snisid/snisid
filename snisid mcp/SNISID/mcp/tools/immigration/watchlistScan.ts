import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { immigrationService } from '../../services/immigration.service.js';
import { nationalIdSchema, passportNumberSchema } from '../../validators/identity.validator.js';

export function registerWatchlistScanTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'immigration.watchlistScan',
    description: 'Scan against authorized immigration watchlists.',
    permission: PERMISSIONS.IMMIGRATION_READ,
    inputShape: {
      nationalId: nationalIdSchema.optional(),
      passportNumber: passportNumberSchema.optional(),
      name: z.string().max(128).optional(),
    },
    handler: async (input, ctx) => {
      return immigrationService.watchlistScan(input, ctx);
    }
  });
}