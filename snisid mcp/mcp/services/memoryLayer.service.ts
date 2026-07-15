import { encryptJson, decryptJson } from '../security/encryption.js';
import { redact } from '../utils/logger.js';
import type { SecurityContext } from '../types/security.types.js';

export interface MemoryRecord {
  key: string;
  encryptedValue: string; // encrypted base64 payload
  owner: string; // actor/agent who created it
  clearanceRequired: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET';
  updatedAt: string;
}

const CLEARANCE_LEVELS = {
  PUBLIC: 0,
  INTERNAL: 1,
  CONFIDENTIAL: 2,
  SECRET: 3,
  TOP_SECRET: 4
};

export class MemoryLayerService {
  private static instance: MemoryLayerService;
  private storage = new Map<string, MemoryRecord>();

  private constructor() {}

  public static getInstance(): MemoryLayerService {
    if (!MemoryLayerService.instance) {
      MemoryLayerService.instance = new MemoryLayerService();
    }
    return MemoryLayerService.instance;
  }

  private hasSufficientClearance(
    actorClearance: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET',
    requiredClearance: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET'
  ): boolean {
    const actorNum = CLEARANCE_LEVELS[actorClearance] ?? 0;
    const requiredNum = CLEARANCE_LEVELS[requiredClearance] ?? 0;
    return actorNum >= requiredNum;
  }

  public async store(
    key: string,
    value: unknown,
    clearance: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET',
    ctx: SecurityContext
  ): Promise<void> {
    // 1. Enforce Clearance check: the storer cannot set memory with clearance higher than their own
    if (!this.hasSufficientClearance(ctx.principal.clearance, clearance)) {
      throw new Error(`INSUFFICIENT_CLEARANCE: Cannot write memory with clearance ${clearance} with your clearance ${ctx.principal.clearance}`);
    }

    // 2. Redact sensitive values first (defense-in-depth)
    const redactedValue = redact(value);

    // 3. Encrypt the value
    const encryptedValue = encryptJson(redactedValue);

    // 4. Save record
    const record: MemoryRecord = {
      key,
      encryptedValue,
      owner: ctx.principal.subject,
      clearanceRequired: clearance,
      updatedAt: new Date().toISOString()
    };

    this.storage.set(key, record);
  }

  public async retrieve<T>(key: string, ctx: SecurityContext): Promise<T | undefined> {
    const record = this.storage.get(key);
    if (!record) {
      return undefined;
    }

    // 1. Verify clearance
    if (!this.hasSufficientClearance(ctx.principal.clearance, record.clearanceRequired)) {
      throw new Error(`INSUFFICIENT_CLEARANCE: Access to memory ${key} requires ${record.clearanceRequired} clearance`);
    }

    // 2. Decrypt the value
    return decryptJson<T>(record.encryptedValue);
  }

  public async delete(key: string, ctx: SecurityContext): Promise<boolean> {
    const record = this.storage.get(key);
    if (!record) {
      return false;
    }

    // Enforce clearance or ownership check
    if (!this.hasSufficientClearance(ctx.principal.clearance, record.clearanceRequired) && record.owner !== ctx.principal.subject) {
      throw new Error(`INSUFFICIENT_CLEARANCE: Cannot delete memory ${key}`);
    }

    return this.storage.delete(key);
  }

  public clearAll(): void {
    this.storage.clear();
  }
}

export const memoryLayerService = MemoryLayerService.getInstance();
