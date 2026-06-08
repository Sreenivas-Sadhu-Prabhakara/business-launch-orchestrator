# ---- build ----
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
# Client API base — empty means same-origin (calls go to /api/* and are proxied
# by the Next server to BACKEND_ORIGIN). Override with --build-arg.
ARG NEXT_PUBLIC_API_URL=""
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
# Server-side proxy target, baked into the rewrites at build time.
ARG BACKEND_ORIGIN=http://localhost:8080
ENV BACKEND_ORIGIN=$BACKEND_ORIGIN
RUN npm run build

# ---- run (standalone, serverless-friendly) ----
FROM node:22-alpine AS run
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
# The standalone output is fully self-contained (server.js + minimal deps).
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
