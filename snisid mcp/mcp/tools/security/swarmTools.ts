import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { PERMISSIONS } from '../../config/permissions.js';
import { registerGovernmentTool } from '../_shared/toolFactory.js';
import { agentRegistryService } from '../../services/agentRegistry.service.js';
import { memoryLayerService } from '../../services/memoryLayer.service.js';
import { securityAuditService } from '../../services/securityAudit.service.js';
import { distributedAgentCommsService } from '../../services/distributedAgentComms.service.js';

export function registerSwarmTools(server: McpServer): void {
  // 1. Agent Registry: Register Agent
  registerGovernmentTool(server, {
    name: 'swarm.registerAgent',
    description: 'Register a new secure AI agent with its SPIFFE ID, role, and fine-grained permissions.',
    permission: PERMISSIONS.SECURITY_ADMIN,
    inputShape: {
      name: z.string().min(3).max(128),
      spiffeId: z.string().startsWith('spiffe://'),
      role: z.string().min(3).max(64),
      permissions: z.array(z.string()),
      status: z.enum(['ACTIVE', 'INACTIVE', 'QUARANTINED']).default('ACTIVE'),
      metadata: z.record(z.unknown()).optional()
    },
    handler: async (input) => {
      return agentRegistryService.registerAgent({
        name: input.name,
        spiffeId: input.spiffeId,
        role: input.role,
        permissions: input.permissions,
        status: input.status,
        metadata: input.metadata
      });
    }
  });

  // 2. Agent Registry: List Agents
  registerGovernmentTool(server, {
    name: 'swarm.listAgents',
    description: 'Retrieve the inventory of registered swarm AI agents and their current lifecycle statuses.',
    permission: PERMISSIONS.SECURITY_ADMIN,
    inputShape: {},
    handler: async () => {
      return agentRegistryService.listAgents();
    }
  });

  // 3. Agent Registry: Quarantine Agent
  registerGovernmentTool(server, {
    name: 'swarm.quarantineAgent',
    description: 'Revoke active status and immediately quarantine a compromised or anomalous AI agent.',
    permission: PERMISSIONS.SECURITY_ADMIN,
    inputShape: {
      id: z.string().uuid()
    },
    handler: async (input) => {
      return agentRegistryService.updateAgentStatus(input.id, 'QUARANTINED');
    }
  });

  // 4. Memory Layer: Store
  registerGovernmentTool(server, {
    name: 'swarm.memoryStore',
    description: 'Store encrypted and redacted agent memory with specific cryptographic security clearance requirements.',
    permission: PERMISSIONS.AI_ORCHESTRATE,
    inputShape: {
      key: z.string().min(3).max(256),
      value: z.unknown(),
      clearanceRequired: z.enum(['PUBLIC', 'INTERNAL', 'CONFIDENTIAL', 'SECRET', 'TOP_SECRET'])
    },
    handler: async (input, ctx) => {
      await memoryLayerService.store(input.key, input.value, input.clearanceRequired, ctx);
      return { success: true, key: input.key, status: 'ENCRYPTED_AND_STORED' };
    }
  });

  // 5. Memory Layer: Retrieve
  registerGovernmentTool(server, {
    name: 'swarm.memoryRetrieve',
    description: 'Decrypt and retrieve agent memory securely, enforcing access control based on cryptographic clearance levels.',
    permission: PERMISSIONS.AI_ORCHESTRATE,
    inputShape: {
      key: z.string().min(3).max(256)
    },
    handler: async (input, ctx) => {
      const data = await memoryLayerService.retrieve(input.key, ctx);
      if (!data) {
        throw new Error(`MEMORY_NOT_FOUND: Memory for key ${input.key} not found`);
      }
      return data;
    }
  });

  // 6. Security Audit: Verify
  registerGovernmentTool(server, {
    name: 'swarm.auditVerify',
    description: 'Verify the cryptographic immutability and complete integrity of the hash-chained Sovereign Audit Ledger.',
    permission: PERMISSIONS.AUDIT_READ,
    inputShape: {},
    handler: async () => {
      return securityAuditService.verifyIntegrity();
    }
  });

  // 7. Security Audit: Analyze Logs
  registerGovernmentTool(server, {
    name: 'swarm.auditAnalyze',
    description: 'Analyze secure system audit logs for potential anomalies, consecutive denials, or threat incidents.',
    permission: PERMISSIONS.AUDIT_READ,
    inputShape: {},
    handler: async () => {
      return securityAuditService.analyzeLogs();
    }
  });

  // 8. Distributed Agent Communication: Publish Message
  registerGovernmentTool(server, {
    name: 'swarm.publishMessage',
    description: 'Transmit a signed, secure swarm communication message (direct or broadcast) over the Sovereign Event Backbone.',
    permission: PERMISSIONS.AI_ORCHESTRATE,
    inputShape: {
      senderSpiffeId: z.string().startsWith('spiffe://'),
      receiverSpiffeId: z.string().startsWith('spiffe://').optional(),
      messageType: z.enum(['Swarm_Incident_Proposed', 'Evidence_Package', 'Containment_Request', 'Broadcast']),
      payload: z.unknown(),
      correlationId: z.string().min(8)
    },
    handler: async (input) => {
      return distributedAgentCommsService.publish({
        senderSpiffeId: input.senderSpiffeId,
        receiverSpiffeId: input.receiverSpiffeId,
        messageType: input.messageType,
        payload: input.payload,
        correlationId: input.correlationId
      });
    }
  });
}
