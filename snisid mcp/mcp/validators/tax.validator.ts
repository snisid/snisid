import { z } from 'zod';

export const nifSchema = z.string().min(5).max(32).regex(/^[A-Z0-9-]+$/i);
export const businessRegistrySchema = z.object({
  registrationNumber: z.string().min(4).max(64).regex(/^[A-Z0-9-]+$/i)
});
