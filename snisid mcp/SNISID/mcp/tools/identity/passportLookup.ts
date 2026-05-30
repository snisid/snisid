import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { passportService } from '../../services/passport.service.js';
import { passportNumberSchema } from '../../validators/identity.validator.js';

export function registerPassportLookupTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.passportLookup',
    description: 'Lookup passport status via passport API.',
    permission: PERMISSIONS.IDENTITY_READ,
    inputShape: {
      passportNumber: passportNumberSchema,
    },
    handler: async (input, ctx) => {
      return passportService.passportLookup(input, ctx);
    }
  });
}