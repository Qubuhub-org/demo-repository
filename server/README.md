# Qubuhub authentication API

Self-hosted Go authentication service for the existing React/Vite dashboard.

## Configuration

Required: `DATABASE_URL`. Optional: `AUTH_ADDR` (defaults to `:8080`) and `COOKIE_SECURE` (set to `true` behind HTTPS; defaults to insecure for local HTTP development).

Apply the schema before starting the service:

```sh
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f server/migrations/001_auth.sql
go mod tidy
go run ./server/cmd/api
```

## Routes

- `GET /healthz`
- `POST /api/auth/register` — JSON `{"email":"...","password":"..."}`
- `POST /api/auth/login` — same JSON shape
- `POST /api/auth/logout`
- `GET /api/auth/me`

Mutation requests require a same-origin `Origin` header. Session tokens are random, stored only as SHA-256 hashes, and issued in HttpOnly, SameSite=Strict cookies. Passwords use Argon2id. The service currently has process-local IP throttling; deploy behind a trusted reverse proxy and replace this limiter with a shared store before running multiple replicas.

## Production requirements

Use HTTPS and set `COOKIE_SECURE=true`. Keep `DATABASE_URL` in the deployment secret manager, never in the repository or frontend bundle. Configure database backups, TLS to PostgreSQL, monitoring, email-based verified password recovery, account verification, MFA, session cleanup, and a distributed rate limiter before production use. This initial foundation intentionally does not claim to provide password reset, email verification, MFA, or role-based authorization yet.
