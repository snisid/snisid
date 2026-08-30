import type { NextFunction, Request, Response } from 'express';
import { writeAuditEvent } from '../audit/auditLogger.js';

export function auditMiddleware(req: Request, res: Response, next: NextFunction): void {
  const started = Date.now();
  res.on('finish', () => {
    void writeAuditEvent({
      actor: res.locals['principal']?.subject ?? 'anonymous',
      action: `${req.method} ${req.path}`,
      resource: 'http.mcp',
      purpose: req.header('x-purpose') ?? 'UNSPECIFIED',
      correlationId: req.header('x-correlation-id') ?? 'UNSPECIFIED',
      outcome: res.statusCode < 400 ? 'ALLOW' : 'DENY',
      severity: res.statusCode >= 500 ? 'HIGH' : 'LOW',
      metadata: { statusCode: res.statusCode, durationMs: Date.now() - started }
    });
  });
  next();
}
