# ---- Build Stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack for yarn/pnpm if needed, using npm here
COPY package.json package-lock.json* ./
RUN npm install --prefer-offline --no-audit

COPY . .
RUN npm run build

# ---- Production Stage ----
FROM nginx:alpine

# Copy build artifacts
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
