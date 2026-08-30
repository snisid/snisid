export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function assertNever(value: never): never {
  throw new Error(`Unexpected value: ${String(value)}`);
}

export function minimalResult<T extends Record<string, unknown>>(payload: T, allowedKeys: (keyof T)[]): Partial<T> {
  return Object.fromEntries(Object.entries(payload).filter(([key]) => allowedKeys.includes(key as keyof T))) as Partial<T>;
}
