# Repository Guidelines

## Project Structure & Module Organization
- `apps/` hosts domain services: `user/`, `social/`, `im/`, `msg/`, each with `api/`, `rpc/`, and domain logic. Place new handlers under the matching service and transport.
- `pkg/` stores shared utilities (`ctxdata`, `middleware`, `resultx`, `wuid`). Extend these before duplicating logic in services.
- `components/` contains Docker volumes and config for Redis, MySQL, Mongo, Kafka, APISIX, etc.; update when infrastructure requirements change.
- `deploy/` provides make fragments (`deploy/mk/`) and helper scripts for release pipelines; mirror any new automation here.
- `test/` captures cross-service Go scenarios, while unit tests live beside their sources in `apps/**/*_test.go`. Consult `doc/` for architecture notes.

## Build, Test, and Development Commands
- `source env.sh` populates `HOST_IP` before using Docker services.
- `docker-compose up -d` starts the local stack defined in `docker-compose.yaml`; stop with `docker-compose down`.
- `make <service>-dev` (e.g., `make user-api-dev`, `make im-ws-dev`) runs the service-specific release flow via `deploy/mk/*.mk`.
- `make release-test` executes every service target and is the pre-merge sanity check.
- `go test ./...` runs package tests; `go test ./test -run TestSocialFriendGroup` exercises integration suites selectively.

## Coding Style & Naming Conventions
- Format Go files with `gofmt` or `goimports`; keep import groups ordered std/lib/internal and rely on gofmt tabs.
- Keep packages lowercase with underscores sparingly (`im_rpc`), exported identifiers in PascalCase, and config structs ending in `Config`.
- HTTP and RPC handlers belong in `apps/<domain>/<api|rpc>/internal/handler/` and should end with `Handler` (e.g., `CreateFriendHandler`).
- Prefer structured logging and error wrapping via helpers in `pkg/resultx` and `pkg/xerr`; avoid bare `fmt` in service flows.

## Testing Guidelines
- Name tests `Test<Service><Behavior>` and place fixtures beside the code under test or in `test/` when they span services.
- Cover both happy-path and failure cases; reuse shared helpers from `test/test_helpers_test.go` for setup and teardown.
- Run `go test ./test` before opening a PR; include relevant `make <service>-dev` results when behavior spans transports.

## Commit & Pull Request Guidelines
- Use short Conventional Commit prefixes seen in history (`feat:`, `fix:`, `chore:`) followed by an imperative summary (≤ 50 chars).
- Reference the impacted service (`feat: social add group audit`) and link issues or tickets where possible.
- PR descriptions must state scope, highlight schema or environment changes, and list validation commands executed; add screenshots for API surface changes.
- Keep automation scripts in `deploy/script/` aligned with code updates and call out required Docker or config changes for reviewers.
