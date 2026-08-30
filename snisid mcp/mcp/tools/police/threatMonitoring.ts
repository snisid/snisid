import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { pnhService } from '../../services/pnh.service.js';

export function registerThreatMonitoringTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'police.threatMonitoring',
    description: 'Submit a threat monitoring query with legal purpose and audit.',
    permission: PERMISSIONS.POLICE_THREAT,
    inputShape: {
      subjectRef: z.string().min(4).max(128),
      timeWindowHours: z.number().int().min(1).max(720).default(24),
    },
    handler: async (input, ctx) => {
      return pnhService.threatMonitoring(input, ctx);
    }
  });
}