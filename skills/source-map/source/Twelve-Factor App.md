# Twelve-Factor App Methodology

**Source**: [12factor.net](https://12factor.net/) by Adam Wiggins
**Purpose**: A methodology for building modern, scalable, maintainable software-as-a-service applications. Applicable to any language and any combination of backing services.

## When to Apply

Use this methodology when designing or reviewing any service that:
- Deploys to cloud infrastructure (Kubernetes, cloud VMs, PaaS)
- Needs to scale horizontally
- Requires continuous deployment
- Must minimize divergence between development and production
- Is part of a microservices or distributed system architecture

## I. Codebase

> One codebase tracked in revision control, many deploys

### What It Means

Exactly **one codebase** (one repo, or repos sharing a root commit) per app, with many **deploys** of it -- production,
staging, each developer's machine -- possibly at different commits. Multiple codebases means it is a distributed
system, not one app; multiple apps sharing code is a violation -- factor shared code into a library.

### How to Apply

- One Git repository per service/application.
- Shared code lives in separate library repositories, pulled via dependency manager (`go get`, `npm install`, etc.).
- Every environment (dev, staging, prod) deploys from the same repository -- different tags/branches, same repo.

## II. Dependencies

> Explicitly declare and isolate dependencies

### What It Means

Never rely on the implicit existence of system-wide packages: declare every dependency completely and exactly in a
manifest, and use an isolation tool so nothing leaks in from the surrounding system. The full dependency specification
applies uniformly to production and development alike.

### How to Apply

- All dependencies in `go.mod` (or equivalent). No implicit reliance on system tools.
- If the app requires a system tool (e.g., `ImageMagick`, `ffmpeg`), vendor it into the app or bundle it in the container image.
- A new developer should need only the language runtime and dependency manager to build and run the app.

## III. Config

> Store config in the environment

### What It Means

Config is everything likely to vary between deploys -- credentials, connection strings, hostnames, feature flags, log
levels -- and never internal application structure (routes, DI wiring), which belongs in code. **Litmus test**: could
the codebase be open-sourced right now without exposing any credentials? Config lives in **environment variables**.

### How to Apply

- All secrets, connection strings, and per-deploy values come from environment variables.
- Use a config struct that reads from env vars at startup (e.g., `env` tags in Go, `os.Getenv`).
- Never hardcode credentials, URLs, or environment-specific values in source code.
- Avoid "environment groups" like `config/production.yaml` -- each variable is independently managed.
- Config files are acceptable for internal application wiring (route definitions, DI configuration) since they don't vary between deploys.

## IV. Backing Services

> Treat backing services as attached resources

### What It Means

A **backing service** is anything the app consumes over the network: datastores, message queues, SMTP, caches,
monitoring, third-party APIs. The code makes **no distinction between local and third-party services** -- both are
attached resources reached via config, so swapping one for another is **only a config change, zero code changes**.

### How to Apply

- All backing service connections are configured via environment variables (Factor III).
- Use interfaces/abstractions for storage, messaging, and external services so implementations are swappable.
- Never hardcode hostnames, ports, or credentials for any backing service.
- Each distinct backing service (e.g., two different databases) is a distinct attached resource with its own config entry.

## V. Build, Release, Run

> Strictly separate build and run stages

### What It Means

Three strictly separated stages: **build** (code repo into an executable bundle at a specific commit), **release**
(build combined with deploy-specific config, stamped with a unique release ID), **run** (start processes against a
selected release). Releases are immutable, code cannot change at runtime, and rollback deploys a previous release.

### How to Apply

- CI/CD pipeline builds a Docker image (build stage), tags it with a release ID, and pushes it to a registry.
- Deployment combines the image with environment-specific config (release stage).
- The orchestrator (Kubernetes) runs the image (run stage).
- Never SSH into a container to edit code. Never modify a running release.
- Keep the run stage as simple as possible -- if it breaks at 3 AM, it should auto-recover without developer intervention.

## VI. Processes

> Execute the app as one or more stateless processes

### What It Means

Processes are **stateless and share-nothing**; anything that must persist goes to a stateful backing service. Memory
and local disk may serve only as a brief single-transaction cache -- never assume a later request hits the same
process. **Sticky sessions are a violation**; session state belongs in a datastore with time-expiration.

### How to Apply

- No in-process session state. Use Redis/Memcached for sessions.
- No local filesystem writes that other processes or future requests depend on. Use object storage (S3) or a database.
- Design every process to be killable and replaceable at any moment without data loss.
- Asset bundling and code generation happen at build time, not runtime.

## VII. Port Binding

> Export services via port binding

### What It Means

The app is **completely self-contained**: it binds to a port and exports its service itself, never relying on a
webserver injected into the execution environment. This applies to **any protocol**, not just HTTP. In production a
routing layer maps the public hostname to the bound port; one app's service can be another's backing service.

### How to Apply

- The app starts its own HTTP/gRPC server. No external webserver container required.
- The port to bind to comes from config (environment variable, e.g., `PORT=8080`).
- In Kubernetes, the container listens on a port; the Service/Ingress routes traffic to it.

## VIII. Concurrency

> Scale out via the process model

### What It Means

Processes are a first-class citizen: each type of work gets its own **process type** (web, worker, clock), and the
array of types with their counts is the **process formation**. The app **scales horizontally by running more
processes**, not larger ones, and should **never daemonize or write PID files** -- the OS process manager does that.

### How to Apply

- Scale by increasing replicas of a process type (Kubernetes `replicas: N`), not by giving one instance more CPU/RAM.
- Different workloads run as different process types (separate Deployments in Kubernetes): `web`, `worker`, `scheduler`.
- The app does not daemonize itself, manage PID files, or trap signals for self-restarting. Let the orchestrator handle it.
- Internal concurrency (goroutines, threads) is fine but not a substitute for horizontal scaling.

## IX. Disposability

> Maximize robustness with fast startup and graceful shutdown

### What It Means

Processes are **disposable** -- startable and stoppable at a moment's notice, ideally a few seconds from launch to
ready-to-serve. On `SIGTERM` a process stops accepting new requests, finishes in-flight ones, then exits; a worker
returns its current job to the queue. Jobs must be **reentrant**, and the architecture embraces **crash-only design**.

### How to Apply

- Optimize startup: defer heavy initialization, use connection pooling, minimize bootstrap work.
- Implement graceful shutdown: handle `SIGTERM`, drain in-flight requests, close connections cleanly.
- Make all background jobs idempotent and reentrant -- safe to retry after a crash.
- Kubernetes readiness/liveness probes should reflect actual readiness, not just process existence.
- Use `preStop` hooks or `terminationGracePeriodSeconds` in Kubernetes to allow graceful drain.

## X. Dev/Prod Parity

> Keep development, staging, and production as similar as possible

### What It Means

Keep three gaps small: **time** (hours from commit to deploy, not weeks), **personnel** (whoever writes the code is
involved in deploying and observing it), **tools** (the same backing services everywhere). The twelve-factor developer
**resists the urge** to swap in a different service for dev, even when an adapter library appears to hide the gap.

### How to Apply

- Use the same database, cache, queue, and message broker in dev as in prod. Docker Compose for local development.
- Deploy continuously -- code merged today should reach production today (or within hours).
- Developers observe their code in production (logs, metrics, alerts).
- Never use SQLite locally and PostgreSQL in production, or in-memory cache locally and Redis in production.
- Infrastructure-as-code ensures staging and production are structurally identical.

## XI. Logs

> Treat logs as event streams

### What It Means

Logs are a **stream of aggregated, time-ordered events**, typically one event per line. A twelve-factor app **never
concerns itself with routing or storage of its output stream**: it writes no logfiles and emits its event stream
**unbuffered to `stdout`**. The execution environment captures and routes it to destinations the app cannot configure.

### How to Apply

- Log to `stdout` / `stderr`. Period. No log files, no log rotation, no log management in the app.
- Use structured logging (JSON lines) for machine-parseable output.
- In Kubernetes, container stdout is automatically captured. Use Fluentd/Fluent Bit or a log collector to route to Elasticsearch, Loki, etc.
- Logging configuration (level, format) can come from environment variables (Factor III).
- Never use `log.Fatal` (or equivalent) in request handlers -- it kills the process. Use error returns.

## XII. Admin Processes

> Run admin/management tasks as one-off processes

### What It Means

One-off administrative tasks -- database migrations, a REPL against live data, one-time data fixes -- run as **one-off
processes** in an environment **identical** to the long-running ones: same **release** (same code and config), same
dependency isolation, same backing services. Admin code **ships with application code**, same codebase, same image.

### How to Apply

- Database migrations are part of the application codebase and run as a step in the deployment pipeline (init container in Kubernetes, or a CI/CD job using the same image).
- One-off admin tasks run via `kubectl exec` or a Kubernetes Job using the same container image and config.
- Never run admin tasks from a developer laptop against production data. Use the same image, same config, same environment.
- Admin scripts live in the same repository as the application code.

## Quick Reference

| # | Factor | Rule | Violation Example |
|---|--------|------|-------------------|
| I | Codebase | 1 repo = 1 app, many deploys | Sharing code by copy-paste between repos |
| II | Dependencies | Declare and isolate all deps | Relying on a globally installed system tool |
| III | Config | Environment variables | Hardcoded DB password in source code |
| IV | Backing Services | Attached resources via config | `if env == "prod" { useRDS() } else { useSQLite() }` |
| V | Build, Release, Run | Immutable releases, strict stages | SSH into prod to hotfix code |
| VI | Processes | Stateless, share-nothing | Sticky sessions, local file uploads |
| VII | Port Binding | Self-contained, export via port | Deploying a WAR file into Tomcat |
| VIII | Concurrency | Scale via process model | Scaling vertically to a bigger VM instead of adding replicas |
| IX | Disposability | Fast startup, graceful shutdown | 5-minute startup, no SIGTERM handling |
| X | Dev/Prod Parity | Same stack everywhere | SQLite in dev, PostgreSQL in prod |
| XI | Logs | stdout, no file management | Writing to `/var/log/app.log` and rotating manually |
| XII | Admin Processes | One-off, same environment | Running migrations from a dev laptop against prod DB |
