# E-commerce Project Plan — Features, Microservices & Tasks

Mercado Livre-style marketplace: multiple sellers, buyers, and the full purchase journey.

---

## 1. Microservices Map

Each service owns a bounded context and its own database. Communication: gRPC (sync) or RabbitMQ (async).

| Service | Responsibility | Sync (gRPC) | Async (Events) |
|---------|----------------|-------------|----------------|
| **users** | Auth, profiles, addresses | - | `user.created`, `user.updated` |
| **catalog** | Products, categories, search | Product lookup, stock check | `product.created`, `product.updated` |
| **inventory** | Stock per seller/SKU | Reserve stock, release stock | `stock.reserved`, `stock.released`, `stock.low` |
| **orders** | Cart, checkout, order lifecycle | - | `order.created`, `order.paid`, `order.shipped` |
| **payments** | Payment processing, refunds | - | `payment.completed`, `payment.failed`, `refund.processed` |
| **shipping** | Shipping options, tracking, delivery | Get shipping rates | `shipment.created`, `shipment.delivered` |
| **sellers** | Seller registry, storefront | Seller info | `seller.registered` |
| **reviews** | Product/seller ratings | - | `review.created` |
| **notifications** | Email, SMS, push, in-app | - | (consumer only) |

---

## 2. Feature Breakdown by Domain

### 2.1 Users
- [x] User registration (email, password, name, phone, CPF with validation)
- [x] Login / JWT tokens
- [x] Profile (GET /me, update name/phone/CPF)
- [x] Addresses (CRUD): billing and shipping separation; default per type (default billing, default shipping)
- [x] Password reset (request token by email; confirm with token + new password)
- [ ] (Later) OAuth (Google, etc.)

### 2.2 Catalog
- [x] Categories (tree)
- [x] Products (title, description, images, price, seller_id, category_id)
- [ ] Product variants (size, color, etc.)
- [ ] Product search (by name, category, price range)
- [ ] Product listing by seller
- [ ] Pagination and filters

### 2.3 Inventory
- [ ] Stock per product/variant
- [ ] Reserve stock on order creation
- [ ] Release stock on order cancellation/payment failure
- [ ] Reduce stock on payment confirmation
- [ ] Low stock alerts (event)
- [ ] Stock history / audit

### 2.4 Orders
- [ ] Cart (add, remove, update quantity)
- [ ] Checkout (cart → order)
- [ ] Order creation (pending payment)
- [ ] Order status: pending → paid → shipped → delivered
- [ ] Order history per user
- [ ] Order detail (items, shipping, payment)
- [ ] Order cancellation (before payment)

### 2.5 Payments
- [ ] Payment methods (card, PIX, boleto — mocked for learning)
- [ ] Process payment (mock gateway)
- [ ] Payment status webhook / event
- [ ] Refund
- [ ] Payment history

### 2.6 Shipping
- [ ] Shipping address validation
- [ ] Calculate shipping cost (mock carrier)
- [ ] Shipping options (economy, express)
- [ ] Shipment creation on order paid
- [ ] Tracking (mock tracking number)
- [ ] Delivery confirmation (event)

### 2.7 Sellers
- [ ] Seller registration
- [ ] Seller profile (store name, logo, description)
- [ ] Storefront (products by seller)
- [ ] Seller dashboard (orders, sales — later)

### 2.8 Reviews
- [ ] Review product (rating 1–5, comment)
- [ ] Review seller
- [ ] List reviews for product/seller
- [ ] Only buyers who purchased can review

### 2.9 Notifications *(last feature to implement)*
- [ ] Notifications service skeleton (Go, RabbitMQ consumer)
- [ ] Consume SMS/email requests from RabbitMQ queue `notifications`; send via provider or log (dev)
- [ ] Order confirmed (email)
- [ ] Payment received
- [ ] Shipment shipped
- [ ] Shipment delivered
- [ ] (Later) In-app, push

