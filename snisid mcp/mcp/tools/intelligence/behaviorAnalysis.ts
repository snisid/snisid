import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { intelligenceService } from '../../services/intelligence.service.js';

export function registerBehaviorAnalysisTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'intelligence.behaviorAnalysis',
    description: 'Analyze behavior patterns with minimization and audit.',
    permission: PERMISSIONS.INTELLIGENCE_ANALYZE,
    inputShape: {
      subjectRef: z.string().min(4).max(128),
      scope: z.enum(["LOW","MEDIUM","HIGH"]).default("LOW"),
    },
    handler: async (input, ctx) => {
      return intelligenceService.behaviorAnalysis(input, ctx);
    }
  });
}