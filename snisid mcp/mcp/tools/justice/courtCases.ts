import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { mjspService } from '../../services/mjsp.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';
import { caseIdSchema } from '../../validators/justice.validator.js';

export function registerCourtCasesTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'justice.courtCases',
    description: 'Search court cases by national identifier or case id.',
    permission: PERMISSIONS.JUSTICE_READ,
    inputShape: {
      nationalId: nationalIdSchema.optional(),
      caseId: caseIdSchema.optional(),
    },
    handler: async (input, ctx) => {
      return mjspService.courtCases(input, ctx);
    }
  });
}