import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { educationService } from '../../services/education.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerDiplomaVerificationTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'education.diplomaVerification',
    description: 'Verify diploma authenticity.',
    permission: PERMISSIONS.EDUCATION_READ,
    inputShape: {
      diplomaNumber: z.string().min(4).max(64),
      nationalId: nationalIdSchema.optional(),
    },
    handler: async (input, ctx) => {
      return educationService.diplomaVerification(input, ctx);
    }
  });
}