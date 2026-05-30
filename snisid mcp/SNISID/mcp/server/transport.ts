import express from 'express';
import cors from 'cors';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { StreamableHTTPServerTransport } from '@modelcontextprotocol/sdk/server/streamableHttp.js';
import { createSnisidMcpServer } from './server.js';
import { env } from '../config/env.js';
import { logger } from '../utils/logger.js';
import { securityHeaders } from '../middleware/securityHeaders.js';
import { rateLimitMiddleware } from '../middleware/rateLimit.js';
import { requestValidator } from '../middleware/requestValidator.js';
import { auditMiddleware } from '../middleware/auditMiddleware.js';
import { threatDetectionMiddleware } from '../middleware/threatDetection.js';
import { errorHandler } from '../middleware/errorHandler.js';
import { registerClientConnection } from './clientConnections.js';

export async function startStdioTransport(): Promise<void> {
  const server = createSnisidMcpServer();
  const transport = new StdioServerTransport();
  registerClientConnection('stdio');
  await server.connect(transport);
  logger.info('SNISID MCP started on STDIO');
}

export async function startHttpTransport(): Promise<void> {
  const app = express();
  app.disable('x-powered-by');
  app.use(securityHeaders);
  app.use(cors({ origin: false }));
  app.use(express.json({ limit: '1mb' }));
  app.use(rateLimitMiddleware);
  app.get('/healthz', (_req, res) => res.json({ ok: true, service: 'snisid-mcp' }));
  app.all('/mcp', requestValidator, auditMiddleware, threatDetectionMiddleware, async (req, res) => {
    const server = createSnisidMcpServer();
    const transport = new StreamableHTTPServerTransport({ sessionIdGenerator: undefined, enableJsonResponse: true });
    const connection = registerClientConnection('http');
    res.on('close', () => {
      void transport.close();
      void server.close();
    });
    try {
      await server.connect(transport);
      await transport.handleRequest(req, res, req.body);
    } catch (error) {
      logger.error('mcp_http_transport_error', { connectionId: connection.id, error: error instanceof Error ? error.message : String(error) });
      if (!res.headersSent) res.status(500).json({ jsonrpc: '2.0', error: { code: -32603, message: 'Internal server error' }, id: null });
    }
  });
  app.use(errorHandler);
  app.listen(env.PORT, () => logger.info('SNISID MCP HTTP listening', { port: env.PORT, path: '/mcp' }));
}
