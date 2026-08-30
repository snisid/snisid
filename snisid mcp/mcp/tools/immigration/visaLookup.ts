import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { immigrationService } from '../../services/immigration.service.js';
import { passportNumberSchema } from '../../validators/identity.validator.js';

export function registerVisaLookupTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'immigration.visaLookup',
    description: 'Lookup visa records.',
    permission: PERMISSIONS.IMMIGRATION_READ,
    inputShape: {
      passportNumber: passportNumberSchema,
      country: z.string().length(2).optional(),
    },
    handler: async (input, ctx) => {
      return immigrationService.visaLookup(input, ctx);
    }
  });
}