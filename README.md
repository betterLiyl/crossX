# crossX

Cross-platform Matching and Liquidity Aggregation engine for CEX and DEX.

## Description
crossX is a high-performance trading system that aggregates liquidity and matches orders across centralized exchanges (CEX) and decentralized exchanges (DEX). It provides a modular architecture with a matching engine, protocol adapters, gateway services, order management, settlement coordination, and optional on-chain components for Solana.

## Key Features
- Minimal matching engine with price-time priority and partial/complete fills
- Unified protocol adapters for CEX/DEX (OKX implemented for market data; Hypeliquld stub)
- Gateway service with JWT auth, RESTful endpoints, and OpenAPI documentation
- Standardized order event ingestion (push HTTP; WS/Kafka planned; pull adapters planned)
- Storage configuration stubs for PostgreSQL and Redis
- Containerization for gateway service

## Technology Stack
- Core services: Go (concurrency, performance)
- On-chain (future): Rust for Solana programs
- Database: PostgreSQL (relational), Redis (cache)
- Messaging: Kafka/RabbitMQ (planned)
- Monitoring: Prometheus + Grafana (planned)

## Module Structure
- `cmd/gateway`: HTTP gateway, JWT middleware, routes, OpenAPI serving
- `internal/engine`: Matching engine, order book, aggregation of top-of-book
- `internal/adapters`: Market data adapters (OKX implemented, Hypeliquld placeholder)
- `internal/ingest/common`: Standardization, validation, signature verification components
- `internal/ingest/push`: HTTP push handler (WS/Kafka planned)
- `internal/storage`: Configuration stubs for Postgres/Redis
- `api/swagger.json`: OpenAPI spec for gateway endpoints

## Installation & Setup
Requirements: Go 1.21+

```bash
git clone <repo-url>
cd crossX
go build ./...
go test ./... -cover
```

Run gateway locally:
```bash
PORT=8080 JWT_SECRET= go run ./cmd/gateway
```

Docker:
```bash
docker build -f cmd/gateway/Dockerfile -t crossx-gateway .
docker run -p 8080:8080 crossx-gateway
```

Environment variables:
- `PORT`: gateway port (default `8080`)
- `JWT_SECRET`: shared secret for HMAC JWT verification (empty allows bearer header without signature check)
- `OKX_BASE_URL`: override OKX API base (optional)
- `POSTGRES_DSN`, `REDIS_ADDR`: storage configuration (stubs)
- `PUSH_SECRET`: HMAC secret for push payload verification (optional)

## Usage & API
Endpoints:
- `GET /health` → health check
- `GET /swagger.json` → OpenAPI document
- `POST /orders` → place order
- `GET /orders/{id}` → get order status
- `GET /aggregate/{symbol}` → aggregate market data across adapters
- `POST /v1/orders/push/http` → push order events (standardized JSON)

Examples:
```bash
curl -H "Authorization: Bearer a.b.c" \
     -H "Content-Type: application/json" \
     -d '{"symbol":"BTC-USDT","side":"BUY","price":100,"size":1}' \
     http://localhost:8080/orders

curl -H "Authorization: Bearer a.b.c" http://localhost:8080/aggregate/BTC-USDT

curl -H "Authorization: Bearer a.b.c" \
     -H "X-Signature: <hex-hmac>" \
     -H "Content-Type: application/json" \
     -d '{"order_id":"o1","symbol":"BTC-USDT","price":100,"size":1,"side":"BUY","ts":"2025-01-01T12:00:00Z","platform":"mock","signature":"<optional>","seq":1}' \
     http://localhost:8080/v1/orders/push/http
```

## Contribution Guidelines
- Use Go 1.21+, follow existing code style and module boundaries
- Add unit tests for new functionality; target coverage ≥80%
- Update `api/swagger.json` when changing endpoints
- Avoid committing secrets/credentials; use environment variables
- Open issues with clear steps to reproduce and expected behavior

## License
MIT (pending confirmation). If contributing, include license headers where appropriate.

## Roadmap & Milestones
Priorities:
1. Push path: WebSocket server and Kafka consumer with signature verification and idempotency
2. Pull path: Scheduler with 30s polling, per-platform rate limiting, checkpoint persistence
3. Metrics & alerts: Prometheus metrics and Grafana dashboards, Kafka consumer lag monitoring
4. Security hardening: mTLS, key rotation, replay protection, audit logging
5. Documentation: Platform onboarding guide, detailed API examples

Implementation Plans:
- Push/WS/Kafka: implement `internal/ingest/push/ws_server.go`, `kafka_consumer.go`; add signature, nonce, and time window checks; wire to processing pipeline
- Pull/scheduler/checkpoint: implement `internal/ingest/pull/scheduler.go`, `checkpoint.go`; define `FetchChanges(since)` for platform adapters; add rate-limiting and retry with backoff
- Metrics: create `internal/ingest/metrics/metrics.go`, expose `/metrics` and instrument ingest pipeline
- Security: configure TLS/mTLS, integrate key management, enforce headers and replay protection
Milestones & Deliverables:
- MVP Push: HTTP + WS + Kafka ingestion with validation, metrics, and basic dashboards
- MVP Pull: OKX/Hypeliquld adapters with checkpointing and retries
- Hardening: security and observability improvements, CI pipelines and tests

Testing & Quality
- Unit tests for normalization, validation, signature verification, scheduler logic
- Integration tests for WS/HTTP/Kafka ingestion and checkpoint persistence
- Static analysis and linting; keep test coverage ≥80%
