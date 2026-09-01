# 1. Build frontend
FROM node:24-alpine AS frontend-builder

WORKDIR /build/frontend

RUN npm install @hexlet/project-url-shortener-frontend


# 2. Build backend
FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine AS backend-builder

WORKDIR /build/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /build/shortener ./cmd/shortener


# 3. Runtime
FROM node:24-alpine

RUN apk add --no-cache \
    ca-certificates \
    caddy

WORKDIR /app

RUN npm install concurrently

COPY --from=backend-builder \
    /build/shortener \
    ./shortener

COPY --from=frontend-builder \
    /build/frontend/node_modules/@hexlet/project-url-shortener-frontend/dist \
    ./public

COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80

CMD ["npx", "concurrently", "--kill-others", "./shortener", "caddy run --config /etc/caddy/Caddyfile --adapter caddyfile"]