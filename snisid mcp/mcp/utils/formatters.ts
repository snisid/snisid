export function formatMcpText(payload: unknown): string {
  return JSON.stringify(payload, null, 2);
}

export function normalizeNationalId(value: string): string {
  return value.replace(/[^A-Z0-9-]/gi, '').toUpperCase();
}
