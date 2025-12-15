# edora — High-Performance Ecosystem

Minimal scaffold implementing the provided FULL_SPEC.md.

Run (using Docker Compose):

```bash
docker compose up --build
```

Services:
- backend: Go + Fiber on :8080
- db: Postgres 16
- redis: Redis 7
- web: Next.js on :3000

Load test with k6:

```bash
k6 run k6/load-test.js
```
