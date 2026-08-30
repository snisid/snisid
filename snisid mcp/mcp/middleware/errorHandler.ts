import type { ErrorRequestHandler } from 'express';
import { logger } from '../utils/logger.js';

export const errorHandler: ErrorRequestHandler = (err, _req, res, _next) => {
  logger.error('http_error', { error: err instanceof Error ? err.message : String(err) });
  if (res.headersSent) return;
  res.status(500).json({ error: 'internal_error' });
};
