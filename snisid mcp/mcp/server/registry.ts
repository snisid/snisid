import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { registerAllTools } from '../tools/index.js';
import { registerAllPrompts } from '../prompts/index.js';
import { registerAllResources } from '../resources/index.js';

export function registerMcpSurface(server: McpServer): void {
  registerAllTools(server);
  registerAllPrompts(server);
  registerAllResources(server);
}
