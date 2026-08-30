import { describe, expect, it } from 'vitest';
import { createSnisidMcpServer } from '../server/server.js';

describe('MCP server', () => {
  it('creates server without direct DB access', () => {
    const server = createSnisidMcpServer();
    expect(server).toBeTruthy();
  });
});
