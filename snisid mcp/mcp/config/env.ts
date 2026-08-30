import 'dotenv/config';
import { z } from 'zod';

const emptyToUndefined = (value: unknown) => value === '' ? undefined : value;

const optionalString = () =>
  z.preprocess(emptyToUndefined, z.string().optional());

const optionalUrl = () =>
  z.preprocess(emptyToUndefined, z.string().url().optional());

const urlWithDefault = (fallback: string) =>
  z.preprocess(emptyToUndefined, z.string().url().default(fallback));

const boolWithDefault = (fallback: boolean) =>
  z.preprocess(emptyToUndefined, z.coerce.boolean().default(fallback));

const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'test', 'staging', 'production']).default('development'),
  PORT: z.coerce.number().int().positive().default(3001),
  MCP_TRANSPORT: z.enum(['stdio', 'http']).default('http'),
  SNISID_DEMO_MODE: boolWithDefault(false),

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
  QDRANT_URL: urlWithDefault('http://localhost:6333'),
  OTEL_EXPORTER_OTLP_ENDPOINT: optionalUrl(),
  LOG_LEVEL: z.enum(['error', 'warn', 'info', 'debug']).default('info'),

  OPENAI_API_KEY: optionalString(),
  OPENAI_BASE_URL: optionalUrl(),
  ANTHROPIC_API_KEY: optionalString(),
  ANTHROPIC_BASE_URL: optionalUrl(),
  CLAUDE_API_KEY: optionalString(),

  MISTRAL_API_KEY: optionalString(),
  MISTRAL_BASE_URL: urlWithDefault('https://api.mistral.ai/v1'),

  DEEPSEEK_API_KEY: optionalString(),
  DEEPSEEK_BASE_URL: urlWithDefault('https://api.deepseek.com/v1'),

  MINIMAX_API_KEY: optionalString(),
  MINIMAX_BASE_URL: urlWithDefault('https://api.minimax.io/v1'),

  NVIDIA_API_KEY: optionalString(),
  NVIDIA_BASE_URL: urlWithDefault('https://integrate.api.nvidia.com/v1'),

  SNISID_DEFAULT_AI_PROVIDER: z.enum(['mistral', 'deepseek', 'minimax', 'nvidia']).default('mistral'),
  SNISID_DEFAULT_AI_MODEL: z.string().default('mistral-small-latest'),

  SNISID_FALLBACK_AI_PROVIDER: z.enum(['mistral', 'deepseek', 'minimax', 'nvidia']).default('deepseek'),
  SNISID_FALLBACK_AI_MODEL: z.string().default('deepseek-chat'),

  SNISID_SECONDARY_AI_PROVIDER: z.enum(['mistral', 'deepseek', 'minimax', 'nvidia']).default('minimax'),
  SNISID_SECONDARY_AI_MODEL: z.string().default('MiniMax-M2.7'),

  SNISID_TERTIARY_AI_PROVIDER: z.enum(['mistral', 'deepseek', 'minimax', 'nvidia']).default('nvidia'),
  SNISID_TERTIARY_AI_MODEL: z.string().default('nvidia/nemotron-3-super-120b-a12b'),

  SNISID_AI_EXTERNAL_ALLOWED: boolWithDefault(false),
  SNISID_AI_SENSITIVE_DATA_ALLOWED: boolWithDefault(false),

  GOOGLE_ANTIGRAVITY_API_KEY: optionalString(),
  GOOGLE_ANTIGRAVITY_BASE_URL: optionalUrl(),
  ARENA_AI_API_KEY: optionalString(),
  ARENA_AI_BASE_URL: optionalUrl(),
  CURSOR_BRIDGE_URL: optionalUrl(),
  VSCODE_BRIDGE_URL: optionalUrl(),
  GITHUB_COPILOT_BRIDGE_URL: optionalUrl()
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

