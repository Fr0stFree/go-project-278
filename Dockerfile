# 1. Build frontend
FROM node:24-alpine AS frontend-builder

WORKDIR /build/frontend

RUN npm install @hexlet/project-url-shortener-frontend


# 2. Build backend
FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine AS backend-builder

RUN apk add --no-cache git

WORKDIR /build/backend

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /build/goose github.com/pressly/goose/v3/cmd/goose

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /build/shortener ./cmd/shortener


# 3. Runtime
FROM node:24-alpine

RUN apk add --no-cache \
    bash \
    ca-certificates \
    caddy

WORKDIR /app

RUN npm install concurrently

COPY --from=backend-builder \
    /build/shortener \
    ./bin/shortener

COPY --from=backend-builder \
    /build/goose \
    /usr/local/bin/goose

COPY --from=backend-builder \
    /build/backend/db/migrations \
    ./db/migrations

COPY --from=frontend-builder \
    /build/frontend/node_modules/@hexlet/project-url-shortener-frontend/dist \
    ./public

COPY Caddyfile /etc/caddy/Caddyfile
COPY bin/run.sh ./bin/run.sh

RUN chmod +x ./bin/run.sh

EXPOSE 80

CMD ["/app/bin/run.sh"]
