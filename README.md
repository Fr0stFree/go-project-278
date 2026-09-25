# URL Shortener

URL Shortener is a Go web service for creating, managing, and tracking shortened links. It exposes a JSON API, redirects short URLs to their original targets, and stores links and redirect visits in PostgreSQL.

[A live deployment](https://shortener-latest.onrender.com/) is available on Render.

## Features

- Creates, reads, updates, and soft-deletes shortened links.
- Generates a random 6-character alphanumeric short name when one is not provided.
- Supports custom short names with validation and uniqueness checks.
- Redirects `/r/:short_name` requests with `302 Found`.
- Records redirect visits with IP address, user agent, referrer, and response status.
- Supports inclusive range pagination, sorting, and `Content-Range` response headers.
- Uses PostgreSQL, GORM, and versioned Goose migrations.
- Includes a production Docker image that serves the Hexlet frontend through Caddy and proxies API requests to the Go backend.

## Pipeline Status

[![Actions Status](https://github.com/Fr0stFree/go-project-278/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/Fr0stFree/go-project-278/actions)
[![CI](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml/badge.svg?branch=master)](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml)
[![Test Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/Fr0stFree/go-project-278/master/.github/badges/coverage-badge.json)](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml)

## Requirements

- Go `1.26.3`
- PostgreSQL
- `make`
- Docker for integration tests and container builds

The lint and formatting targets use `golangci-lint` `v2.12.2`. Install it once with `make install-lint`.

## Local Setup

Clone the repository and create a local configuration file:

```bash
git clone git@github.com:Fr0stFree/go-project-278.git
cd go-project-278
cp .env.example .env
make install
```

Edit `DATABASE_URL` if your local PostgreSQL credentials differ, create the database, export the example configuration, and apply the migrations. The locally built binary does not run migrations automatically.

```bash
set -a
source .env
set +a
go tool goose -dir db/migrations postgres "$DATABASE_URL" up
```

Build and run the backend:

```bash
make start
```

`make start` builds the current source into `bin/shortener` and then runs it. The backend is available at `http://localhost:8080` with the example configuration.

For live reload, install [Air](https://github.com/air-verse/air) and run:

```bash
go install github.com/air-verse/air@latest
make dev
```

## Configuration

Configuration is read from environment variables. A local `.env` file is loaded automatically when present; process environment variables take precedence. Start with the committed [`.env.example`](.env.example).

Durations use Go duration syntax such as `2s`, `5m`, or `12h`. `HTTP_MAX_BODY_SIZE` is specified in bytes. Comma-separated values are used for `HTTP_CORS_ALLOW_ORIGINS`.

| Variable | Description | Default |
| --- | --- | --- |
| `APP_BASE_URL` | Public base URL used to build `short_url` in API responses | `http://localhost:8080` |
| `HTTP_PORT` | Backend HTTP port | `8080` |
| `HTTP_READ_TIMEOUT` | HTTP request read timeout | `10s` |
| `HTTP_WRITE_TIMEOUT` | HTTP response write timeout | `10s` |
| `HTTP_IDLE_TIMEOUT` | HTTP keep-alive idle timeout | `10s` |
| `HTTP_MAX_BODY_SIZE` | Maximum request body size in bytes | `16384` |
| `HTTP_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `10s` |
| `HTTP_HEALTHCHECK_TIMEOUT` | Maximum wait for the PostgreSQL readiness check | `2s` |
| `HTTP_CORS_ALLOW_ORIGINS` | Comma-separated allowed CORS origins | `*` |
| `HTTP_CORS_MAX_AGE` | Browser CORS preflight cache duration | `12h` |
| `DATABASE_URL` | PostgreSQL connection URL | required |
| `DB_MAX_OPEN_CONNECTIONS` | Maximum open database connections | `10` |
| `DB_MAX_IDLE_CONNECTIONS` | Maximum idle database connections | `5` |
| `DB_CONNECTION_MAX_LIFETIME` | Maximum lifetime of a reused database connection | `5m` |
| `SENTRY_ENABLED` | Enables Sentry initialization and middleware | `false` |
| `SENTRY_DSN` | Sentry project DSN; needed when Sentry is enabled | empty |
| `SENTRY_ENVIRONMENT` | Environment label sent to Sentry | `development` |
| `SENTRY_FLUSH_TIMEOUT` | Maximum Sentry event delivery wait | `2s` |

The default CORS origin `*` is intended for local development. In a deployed environment, set `HTTP_CORS_ALLOW_ORIGINS` to the exact frontend origins, separated by commas. Cross-origin cookie credentials are disabled.

## Docker

Build the production image:

```bash
make docker-build
```

The default image is `frostfree/shortener:latest`. Override the tag with `DOCKER_TAG`:

```bash
make docker-build DOCKER_TAG=1.0.0
make docker-push DOCKER_TAG=1.0.0
```

The container listens on port `80`. Caddy serves the frontend from `/app/public` and proxies backend requests to the Go process on port `8080` inside the container.

The image does not include PostgreSQL. Supply a `DATABASE_URL` that is reachable from the container. On every container start, `/app/bin/run.sh` applies all pending Goose migrations before launching the backend and Caddy; startup stops immediately if a migration fails.

Keep `HTTP_PORT=8080` inside the container because the bundled Caddy configuration proxies to that port. Set `APP_BASE_URL` to the public URL through which clients reach the container.

## API

### Health Checks

The liveness endpoint reports whether the HTTP process is running:

```http
GET /ping
```

Returns `200 OK` with the plain-text body `pong`.

The readiness endpoint checks whether PostgreSQL is reachable within `HTTP_HEALTHCHECK_TIMEOUT`:

```http
GET /health
```

It returns `200 OK` with `{"ok":true}` when the application is ready to serve data-dependent requests. If PostgreSQL is unavailable or the check times out, it returns `503 Service Unavailable` with `{"ok":false}`. Database error details are written to the server log and are not exposed to the client.

Configure the deployment readiness probe or load balancer health check to use `/health`; `/ping` should be used only as a liveness probe.

### Create Link

```http
POST /api/links
Content-Type: application/json
```

```json
{
  "original_url": "https://example.com",
  "short_name": "example"
}
```

`original_url` must be an HTTP(S) URL. `short_name` is optional; when omitted or empty, a random 6-character alphanumeric value is generated. A custom value must contain 3–32 letters, digits, hyphens, or underscores.

Successful response: `201 Created`.

```json
{
  "id": 1,
  "original_url": "https://example.com",
  "short_name": "example",
  "short_url": "http://localhost:8080/r/example"
}
```

### List Links

```http
GET /api/links
```

| Parameter | Description | Default |
| --- | --- | --- |
| `range` | Inclusive range `[from,to]`; at most 1001 records | `[0,9]` |
| `sort` | JSON pair `["field","ASC"]` or `["field","DESC"]` | `["id","DESC"]` |

Supported sort fields are `id`, `original_url`, `short_name`, `short_url`, and `created_at`. Sorting by `short_url` is equivalent to sorting by `short_name`.

```bash
curl 'http://localhost:8080/api/links?range=[0,9]&sort=["id","DESC"]'
```

The response is a JSON array of link objects. The total unfiltered link count is returned in the header:

```http
Content-Range: links 0-9/42
```

### Get Link

```http
GET /api/links/:id
```

Returns `200 OK` and the same link representation used by the create endpoint.

### Update Link

```http
PUT /api/links/:id
Content-Type: application/json
```

Both fields are required:

```json
{
  "original_url": "https://example.org",
  "short_name": "docs"
}
```

Returns `200 OK` and the updated link representation.

### Delete Link

```http
DELETE /api/links/:id
```

Successful deletion returns `204 No Content`. Links are soft-deleted, so they are excluded from lookups, lists, and redirects. Their visits remain stored, and the database uniqueness constraint keeps the deleted short name reserved.

### Redirect

```http
GET /r/:short_name
```

Returns `302 Found` with the original URL in the `Location` header. The redirect continues even if recording its visit fails.

```bash
curl -I http://localhost:8080/r/example
```

### List Link Visits

```http
GET /api/link_visits
```

The endpoint supports the inclusive `range=[from,to]` query parameter, defaults to `[0,9]`, and allows at most 100 records. The `sort` parameter is a JSON pair `["field","ASC"]` or `["field","DESC"]`; supported fields are `id`, `link_id`, `created_at`, `ip`, `user_agent`, and `status`. The default order is `created_at DESC`.

```bash
curl 'http://localhost:8080/api/link_visits?range=[0,9]&sort=["status","DESC"]'
```

```http
Content-Range: link_visits 0-9/42
```

```json
[
  {
    "id": 1,
    "link_id": 1,
    "created_at": "2026-09-04T12:00:00Z",
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0",
    "status": 302
  }
]
```

## Error Responses

General errors use `error`:

```json
{
  "error": "link not found"
}
```


Validation errors use `errors`:

```json
{
  "errors": {
    "short_name": "short name must be between 3 and 32 characters long"
  }
}
```

| Status | When it is used |
| --- | --- |
| `400 Bad Request` | Malformed, empty, or type-invalid JSON |
| `404 Not Found` | Link or route does not exist |
| `405 Method Not Allowed` | Route exists but does not support the HTTP method |
| `409 Conflict` | `short_name` is already in use |
| `413 Request Entity Too Large` | Request body exceeds `HTTP_MAX_BODY_SIZE` |
| `422 Unprocessable Entity` | Binding validation, path validation, range, or sort error |
| `500 Internal Server Error` | Unexpected internal error |

## Development Commands

| Command | Description |
| --- | --- |
| `make build` | Build `bin/shortener` from `cmd/shortener`. |
| `make run` | Run the compiled binary; accepts `ARGS="..."`. |
| `make dev` | Run the application with Air live reload. |
| `make test` | Run unit tests with verbose output. |
| `make test-integration` | Run PostgreSQL integration tests using Testcontainers; requires Docker. |
| `make test-coverage` | Write `coverage.out` and print coverage by function. |
| `make install-lint` | Install the pinned `golangci-lint` version. |
| `make lint` | Run configured linters. |
| `make fmt` | Format code through `golangci-lint fmt`. |
| `make fmt-check` | Check formatting without modifying files. |
| `make tidy-check` | Verify that `go.mod` and `go.sum` are tidy. |
| `make lint-fix` | Format code and apply supported lint fixes. |
| `make docker-build` | Build the `linux/amd64` Docker image. |
| `make docker-push` | Build and push the `linux/amd64` Docker image. |

Typical verification before pushing:

```bash
make build
make test
make test-integration
make test-coverage
make fmt-check
make tidy-check
make lint
```

## Project Structure

```text
.
|-- cmd/shortener                 # Application entry point
|-- db/migrations                 # Goose SQL migrations
|-- internal/app                  # Application lifecycle
|-- internal/config               # Environment-based configuration
|-- internal/db                   # PostgreSQL connection and repositories
|-- internal/httpserver           # Gin server, handlers, DTOs, and middleware
|-- internal/services/shortener   # Use cases and application contracts
|-- internal/tests/integration    # Testcontainers integration tests
|-- .github                       # CI workflows and generated badge data
|-- .env.example                  # Complete local configuration example
|-- bin/run.sh                    # Container migrations and process entrypoint
|-- Caddyfile                     # Frontend/static proxy configuration
|-- Dockerfile                    # Multi-stage production image
|-- Makefile                      # Development commands
`-- README.md                     # Project documentation
```
