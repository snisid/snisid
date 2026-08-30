import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['mcp/tests/**/*.test.ts'],
    exclude: [
      'node_modules/**',
      'dist/**',
      'SNISID/**',
      'coverage/**',
      'tmp/**',
      'uploads/**'
    ],
    environment: 'node'
  }
});
