import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { oniService } from '../../services/oni.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerBirthCertificateTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.birthCertificate',
    description: 'Validate birth certificate metadata.',
    permission: PERMISSIONS.IDENTITY_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      certificateNumber: z.string().min(4).max(64),
    },
    handler: async (input, ctx) => {
      return oniService.birthCertificate(input, ctx);
    }
  });
}