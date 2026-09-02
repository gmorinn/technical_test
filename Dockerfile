# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# ── Dev stage ─────────────────────────────────────────────────────────────────
FROM golang:1.25-bookworm AS dev

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@v1.67.0

WORKDIR /app

CMD ["air", "-c", "air.toml"]

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates

RUN addgroup -S api && adduser -S -G api api

WORKDIR /app

COPY --from=builder --chown=api:api /out/api ./api

COPY --from=builder --chown=api:api /build/internal/db/migrations ./internal/db/migrations

RUN mkdir -p ./public/uploads && chown -R api:api ./public

USER api

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- "http://localhost:${API_PORT:-8080}/health" >/dev/null || exit 1

CMD ["./api"]
