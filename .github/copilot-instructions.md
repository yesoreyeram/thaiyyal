# Copilot coding agent instructions — thaiyyal

Purpose: give a new coding agent the essential, high‑value facts and actionable steps needed to develop, build, test, and validate changes in this repository without long exploratory searches. Trust these instructions first; search only if something here is missing or a command fails.

## Quick summary

- This repository is a monorepo with a Next.js frontend (Node) and a Go backend. The frontend builds static output which the Go server serves from backend/pkg/server/static.
- Primary languages: Go (backend) and TypeScript/JavaScript (Next.js frontend). There are also Docker + shell files for building and running.
- Useful entry points you will use often:
  - package.json (root) — frontend scripts (build, dev, copy:static, lint, playwright test)
  - go.mod (root) — Go module and Go version
  - backend/pkg/server — Go server sources and static assets location
  - Dockerfile — canonical multi‑stage build showing the intended build order and commands
  - CONTRIBUTING.md — contribution guidance (project conventions)

## Expected tool/runtime versions (use these or newer compatible)

- Go: use Go 1.24.x (go.mod indicates Go 1.24 series).
- Node: prefer Node 20.x (project dev deps target Node 20+; Node 18+ may work).
- npm / pnpm / yarn: the repo uses npm scripts; running `npm ci` or `npm install` in the root is the standard flow.
- Playwright: e2e tests use Playwright (dev dependency). Ensure Playwright system dependencies are available if running headed tests.

If your environment differs, run the build locally once and stop on the first error before making code changes; document any version mismatches.

## High‑level build & validation sequence (canonical order)

Follow this order when making changes and preparing a PR — it mirrors the Dockerfile and avoids common CI failures.

1. Prepare environment

   - Always: install Go 1.24.x and Node (20.x recommended).
   - Always: run `go env` and `node -v` to confirm versions.

2. Install Node deps (root)

   - From repo root:
     - `npm ci` (preferred for CI) or `npm install` (local dev).
   - Reason: frontend build and scripts rely on these deps.

3. Build frontend (Next.js) and produce static output

   - From repo root:
     - `npm run build` — builds the Next.js app (produces `out/` or `.next` per Next config).
     - Alternatively `npm run build:frontend` if present (it runs build + copy).
   - After a successful Next build:
     - `npm run copy:static` — this project expects static front output to be copied into backend/pkg/server/static (the Dockerfile and package.json use this path). Always run this before building the Go server.

   Note: If you prefer a single step, `npm run build:frontend` is the canonical combined step if available.

4. Prepare Go backend

   - From repo root:
     - `go mod download`
     - `go test ./...` — run Go unit tests (do this before building; fix any failures).
     - Build the server binary (example):
       - `cd backend/pkg/server` (or from repo root: `go build -o ./bin/thaiyyal ./backend/pkg/server`)
     - If your change touches both frontend and backend, ensure `npm run copy:static` ran first so static files are included in the build.

5. Lint and static checks

   - Frontend: `npm run lint` (runs eslint).
   - Backend: run any Go linter used by the project if present (e.g., golangci‑lint). If none is configured, run `go vet ./...`.

6. Run server locally
   - After build, run the Go server binary:
     - e.g., `./bin/thaiyyal` (or `cd backend/pkg/server && go run .` for dev).
   - The server listens on port 8080 by default (Dockerfile exposes 8080). Use `curl http://localhost:8080/health` or root path to confirm.

## Common pitfalls and mitigations

- Missing static frontend assets in the Go build:

  - Symptom: server returns 404s for frontend paths or tests fail.
  - Fix: always run `npm run build` then `npm run copy:static` (or `npm run build:frontend`) before `go build`.

- Node version mismatches:

  - Symptom: next build errors related to Node features.
  - Fix: switch to Node 20 (use nvm) and reinstall node_modules.

- Playwright failing in CI due to missing browsers:

  - Fix: `npx playwright install` and ensure CI installs browsers or use Playwright CLI setup used by the repo.

- Lint or formatting errors:
  - Run `npm run lint` and `go vet ./...` locally and fix before opening a PR.

## Project layout (high priority paths)

- package.json (repo root) — frontend scripts, build/copy steps, lint, e2e test scripts.
- go.mod (repo root) — Go module root and Go version.
- Dockerfile (repo root) — canonical two‑stage build that shows:
  - Stage 1: Node build (Next.js)
  - Stage 2: Go build and server image using built frontend static files
- backend/pkg/server — Go server code and static assets directory: backend/pkg/server/static
- src (or `src/`) and public — frontend code for the Next.js app (standard Next layout)
- CONTRIBUTING.md — repo contribution notes and expectations

Files you should inspect first when changing behavior:

1. package.json — understand npm scripts and copy paths
2. backend/pkg/server/main.go (or equivalent) — server entrypoint and static file handling
3. go.mod — confirm module and Go version
4. Dockerfile — to see the exact build steps CI will perform

## CI / workflows (what to replicate locally)

- Typical checks you should run locally before PR:
  - `npm ci && npm run build && npm run copy:static`
  - `go mod download && go test ./...`
  - `npm run lint`
  - `npm run test:e2e` (if your change touches end‑to‑end behavior)
- Replicate the Dockerfile build locally with `docker build .` if your change affects packaging.

## How to make minimal, safe changes

- Small, focused PRs that include:
  - Updated unit tests (Go or JS)
  - Lint fixes if new code triggers linter
  - Repro steps in PR description
- Run the full build flow locally and include brief CI reproduction steps in the PR description.

## When to search the repo

- Trust this file as the first source of truth.
- Only search the repo if:
  - A command here fails with a reproducible error.
  - You need to find a specific file referenced by a failing test or CI log.

## Code quality

- The code MUST be enterprise grade and follow best practices
- The code MUST follow security first principles
- The UI MUST be enterprise grade and consistent

## Final notes for the agent

- Always run installs and builds in the order above (install → frontend build → copy static → go build → tests → lint).
- Document any deviation you make and why (version pin, new dependency, changed copy path).
- If a test or build fails and the failure isn't obviously a missing install or static assets issue, run focused searches for the failing symbol/file and include that trace in your PR description.

If anything in the repository changes (new CI file, different build scripts), update this instructions file so future agents save time.
