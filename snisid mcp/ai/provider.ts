export type AiRole = 'system' | 'user' | 'assistant';
export interface AiMessage { role: AiRole; content: string }
export interface AiCompletionRequest {
  provider?: string;
  model?: string;
  messages: AiMessage[];
  correlationId: string;
  maxTokens?: number;
}
export interface AiProvider {
  name: string;
  available(): boolean;
  complete(request: AiCompletionRequest): Promise<string>;
}
