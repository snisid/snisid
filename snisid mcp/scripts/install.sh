#!/usr/bin/env bash
set -euo pipefail

npm install \
  @modelcontextprotocol/sdk \
  express \
  axios \
  dotenv \
  jsonwebtoken \
  zod \
  winston \
  bcrypt \
  cors \
  helmet \
  express-rate-limit \
  ioredis \
  @qdrant/js-client-rest \
  @opentelemetry/api \
  @opentelemetry/sdk-node \
  @opentelemetry/auto-instrumentations-node \
  uuid

npm install -D \
  typescript \
  tsx \
  vitest \
  supertest \
  eslint \
  prettier \
  @types/node \
  @types/express \
  @types/cors \
  @types/jsonwebtoken \
  @types/bcrypt \
  @types/supertest
