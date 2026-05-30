import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { intelligenceService } from '../../services/intelligence.service.js';

export function registerFusionAnalysisTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'intelligence.fusionAnalysis',
    description: 'Perform multi-source fusion analysis with human oversight.',
    permission: PERMISSIONS.INTELLIGENCE_ANALYZE,
    inputShape: {
      subjectRefs: z.array(z.string().min(4).max(128)).min(1).max(20),
      hypothesis: z.string().min(8).max(1024),
    },
    handler: async (input, ctx) => {
      return intelligenceService.fusionAnalysis(input, ctx);
    }
  });
}