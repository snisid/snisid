import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { educationService } from '../../services/education.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerAcademicHistoryTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'education.academicHistory',
    description: 'Retrieve academic history with minimization.',
    permission: PERMISSIONS.EDUCATION_READ,
    inputShape: {
      nationalId: nationalIdSchema,
    },
    handler: async (input, ctx) => {
      return educationService.academicHistory(input, ctx);
    }
  });
}