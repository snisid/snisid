import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { intelligenceService } from '../../services/intelligence.service.js';

export function registerNetworkAnalysisTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'intelligence.networkAnalysis',
    description: 'Analyze authorized relationship network.',
    permission: PERMISSIONS.INTELLIGENCE_ANALYZE,
    inputShape: {
      seedRefs: z.array(z.string().min(4).max(128)).min(1).max(20),
      depth: z.number().int().min(1).max(3).default(1),
    },
    handler: async (input, ctx) => {
      return intelligenceService.networkAnalysis(input, ctx);
    }
  });
}