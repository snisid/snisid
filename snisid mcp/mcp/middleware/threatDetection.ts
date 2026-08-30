import type { NextFunction, Request, Response } from 'express';
import { recordSecurityEvent } from '../audit/securityEvents.js';

const suspicious = [/ignore previous/i, /jailbreak/i, /bypass/i, /<script/i, /\.\.\//, /\$\(/];

export function threatDetectionMiddleware(req: Request, res: Response, next: NextFunction): void {
  const body = JSON.stringify(req.body ?? {});
  if (suspicious.some((pattern) => pattern.test(body))) {
    void recordSecurityEvent({
      action: 'threat.detected.http_payload',
      resource: req.path,
      outcome: 'DENY',
      severity: 'HIGH',
      ...(req.header('x-purpose') ? { purpose: req.header('x-purpose')! } : {}),
      ...(req.header('x-correlation-id') ? { correlationId: req.header('x-correlation-id')! } : {}),
      metadata: { pattern: 'prompt_or_injection_attempt' }
    });
    res.status(400).json({ error: 'suspicious_payload' });
    return;
  }
  next();
}
