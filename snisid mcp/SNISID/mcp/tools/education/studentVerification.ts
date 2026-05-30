import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { educationService } from '../../services/education.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerStudentVerificationTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'education.studentVerification',
    description: 'Verify student status.',
    permission: PERMISSIONS.EDUCATION_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      institutionId: z.string().min(3).max(64).optional(),
    },
    handler: async (input, ctx) => {
      return educationService.studentVerification(input, ctx);
    }
  });
}