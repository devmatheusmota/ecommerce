# Kong API Gateway

Declarative config: **kong.yml.template** is the source for seed. At startup `kong-seed` substitutes `__JWT_SECRET__` with `JWT_SECRET` from `services/users/.env` (default `dev`), then runs `kong config db_import`. Static **kong.yml** is a fallback for manual import.

## Dev vs prod: who is the source of truth?

**Dev (local):** `kong-seed` runs when you bring the stack up. The YAML is the source of truth; `make kong-reset` wipes Kong’s DB and re-imports so you can try config changes. Any change you make in the Kong UI will be **overwritten** the next time the seed runs (e.g. after `kong-reset` or a full down/up with a fresh volume).

**Prod:** You will **not** run `kong-reset` in production. Use `kong.yml` + `kong-seed` only for **initial bootstrap** (first deploy): that loads the base routes, consumer, and JWT plugin. After that, the **Kong DB (and UI) is the source of truth**. Create or change routes, plugins, rate limits, etc. in the Kong Manager UI (or Admin API); those changes persist. So in prod: bootstrap once from YAML, then manage everything via the UI. On later deploys, either **skip running kong-seed** (so the seed doesn’t re-import and overwrite UI changes) or run it only when the Kong DB is empty (e.g. new environment).

**If `kong-seed` fails with `UNIQUE violation on key="users"`:** the Kong DB already has the config from a previous run. In dev, run `make kong-reset` then `make dev`. In prod you normally don’t re-run the seed after the first deploy.

## JWT for protected routes

- **GET /v1/users/me** is protected: Kong validates the Bearer JWT and forwards the user ID to the upstream in the `X-User-ID` header.
- The users service issues JWTs with claim **`iss`: `"users"`** so Kong can match the consumer.
- **Secret**: Kong reads the same `JWT_SECRET` as the users service from `services/users/.env` (seed uses `kong.yml.template` and substitutes it). It **must match** or token verification will fail.

**Local dev:** Set `JWT_SECRET=dev` in `services/users/.env` (see `.env.example`). Kong seed uses that value automatically.

**Production:** Do not commit the real secret in `kong.yml`. Use Kong Admin API to create/update the consumer’s JWT credential, or generate the config from a secret store.

## Per-route plugins (e.g. rate limit on login)

Each path has its **own route** in Kong (e.g. `users-login-route`, `users-register-route`). So in the Kong Manager UI you can attach plugins to a single path: e.g. open **Routes → users-login-route → Plugins** and add a rate-limit plugin only for login, without affecting register or other paths.

## Path routing (strip_path pattern)

Kong uses `strip_path: true` on all routes. The catch-all `/v1/users` strips the prefix and forwards the path suffix to the upstream — no request-transformer needed for normal routes:

```
GET  /v1/users/login    → upstream /login
POST /v1/users/register → upstream /register
GET  /v1/users/health   → upstream /health
GET  /v1/users/docs     → upstream /docs
```

`/v1/users/me` is a separate route (for JWT protection). It also uses `strip_path: true`, but since the full path is stripped on an exact match, a request-transformer sets the upstream URI to `/me`.

## Per-route plugins (e.g. rate limit on login)

When you need a plugin on a specific path (e.g. rate limit only on login), add a dedicated route for that path in `kong.yml` or via the Kong UI:

```yaml
- name: users-login-route
  paths:
    - /v1/users/login
  strip_path: true
  methods:
    - POST
```

Then attach the plugin to `users-login-route` only. Add a request-transformer to that route to set `uri: "/login"` (exact match strips the full path). The catch-all `users-route` remains unchanged for all other paths.

## Route summary

| Route name            | Path                | Method | Auth | Upstream path |
|-----------------------|---------------------|--------|------|----------------|
| health-route          | /health             | *      | no   | (httpbin)      |
| users-me-route        | /v1/users/me        | GET    | JWT  | /me            |
| users-login-route     | /v1/users/login     | POST   | no   | /login         |
| users-register-route  | /v1/users/register  | POST   | no   | /register      |
| users-route           | /v1/users           | *      | no   | /{suffix}      |

`users-login-route` and `users-register-route` exist so plugins (e.g. rate limit) can be attached to each individually in the Kong UI. Other paths fall through to `users-route` (catch-all). When another endpoint needs a per-route plugin, add a specific route for it following the same pattern (strip_path: true + request-transformer to set the upstream URI).
