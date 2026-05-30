import { z } from 'zod';

export const nationalIdSchema = z.string().min(5).max(64).regex(/^[A-Z0-9-]+$/i);
export const passportNumberSchema = z.string().min(5).max(32).regex(/^[A-Z0-9-]+$/i);
export const biometricProbeSchema = z.object({
  templateRef: z.string().min(8).max(256),
  modality: z.enum(['FACE', 'FINGERPRINT', 'IRIS', 'MULTI']).default('FACE')
});
export const identityQuerySchema = z.object({
  nationalId: nationalIdSchema,
  consentReference: z.string().min(4).max(128).optional()
});
