import { randomUUID } from 'node:crypto';
import { agentRegistryService } from './agentRegistry.service.js';
import { securityAuditService } from './securityAudit.service.js';
import { sha256 } from '../utils/crypto.js';

export interface SwarmMessage {
  id: string;
  senderSpiffeId: string;
  receiverSpiffeId?: string; // undefined for broadcast
  messageType: 'Swarm_Incident_Proposed' | 'Evidence_Package' | 'Containment_Request' | 'Broadcast';
  payload: unknown;
  correlationId: string;
  timestamp: string;
  signature: string; // Cryptographic integrity proof
}

export type SwarmMessageHandler = (message: SwarmMessage) => void | Promise<void>;

export class DistributedAgentCommsService {
  private static instance: DistributedAgentCommsService;
  private subscriptions = new Map<string, Set<SwarmMessageHandler>>();

  private constructor() {}

  public static getInstance(): DistributedAgentCommsService {
    if (!DistributedAgentCommsService.instance) {
      DistributedAgentCommsService.instance = new DistributedAgentCommsService();
    }
    return DistributedAgentCommsService.instance;
  }

  // Create a cryptographic signature for a message
  private generateSignature(senderSpiffeId: string, payload: unknown, correlationId: string): string {
    const serializedPayload = JSON.stringify(payload);
    // Combine fields with a internal secret seed/salt to sign
    const signInput = `${senderSpiffeId}:${serializedPayload}:${correlationId}:snisid-swarm-secret-key-2026`;
    return sha256(signInput);
  }

  public async publish(msgInput: {
    senderSpiffeId: string;
    receiverSpiffeId?: string;
    messageType: 'Swarm_Incident_Proposed' | 'Evidence_Package' | 'Containment_Request' | 'Broadcast';
    payload: unknown;
    correlationId: string;
  }): Promise<SwarmMessage> {
    // 1. Zero Trust: Verify sender identity in AgentRegistry
    const senderVerified = await agentRegistryService.verifyAgentIdentity(msgInput.senderSpiffeId);
    if (!senderVerified) {
      // Log the security violation
      await securityAuditService.logEvent({
        actor: msgInput.senderSpiffeId,
        action: 'swarm.message.blocked',
        resource: 'swarm.backbone',
        purpose: 'PUBLISH',
        correlationId: msgInput.correlationId,
        outcome: 'DENY',
        severity: 'HIGH',
        metadata: { reason: 'SENDER_NOT_ACTIVE_OR_UNREGISTERED', input: msgInput }
      });
      throw new Error(`UNAUTHORIZED_SENDER: SPIFFE ID ${msgInput.senderSpiffeId} is not active or registered`);
    }

    // 2. Zero Trust: Verify receiver identity if specified
    if (msgInput.receiverSpiffeId) {
      const receiverVerified = await agentRegistryService.verifyAgentIdentity(msgInput.receiverSpiffeId);
      if (!receiverVerified) {
        await securityAuditService.logEvent({
          actor: msgInput.senderSpiffeId,
          action: 'swarm.message.blocked',
          resource: 'swarm.backbone',
          purpose: 'PUBLISH',
          correlationId: msgInput.correlationId,
          outcome: 'DENY',
          severity: 'HIGH',
          metadata: { reason: 'RECEIVER_NOT_ACTIVE_OR_UNREGISTERED', input: msgInput }
        });
        throw new Error(`UNAUTHORIZED_RECEIVER: SPIFFE ID ${msgInput.receiverSpiffeId} is not active or registered`);
      }
    }

    // 3. Construct Swarm Message with Cryptographic Signature
    const message: SwarmMessage = {
      id: randomUUID(),
      timestamp: new Date().toISOString(),
      senderSpiffeId: msgInput.senderSpiffeId,
      receiverSpiffeId: msgInput.receiverSpiffeId,
      messageType: msgInput.messageType,
      payload: msgInput.payload,
      correlationId: msgInput.correlationId,
      signature: this.generateSignature(msgInput.senderSpiffeId, msgInput.payload, msgInput.correlationId)
    };

    // 4. Log publish success to Sovereign Audit Trail
    await securityAuditService.logEvent({
      actor: msgInput.senderSpiffeId,
      action: 'swarm.message.publish',
      resource: 'swarm.backbone',
      purpose: 'COMMUNICATION',
      correlationId: msgInput.correlationId,
      outcome: 'ALLOW',
      severity: 'LOW',
      metadata: { id: message.id, messageType: message.messageType, receiverSpiffeId: message.receiverSpiffeId }
    });

    // 5. Deliver message to subscribers/receivers
    await this.deliver(message);

    return message;
  }

  public subscribe(agentSpiffeId: string, handler: SwarmMessageHandler): void {
    if (!this.subscriptions.has(agentSpiffeId)) {
      this.subscriptions.set(agentSpiffeId, new Set());
    }
    this.subscriptions.get(agentSpiffeId)!.add(handler);
  }

  public unsubscribe(agentSpiffeId: string, handler?: SwarmMessageHandler): void {
    if (!handler) {
      this.subscriptions.delete(agentSpiffeId);
    } else {
      const handlers = this.subscriptions.get(agentSpiffeId);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          this.subscriptions.delete(agentSpiffeId);
        }
      }
    }
  }

  private async deliver(message: SwarmMessage): Promise<void> {
    const recipients: { spiffeId: string; handlers: Set<SwarmMessageHandler> }[] = [];

    if (message.receiverSpiffeId) {
      // 1-to-1 delivery
      const handlers = this.subscriptions.get(message.receiverSpiffeId);
      if (handlers) {
        recipients.push({ spiffeId: message.receiverSpiffeId, handlers });
      }
    } else {
      // Broadcast delivery to all other agents (excluding sender)
      for (const [spiffeId, handlers] of this.subscriptions.entries()) {
        if (spiffeId !== message.senderSpiffeId) {
          recipients.push({ spiffeId, handlers });
        }
      }
    }

    for (const recipient of recipients) {
      for (const handler of recipient.handlers) {
        try {
          // Log dispatch event to Sovereign Audit Trail
          await securityAuditService.logEvent({
            actor: recipient.spiffeId,
            action: 'swarm.message.receive',
            resource: 'swarm.backbone',
            purpose: 'COMMUNICATION',
            correlationId: message.correlationId,
            outcome: 'ALLOW',
            severity: 'LOW',
            metadata: { messageId: message.id, senderSpiffeId: message.senderSpiffeId }
          });

          await handler(message);
        } catch (error) {
          // Log processing failure
          await securityAuditService.logEvent({
            actor: recipient.spiffeId,
            action: 'swarm.message.process_failed',
            resource: 'swarm.backbone',
            purpose: 'COMMUNICATION',
            correlationId: message.correlationId,
            outcome: 'ERROR',
            severity: 'HIGH',
            metadata: { messageId: message.id, error: error instanceof Error ? error.message : String(error) }
          });
        }
      }
    }
  }

  public clearAll(): void {
    this.subscriptions.clear();
  }
}

export const distributedAgentCommsService = DistributedAgentCommsService.getInstance();
