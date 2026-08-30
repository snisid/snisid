import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { investigationPrompt } from './investigation.prompt.js';
import { identityPrompt } from './identity.prompt.js';
import { securityPrompt } from './security.prompt.js';
import { judicialPrompt } from './judicial.prompt.js';
import { intelligencePrompt } from './intelligence.prompt.js';

const prompts = { investigationPrompt, identityPrompt, securityPrompt, judicialPrompt, intelligencePrompt };

export function registerAllPrompts(server: McpServer): void {
  for (const [name, prompt] of Object.entries(prompts)) {
    server.prompt(name, prompt.title, {}, async () => ({
      messages: [{ role: 'user' as const, content: { type: 'text' as const, text: prompt.system } }]
    }));
  }
}
