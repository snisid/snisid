import { env } from './env.js';

export const aiProviders = {
  openai: { baseUrl: env.OPENAI_BASE_URL, apiKey: env.OPENAI_API_KEY },
  deepseek: { baseUrl: env.DEEPSEEK_BASE_URL, apiKey: env.DEEPSEEK_API_KEY },
  anthropic: { baseUrl: env.ANTHROPIC_BASE_URL, apiKey: env.ANTHROPIC_API_KEY ?? env.CLAUDE_API_KEY },
  googleAntigravity: { baseUrl: env.GOOGLE_ANTIGRAVITY_BASE_URL, apiKey: env.GOOGLE_ANTIGRAVITY_API_KEY },
  arena: { baseUrl: env.ARENA_AI_BASE_URL, apiKey: env.ARENA_AI_API_KEY },
  cursor: { baseUrl: env.CURSOR_BRIDGE_URL },
  vscode: { baseUrl: env.VSCODE_BRIDGE_URL },
  githubCopilot: { baseUrl: env.GITHUB_COPILOT_BRIDGE_URL }
} as const;

export type AiProviderName = keyof typeof aiProviders;
