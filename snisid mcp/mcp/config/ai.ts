import { env } from './env.js';

export const aiProviders = {
  mistral: {
    baseUrl: env.MISTRAL_BASE_URL,
    apiKey: env.MISTRAL_API_KEY,
    defaultModel:
      env.SNISID_DEFAULT_AI_PROVIDER === 'mistral'
        ? env.SNISID_DEFAULT_AI_MODEL
        : 'mistral-small-latest'
  },
  deepseek: {
    baseUrl: env.DEEPSEEK_BASE_URL,
    apiKey: env.DEEPSEEK_API_KEY,
    defaultModel:
      env.SNISID_DEFAULT_AI_PROVIDER === 'deepseek'
        ? env.SNISID_DEFAULT_AI_MODEL
        : 'deepseek-chat'
  },
  minimax: {
    baseUrl: env.MINIMAX_BASE_URL,
    apiKey: env.MINIMAX_API_KEY,
    defaultModel:
      env.SNISID_SECONDARY_AI_PROVIDER === 'minimax'
        ? env.SNISID_SECONDARY_AI_MODEL
        : 'minimax-m2.5'
  },
  nvidia: {
    baseUrl: env.NVIDIA_BASE_URL,
    apiKey: env.NVIDIA_API_KEY,
    defaultModel:
      env.SNISID_TERTIARY_AI_PROVIDER === 'nvidia'
        ? env.SNISID_TERTIARY_AI_MODEL
        : 'nvidia/nemotron-3-super-120b-a12b'
  }
} as const;

export type AiProviderName = keyof typeof aiProviders;

export const aiPolicy = {
  defaultProvider: env.SNISID_DEFAULT_AI_PROVIDER,
  defaultModel: env.SNISID_DEFAULT_AI_MODEL,
  fallbackProvider: env.SNISID_FALLBACK_AI_PROVIDER,
  fallbackModel: env.SNISID_FALLBACK_AI_MODEL,
  secondaryProvider: env.SNISID_SECONDARY_AI_PROVIDER,
  secondaryModel: env.SNISID_SECONDARY_AI_MODEL,
  tertiaryProvider: env.SNISID_TERTIARY_AI_PROVIDER,
  tertiaryModel: env.SNISID_TERTIARY_AI_MODEL,
  externalAllowed: env.SNISID_AI_EXTERNAL_ALLOWED,
  sensitiveDataAllowed: env.SNISID_AI_SENSITIVE_DATA_ALLOWED
} as const;
