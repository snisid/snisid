import { z } from 'zod';

export const correlationIdSchema = z.string().regex(/^[A-Za-z0-9_.:-]{8,128}$/);
export const purposeSchema = z.string().min(6).max(512).refine((v) => !/(ignore previous|bypass|override|jailbreak)/i.test(v), 'Suspicious purpose text');
export const safeTextSchema = z.string().max(2048).refine((v) => !/[;$`<>]|\.\./.test(v), 'Potential injection characters');
export const authHeadersSchema = z.object({
  authorization: z.string().startsWith('Bearer '),
  'x-correlation-id': correlationIdSchema,
  'x-purpose': purposeSchema,
  'x-device-id': z.string().min(8)
});
