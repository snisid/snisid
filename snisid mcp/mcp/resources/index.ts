import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { nationalSchemas } from './nationalSchemas.js';
import { judicialProcedures } from './judicialProcedures.js';
import { laws } from './laws.js';
import { governmentPolicies } from './governmentPolicies.js';
import { securityProtocols } from './securityProtocols.js';

const resources = {
  'snisid://schemas/national': nationalSchemas,
  'snisid://procedures/judicial': judicialProcedures,
  'snisid://laws/core': laws,
  'snisid://policies/government': governmentPolicies,
  'snisid://protocols/security': securityProtocols
};

export function registerAllResources(server: McpServer): void {
  for (const [uri, data] of Object.entries(resources)) {
    server.resource(uri, uri, async () => ({
      contents: [{ uri, mimeType: 'application/json', text: JSON.stringify(data, null, 2) }]
    }));
  }
}
