import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { dgiService } from '../../services/dgi.service.js';

export function registerBusinessRegistryTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'tax.businessRegistry',
    description: 'Lookup business registry.',
    permission: PERMISSIONS.TAX_READ,
    inputShape: {
      registrationNumber: z.string().min(4).max(64),
    },
    handler: async (input, ctx) => {
      return dgiService.businessRegistry(input, ctx);
    }
  });
}