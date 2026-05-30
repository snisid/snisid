import type { NextFunction, Request, Response } from 'express';
import { authHeadersSchema } from '../validators/security.validator.js';

export function requestValidator(req: Request, res: Response, next: NextFunction): void {
  const parsed = authHeadersSchema.safeParse({
    authorization: req.header('authorization'),
    'x-correlation-id': req.header('x-correlation-id'),
    'x-purpose': req.header('x-purpose'),
    'x-device-id': req.header('x-device-id')
  });
  if (!parsed.success) {
    res.status(400).json({ error: 'invalid_headers', details: parsed.error.flatten().fieldErrors });
    return;
  }
  next();
}
