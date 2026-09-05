# URL Shortener

URL Shortener is a Go web service for creating, managing, and tracking shortened links.

The application exposes a JSON API for link CRUD operations, redirects short URLs to their original targets, and stores redirect visits in PostgreSQL.

[A live deployment](https://shortener-latest.onrender.com/) is available on Render.

## Features

- Creates short links from original URLs.
- Generates a deterministic 6-character short name when a custom one is not provided.
- Supports custom short names with uniqueness checks.
- Redirects users from `/r/:short_name` to the original URL.
- Records redirect visits with IP address, user agent, referrer, and response status.
- Provides paginated API responses with `Content-Range` headers.
- Uses PostgreSQL through GORM and SQL migrations.
- Includes a Docker image that serves the Hexlet frontend through Caddy and proxies API requests to the Go backend.

## Pipeline Status

[![Actions Status](https://github.com/Fr0stFree/go-project-278/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/Fr0stFree/go-project-278/actions)
[![CI](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml/badge.svg?branch=master)](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml)

## Quality

[![Test Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/Fr0stFree/go-project-278/master/.github/badges/coverage-badge.json)](https://github.com/Fr0stFree/go-project-278/actions/workflows/test-and-lint.yml)

## Requirements

- Go `1.26.3`
- PostgreSQL
- `make`
- Docker, for container builds

The linter target uses `golangci-lint` `v2.12.2`. Run `make install-lint` once before `make lint`, `make fmt`, or `make lint-fix`.

## Installation

Clone the repository and build the binary:

```bash
git clone git@github.com:Fr0stFree/go-project-278.git
cd go-project-278
make build
```

The compiled program will be available at:

```bash
./bin/shortener
```

## Configuration

The application reads configuration from environment variables. A local `.env` file is optional and is loaded automatically when present.

Example local configuration:

```env
# APP
APP_BASE_URL=http://localhost:8080

# HTTP
HTTP_PORT=8080
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s

# DB
DATABASE_URL=postgres://shortener:password@localhost:5432/shortener?sslmode=disable
DB_MAX_OPEN_CONNECTIONS=10
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_MAX_LIFETIME=5m
```

Environment variables:

| Variable | Description | Default |
| --- | --- | --- |
| `APP_BASE_URL` | Public base URL used to build `short_url` values |  |
| `HTTP_PORT` | HTTP server port | `8080` |
| `HTTP_READ_TIMEOUT` | HTTP read timeout | `10s` |
| `HTTP_WRITE_TIMEOUT` | HTTP write timeout | `10s` |
| `DATABASE_URL` | PostgreSQL connection URL | required |
| `DB_MAX_OPEN_CONNECTIONS` | Maximum open database connections | `10` |
| `DB_MAX_IDLE_CONNECTIONS` | Maximum idle database connections | `5` |
| `DB_CONNECTION_MAX_LIFETIME` | Maximum lifetime of reused database connections | `5m` |

Database schema migrations are stored in `db/migrations` and use `goose` annotations.

## Running Locally

Build and run the server:

```bash
make build
make run
```

The API will be available at:

```text
http://localhost:8080
```

For live reload development, install Air and run:

```bash
make dev
```

## Docker

Build the Docker image:

```bash
make docker-build
```

The image is tagged as:

```text
frostfree/shortener:latest
```

You can override the tag:

```bash
make docker-build DOCKER_TAG=1.0.0
```

Push the image:

```bash
make docker-push DOCKER_TAG=1.0.0
```

The container listens on port `80`. It serves the frontend from `/app/public` with Caddy and reverse-proxies all non-static requests to the Go backend on port `8080`.

The Docker image does not start PostgreSQL. Provide database environment variables that point to an existing PostgreSQL instance.

## API

### Health Check

```http
GET /ping
```

Response:

```text
pong
```

### Create Link

```http
POST /api/links
Content-Type: application/json
```

Request body:

```json
{
  "original_url": "https://example.com",
  "short_name": "example"
}
```

`short_name` is optional. When it is omitted or empty, the service generates one from the original URL.

Response:

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

Supported query parameters:

| Parameter | Description | Default |
| --- | --- | --- |
| `range` | Inclusive pagination range in the form `[from,to]` | `[0,9]` |
| `sort` | Sort field and order in the form `["field","ASC"]` or `["field","DESC"]` | `["id","DESC"]` |

Supported sort fields are `id`, `original_url`, `short_name`, `short_url`, and `created_at`.

Example:

```bash
curl 'http://localhost:8080/api/links?range=[0,9]&sort=["id","DESC"]'
```

Response headers include the total count:

```http
Content-Range: links 0-9/42
```

Response body:

```json
[
  {
    "id": 1,
    "original_url": "https://example.com",
    "short_name": "example",
    "short_url": "http://localhost:8080/r/example"
  }
]
```

### Get Link

```http
GET /api/links/:id
```

Response:

```json
{
  "id": 1,
  "original_url": "https://example.com",
  "short_name": "example",
  "short_url": "http://localhost:8080/r/example"
}
```

### Update Link

```http
PUT /api/links/:id
Content-Type: application/json
```

Request body:

```json
{
  "original_url": "https://example.org",
  "short_name": "docs"
}
```

Response:

```json
{
  "id": 1,
  "original_url": "https://example.org",
  "short_name": "docs",
  "short_url": "http://localhost:8080/r/docs"
}
```

### Delete Link

```http
DELETE /api/links/:id
```

Successful deletion returns `204 No Content`.

### Redirect

```http
GET /r/:short_name
```

The service records a visit and returns `302 Found` with a `Location` header pointing to the original URL.

Example:

```bash
curl -I http://localhost:8080/r/example
```

### List Link Visits

```http
GET /api/link_visits
```

Supported query parameters:

| Parameter | Description | Default |
| --- | --- | --- |
| `range` | Inclusive pagination range in the form `[from,to]` | `[0,9]` |

Example:

```bash
curl 'http://localhost:8080/api/link_visits?range=[0,9]'
```

Response headers include the total count:

```http
Content-Range: link_visits 0-9/42
```

Response body:

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

Errors are returned as JSON:

```json
{
  "error": "link not found"
}
```

Validation and conflict errors are keyed by field:

```json
{
  "error": {
    "short_name": "shortname already in use"
  }
}
```

Common statuses:

| Status | When it is used |
| --- | --- |
| `400 Bad Request` | Invalid JSON or missing required JSON fields |
| `404 Not Found` | Requested link does not exist |
| `409 Conflict` | `short_name` is already in use |
| `422 Unprocessable Entity` | Invalid path or query parameter |
| `500 Internal Server Error` | Unexpected server error |

## Makefile Commands

| Command | Description |
| --- | --- |
| `make build` | Builds the CLI binary into `bin/shortener`. |
| `make run` | Runs the compiled binary. Additional arguments can be passed with `ARGS="..."`. |
| `make dev` | Runs the app with Air live reload. |
| `make test` | Runs all Go tests with verbose output. |
| `make test-coverage` | Runs tests, writes `coverage.out`, and prints coverage by function. |
| `make install-lint` | Installs the configured `golangci-lint` version. |
| `make lint` | Runs `golangci-lint` with the project config. |
| `make fmt` | Formats code through `golangci-lint fmt`. |
| `make lint-fix` | Formats code and applies automatic lint fixes. |
| `make docker-build` | Builds the Docker image for `linux/amd64`. |
| `make docker-push` | Builds and pushes the Docker image for `linux/amd64`. |

Run database migrations with `goose` when preparing a new PostgreSQL database:

```bash
go tool goose -dir db/migrations postgres "$DATABASE_URL" up
```

To update the coverage badge manually:

```bash
make test-coverage
.github/scripts/generate-coverage-badge.sh
```

Typical local check before pushing changes:

```bash
make install-lint
make test
make test-coverage
make lint
```

## Project Structure

```text
.
|-- cmd/shortener                 # Application entry point
|-- internal/app                  # Dependency wiring and application runner
|-- internal/config               # Environment-based configuration
|-- internal/db                   # PostgreSQL and GORM setup
|-- internal/db/models            # Repository layer and storage models
|-- db/migrations                 # SQL database migrations
|-- internal/httpserver           # Gin router, HTTP server, and error responses
|-- internal/services/shortener   # Business logic for links and visits
|-- .github/badges                # Generated badge data
|-- .github/scripts               # Maintenance scripts
|-- .github/workflows             # GitHub Actions workflows
|-- Caddyfile                     # Runtime frontend/static proxy config
|-- Dockerfile                    # Multi-stage production image
|-- Makefile                      # Development commands
`-- README.md                     # Project documentation
```
