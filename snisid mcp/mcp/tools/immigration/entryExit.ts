import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { immigrationService } from '../../services/immigration.service.js';
import { passportNumberSchema } from '../../validators/identity.validator.js';

export function registerEntryExitTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'immigration.entryExit',
    description: 'Query entry/exit records.',
    permission: PERMISSIONS.IMMIGRATION_READ,
    inputShape: {
      passportNumber: passportNumberSchema,
      fromDate: z.string().date().optional(),
      toDate: z.string().date().optional(),
    },
    handler: async (input, ctx) => {
      return immigrationService.entryExit(input, ctx);
    }
  });
}