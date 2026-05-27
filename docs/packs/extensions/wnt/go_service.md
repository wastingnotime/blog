# WNT Pack Extension: go_service

## Purpose

Use this pack when a repository is centered on a Go backend service with explicit domain, application, infrastructure, runtime, and HTTP adapter boundaries.

This pack is useful for API-first services that need deterministic local validation, a packaged runtime artifact, and clear separation between business behavior and transport or deployment concerns.

## Scope

This pack extension defines WNT-owned implementation defaults.
It does not define MRL core behavior or project-specific domain behavior.

This pack exists to define reusable WNT defaults for MRL repositories that share a Go service runtime shape.

---

## Shape

- language: Go
- runtime topology: single service runtime plus supporting command entrypoints
- architecture mode: service / layered application

---

## Good Fit

- API-first backend services
- repositories that expose one long-running service runtime
- systems where use cases should stay independent from HTTP, persistence, and deployment details
- repositories that need deterministic local scenarios plus packaged runtime validation
- services with explicit persistence adapters and runtime configuration

---

## Less Suitable

- repositories that require a separate frontend and backend model in the same pack
- event-sourced systems where an event log is the primary source of truth
- repositories with multiple independent application runtimes in different languages
- libraries that do not produce a runnable service artifact

---

## Layout

```text
cmd/
  api/
  run-scenario/
  run-packaged-runtime/

internal/
  domain/
  application/
  httpapi/
  infrastructure/
  runtime/

tests/
  unit/
  integration/
```

Optional release or lifecycle helper commands may live under `cmd/` when the repository has explicit release or expose requirements, but those helpers should not redefine the service architecture.

---

## Boundaries

- business rules live in `internal/domain/`
- orchestration lives in `internal/application/`
- persistence and external adapters live in `internal/infrastructure/`
- HTTP request and response translation lives in `internal/httpapi/`
- runtime configuration and dependency composition live in `internal/runtime/`
- long-running service startup lives in `cmd/api/`
- deterministic scenario execution lives in `cmd/run-scenario/`
- packaged runtime smoke validation lives in `cmd/run-packaged-runtime/`

HTTP handlers should translate requests, invoke application behavior, map errors, and write responses. They should not own business workflow logic.

---

## Testing Defaults

- unit/spec focus: domain invariants and use-case behavior
- integration focus: runtime wiring, persistence adapters, HTTP behavior, and packaged runtime behavior
- determinism rules: keep time, IDs, external responses, persistence state, and event order deterministic in tests
- preferred test command: `go test ./...`

Integration tests should favor local fakes, in-memory adapters, temporary directories, SQLite, or equivalent local-safe dependencies before requiring external services.

---

## Scenario Runner Shape

- preferred runner location: `cmd/run-scenario/`
- preferred command: `go run ./cmd/run-scenario`
- runner responsibility: exercise one deterministic service scenario and print structured evidence, preferably JSON
- runner must not: depend on external services for its default path

When a packaged artifact exists, add a separate smoke path:

- preferred runner location: `cmd/run-packaged-runtime/`
- preferred command: `go run ./cmd/run-packaged-runtime`
- runner responsibility: validate the configured runtime artifact through the same application behavior used by the service
- runner must not: introduce business behavior that is unavailable through the application layer

---

## Runtime Targeting Rules

- slices using this pack should declare runtime targets as:
  - `Go service runtime`
  - `Go command entrypoints` when scenario, packaging, release, or expose helpers are affected
  - `container image` when the slice changes the packaged runtime contract
- API contracts should be described in the slice, contract artifact, or repository-specific API documentation
- event/message boundaries should be expressed through application ports and infrastructure adapters
- deployment or promotion mechanics should live in expose extensions, lifecycle scripts, or repository decisions, not inside the pack itself

---

## API Surface Defaults

- provide a health endpoint when the service is intended to run as a long-lived API
- keep HTTP status mapping explicit and deterministic
- reject malformed request payloads clearly
- keep CORS, auth claims, tracing, and other interface concerns at the HTTP or runtime boundary
- avoid framework magic that hides route behavior, request decoding, or dependency wiring from tests

The pack does not require a specific Go HTTP framework. The standard library is an acceptable default until a framework provides clear value.

---

## Runtime Configuration Defaults

- parse environment-driven configuration in `internal/runtime/`
- fail clearly and early on invalid configuration
- keep local defaults safe and deterministic
- support an in-memory or local persistence mode when practical
- support a production persistence mode only when a slice or repository decision requires it
- centralize runtime cleanup so tests and command entrypoints can close resources consistently

Repository-specific environment variable names, secret names, ports, and persistence engines belong in project artifacts, not in this generic pack.

---

## Packaging Defaults

- preferred packaged artifact: container image
- preferred container shape: multi-stage Go build with a small runtime image
- packaged runtime should start without external services by default when practical
- packaged runtime should expose the same application behavior validated by unit and integration tests

Publishing targets, registries, repository dispatch events, and infrastructure pull requests are lifecycle concerns. Record them in expose extensions, change artifacts, or `decisions.md` rather than making them part of this pack.

---

## Naming Conventions

- use cases or commands: verb-led names such as `CreateContact`, `PlaceOrder`, or `PublishInvoice`
- domain objects: singular business nouns such as `Contact`, `Order`, or `Invoice`
- repositories/gateways: intention-revealing names such as `ContactRepository` or `PaymentGateway`
- events: past-tense or fact-shaped names such as `contact.created` or `invoice.published`
- packages: short responsibility names such as `domain`, `application`, `httpapi`, `runtime`, and adapter-specific infrastructure package names

---

## Rules

- keep the application layer free of HTTP, SQL, framework, and deployment concerns
- keep infrastructure adapters small and explicit
- keep runtime composition in one obvious boundary
- prefer deterministic local validation before external integration
- make API contract changes explicit in slice or contract artifacts
- record meaningful pack deviations in `decisions.md`

---

## Starter Tooling

- package manager: `go modules`
- formatter/linter: `gofmt` and repository-selected linting if needed
- test framework: Go `testing` package
- persistence default: in-memory or SQLite for local deterministic validation
- packaged runtime default: container image when the repository needs `release` or `expose`

---

## Notes

This pack extension is a reusable WNT default, not a law. If a repository diverges meaningfully, record that in `decisions.md`.