### 2.10 Search
- [ ] Full-text search products
- [ ] Filters (category, price, seller)
- [ ] Sorting (price, date, rating)
- [ ] (Later) Elasticsearch/OpenSearch for scale

### 2.11 Cross-cutting
- [ ] Kong API Gateway (routes, rate limiting, auth)
- [ ] JWT validation at gateway or services
- [ ] OpenAPI/Swagger for HTTP APIs
- [ ] OpenTelemetry traces
- [ ] Structured logging
- [ ] Prometheus metrics
- [ ] Health checks (readiness, liveness)
- [ ] CI/CD (build, test, deploy)

---

## 3. Infrastructure & Platform Tasks

- [ ] Kong in Docker Compose
- [ ] Kong in K8s
- [ ] Kong routes for each service
- [ ] Kong JWT plugin (or auth service integration)
- [ ] Redis for session/cache (optional, later)
- [ ] Prometheus + Grafana
- [ ] OpenTelemetry instrumentation
- [ ] GitHub Actions (or similar) for CI/CD
- [ ] K8s manifests for each service
- [ ] Helm charts (optional, for packaging)

---

## 4. Shared / Proto Contracts

- [ ] `shared/proto` — gRPC definitions
  - [ ] catalog.proto (GetProduct, CheckStock)
  - [ ] inventory.proto (ReserveStock, ReleaseStock)
  - [ ] shipping.proto (GetRates)
  - [ ] users.proto (GetUser, ValidateToken — if needed)
- [ ] Event schemas (JSON or protobuf) for RabbitMQ
  - [ ] order.created, order.paid, order.shipped
  - [ ] payment.completed, payment.failed
  - [ ] stock.reserved, stock.released
  - [ ] etc.
- [ ] Shared Go packages
  - [ ] Event publishing/consuming
  - [ ] DB migrations setup
  - [ ] Logging, tracing helpers

---

## 5. Suggested Implementation Order

### Phase 1 — Foundation
1. Kong (Docker Compose + basic route)
2. **users** service — registration, login, JWT
3. Kong auth (JWT validation or pass-through to users)
4. **catalog** service — categories, products (CRUD)
5. **sellers** service — seller registration, link to catalog

### Phase 2 — Core Purchase Flow
6. **orders** service — cart, checkout, order creation
7. **inventory** service — stock, reserve/release
8. **orders** ↔ **inventory** gRPC (reserve on checkout)
9. **payments** service — mock payment, emit `payment.completed`
10. **orders** — listen to `payment.completed`, update status
11. **inventory** — listen to `payment.completed`, reduce stock
12. **shipping** service — mock rates, create shipment on `order.paid`

### Phase 3 — Completion & Polish
13. **reviews** service — rate products/sellers
14. Search improvements (DB full-text or dedicated search service)
15. Admin/seller dashboard (simplified)

### Phase 4 — Scale & Ops
16. Observability (Prometheus, Grafana, OpenTelemetry)
17. Kong in K8s, K8s manifests for all services
18. CI/CD pipeline
19. Load testing, chaos engineering (optional)

### Phase 5 — Notifications (last)
20. **notifications** service — consume events from RabbitMQ, send SMS/email (mock)

---

## 6. Task Checklist (Flat)

Use this as a backlog. Check off as you go.

### Infra
- [x] Add Kong to Docker Compose
- [x] Kong JWT plugin on GET /v1/users/me (validates Bearer, sets X-User-ID from claim `sub`; secret from env via kong.yml.template)
- [ ] Add Kong to K8s
- [ ] Define Kong routes (placeholder for each service)
- [ ] Prometheus + Grafana
- [ ] OpenTelemetry in at least one service

### Users *(addresses + password reset completed with tests)*
- [x] users service skeleton (Go, chi, Postgres)
- [x] POST /register (email, password, name, phone, CPF; validation + duplicate email 409)
- [x] POST /login (returns JWT)
- [x] GET /me (profile, requires JWT)
- [x] PATCH /me (update name, phone, CPF; at least one field required)
- [x] CRUD addresses (POST/GET/PATCH/DELETE /me/addresses; type billing/shipping; default flags)
- [x] POST /password-reset/request, POST /password-reset/confirm
- [ ] K8s Deployment + Service

