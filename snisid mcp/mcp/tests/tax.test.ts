import { describe, expect, it } from 'vitest';
import { nifSchema } from '../validators/tax.validator.js';

describe('tax validators', () => {
  it('validates NIF', () => {
    expect(nifSchema.parse('NIF-123456')).toBe('NIF-123456');
  });
});
