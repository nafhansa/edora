# MASTER PROJECT SPEC: High-Performance Ecosystem (Go + Next.js + Flutter)

**Version:** 1.0.0
**Status:** Development Ready
**Performance Target:** 1000+ Concurrent Users | Google Lighthouse Score 90+ | <200ms API Latency

---

## 1. High-Level Architecture Overview

This system is designed as a high-throughput micro-service architecture capable of serving web and mobile clients simultaneously via a unified RESTful API.

### The Stack
| Component | Technology | Reasoning |
| :--- | :--- | :--- |
| **Backend** | **Go (Golang) 1.22+** | Selected for raw speed, concurrency handling, and low memory footprint. |
| **Framework** | **Fiber v2** | Express-like syntax but built on `fasthttp` (fastest Go HTTP engine). |
| **Database** | **PostgreSQL 16** | Robust relational data with `pgx` driver for high-performance connection pooling. |
| **Caching** | **Redis 7** | Mandatory for session storage and caching hot data to handle 1000+ users. |
| **Web Frontend** | **Next.js 14 (App Router)** | Server Components for SEO and LCP/CLS optimization (Lighthouse 90+). |
| **Mobile App** | **Flutter** | Compiled native performance (60fps) for Android & iOS. |
| **Infrastructure** | **Docker Compose** | Containerized environment for consistent Dev/Prod parity. |

---

## 2. Backend Implementation (Go)

### A. Folder Structure (Standard Go Layout)
```text
/backend
  /cmd/api           # Application entry point (main.go)
  /internal
    /handler         # HTTP Handlers (Fiber Context)
    /service         # Business Logic (Decoupled from HTTP)
    /repository      # Database interactions (SQL)
    /middleware      # Auth, Logger, Rate Limiter
    /models          # Struct definitions (Request/Response)
  /pkg
    /database        # Postgres connection setup (pgxpool)
    /redis           # Redis client setup
    /utils           # Helper functions (Validator, JWT)
  /config            # Env loader


B. Core Performance Logic

    Connection Pooling: Use pgxpool. Configure MaxConns to 50-100 depending on CPU to prevent connection starvation.

    Concurrency (Worker Pools): Do NOT run heavy tasks (e.g., image processing, email sending) in the main HTTP handler. Dispatch them to a Go Routine Worker Pool.

    Caching Strategy (Look-Aside):

        Read: Check Redis key -> If exists, return JSON -> If nil, Query DB -> Write to Redis (TTL 5 min) -> Return.

        Write: Write to DB -> Invalidate/Delete Redis key immediately.

C. Database Schema (Optimization Ready)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users: Indexed on email for fast login lookups
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

-- Sessions: High churn, managed via Redis, but strictly logged here if needed

3. Web Frontend Implementation (Next.js)
A. Achieving Google Lighthouse 90+

    Image Optimization:

        Use <Image /> component strictly.

        Set priority={true} for LCP elements (Hero images).

        Explicitly set sizes prop.

    Font Loading: Use next/font/google with subset: ['latin'] to prevent Layout Shift (CLS).

    Bundle Splitting:

        Use Server Components by default.

        Only use 'use client' for interactive leaves (buttons, forms).

        Dynamic imports: const HeavyComponent = dynamic(() => import('./Heavy')).

B. State Management

    Global: Zustand (Minimal footprint).

    Server State: TanStack Query (React Query) v5.

        Why? Automatic background refetching, caching, and deduping requests keep the UI snappy.

4. Mobile App Implementation (Flutter)
A. Architecture (Feature-First + Riverpod)
/lib
  /src
    /features
      /auth          # Login/Register Screens & Logic
      /dashboard     # Home Logic
    /core
      /api           # Dio Client with Interceptors
      /storage       # Hive (Local DB)

B. Performance & Offline Mode

    Dio Interceptors: Handle JWT Refresh tokens automatically silently.

    Hive Storage: Cache critical JSON responses locally.

        Logic: App Start -> Load data from Hive (Instant UI) -> Fetch API in background -> Update UI & Hive.

    Image Caching: Use cached_network_image.

5. Testing Strategy (QA & Load Testing)
A. Load Testing (Proof of 1000+ Users)

Tool: k6. Save as load-test.js.
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 500 },  // Ramp up
    { duration: '1m', target: 1000 },  // 1000 concurrent users
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p95<200'],    // 95% of requests must be < 200ms
    http_req_failed: ['rate<0.01'],    // < 1% errors
  },
};

export default function () {
  const res = http.get('http://localhost:8080/api/v1/products');
  check(res, { 'status was 200': (r) => r.status == 200 });
  sleep(1);
}

6. End-to-End Development Setup (Docker)

File: docker-compose.yml
version: '3.8'
services:
  # The Brain
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_URL=postgres://user:pass@db:5432/appdb
      - REDIS_ADDR=redis:6379
    depends_on:
      - db
      - redis

  # The Data
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: appdb
    volumes:
      - pgdata:/var/lib/postgresql/data

  # The Speed (Cache)
  redis:
    image: redis:7-alpine

  # The Face
  web:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080

volumes:
  pgdata:

