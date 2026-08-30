import type { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z, type ZodRawShape } from 'zod';
import type { Permission } from '../../config/permissions.js';
import type { SecurityContext } from '../../types/security.types.js';
import { authenticateAndAuthorize, authContextSchema } from '../../security/auth.js';
import { writeAuditEvent } from '../../audit/auditLogger.js';
import { formatMcpText } from '../../utils/formatters.js';
import { toolCounter } from '../../utils/monitoring.js';
import { redact } from '../../utils/logger.js';

export interface GovernmentToolDefinition<T extends ZodRawShape> {
  name: string;
  description: string;
  permission: Permission;
  inputShape: T;
  handler: (input: z.infer<z.ZodObject<T>>, ctx: SecurityContext) => Promise<unknown>;
}

export function registerGovernmentTool<T extends ZodRawShape>(server: McpServer, def: GovernmentToolDefinition<T>): void {
  const fullShape = { ...def.inputShape, auth: authContextSchema };
  const parser = z.object(fullShape);

  (server.tool as any)(def.name, def.description, fullShape, async (rawInput: unknown) => {
    const input = parser.parse(rawInput) as any;
    let ctx: SecurityContext | undefined;
    try {
      ctx = await authenticateAndAuthorize(input.auth, def.permission);
      await writeAuditEvent({
        actor: ctx.principal.subject,
        action: `tool.call.${def.name}`,
        resource: def.permission,
        purpose: ctx.purpose,
        correlationId: ctx.correlationId,
        outcome: 'ALLOW',
        severity: 'LOW',
        metadata: { input: redact(input) }
      });
      toolCounter.add(1, { tool: def.name, permission: def.permission });
      const { auth: _auth, ...businessInput } = input;
      const result = await def.handler(businessInput as z.infer<z.ZodObject<T>>, ctx);
      await writeAuditEvent({
        actor: ctx.principal.subject,
        action: `tool.success.${def.name}`,
        resource: def.permission,
        purpose: ctx.purpose,
        correlationId: ctx.correlationId,
        outcome: 'ALLOW',
        severity: 'LOW',
        metadata: { resultSummary: redact(result) }
      });
      return { content: [{ type: 'text' as const, text: formatMcpText({ ok: true, data: result }) }] };
    } catch (error) {
      await writeAuditEvent({
        actor: ctx?.principal.subject ?? 'anonymous',
        action: `tool.denied_or_error.${def.name}`,
        resource: def.permission,
        purpose: ctx?.purpose ?? input.auth?.purpose ?? 'UNSPECIFIED',
        correlationId: ctx?.correlationId ?? input.auth?.correlationId ?? 'UNSPECIFIED',
        outcome: error instanceof Error && /FORBIDDEN|MFA|RISK|PURPOSE|SESSION|TOKEN/.test(error.message) ? 'DENY' : 'ERROR',
        severity: 'HIGH',
        metadata: { error: error instanceof Error ? error.message : String(error) }
      });
      throw error;
    }
  });
}
