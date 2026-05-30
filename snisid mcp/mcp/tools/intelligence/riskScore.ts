import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { intelligenceService } from '../../services/intelligence.service.js';

export function registerRiskScoreTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'intelligence.riskScore',
    description: 'Compute intelligence risk score with explainability.',
    permission: PERMISSIONS.INTELLIGENCE_ANALYZE,
    inputShape: {
      subjectRef: z.string().min(4).max(128),
      modelVersion: z.string().max(32).optional(),
    },
    handler: async (input, ctx) => {
      return intelligenceService.riskScore(input, ctx);
    }
  });
}