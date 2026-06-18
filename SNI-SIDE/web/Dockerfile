FROM node:20-alpine AS build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN adduser -D -u 1001 sniside && chown -R sniside:sniside /usr/share/nginx/html
USER sniside

EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
