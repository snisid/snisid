import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { pnhService } from '../../services/pnh.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerWeaponPermitTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'police.weaponPermit',
    description: 'Verify weapon permit status.',
    permission: PERMISSIONS.POLICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      permitNumber: z.string().min(4).max(64).optional(),
    },
    handler: async (input, ctx) => {
      return pnhService.weaponPermit(input, ctx);
    }
  });
}