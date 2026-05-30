import axios from 'axios';
import { aiProviders } from '../mcp/config/ai.js';
import type { AiCompletionRequest, AiProvider } from './provider.js';
import { AiRouter } from './router.js';
import { sanitizeAgentPrompt, systemGuardrail } from './orchestrator.js';

class OpenAiCompatibleProvider implements AiProvider {
  constructor(public readonly name: string, private readonly baseUrl?: string, private readonly apiKey?: string) {}
  available(): boolean { return Boolean(this.baseUrl && this.apiKey); }
  async complete(request: AiCompletionRequest): Promise<string> {
    const messages = [{ role: 'system', content: systemGuardrail() }, ...request.messages.map((m) => ({ ...m, content: sanitizeAgentPrompt(m.content) }))];
    const response = await axios.post(`${this.baseUrl}/chat/completions`, {
      model: request.model ?? 'default', messages, max_tokens: request.maxTokens ?? 1024
    }, { headers: { authorization: `Bearer ${this.apiKey}`, 'x-correlation-id': request.correlationId }, timeout: 30_000 });
    return response.data?.choices?.[0]?.message?.content ?? '';
  }
}

class BridgeProvider implements AiProvider {
  constructor(public readonly name: string, private readonly baseUrl?: string) {}
  available(): boolean { return Boolean(this.baseUrl); }
  async complete(request: AiCompletionRequest): Promise<string> {
    const response = await axios.post(`${this.baseUrl}/complete`, request, { timeout: 30_000 });
    return response.data?.text ?? '';
  }
}

const providers: AiProvider[] = [
  new OpenAiCompatibleProvider('openai', aiProviders.openai.baseUrl, aiProviders.openai.apiKey),
  new OpenAiCompatibleProvider('deepseek', aiProviders.deepseek.baseUrl, aiProviders.deepseek.apiKey),
  new OpenAiCompatibleProvider('anthropic', aiProviders.anthropic.baseUrl, aiProviders.anthropic.apiKey),
  new OpenAiCompatibleProvider('googleAntigravity', aiProviders.googleAntigravity.baseUrl, aiProviders.googleAntigravity.apiKey),
  new OpenAiCompatibleProvider('arena', aiProviders.arena.baseUrl, aiProviders.arena.apiKey),
  new BridgeProvider('cursor', aiProviders.cursor.baseUrl),
  new BridgeProvider('vscode', aiProviders.vscode.baseUrl),
  new BridgeProvider('githubCopilot', aiProviders.githubCopilot.baseUrl)
];

export const aiGateway = new AiRouter(providers);
