# Stage 1: Build
FROM node:22-bookworm-slim AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=optional
COPY . .
RUN npx tsc -p tsconfig.json

# Stage 2: Production runtime
FROM node:22-bookworm-slim
WORKDIR /app
ENV NODE_ENV=production
RUN groupadd -r snisid && useradd -r -g snisid -d /app -s /sbin/nologin snisid
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/mcp/resources ./mcp/resources
COPY --from=builder /app/mcp/prompts ./mcp/prompts
EXPOSE 3001
USER snisid
CMD ["node", "dist/mcp/server/index.js", "--http"]