import { describe, expect, it } from 'vitest';
import { nationalIdSchema } from '../validators/identity.validator.js';

describe('identity validators', () => {
  it('accepts safe national ids', () => {
    expect(nationalIdSchema.parse('HT-12345')).toBe('HT-12345');
  });
  it('rejects injection characters', () => {
    expect(() => nationalIdSchema.parse('../etc/passwd')).toThrow();
  });
});
