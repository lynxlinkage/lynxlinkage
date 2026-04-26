# lynxlinkage landing page

A monorepo containing the lynxlinkage marketing site.

- `backend/` — Go 1.22+ + Gin, SQLite via `modernc.org/sqlite` (pure Go, no CGO),
  serves `/api/v1/*` and (when built with `-tags=embed`) the prerendered
  SvelteKit frontend on every other path.
- `frontend/` — SvelteKit 2 + Svelte 5 (runes), TypeScript, `adapter-static`.
  Pages are prerendered to static HTML at build time.

## Quick start

```sh
# 1. Install backend deps (Go modules)
cd backend && go mod download && cd ..

# 2. Install frontend deps
cd frontend && pnpm install   # or `npm install`
cd ..

# 3. Seed initial content into SQLite
make seed

# 4. Run backend and frontend dev servers in parallel
make dev
```

- Backend listens on http://localhost:8080
- Frontend Vite dev server on http://localhost:5173 (Vite proxies `/api/*` to 8080)

## Production build (single binary)

```sh
make build         # builds frontend + Go binary with frontend embedded
./bin/lynxlinkage  # serves both the static site and the API on :8080
```

The production binary contains:

- The SvelteKit static export (under `backend/internal/static/dist/`,
  embedded via `//go:embed`).
- All API handlers and the SQLite driver. SQLite uses a file at
  `./data/lynxlinkage.db` by default; override with `DATABASE_URL`.

## Project structure

```
backend/
  cmd/server/         # main HTTP server entry point
  cmd/seed/           # `go run ./cmd/seed` loads YAML into the DB
  internal/api/       # /api/v1 handlers
  internal/domain/    # entities and DTOs
  internal/store/     # sqlx-based repos + embedded migrations
  internal/middleware # logger, recover, CORS, IP rate limiter
  internal/static/    # embed.FS wrapper for the SvelteKit build
  seed/               # initial content (research cards, jobs, partners)

frontend/
  src/routes/         # / · /about · /researches · /hiring · /partners
  src/lib/components/ # Header, Footer, Hero, ResearchCard, JobCard, …
  src/lib/api/        # server.ts (build-time) + client.ts (browser)
  src/lib/styles/     # tokens.css (design tokens)
  static/             # favicon, og image, robots.txt
```

## Public API

All endpoints are JSON. Read endpoints are safe to call at prerender time.

| Method | Path                              | Purpose                              |
| ------ | --------------------------------- | ------------------------------------ |
| GET    | `/api/v1/health`                  | Liveness                             |
| GET    | `/api/v1/researches?tag=&limit=`  | Public research cards                |
| GET    | `/api/v1/jobs`                    | Active job postings                  |
| GET    | `/api/v1/jobs/:id`                | Single job posting                   |
| GET    | `/api/v1/partners`                | Partners (logo wall)                 |
| POST   | `/api/v1/contact`                 | Contact submission (rate-limited)    |
| POST   | `/api/v1/auth/login`              | Sign in with email + password        |
| POST   | `/api/v1/auth/logout`             | Clear session cookie                 |
| GET    | `/api/v1/auth/me`                 | Current authenticated user           |
| GET    | `/api/v1/admin/jobs`              | All postings (HR only)               |
| POST   | `/api/v1/admin/jobs`              | Create posting (HR only)             |
| PUT    | `/api/v1/admin/jobs/:id`          | Update posting (HR only)             |
| DELETE | `/api/v1/admin/jobs/:id`          | Hard-delete posting (HR only)        |

## HR / admin

Recruiters sign in at `/login` and manage job postings at `/admin`. Both
pages are client-rendered SPAs (`prerender = false`, `ssr = false`) served
through the SvelteKit `200.html` fallback so the rest of the site remains
fully prerendered.

Authentication is HMAC-signed session cookies (HttpOnly, SameSite=Strict,
7-day TTL by default). Rotate `SESSION_SECRET` to invalidate every
outstanding session.

### Bootstrap an HR user

```sh
# Interactive (password prompted, not echoed):
make createuser EMAIL=hr@example.com

# Non-interactive:
make createuser EMAIL=hr@example.com PASSWORD='choose-a-strong-one'
```

Behind the scenes this runs `go run ./cmd/createuser`, which inserts a row
into the `users` table with a bcrypt hash. Currently only the `hr` role is
supported.

### Production checklist

- Set `SESSION_SECRET` to a long random value
  (`openssl rand -base64 48`). The server refuses to start in
  `APP_ENV=production` without it.
- Serve the binary behind TLS — the `Secure` cookie flag is set when
  `APP_ENV=production`, so the cookie won't be sent over plain HTTP.

## Configuration

The backend reads everything from environment variables. See
[`backend/.env.example`](backend/.env.example) for the full list and defaults.
The most relevant ones:

- `APP_ENV` — `development` or `production`
- `HTTP_ADDR` — listen address (default `:8080`)
- `DATABASE_URL` — SQLite DSN (default `file:./data/lynxlinkage.db?...`)
- `CORS_ALLOW_ORIGIN` — comma-separated origins (default `http://localhost:5173`)
- `CONTACT_RPS` / `CONTACT_BURST` — per-IP rate limit on the contact endpoint

The frontend uses `BACKEND_URL` (see [`frontend/.env.example`](frontend/.env.example))
during prerender to point load functions at the backend.

## Design notes

- **Static-first.** The frontend is prerendered. Updating cards/jobs/partners
  requires a rebuild — trigger this from CI on a webhook, or run
  `make build` on the host. The contact form is the only runtime call.
- **One binary.** In production the Go server serves both the static frontend
  and the API on the same origin; no CORS, no separate frontend host.
- **Lean store.** SQLite is sufficient for the scale of a marketing site.
  The `store/` package is small and behind interfaces, so swapping in
  Postgres is mostly a wiring change.
- **Style.** Vanilla CSS with design tokens (no Tailwind). Single navy
  accent, generous whitespace, scoped Svelte styles.

## CI

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) runs:

- `go vet`, `go test` and `go build` on the backend.
- `pnpm lint`, `pnpm check`, `pnpm build` on the frontend.
- A combined embedded production binary build, uploaded as an artifact.

## License

Proprietary &mdash; lynxlinkage.
