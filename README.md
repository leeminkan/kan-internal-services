# kan-internal-services

Small internal service to automate check-in/out against an external Blueprint SSO/API.

This repository provides a minimal HTTP server exposing a single endpoint that performs a login flow and calls a check-in/out API. The service is intentionally lightweight and currently implements the core flow used in the `features/checkinout` package.

## What it does

- Exposes POST /api/checkinout
- Performs an SSO/login flow against Blueprint (the code uses `LOGIN_URL` and follows redirects), then calls a CHECKINOUT API endpoint.
- The request body must include `username` and `password`. Both fields are required; if missing the endpoint returns a 400 validation error.

## Repository layout

- `main.go` - entrypoint, loads environment and registers routes with Fiber.
- `config/` - placeholder for configuration logic (`config/config.go`).
- `features/checkinout/` - implementation of the checkin/out handler and the HTTP scraping/requests flow.
  - `handler.go` - HTTP handler and route registration.
  - `service.go` - low-level login, redirect handling and check-in/out HTTP calls.
- `pkg/` - package-level helpers (currently `pkg/telegram.go` is a placeholder for notification logic).

## Requirements

- Go 1.21

The project uses these notable modules (see `go.mod`):

- `github.com/gofiber/fiber/v2` - HTTP server
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/joho/godotenv` - environment variable loading from `.env` (used in `main.go`)

## Environment variables

- `PORT` - server listen port (default `8080`).

You can add a `.env` file in the project root for local development to set `PORT` and other runtime variables.

## Running locally

1. Build

```sh
go build -o kan-internal-services .
```

2. Run

```sh
PORT=8080 ./kan-internal-services
```

Or use `go run` during development:

```sh
go run main.go
```

### Docker

You can build and run the service in a container using the provided `Dockerfile`.

Build the image locally:

```sh
docker build -t kan-internal-services:local .
```

Run the container:

```sh
docker run --rm -p 8080:8080 -e PORT=8080 kan-internal-services:local
```

Or use Docker Compose to build and run (recommended during development):

```sh
docker compose up --build -d
```

The service will be available on `http://localhost:8080`.

Notes:

- The `/api/checkinout` endpoint still requires `username` and `password` to be provided in the request body (JSON). The container does not inject credentials automatically.
- If you wish to pass any runtime environment variables, add them to the `docker-compose.yml` `environment` section or provide an `.env` file consumed by Compose.

## API

POST /api/checkinout

Request body (JSON):

```json
{
  "username": "required_username",
  "password": "required_password"
}
```

If `username` or `password` are omitted from the JSON body, the endpoint returns a 400 validation error. Provide both fields in the request body.

Example curl (explicit credentials):

```sh
curl -X POST http://localhost:8080/api/checkinout \
  -H 'Content-Type: application/json' \
  -d '{"username":"phuchoangnguyen","password":"secret"}'
```

Response (JSON):

```json
{ "success": true, "message": "Check-in/out completed" }
```

or on error:

```json
{ "success": false, "error": "..." }
```

## Important implementation notes

-- The check-in/out flow is implemented in `features/checkinout/service.go`. It:

- Loads login page and follows a series of redirects
- Parses an HTML form (`form#login-form`) to find the form `action`
- Submits username/password as form data and follows redirects
- Finally posts to `CHECKINOUT_URL` and prints the response body

- The HTTP client uses a cookie jar and disables automatic redirect following for certain steps. There is currently no timeout set on the HTTP client.

## Security

- Credentials are sensitive. Do not include them in code, logs, or commits. For production consider using a secrets manager instead of placing credentials in files or environment variables accessible to many systems.

## Suggested improvements (next steps)

1. Add structured logging and configurable log levels.
2. Extract configuration with a proper loader (Viper or similar) and avoid globals.
3. Add timeouts and better error handling for HTTP requests. Use context.Context.
4. Make the HTTP interactions mockable and add unit tests for the login/check-in flow.
5. Add retries/backoff and graceful rate limiting if running periodic automation.
6. Implement notification integration (e.g., Telegram) in `pkg/telegram.go` and attach success/failure hooks.

## License

This repository currently contains no license. Add an appropriate LICENSE file if you intend to share or publish this code.

## Contact / Maintainer

Repository owner / author: (update README with appropriate contact info)

---

If you want, I can also:

- run `go build` here and report results, or
- add a small `.github/workflows` CI file to build/check the repo, or
- create unit test scaffolding for `Run()` (with an HTTP client interface) to make the code testable.
