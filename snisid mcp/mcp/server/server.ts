import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { SNISID } from '../config/constants.js';
import { registerMcpSurface } from './registry.js';

export function createSnisidMcpServer(): McpServer {
  const server = new McpServer({
    name: SNISID.mcpName,
    version: SNISID.version
  });
  registerMcpSurface(server);
  return server;
}
