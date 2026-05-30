import { ROLE_PERMISSIONS, type Role } from '../config/roles.js';
import type { Permission } from '../config/permissions.js';

export function permissionsForRoles(roles: Role[]): Permission[] {
  return [...new Set(roles.flatMap((role) => ROLE_PERMISSIONS[role] ?? []))];
}

export function hasPermission(roles: Role[], explicit: Permission[], required: Permission): boolean {
  return explicit.includes(required) || permissionsForRoles(roles).includes(required);
}
