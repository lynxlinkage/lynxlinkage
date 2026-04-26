syntax=docker/dockerfile:1.6

# ---- Stage 1: build the SvelteKit frontend ----
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* frontend/pnpm-lock.yaml* ./
RUN if [ -f pnpm-lock.yaml ]; then \
        corepack enable && corepack prepare pnpm@latest --activate && pnpm install --frozen-lockfile; \
    elif [ -f package-lock.json ]; then \
        npm ci; \
    else \
        npm install; \
    fi

COPY frontend/ .
RUN if [ -f pnpm-lock.yaml ]; then pnpm build; else npm run build; fi

# ---- Stage 2: build the Go backend with embedded frontend ----
FROM golang:1.22-alpine AS backend-builder
RUN apk add --no-cache make
WORKDIR /src

COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

COPY backend/ ./backend/
COPY --from=frontend-builder /app/frontend/build/ ./backend/internal/static/dist/

WORKDIR /src/backend
RUN CGO_ENABLED=0 go build -tags=embed -ldflags="-s -w" -o /out/lynxlinkage ./cmd/server

# ---- Stage 3: runtime ----
FROM alpine:3.20
RUN adduser -D -u 10001 app
WORKDIR /app
COPY --from=backend-builder /out/lynxlinkage /app/lynxlinkage

USER app
ENV APP_ENV=production HTTP_ADDR=:8080 \
    DATABASE_URL="file:/data/lynxlinkage.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/app/lynxlinkage"]
