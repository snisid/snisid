import { env } from '../config/env.js';
import { startHttpTransport, startStdioTransport } from './transport.js';
import { logger } from '../utils/logger.js';

async function main(): Promise<void> {
  const args = new Set(process.argv.slice(2));
  const transport = args.has('--stdio') ? 'stdio' : args.has('--http') ? 'http' : env.MCP_TRANSPORT;
  if (transport === 'stdio') await startStdioTransport();
  else await startHttpTransport();
}

main().catch((error) => {
  logger.error('SNISID MCP fatal startup error', { error: error instanceof Error ? error.message : String(error) });
  process.exit(1);
});
