import { randomUUID } from 'node:crypto';

interface ClientConnection {
  id: string;
  transport: 'stdio' | 'http';
  createdAt: string;
  lastSeenAt: string;
  principal?: string;
}

const connections = new Map<string, ClientConnection>();

export function registerClientConnection(transport: 'stdio' | 'http', principal?: string): ClientConnection {
  const connection: ClientConnection = { id: randomUUID(), transport, principal, createdAt: new Date().toISOString(), lastSeenAt: new Date().toISOString() };
  connections.set(connection.id, connection);
  return connection;
}

export function touchClientConnection(id: string): void {
  const c = connections.get(id);
  if (c) c.lastSeenAt = new Date().toISOString();
}

export function listClientConnections(): ClientConnection[] {
  return [...connections.values()];
}
