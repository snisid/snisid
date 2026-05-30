import { metrics, trace } from '@opentelemetry/api';

export const tracer = trace.getTracer('snisid-mcp');
export const meter = metrics.getMeter('snisid-mcp');
export const toolCounter = meter.createCounter('snisid_mcp_tool_calls');
export const securityCounter = meter.createCounter('snisid_mcp_security_events');

export async function withSpan<T>(name: string, fn: () => Promise<T>): Promise<T> {
  return tracer.startActiveSpan(name, async (span) => {
    try {
      return await fn();
    } catch (error) {
      span.recordException(error as Error);
      throw error;
    } finally {
      span.end();
    }
  });
}