### Catalog
- [x] catalog service skeleton
- [x] CRUD categories
- [x] CRUD products (with seller_id)
- [ ] GET /products (list, pagination, filters)
- [ ] gRPC GetProduct, ListProducts
- [ ] K8s Deployment + Service

### Sellers
- [ ] sellers service skeleton
- [ ] POST /sellers (register)
- [ ] GET /sellers/:id
- [ ] GET /sellers/:id/products (or via catalog)
- [ ] K8s Deployment + Service

### Inventory
- [ ] inventory service skeleton
- [ ] Stock table (product_id, seller_id, quantity)
- [ ] gRPC ReserveStock, ReleaseStock
- [ ] Consume order.created → reserve (or orders calls gRPC)
- [ ] Consume payment.completed → reduce
- [ ] Consume order.cancelled → release
- [ ] K8s Deployment + Service

### Orders
- [ ] orders service skeleton
- [ ] Cart (in-memory or DB per user)
- [ ] POST /checkout (cart → order)
- [ ] gRPC call to inventory to reserve
- [ ] Emit order.created
- [ ] Listen payment.completed → update status, emit order.paid
- [ ] Listen order.paid → (or shipping listens)
- [ ] GET /orders (user history)
- [ ] K8s Deployment + Service

### Payments
- [ ] payments service skeleton
- [ ] POST /payments (create payment for order)
- [ ] Mock gateway (always success for dev)
- [ ] Emit payment.completed / payment.failed
- [ ] Refund endpoint
- [ ] K8s Deployment + Service

### Shipping
- [ ] shipping service skeleton
- [ ] gRPC GetShippingRates (for checkout)
- [ ] Listen order.paid → create shipment
- [ ] Emit shipment.shipped, shipment.delivered (mock)
- [ ] GET /shipments/:id/tracking
- [ ] K8s Deployment + Service

### Notifications *(implement last)*
- [ ] notifications service skeleton (Go, RabbitMQ consumer)
- [ ] Consume queue `notifications`; on `type:sms` send SMS
- [ ] Consume order.created, payment.completed, shipment.shipped, etc. (email)
- [ ] K8s Deployment + Service

### Reviews
- [ ] reviews service skeleton
- [ ] POST /reviews (product or seller)
- [ ] GET /reviews (by product, by seller)
- [ ] K8s Deployment + Service

---

## 7. Event Flow (Key Flows)

### Order creation
```
Client → Kong → orders (POST /checkout)
  orders → inventory (gRPC ReserveStock)
  orders → RabbitMQ (order.created)
  notifications ← order.created (send email)
```

### Payment
```
Client → Kong → payments (POST /payments)
  payments → RabbitMQ (payment.completed or payment.failed)
  orders ← payment.completed (update status, emit order.paid)
  inventory ← payment.completed (reduce stock)
  shipping ← order.paid (create shipment)
  notifications ← payment.completed (send email)
```

### Shipment
```
shipping → RabbitMQ (shipment.shipped, shipment.delivered)
  orders ← shipment.delivered (update status)
  notifications ← shipment.shipped, shipment.delivered
```

---

## 8. Data Model Hints (per service)

| Service | Main entities |
|---------|---------------|
| users | users, addresses (billing/shipping type, default per type) |
| catalog | categories, products, product_variants |
| sellers | sellers |
| inventory | stock (product_id, seller_id, quantity, reserved) |
| orders | orders, order_items, carts |
| payments | payments, refunds |
| shipping | shipments |
| reviews | reviews (product_id or seller_id, user_id, rating, comment) |
| notifications | (stateless; optional: notification_log) |

---

*This document is a living plan. Update phases and checkboxes as you progress.*
