# AGENTS.md

## 1. Purpose of this Repository

This repo is based on **evrone/go-clean-template**. It is a **Clean Architecture Golang service** that exposes:

- REST API (via **Fiber**)  
- gRPC (via protobuf)  
- AMQP RPC (RabbitMQ)  
- NATS RPC  

We will use this template to build a **multi-exam learning platform backend** (NEET PG, NEET UG, JEE, UPSC), including:

- User-facing APIs (practice, exams, podcasts, wallet, revision, referral, etc.)
- Admin APIs (CRUD for questions, subjects, topics, exams, podcasts, coupons, users, AI settings, analytics)

The **key requirement**:  
> **Always respect Clean Architecture: inner layers (entities, usecases) must NOT depend on outer layers (controller, DB, HTTP, Fiber, RabbitMQ, etc.).**

---

## 2. Tech Stack & Tools

- **Language**: Go
- **HTTP Framework**: [Fiber](https://github.com/gofiber/fiber)
- **Swagger / Docs**: [swag](https://github.com/swaggo/swag)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **JSON**: [goccy/go-json](https://github.com/goccy/go-json)
- **Query Builder**: [Squirrel](https://github.com/Masterminds/squirrel)
- **Migrations**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Metrics**: Prometheus via fiberprometheus
- **Testing**: testify
- **Mocking**: go.uber.org/mock

The template already knows how to:

- Run with `make run`
- Start infra with `make compose-up`
- Serve swagger at `/swagger`

**Agents should fit into this setup, not replace it.**

---

## 3. Project Structure & Layering Rules

Important directories:

- `cmd/app/main.go`
  - Entry point: loads config, logger, and calls `internal/app.Run()`.
- `internal/app`
  - Wires dependencies (repositories, usecases, controllers).
  - Starts servers (HTTP, gRPC, AMQP, NATS).
- `internal/controller`
  - Entry points for external protocols:
    - `internal/controller/http` – Fiber routes (REST).
    - `internal/controller/grpc` – gRPC handlers.
    - `internal/controller/amqp_rpc` – RabbitMQ.
  - **This is the only layer that should know about Fiber, HTTP, etc.**
- `internal/entity`
  - Pure domain **entities** (models). No imports from Fiber, pgx, etc.
- `internal/usecase`
  - **Business logic layer** (use cases).
  - Defines usecase structs that depend on interfaces (repositories, external services).
- `internal/repo/persistent`
  - Concrete DB implementations (Postgres via Squirrel / pgx).
- `internal/repo/webapi`
  - Integrations with external HTTP APIs or other services (if any).
- `pkg/*`
  - Shared utilities like RabbitMQ client, logging, etc.
- `config`
  - Config structs & loading from env vars.

### 3.1 Clean Architecture Rules (hard constraints)

**Agents MUST follow these rules**:

- `internal/entity`:
  - Can only use Go standard library (no Fiber, no DB).
  - Contains domain structs, enums, and maybe simple validation methods.

- `internal/usecase`:
  - No direct imports of Fiber, HTTP, DB drivers, or external frameworks.
  - Depends on interfaces defined in its own package or in `internal/usecase/...` subpackages.
  - Implements business logic methods that controllers call.

- `internal/controller`:
  - Can import `fiber`, `zerolog`, etc.
  - Accepts HTTP requests, parses params, calls usecases, and writes HTTP responses.
  - Translates between HTTP DTOs and `entity` structs.

- `internal/repo/persistent`:
  - Implements repository interfaces defined in `usecase`.
  - Knows about DB connection, SQL queries, and migrations.

**Direction of dependencies** (must be respected):

```text
controller (HTTP / gRPC / AMQP / NATS)
        ↓
      usecase
        ↓
 repository (persistent / webapi)
        ↓
  database / external APIs
