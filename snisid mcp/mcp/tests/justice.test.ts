import { describe, expect, it } from 'vitest';
import { caseIdSchema } from '../validators/justice.validator.js';

describe('justice validators', () => {
  it('validates case ids', () => {
    expect(caseIdSchema.parse('CASE-2026-001')).toBe('CASE-2026-001');
  });
});
