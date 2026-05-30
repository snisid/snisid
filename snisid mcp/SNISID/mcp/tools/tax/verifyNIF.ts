import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { dgiService } from '../../services/dgi.service.js';
import { nifSchema } from '../../validators/tax.validator.js';

export function registerVerifyNifTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'tax.verifyNIF',
    description: 'Verify NIF through DGI API.',
    permission: PERMISSIONS.TAX_READ,
    inputShape: {
      nif: nifSchema,
    },
    handler: async (input, ctx) => {
      return dgiService.verifyNif(input, ctx);
    }
  });
}