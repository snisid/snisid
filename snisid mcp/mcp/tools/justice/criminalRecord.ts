import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { mjspService } from '../../services/mjsp.service.js';
import { nationalIdSchema } from '../../validators/identity.validator.js';

export function registerCriminalRecordTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'justice.criminalRecord',
    description: 'Retrieve criminal record summary with judicial authorization.',
    permission: PERMISSIONS.JUSTICE_READ,
    inputShape: {
      nationalId: nationalIdSchema,
      caseScope: z.enum(["SUMMARY","FULL"]).default("SUMMARY"),
    },
    handler: async (input, ctx) => {
      return mjspService.criminalRecord(input, ctx);
    }
  });
}