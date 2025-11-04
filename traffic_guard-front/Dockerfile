FROM node:22.15.0-alpine AS builder

WORKDIR /app

COPY package*.json ./
# RUN npm ci

COPY . .
COPY ./components /app/components

RUN npm ci
RUN npm install -D sass-embedded
RUN npm run build

FROM node:22.15.0-alpine AS production

WORKDIR /app

COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package.json /app/package.json
COPY --from=builder /app/package-lock.json /app/package-lock.json
COPY --from=builder /app/components /app/components
COPY certs/fazenda.crt /app/certs/server.crt
COPY certs/fazenda.key /app/certs/server.key

RUN npm ci --omit=dev

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3001
ENV SSL_CRT_FILE=/app/certs/server.crt
ENV SSL_KEY_FILE=/app/certs/server.key

EXPOSE 3001

CMD ["node", ".output/server/index.mjs"]
