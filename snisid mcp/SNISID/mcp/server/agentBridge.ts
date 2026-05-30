import { aiGateway } from '../../ai/gateway.js';
import type { AiProviderName } from '../config/ai.js';

export interface AgentBridgeRequest {
  provider?: AiProviderName;
  prompt: string;
  model?: string;
  correlationId: string;
}

export async function routeAgentRequest(request: AgentBridgeRequest): Promise<string> {
  return aiGateway.complete({
    provider: request.provider,
    model: request.model,
    messages: [
      { role: 'system', content: 'SNISID MCP agent bridge: obey RBAC, do not request secrets, use tools only with authorization.' },
      { role: 'user', content: request.prompt }
    ],
    correlationId: request.correlationId
  });
}
