export function sanitizeAgentPrompt(prompt: string): string {
  return prompt
    .replace(/ignore previous instructions/gi, '[blocked-injection]')
    .replace(/reveal.*(secret|token|key)/gi, '[blocked-secret-request]')
    .slice(0, 8000);
}

export function systemGuardrail(): string {
  return [
    'You are operating in SNISID sovereign AI environment.',
    'Never bypass RBAC, MFA, audit, device trust or legal purpose requirements.',
    'Do not infer guilt or make coercive decisions; provide decision support only.',
    'Minimize personal data and explain uncertainty.'
  ].join('\n');
}
