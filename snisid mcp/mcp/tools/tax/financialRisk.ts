import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { dgiService } from '../../services/dgi.service.js';
import { nifSchema } from '../../validators/tax.validator.js';

export function registerFinancialRiskTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'tax.financialRisk',
    description: 'Compute financial risk through DGI risk API.',
    permission: PERMISSIONS.TAX_RISK,
    inputShape: {
      nif: nifSchema,
      signals: z.array(z.string().max(64)).max(20).optional(),
    },
    handler: async (input, ctx) => {
      return dgiService.financialRisk(input, ctx);
    }
  });
}