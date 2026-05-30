# ---- Build Stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack for yarn/pnpm if needed, using npm here
COPY package.json package-lock.json* ./
RUN npm ci --prefer-offline --no-audit

COPY . .
RUN npm run build

# ---- Production Stage ----
# NGINX unprivileged image running as non-root (UID 101)
FROM nginxinc/nginx-unprivileged:alpine

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy build artifacts
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port (must be >1024 for non-root)
EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
