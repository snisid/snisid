import axios from 'axios';
import { aiPolicy, aiProviders } from '../mcp/config/ai.js';
import type { AiCompletionRequest, AiProvider } from './provider.js';
import { AiRouter } from './router.js';
import { sanitizeAgentPrompt, systemGuardrail } from './orchestrator.js';

function containsSensitiveSnisidData(text: string): boolean {
  return /biometric|fingerprint|iris|face-template|criminal|warrant|detention|taxpayer|nif|passport|watchlist|gang|weapon|nationalId|HT-NID|HT-PASSPORT/i.test(text);
}

function assertAiPolicy(request: AiCompletionRequest): void {
  const joined = request.messages.map((m) => m.content).join('\n');

  if (!aiPolicy.externalAllowed && process.env['NODE_ENV'] === 'production') {
    throw new Error('EXTERNAL_AI_DISABLED_BY_POLICY');
  }

  if (!aiPolicy.sensitiveDataAllowed && containsSensitiveSnisidData(joined)) {
    throw new Error('SENSITIVE_SNISID_DATA_NOT_ALLOWED_FOR_EXTERNAL_AI');
  }
}

class OpenAiCompatibleProvider implements AiProvider {
  constructor(
    public readonly name: string,
    private readonly baseUrl?: string,
    private readonly apiKey?: string,
    private readonly defaultModel?: string
  ) {}

  available(): boolean {
    return Boolean(this.baseUrl && this.apiKey);
  }

  async complete(request: AiCompletionRequest): Promise<string> {
    assertAiPolicy(request);

    const messages = [
      { role: 'system', content: systemGuardrail() },
      ...request.messages.map((m) => ({
        ...m,
        content: sanitizeAgentPrompt(m.content)
      }))
    ];

    const response = await axios.post(
      `${this.baseUrl}/chat/completions`,
      {
        model: request.model ?? this.defaultModel ?? aiPolicy.defaultModel,
        messages,
        max_tokens: request.maxTokens ?? 1024
      },
      {
        headers: {
          authorization: `Bearer ${this.apiKey}`,
          'x-correlation-id': request.correlationId
        },
        timeout: 30_000
      }
    );

    return response.data?.choices?.[0]?.message?.content ?? '';
  }
}

const providers: AiProvider[] = [
  new OpenAiCompatibleProvider(
    'mistral',
    aiProviders.mistral.baseUrl,
    aiProviders.mistral.apiKey,
    aiProviders.mistral.defaultModel
  ),
  new OpenAiCompatibleProvider(
    'deepseek',
    aiProviders.deepseek.baseUrl,
    aiProviders.deepseek.apiKey,
    aiProviders.deepseek.defaultModel
  ),
  new OpenAiCompatibleProvider(
    'minimax',
    aiProviders.minimax.baseUrl,
    aiProviders.minimax.apiKey,
    aiProviders.minimax.defaultModel
  ),
  new OpenAiCompatibleProvider(
    'nvidia',
    aiProviders.nvidia.baseUrl,
    aiProviders.nvidia.apiKey,
    aiProviders.nvidia.defaultModel
  )
];

export const aiGateway = new AiRouter(providers);
