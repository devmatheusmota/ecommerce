# E-commerce Web Frontend

Next.js frontend for the marketplace. Consumes APIs via Kong gateway.

## Prerequisites

- Backend running: `make up` (from project root)
- Node.js 18+

## Run locally

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Environment

- `NEXT_PUBLIC_API_URL` — Kong gateway URL (default: `http://localhost:8000`)

API requests are proxied via Next.js rewrites to avoid CORS.
