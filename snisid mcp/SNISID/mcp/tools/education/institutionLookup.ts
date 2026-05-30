import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { educationService } from '../../services/education.service.js';

export function registerInstitutionLookupTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'education.institutionLookup',
    description: 'Lookup accredited institution.',
    permission: PERMISSIONS.EDUCATION_READ,
    inputShape: {
      institutionId: z.string().min(3).max(64).optional(),
      name: z.string().min(2).max(128).optional(),
    },
    handler: async (input, ctx) => {
      return educationService.institutionLookup(input, ctx);
    }
  });
}