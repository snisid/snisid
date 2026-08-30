import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { biometricService } from '../../services/biometric.service.js';
import { nationalIdSchema, biometricProbeSchema } from '../../validators/identity.validator.js';

export function registerFaceVerificationTool(server: McpServer): void {
  registerGovernmentTool(server, {
    name: 'identity.faceVerification',
    description: 'Verify face template reference against citizen record.',
    permission: PERMISSIONS.IDENTITY_BIOMETRIC,
    inputShape: {
      nationalId: nationalIdSchema,
      probe: biometricProbeSchema,
    },
    handler: async (input, ctx) => {
      return biometricService.faceVerification(input, ctx);
    }
  });
}