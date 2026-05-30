const usage = new Map<string, { tokens: number; resetAt: number }>();

export function reserveTokens(provider: string, requested: number, limit = 200_000): boolean {
  const now = Date.now();
  const current = usage.get(provider);
  if (!current || current.resetAt < now) {
    usage.set(provider, { tokens: requested, resetAt: now + 60_000 });
    return true;
  }
  if (current.tokens + requested > limit) return false;
  current.tokens += requested;
  return true;
}
