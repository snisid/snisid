FROM node:22-bookworm-slim

WORKDIR /app

ENV NODE_ENV=development

COPY package*.json ./

RUN npm install --omit=optional

COPY . .

EXPOSE 3001
EXPOSE 4000

CMD ["npm", "run", "mcp:http"]
