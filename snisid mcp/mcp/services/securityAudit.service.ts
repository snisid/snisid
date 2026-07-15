import { readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { SNISID } from '../config/constants.js';
import { writeAuditEvent } from '../audit/auditLogger.js';
import { verifyAuditChain } from '../audit/forensics.ts';
import type { AuditEvent } from '../types/security.types.js';

export interface AuditLogAnalysis {
  totalEvents: number;
  deniedEvents: number;
  anomalyDetected: boolean;
  warnings: string[];
}

export class SecurityAuditService {
  private static instance: SecurityAuditService;

  private constructor() {}

  public static getInstance(): SecurityAuditService {
    if (!SecurityAuditService.instance) {
      SecurityAuditService.instance = new SecurityAuditService();
    }
    return SecurityAuditService.instance;
  }

  public async logEvent(event: Omit<AuditEvent, 'id' | 'timestamp' | 'previousHash' | 'hash'>): Promise<AuditEvent> {
    return writeAuditEvent(event);
  }

  public async getLogs(): Promise<AuditEvent[]> {
    const path = join(process.cwd(), SNISID.auditLogPath);
    try {
      const data = await readFile(path, 'utf8');
      return data
        .split('\n')
        .filter(Boolean)
        .map((line) => JSON.parse(line) as AuditEvent);
    } catch {
      return [];
    }
  }

  public async queryLogs(filter: {
    actor?: string;
    action?: string;
    outcome?: 'ALLOW' | 'DENY' | 'ERROR';
    severity?: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
    correlationId?: string;
  }): Promise<AuditEvent[]> {
    const logs = await this.getLogs();
    return logs.filter((log) => {
      if (filter.actor && log.actor !== filter.actor) return false;
      if (filter.action && !log.action.includes(filter.action)) return false;
      if (filter.outcome && log.outcome !== filter.outcome) return false;
      if (filter.severity && log.severity !== filter.severity) return false;
      if (filter.correlationId && log.correlationId !== filter.correlationId) return false;
      return true;
    });
  }

  public async verifyIntegrity(): Promise<{ valid: boolean; brokenAt?: string }> {
    return verifyAuditChain();
  }

  public async analyzeLogs(): Promise<AuditLogAnalysis> {
    const logs = await this.getLogs();
    const result: AuditLogAnalysis = {
      totalEvents: logs.length,
      deniedEvents: 0,
      anomalyDetected: false,
      warnings: []
    };

    // Keep track of sequential denials by actor to detect potential brute forcing / attack
    const sequentialDenies = new Map<string, number>();

    for (const log of logs) {
      if (log.outcome === 'DENY') {
        result.deniedEvents++;
        const currentCount = (sequentialDenies.get(log.actor) ?? 0) + 1;
        sequentialDenies.set(log.actor, currentCount);

        if (currentCount >= 3) {
          result.anomalyDetected = true;
          result.warnings.push(`Brute-force / unauthorized access pattern detected for actor: ${log.actor} (3+ sequential DENY events)`);
        }
      } else {
        // Reset sequential denials on ALLOW
        sequentialDenies.set(log.actor, 0);
      }

      if (log.severity === 'CRITICAL') {
        result.anomalyDetected = true;
        result.warnings.push(`CRITICAL event found: [${log.action}] by [${log.actor}]`);
      }
    }

    return result;
  }
}

export const securityAuditService = SecurityAuditService.getInstance();
