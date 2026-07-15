import { randomUUID } from 'node:crypto';

export interface Agent {
  id: string;
  name: string;
  spiffeId: string;
  role: string;
  permissions: string[];
  status: 'ACTIVE' | 'INACTIVE' | 'QUARANTINED';
  registeredAt: string;
  metadata?: Record<string, unknown>;
}

export class AgentRegistryService {
  private static instance: AgentRegistryService;
  private agents = new Map<string, Agent>();

  private constructor() {
    // Seed default agents according to the SNISID SOC Agent Swarm Architecture
    this.seedDefaultAgents();
  }

  public static getInstance(): AgentRegistryService {
    if (!AgentRegistryService.instance) {
      AgentRegistryService.instance = new AgentRegistryService();
    }
    return AgentRegistryService.instance;
  }

  private seedDefaultAgents(): void {
    const defaultAgents: Omit<Agent, 'id' | 'registeredAt'>[] = [
      {
        name: 'SOC Hunter Agent',
        spiffeId: 'spiffe://snisid.gov.ht/ns/soc/sa/hunter-01',
        role: 'HUNTER',
        permissions: ['police:read', 'intelligence:analyze'],
        status: 'ACTIVE',
        metadata: { objective: 'Proactive Search for anomalies' }
      },
      {
        name: 'SOC Investigator Agent',
        spiffeId: 'spiffe://snisid.gov.ht/ns/soc/sa/investigator-01',
        role: 'INVESTIGATOR',
        permissions: ['police:read', 'justice:read', 'identity:verify'],
        status: 'ACTIVE',
        metadata: { objective: 'Forensic Context and evidence gathering' }
      },
      {
        name: 'SOC Response Agent',
        spiffeId: 'spiffe://snisid.gov.ht/ns/soc/sa/responder-01',
        role: 'RESPONDER',
        permissions: ['security:admin'],
        status: 'ACTIVE',
        metadata: { objective: 'Active containment and rollback actions' }
      }
    ];

    for (const agent of defaultAgents) {
      void this.registerAgent(agent);
    }
  }

  public async registerAgent(agentInput: Omit<Agent, 'id' | 'registeredAt'>): Promise<Agent> {
    if (!agentInput.spiffeId.startsWith('spiffe://')) {
      throw new Error('INVALID_SPIFFE_ID: Must start with spiffe://');
    }
    // Check for duplicate spiffeId
    for (const existingAgent of this.agents.values()) {
      if (existingAgent.spiffeId === agentInput.spiffeId) {
        throw new Error(`DUPLICATE_SPIFFE_ID: Agent with SPIFFE ID ${agentInput.spiffeId} is already registered`);
      }
    }

    const agent: Agent = {
      id: randomUUID(),
      registeredAt: new Date().toISOString(),
      ...agentInput
    };
    this.agents.set(agent.id, agent);
    return agent;
  }

  public async getAgent(id: string): Promise<Agent | undefined> {
    return this.agents.get(id);
  }

  public async getAgentBySpiffeId(spiffeId: string): Promise<Agent | undefined> {
    for (const agent of this.agents.values()) {
      if (agent.spiffeId === spiffeId) {
        return agent;
      }
    }
    return undefined;
  }

  public async updateAgentStatus(id: string, status: 'ACTIVE' | 'INACTIVE' | 'QUARANTINED'): Promise<Agent> {
    const agent = this.agents.get(id);
    if (!agent) {
      throw new Error(`AGENT_NOT_FOUND: Agent with ID ${id} not found`);
    }
    const updated = { ...agent, status };
    this.agents.set(id, updated);
    return updated;
  }

  public async listAgents(): Promise<Agent[]> {
    return Array.from(this.agents.values());
  }

  public async verifyAgentIdentity(spiffeId: string): Promise<boolean> {
    const agent = await this.getAgentBySpiffeId(spiffeId);
    return agent !== undefined && agent.status === 'ACTIVE';
  }

  public clearAll(): void {
    this.agents.clear();
    this.seedDefaultAgents();
  }
}

export const agentRegistryService = AgentRegistryService.getInstance();
