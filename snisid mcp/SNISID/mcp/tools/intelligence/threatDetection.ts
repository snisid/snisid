import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { intelligenceService } from '../../services/intelligence.service.js';

export function registerThreatDetectionTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'intelligence.threatDetection',
    description: 'Detect threat patterns from authorized signals.',
    permission: PERMISSIONS.INTELLIGENCE_ANALYZE,
    inputShape: {
      signalRefs: z.array(z.string().min(4).max(128)).min(1).max(100),
      timeWindowHours: z.number().int().min(1).max(720).default(24),
    },
    handler: async (input, ctx) => {
      return intelligenceService.threatDetection(input, ctx);
    }
  });
}