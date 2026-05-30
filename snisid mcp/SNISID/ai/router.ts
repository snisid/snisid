import type { AiProvider, AiCompletionRequest } from './provider.js';
import { reserveTokens } from './tokenManager.js';

export class AiRouter {
  constructor(private readonly providers: AiProvider[]) {}

  select(request: AiCompletionRequest): AiProvider {
    const ordered = request.provider
      ? this.providers.filter((p) => p.name === request.provider)
      : this.providers;
    const provider = ordered.find((p) => p.available() && reserveTokens(p.name, request.maxTokens ?? 2048));
    if (!provider) throw new Error('NO_AI_PROVIDER_AVAILABLE');
    return provider;
  }

  async completeWithFailover(request: AiCompletionRequest): Promise<string> {
    const preferred = request.provider ? [this.select(request)] : this.providers.filter((p) => p.available());
    let lastError: unknown;
    for (const provider of preferred) {
      try {
        return await provider.complete(request);
      } catch (error) {
        lastError = error;
      }
    }
    throw lastError instanceof Error ? lastError : new Error('AI_FAILOVER_EXHAUSTED');
  }
}
