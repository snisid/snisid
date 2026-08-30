import { z } from 'zod';
import { nationalIdSchema } from './identity.validator.js';

export const caseIdSchema = z.string().min(4).max(64).regex(/^[A-Z0-9-]+$/i);
export const warrantQuerySchema = z.object({ nationalId: nationalIdSchema, warrantId: z.string().min(4).max(64).optional() });
