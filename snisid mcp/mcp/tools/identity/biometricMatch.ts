import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { biometricService } from '../../services/biometric.service.js';
import { nationalIdSchema, biometricProbeSchema } from '../../validators/identity.validator.js';

export function registerBiometricMatchTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.biometricMatch',
    description: 'Perform controlled biometric template matching via biometric API.',
    permission: PERMISSIONS.IDENTITY_BIOMETRIC,
    inputShape: {
      nationalId: nationalIdSchema.optional(),
      probe: biometricProbeSchema,
    },
    handler: async (input, ctx) => {
      return biometricService.biometricMatch(input, ctx);
    }
  });
}