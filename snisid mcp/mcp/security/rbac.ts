import type { Permission } from '../config/permissions.js';
import type { AuthenticatedPrincipal } from '../types/security.types.js';
import { hasPermission } from './permissions.js';

export class AuthorizationError extends Error {
  constructor(message = 'FORBIDDEN') {
    super(message);
    this.name = 'AuthorizationError';
  }
}

export function requirePermission(principal: AuthenticatedPrincipal, permission: Permission): void {
  if (!hasPermission(principal.roles, principal.permissions, permission)) {
    throw new AuthorizationError(`Missing permission: ${permission}`);
  }
}
