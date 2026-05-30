import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { oniService } from '../../services/oni.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerCitizenProfileTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.citizenProfile',
    description: 'Read minimized citizen profile through ONI API.',
    permission: PERMISSIONS.IDENTITY_READ,
    inputShape: {
      nationalId: nationalIdSchema,
    },
    handler: async (input, ctx) => {
      return oniService.citizenProfile(input, ctx);
    }
  });
}