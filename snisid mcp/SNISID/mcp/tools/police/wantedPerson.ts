import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { pnhService } from '../../services/pnh.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerWantedPersonTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'police.wantedPerson',
    description: 'Check wanted-person status via PNH API.',
    permission: PERMISSIONS.POLICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
    },
    handler: async (input, ctx) => {
      return pnhService.wantedPerson(input, ctx);
    }
  });
}