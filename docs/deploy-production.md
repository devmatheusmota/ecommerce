# Deploy to production — checklist

When promoting the current setup (Kong as single HTTP entry point, no direct client access to backends) to production, change the following.

---

## 1. **Config and secrets**

| Item | Local (today) | Production |
|------|----------------|------------|
| **`.env` / `services/users/.env`** | Dev credentials, often in repo or `.env.example` | Do **not** commit prod values. Use a secret manager (e.g. AWS Secrets Manager, Vault) or env vars injected by the platform. |
| **Kong admin** | `KONG_PASSWORD: kong_admin` | Strong password; restrict who can call Admin API (see below). |
| **Postgres / RabbitMQ** | `ecommerce_dev`-style passwords | Strong passwords; consider managed Postgres/RabbitMQ. |
| **`infra/k8s/*/secret.yaml`** | Placeholders / dev values | Replace with sealed-secrets, external-secrets, or platform secrets. |

---

## 2. **Kong — single public entry point**

- **HTTP/HTTPS**: In prod, expose **only** Kong’s proxy ports (e.g. 8000 HTTP, 8443 HTTPS) to the internet (or put Kong behind a load balancer / ingress that terminates TLS).
- **Admin API / Manager**: Do **not** expose Kong Admin (8001) or Kong Manager (8002) to the internet. Restrict to internal network / VPN / bastion, or disable if not needed in prod.
- **`infra/kong/kong.yml`**: Same routing logic (e.g. `/v1/users` → users service). Only base URLs / upstreams might change if hostnames differ in prod (e.g. `http://users:8080` stays if the service is still named `users` in the same network).

No change is required to the “no direct backend ports” rule: in prod, backends still must **not** expose HTTP to the client; only Kong (or the LB in front of Kong) should be public.

---

## 3. **Docker Compose in production**

If you run prod with Docker Compose on a server:

- **Backend services (e.g. users)**: Already correct — no `ports:` for HTTP, only `expose` for the internal network. Kong is the only service that needs proxy ports (8000/8443) published.
- **Firewall**: Open only what clients and ops need (e.g. 8000, 8443, SSH). Do not open 8081, 8082, etc. for backend apps.
- **Compose files**: Use prod env files or env vars that point to prod DB/broker and strong secrets; do not use `docker-compose.dev.yml` in prod.

---

## 4. **Kubernetes in production**

When you add Kong and app services (e.g. users) to `infra/k8s/`:

- **Backend Services (users, catalog, orders, …)**: Use `type: ClusterIP` for HTTP. No `LoadBalancer` and no `NodePort` for those — only Kong (or Ingress) should be reachable from outside the cluster.
- **Kong**: Expose via `LoadBalancer`, `NodePort`, or (recommended) an **Ingress** (e.g. NGINX Ingress or platform Ingress) with TLS. Kong’s proxy port is the only public HTTP/HTTPS entry point.
- **Admin/Manager**: Do not expose Kong Admin/Manager publicly; use internal Services or restrict by network policy.

Same principle as Compose: clients only hit Kong; backends are internal.

---

## 5. **Summary**

| Topic | What to do for prod |
|-------|----------------------|
| Backend HTTP ports | No change: keep them internal (Compose: no `ports:`; K8s: `ClusterIP`). |
| Client access | Only through Kong (or LB/Ingress in front of Kong). |
| Secrets & config | Prod credentials and URLs from secret manager or platform env; never commit. |
| Kong Admin/Manager | Not exposed to the internet; internal only or disabled. |
| Kong proxy (8000/8443) | Exposed (or behind LB/Ingress with TLS). |
| Firewall / network | Only Kong (and SSH/ops) open; no backend ports. |

The change you made locally (no direct client access to backends) already matches production: you only need to lock down secrets, Kong admin, and how Kong is exposed (and, in K8s, use `ClusterIP` for backends and expose only Kong).
