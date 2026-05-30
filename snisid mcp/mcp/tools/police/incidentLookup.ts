import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { pnhService } from '../../services/pnh.service.js';

export function registerIncidentLookupTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'police.incidentLookup',
    description: 'Lookup police incident metadata.',
    permission: PERMISSIONS.POLICE_READ,
    inputShape: {
      incidentId: z.string().min(4).max(64),
      district: z.string().max(64).optional(),
    },
    handler: async (input, ctx) => {
      return pnhService.incidentLookup(input, ctx);
    }
  });
}