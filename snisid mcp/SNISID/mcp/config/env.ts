import 'dotenv/config';
import { z } from 'zod';

const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'test', 'staging', 'production']).default('development'),
  PORT: z.coerce.number().int().positive().default(3001),
  MCP_TRANSPORT: z.enum(['stdio', 'http']).default('http'),
  SNISID_DEMO_MODE: z.coerce.boolean().default(false),
  JWT_ISSUER: z.string().min(3),
  JWT_AUDIENCE: z.string().min(3),
  JWT_SECRET: z.string().min(32),
  JWT_EXPIRES_IN: z.string().default('15m'),
  ENCRYPTION_KEY_B64: z.string().min(32),
  API_KEY_PEPPER: z.string().min(16),
  API_GATEWAY_BASE_URL: z.string().url(),
  ONI_API_BASE_URL: z.string().url(),
  DGI_API_BASE_URL: z.string().url(),
  PNH_API_BASE_URL: z.string().url(),
  MJSP_API_BASE_URL: z.string().url(),
  IMMIGRATION_API_BASE_URL: z.string().url(),
  ANH_API_BASE_URL: z.string().url(),
  PASSPORT_API_BASE_URL: z.string().url(),
  BIOMETRIC_API_BASE_URL: z.string().url(),
  EDUCATION_API_BASE_URL: z.string().url(),
  INTELLIGENCE_API_BASE_URL: z.string().url(),
  GOV_API_KEY: z.string().min(8),
  REDIS_URL: z.string().default('redis://localhost:6379'),
  QDRANT_URL: z.string().url().default('http://localhost:6333'),
  OTEL_EXPORTER_OTLP_ENDPOINT: z.string().url().optional(),
  LOG_LEVEL: z.enum(['error', 'warn', 'info', 'debug']).default('info'),
  OPENAI_API_KEY: z.string().optional(),
  OPENAI_BASE_URL: z.string().url().default('https://api.openai.com/v1'),
  DEEPSEEK_API_KEY: z.string().optional(),
  DEEPSEEK_BASE_URL: z.string().url().default('https://api.deepseek.com/v1'),
  ANTHROPIC_API_KEY: z.string().optional(),
  ANTHROPIC_BASE_URL: z.string().url().default('https://api.anthropic.com'),
  GOOGLE_ANTIGRAVITY_API_KEY: z.string().optional(),
  GOOGLE_ANTIGRAVITY_BASE_URL: z.string().optional(),
  ARENA_AI_API_KEY: z.string().optional(),
  ARENA_AI_BASE_URL: z.string().optional(),
  CLAUDE_API_KEY: z.string().optional(),
  CURSOR_BRIDGE_URL: z.string().url().optional(),
  VSCODE_BRIDGE_URL: z.string().url().optional(),
  GITHUB_COPILOT_BRIDGE_URL: z.string().url().optional()
});

const parsed = envSchema.safeParse(process.env);
if (!parsed.success) {
  console.error('Invalid SNISID environment configuration', parsed.error.flatten().fieldErrors);
  throw new Error('SNISID_ENV_VALIDATION_FAILED');
}

if (parsed.data.NODE_ENV === 'production') {
  for (const [key, value] of Object.entries(parsed.data)) {
    if (typeof value === 'string' && value.includes('REPLACE_WITH')) {
      throw new Error(`Production secret/config placeholder not replaced: ${key}`);
    }
  }
}

export const env = parsed.data;
export type Env = typeof env;
