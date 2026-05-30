import type { NextFunction, Request, Response } from 'express';
import { verifyAccessToken } from '../security/jwt.js';

export function authMiddleware(req: Request, res: Response, next: NextFunction): void {
  try {
    const header = req.header('authorization');
    if (!header?.startsWith('Bearer ')) {
      res.status(401).json({ error: 'missing_bearer_token' });
      return;
    }
    res.locals.principal = verifyAccessToken(header.slice('Bearer '.length));
    next();
  } catch {
    res.status(401).json({ error: 'invalid_token' });
  }
}
