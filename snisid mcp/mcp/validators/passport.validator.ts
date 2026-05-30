import { z } from 'zod';

export const travelDocumentSchema = z.object({
  passportNumber: z.string().min(5).max(32).regex(/^[A-Z0-9-]+$/i),
  country: z.string().length(2).optional()
});
