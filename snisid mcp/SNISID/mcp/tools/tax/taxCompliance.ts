import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { dgiService } from '../../services/dgi.service.js';
import { nifSchema } from '../../validators/tax.validator.js';

export function registerTaxComplianceTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'tax.taxCompliance',
    description: 'Check tax compliance.',
    permission: PERMISSIONS.TAX_READ,
    inputShape: {
      nif: nifSchema,
    },
    handler: async (input, ctx) => {
      return dgiService.taxCompliance(input, ctx);
    }
  });
}