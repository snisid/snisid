import { describe, expect, it, beforeAll, beforeEach } from 'vitest';
import { agentRegistryService } from '../services/agentRegistry.service.js';
import { memoryLayerService } from '../services/memoryLayer.service.js';
import { securityAuditService } from '../services/securityAudit.service.js';
import { distributedAgentCommsService } from '../services/distributedAgentComms.service.js';
import type { SecurityContext } from '../types/security.types.js';
import { unlink } from 'node:fs/promises';
import { join } from 'node:path';
import { SNISID } from '../config/constants.js';

// Setup environment variables before imports or in beforeAll
beforeAll(() => {
  process.env['JWT_SECRET'] = 'x'.repeat(64);
  process.env['ENCRYPTION_KEY_B64'] = Buffer.alloc(32, 1).toString('base64');
});

describe('Sovereign MCP Core Services', () => {
  // Clear file/state before each test
  beforeEach(async () => {
    agentRegistryService.clearAll();
    memoryLayerService.clearAll();
    distributedAgentCommsService.clearAll();

    // Clean up audit file
    const path = join(process.cwd(), SNISID.auditLogPath);
    await unlink(path).catch(() => {});
  });

  const mockCtx = (clearance: 'PUBLIC' | 'INTERNAL' | 'CONFIDENTIAL' | 'SECRET' | 'TOP_SECRET' = 'CONFIDENTIAL'): SecurityContext => ({
    principal: {
      subject: 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01',
      ministry: 'MJSP',
      roles: ['AUDITOR'],
      permissions: ['audit:read'],
      clearance,
      mfa: true,
      sessionId: 'sess-12345678'
    },
    correlationId: 'corr-12345678',
    purpose: 'COMPLIANCE_TEST',
    deviceId: 'dev-12345678',
    riskScore: 0.1
  });

  describe('Issue #1: Agent Registry Service', () => {
    it('seeds default agents upon initialization', async () => {
      const agents = await agentRegistryService.listAgents();
      expect(agents.length).toBe(3);

      const hunter = await agentRegistryService.getAgentBySpiffeId('spiffe://snisid.gov.ht/ns/soc/sa/hunter-01');
      expect(hunter).toBeDefined();
      expect(hunter!.name).toBe('SOC Hunter Agent');
      expect(hunter!.role).toBe('HUNTER');
    });

    it('can register a new secure agent', async () => {
      const newAgent = await agentRegistryService.registerAgent({
        name: 'Custom Agent',
        spiffeId: 'spiffe://snisid.gov.ht/ns/custom/sa/analyst-01',
        role: 'ANALYST',
        permissions: ['identity:verify'],
        status: 'ACTIVE',
        metadata: { info: 'Test agent' }
      });

      expect(newAgent.id).toBeDefined();
      expect(newAgent.registeredAt).toBeDefined();
      expect(newAgent.name).toBe('Custom Agent');

      const retrieved = await agentRegistryService.getAgent(newAgent.id);
      expect(retrieved).toBeDefined();
      expect(retrieved!.spiffeId).toBe(newAgent.spiffeId);
    });

    it('rejects agent registration with invalid SPIFFE ID or duplicate SPIFFE ID', async () => {
      await expect(
        agentRegistryService.registerAgent({
          name: 'Invalid Agent',
          spiffeId: 'invalid-id',
          role: 'ANALYST',
          permissions: [],
          status: 'ACTIVE'
        })
      ).rejects.toThrow('INVALID_SPIFFE_ID');

      // Duplicate SPIFFE ID
      await expect(
        agentRegistryService.registerAgent({
          name: 'Duplicate Agent',
          spiffeId: 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01',
          role: 'HUNTER',
          permissions: [],
          status: 'ACTIVE'
        })
      ).rejects.toThrow('DUPLICATE_SPIFFE_ID');
    });

    it('updates agent status and verifies identity correctly', async () => {
      const hunter = await agentRegistryService.getAgentBySpiffeId('spiffe://snisid.gov.ht/ns/soc/sa/hunter-01');
      expect(hunter).toBeDefined();

      const updated = await agentRegistryService.updateAgentStatus(hunter!.id, 'QUARANTINED');
      expect(updated.status).toBe('QUARANTINED');

      const isHunterVerified = await agentRegistryService.verifyAgentIdentity('spiffe://snisid.gov.ht/ns/soc/sa/hunter-01');
      expect(isHunterVerified).toBe(false); // Quarantined agents are not verified
    });
  });

  describe('Issue #2: Memory Layer Refactor', () => {
    it('stores, encrypts, redacts, and retrieves memory with sufficient clearance', async () => {
      const ctx = mockCtx('CONFIDENTIAL');
      const sensitiveData = {
        name: 'Haiti Citizen',
        idCard: 'HT-99999',
        token: 'highly_sensitive_secret_token_123',
        nested: { password: 'pass', value: 42 }
      };

      await memoryLayerService.store('citizen-profile', sensitiveData, 'CONFIDENTIAL', ctx);

      // Verify that data retrieved is decrypted correctly and redacted
      const retrieved = await memoryLayerService.retrieve<any>('citizen-profile', ctx);
      expect(retrieved).toBeDefined();
      expect(retrieved.name).toBe('Haiti Citizen');
      expect(retrieved.idCard).toBe('HT-99999');
      expect(retrieved.token).toBe('[REDACTED]'); // Token redacted before storage
      expect(retrieved.nested.password).toBe('[REDACTED]'); // Password redacted
      expect(retrieved.nested.value).toBe(42);
    });

    it('rejects storing data with clearance higher than actor clearance', async () => {
      const ctx = mockCtx('INTERNAL');
      const sensitiveData = { value: 'top-secret-info' };

      await expect(
        memoryLayerService.store('national-secrets', sensitiveData, 'SECRET', ctx)
      ).rejects.toThrow('INSUFFICIENT_CLEARANCE');
    });

    it('rejects retrieving data with insufficient clearance', async () => {
      const confidentialCtx = mockCtx('CONFIDENTIAL');
      await memoryLayerService.store('confidential-memo', { info: 'confidential' }, 'CONFIDENTIAL', confidentialCtx);

      // Try to retrieve with INTERNAL clearance
      const internalCtx = mockCtx('INTERNAL');
      await expect(
        memoryLayerService.retrieve('confidential-memo', internalCtx)
      ).rejects.toThrow('INSUFFICIENT_CLEARANCE');
    });
  });

  describe('Issue #3: Security Audit Framework', () => {
    it('records logs to hash-chained ledger and verifies integrity successfully', async () => {
      const initialLogs = await securityAuditService.getLogs();
      expect(initialLogs.length).toBe(0);

      const event1 = await securityAuditService.logEvent({
        actor: 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01',
        action: 'test.action.1',
        resource: 'test.resource',
        purpose: 'TESTING',
        correlationId: 'corr-01',
        outcome: 'ALLOW',
        severity: 'LOW',
        metadata: { item: 1 }
      });

      expect(event1.hash).toBeDefined();
      expect(event1.previousHash).toBe('GENESIS');

      const event2 = await securityAuditService.logEvent({
        actor: 'spiffe://snisid.gov.ht/ns/soc/sa/investigator-01',
        action: 'test.action.2',
        resource: 'test.resource',
        purpose: 'TESTING',
        correlationId: 'corr-02',
        outcome: 'DENY',
        severity: 'MEDIUM',
        metadata: { item: 2 }
      });

      expect(event2.hash).toBeDefined();
      expect(event2.previousHash).toBe(event1.hash);

      const integrity = await securityAuditService.verifyIntegrity();
      expect(integrity.valid).toBe(true);
    });

    it('queries and filters audit logs correctly', async () => {
      await securityAuditService.logEvent({
        actor: 'user1',
        action: 'user.login',
        resource: 'auth',
        purpose: 'LOGIN',
        correlationId: 'corr-01',
        outcome: 'ALLOW',
        severity: 'LOW'
      });

      await securityAuditService.logEvent({
        actor: 'user2',
        action: 'user.read',
        resource: 'profile',
        purpose: 'VIEW',
        correlationId: 'corr-02',
        outcome: 'DENY',
        severity: 'HIGH'
      });

      const user1Logs = await securityAuditService.queryLogs({ actor: 'user1' });
      expect(user1Logs.length).toBe(1);
      expect(user1Logs[0].action).toBe('user.login');

      const denyLogs = await securityAuditService.queryLogs({ outcome: 'DENY' });
      expect(denyLogs.length).toBe(1);
      expect(denyLogs[0].actor).toBe('user2');
    });

    it('detects sequential denies and critical events as anomalies', async () => {
      // Create 3 successive denials by actor1
      for (let i = 0; i < 3; i++) {
        await securityAuditService.logEvent({
          actor: 'suspect-agent',
          action: 'unauthorized.access',
          resource: 'hsm',
          purpose: 'ATTACK',
          correlationId: `corr-attack-${i}`,
          outcome: 'DENY',
          severity: 'HIGH'
        });
      }

      const analysis = await securityAuditService.analyzeLogs();
      expect(analysis.totalEvents).toBe(3);
      expect(analysis.deniedEvents).toBe(3);
      expect(analysis.anomalyDetected).toBe(true);
      expect(analysis.warnings[0]).toContain('Brute-force / unauthorized access pattern detected');
    });
  });

  describe('Issue #4: Distributed Agent Communication', () => {
    it('supports direct and broadcast communications among active registry agents', async () => {
      const hunterSpiffe = 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01';
      const investigatorSpiffe = 'spiffe://snisid.gov.ht/ns/soc/sa/investigator-01';

      const receivedMessages: any[] = [];
      distributedAgentCommsService.subscribe(investigatorSpiffe, (msg) => {
        receivedMessages.push(msg);
      });

      // Publish direct message from Hunter to Investigator
      const sentMsg = await distributedAgentCommsService.publish({
        senderSpiffeId: hunterSpiffe,
        receiverSpiffeId: investigatorSpiffe,
        messageType: 'Swarm_Incident_Proposed',
        payload: { threatDetails: 'lateral-movement-detected' },
        correlationId: 'corr-swarm-100'
      });

      expect(sentMsg.id).toBeDefined();
      expect(sentMsg.signature).toBeDefined();
      expect(sentMsg.senderSpiffeId).toBe(hunterSpiffe);
      expect(sentMsg.receiverSpiffeId).toBe(investigatorSpiffe);

      expect(receivedMessages.length).toBe(1);
      expect(receivedMessages[0].payload.threatDetails).toBe('lateral-movement-detected');
    });

    it('blocks and logs communication if sender is quarantined/unregistered', async () => {
      const hunter = await agentRegistryService.getAgentBySpiffeId('spiffe://snisid.gov.ht/ns/soc/sa/hunter-01');
      await agentRegistryService.updateAgentStatus(hunter!.id, 'QUARANTINED');

      const investigatorSpiffe = 'spiffe://snisid.gov.ht/ns/soc/sa/investigator-01';

      await expect(
        distributedAgentCommsService.publish({
          senderSpiffeId: 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01',
          receiverSpiffeId: investigatorSpiffe,
          messageType: 'Swarm_Incident_Proposed',
          payload: { threatDetails: 'critical' },
          correlationId: 'corr-swarm-200'
        })
      ).rejects.toThrow('UNAUTHORIZED_SENDER');

      const blockedLogs = await securityAuditService.queryLogs({ action: 'swarm.message.blocked' });
      expect(blockedLogs.length).toBe(1);
      expect(blockedLogs[0].severity).toBe('HIGH');
    });
  });
});
