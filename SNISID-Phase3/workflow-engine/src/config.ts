/**
 * SNISID — Centralised env-config (12-factor compliant)
 */
import { z } from 'zod';

const Env = z.object({
  NODE_ENV: z.enum(['development','staging','production']).default('production'),

  // Zeebe / Camunda 8
  ZEEBE_GRPC_ADDRESS: z.string().default('zeebe-gateway:26500'),
  ZEEBE_CLIENT_ID: z.string().optional(),
  ZEEBE_CLIENT_SECRET: z.string().optional(),
  ZEEBE_TLS_CA: z.string().default('/etc/snisid/tls/ca.crt'),

  // Temporal
  TEMPORAL_ADDRESS: z.string().default('temporal-frontend:7233'),
  TEMPORAL_NAMESPACE: z.string().default('snisid-prod'),
  TEMPORAL_TASK_QUEUE: z.string().default('snisid-default'),
  TEMPORAL_TLS_CA: z.string().default('/etc/snisid/tls/ca.crt'),

  // Kafka
  KAFKA_BROKERS: z.string().default('kafka-1:9093,kafka-2:9093,kafka-3:9093'),
  KAFKA_CLIENT_ID: z.string().default('snisid-workflow-engine'),
  SCHEMA_REGISTRY_URL: z.string().default('https://schema-registry.snisid.ht'),

  // PKI / TSA
  PKI_SIGN_ENDPOINT: z.string().default('https://pki.snisid.ht/sign'),
  PKI_TSA_ENDPOINT:  z.string().default('https://tsa.snisid.ht/rfc3161'),

  // OTel
  OTEL_EXPORTER_OTLP_ENDPOINT: z.string().default('http://otel-collector:4318'),

  // Workflow Engine
  ENGINE_REGION: z.string().default('DC1-PAP'),
  ENGINE_VERSION: z.string().default(process.env.npm_package_version ?? '1.0.0')
});

export const config = Env.parse(process.env);
export type AppConfig = typeof config;
