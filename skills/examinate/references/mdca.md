# Modular Domain-Centric Architecture (MDCA) — Standard

- Version: 1.0
- Status: Active

---

## Abstract

This document defines **Modular Domain-Centric Architecture (MDCA)** — a software architecture standard for building business applications, optimized for Go services and clients but applicable to any modern statically-typed language. MDCA is a synthesis and refinement of three established schools of thought: **Domain-Driven Design (DDD)**, **Hexagonal Architecture (Ports and Adapters)**, and **Clean Architecture**. It distills their durable contributions, removes ceremony that does not pay for itself in practice, and prescribes a consistent module layout that scales across teams and microservices.

MDCA does not invent new ideas. It standardizes a pragmatic application of existing ones.

---

## 1. Conformance and Requirement Levels

This document is the single normative source for MDCA. It uses the key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** as defined in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119), to indicate requirement levels.

A codebase is **MDCA-conformant** if it satisfies all **MUST** clauses in §6 and §7. Codebases that satisfy the **SHOULD** clauses additionally are **MDCA-recommended**.

Section, clause, and rule numbers are stable citation anchors: cite a principle clause as `P3.2` and a tactical rule as `§8.6`. §2, §5, and §14 were retired; their absence is intentional and the surviving numbers were not shifted.

---

## 3. Terminology

| Term | Definition |
|------|------------|
| **Domain** | The problem space — the business operations the software exists to support. |
| **Bounded Context** | A delimited region of the domain in which a single ubiquitous language and model apply consistently. |
| **Ubiquitous Language** | The vocabulary shared by domain experts and engineers, expressed identically in conversation, code, and documentation. |
| **Module** | An independently replaceable unit of code that owns its data and exposes a stable contract. In MDCA, the canonical module is one bounded context. |
| **Port** | An interface declared by the domain that names a capability it requires from infrastructure. |
| **Adapter** | A concrete implementation of a port; lives outside the domain. |
| **Composition Root** | The single location where concrete adapters are instantiated and injected into the domain. |
| **Domain Event** | A past-tense fact about something that has occurred in the domain, immutable after publication. |
| **Application Service** | A coordinator that translates external requests into domain operations and manages the boundaries of a single use case. |

---

## 4. Design Goals

MDCA-conformant systems prioritize, in order:

1. **Performance efficiency** — the architecture imposes no overhead the workload does not require.
2. **Reliability** — the system behaves predictably under stress and failure.
3. **Maintainability** — change in one place does not ripple into unrelated places.

When tension arises between these goals and *flexibility*, *transferability*, or *compatibility*, the goals above prevail unless an architectural exception is documented in the initiative's design record.

---

## 6. Principles (Normative)

A conformant codebase **MUST** satisfy every clause in this section. Each clause carries a citable anchor of the form `P<principle>.<clause>`.

### P1 — Domain Centricity

| Clause | Requirement |
|--------|-------------|
| `P1.1` | Code **MUST** be organized around business capabilities, not technical layers. |
| `P1.2` | Top-level groupings such as `controllers/`, `models/`, or a flat `services/` dump **MUST NOT** be the primary organizing structure. |
| `P1.3` | The primary organizing unit **MUST** be the bounded context. |

### P2 — Ubiquitous Language

| Clause | Requirement |
|--------|-------------|
| `P2.1` | Types, methods, packages, and tests **MUST** use the vocabulary of domain experts. |
| `P2.2` | When the business term is awkward in code, the awkward term **MUST** be preferred over a translated euphemism. |
| `P2.3` | Generic names such as `Manager`, `Processor`, `Data`, `Info`, `Util` **MUST NOT** be used as type names within the domain. |

### P3 — Dependency Inversion

| Clause | Requirement |
|--------|-------------|
| `P3.1` | Source-code dependencies **MUST** flow inward: presentation depends on application, application depends on domain. |
| `P3.2` | The domain layer **MUST NOT** import infrastructure. Whether a given package counts as infrastructure is decided by the purity test in §7.1, not by the folder it lives in. |
| `P3.3` | Adapters **MUST** import the domain to satisfy its ports; the reverse is forbidden. |

### P4 — Modular Independence

| Clause | Requirement |
|--------|-------------|
| `P4.1` | Each bounded context **MUST** be replaceable without modification to other bounded contexts. |
| `P4.2` | Cross-context type sharing **MUST NOT** occur — when two contexts model "the same" concept, each defines its own struct. |

### P5 — Pragmatic Abstraction

| Clause | Requirement |
|--------|-------------|
| `P5.1` | Abstractions **MUST NOT** be introduced before a second concrete need exists. |
| `P5.2` | An interface with a single implementation **SHOULD** be deleted, with the concrete type used directly, until a second implementation appears. |

### P6 — Explicit Composition

| Clause | Requirement |
|--------|-------------|
| `P6.1` | Concrete adapters **MUST** be instantiated in a single composition root (in Go services: `cmd/<svc>/main.go` and `internal/app/`). |
| `P6.2` | Service locators, global containers, and `init()`-based dependency wiring **MUST NOT** be used. |

### P7 — Event-Driven Internal Communication

| Clause | Requirement |
|--------|-------------|
| `P7.1` | Within a system of services, cross-context communication **SHOULD** occur via domain events on a message broker (NATS in the reference implementation). |
| `P7.2` | Synchronous cross-context calls **MAY** be used where latency or atomicity demands them, but **MUST** be documented in the initiative's design record. |

### P8 — Anticorruption at External Boundaries

| Clause | Requirement |
|--------|-------------|
| `P8.1` | External integrations (third-party SDKs, vendor APIs) **MUST** be encapsulated in adapter packages that translate between vendor types and domain types. |
| `P8.2` | Vendor types **MUST NOT** appear in domain signatures. |

### P9 — DTO/Model Separation

| Clause | Requirement |
|--------|-------------|
| `P9.1` | Types crossing the transport boundary (HTTP, gRPC, events) **MUST** be distinct from domain entities and value objects. |
| `P9.2` | DTOs **MUST NOT** carry persistence tags (`db:""`); domain models **MUST NOT** carry transport tags (`json:""`, protobuf tags). |

### P10 — Quality-Goal Discipline

| Clause | Requirement |
|--------|-------------|
| `P10.1` | When a design choice trades the goals in §4 against flexibility/transferability/compatibility, the §4 goals **MUST** prevail unless a deliberate exception is recorded in the initiative's Architecture Inception Canvas. |

---

## 7. Layering and Layout (Normative)

### 7.1 Layers

A conformant codebase **MUST** distinguish the following layers:

| Layer | Contents | Imports allowed |
|-------|----------|-----------------|
| **Domain** | Entities, value objects, ports, domain events, domain services | Pure packages only — see the purity test below |
| **Application** | Composition root, dependency wiring, use-case orchestration when not collapsed onto domain methods | Domain, infrastructure adapters |
| **Infrastructure** | Repository implementations, brokers, external clients, ACLs | Domain (to satisfy ports), shared infra packages |
| **Presentation** | Handlers, DTOs, routers, UI input systems | Application, domain (read-only) |
| **Shared Infrastructure** | Cross-cutting utilities (logger, config, HTTP framework wrapper) | No business logic; no service-internal imports |

**Purity test for domain imports** (referenced by `P3.2`). A package is importable by the domain when **both** hold:

1. its own imports are confined to the standard library, transitively — no third-party SDK, and no other project package that itself fails this test; and
2. it performs no I/O and holds no mutable global state — no network, no filesystem, no `database/sql`, no process or environment access.

The folder a package lives in is not the test. A package under `pkg/` **MAY** be imported by the domain when it passes both conditions, and a package anywhere **MUST NOT** be imported by the domain when it fails either. Purity is transitive and **MUST** be re-checked when an imported package's own import list changes: a pure package that later grows an HTTP client becomes an illegal domain import.

**Controllable time.** A domain that needs controllable time **MUST** declare a clock port and receive an implementation from the composition root. Importing a concrete clock package that passes the purity test is a testability decision, not a layering violation.

### 7.2 Reference Layout (Go Server)

```
services/<svc>/
  cmd/<svc>/main.go               # composition root entry
  internal/
    app/                          # DI container
    domain/<feature>/             # one bounded context per package
      model.go                    # entities, value objects (db:"" only)
      port.go                     # interfaces the domain requires
      service.go                  # use cases
      service_test.go
      events.go                   # domain events
    dto/<feature>/                # request.go, response.go (json:"" only)
    repository/<feature>.go       # adapter implementing port
    handler/<feature>.go          # presentation
pkg/                              # shared infrastructure
```

Single-domain services **MAY** flatten the `<feature>/` subdirectory.

### 7.3 Module Granularity

The canonical module unit is **one bounded context per package**. A bounded context **MUST**:

1. Own its data; expose access only via its public API or domain events.
2. Be deletable, rewritable, or replaceable without touching unrelated contexts.
3. Be testable without instantiating unrelated contexts.

A package that fails any of these tests is not a valid bounded context.

---

## 8. Tactical Rules (Normative)

### 8.1 Entities

- An entity **MUST** be identified by a typed ID field (e.g., `OrderID`, not bare `string`).
- Methods that mutate state **MUST** use pointer receivers; methods that read state **SHOULD** use value receivers.
- Entity methods **MUST** express domain operations (`Cancel`, `Approve`, `Reschedule`); naive setters (`SetStatus`) **MUST NOT** be exposed when they bypass invariants.
- Construction outside the entity's package **MUST** go through a constructor (`New<Type>`) that enforces all invariants.

### 8.2 Value Objects

- Value objects **MUST** be immutable after construction.
- Value objects **MUST** validate at construction; an invalid value object **MUST NOT** be representable.
- Domain concepts represented as primitives in transport (Money, Email, ID) **SHOULD** be wrapped in value objects within the domain.

### 8.3 Aggregates

- Aggregates **SHOULD NOT** be introduced unless a cross-entity invariant exists that can be violated under concurrent updates.
- When used, the aggregate **MUST** have one root entity, and external code **MUST** access inner entities only via the root.
- Aggregates **MUST** reference other aggregates by ID, not by pointer.

### 8.4 Repositories

- Repository interfaces **MUST** be declared in the domain package alongside the type they persist.
- Repositories **MUST** accept and return domain types. ORM rows, DB driver types, and DTOs **MUST NOT** appear in repository signatures.
- One repository **SHOULD** correspond to one aggregate root (or one entity, if no aggregate exists).
- Generic `Save(any)` repositories **MUST NOT** be used.

### 8.5 Domain Events

- Domain events **MUST** be named in the past tense (`OrderPlaced`, `PaymentCompleted`).
- Events **MUST** be published only after the corresponding state change is durably persisted. The transactional outbox pattern **SHOULD** be used when delivery guarantees matter.
- Event payloads **MUST** carry the minimum data consumers need; entire entities **MUST NOT** be embedded.
- Breaking changes to event shape **MUST** be expressed as new event types (`OrderPlacedV2`).

### 8.6 Application Services

- Application services **MAY** be implemented as methods on a domain service struct when the use-case has a single coordinator.
- Application services **MUST** be the only place where transactions are opened.
- Application services **MUST NOT** contain business invariants that belong on entities.

### 8.7 Anticorruption Layers

- Every external SDK or vendor API **MUST** be wrapped in an adapter package that exposes only domain types.
- The adapter package **MUST** be the only place that imports the vendor SDK.

---

## 9. Examples (Informative)

### 9.3 Composition Root (Go)

```go
// cmd/orderapi/main.go
func main() {
    cfg := config.Load()
    db  := postgres.MustOpen(cfg.PG)
    bus := nats.MustConnect(cfg.NATS)

    repo   := repository.NewOrderRepo(db)
    outbox := repository.NewOutbox(db)
    tx     := postgres.NewTxRunner(db)
    svc    := app.NewOrderService(tx, repo, outbox)

    go messaging.NewOutboxRelay(outbox, messaging.NewPublisher(bus)).Run()

    h := handler.NewOrder(svc)
    httpserver.Run(cfg.HTTP, h.Routes())
}
```

### 9.4 Domain Event with Outbox

The entity enforces its invariants; the transaction boundary follows §8.6.

```go
// internal/app/order.go — application service
func (s *OrderService) Place(ctx context.Context, req PlaceRequest) (order.Order, error) {
    o, err := order.NewOrder(req.UserID, req.Items)
    if err != nil {
        return order.Order{}, err
    }
    err = s.tx.Run(ctx, func(ctx context.Context) error {
        if err := s.orders.Save(ctx, o); err != nil {
            return err
        }
        return s.outbox.Append(ctx, order.OrderPlaced{
            OrderID: o.ID, UserID: o.UserID, Total: o.Total,
        })
    })
    if err != nil {
        return order.Order{}, fmt.Errorf("place order: %w", err)
    }
    return o, nil
}
```

---

## 10. How to Adopt MDCA (Informative)

A team migrating to MDCA from a layered ("controllers/services/repositories") codebase **SHOULD** follow this sequence:

1. **Identify bounded contexts.** Listen for language shifts in stakeholder conversations. Group existing types by which stakeholder cares about them. Each group is a candidate context.
2. **Carve packages by context.** Move types into `domain/<context>/` packages. Resist the urge to share types across new packages.
3. **Extract ports.** For each external dependency the domain has (DB, broker, vendor SDK), declare an interface in the domain package.
4. **Move infrastructure outward.** Implement each port in a peer package (`internal/repository/`, `internal/messaging/`). Replace direct imports of `database/sql` etc. from the domain.
5. **Establish a composition root.** Centralize all `New*` calls in `cmd/<svc>/main.go` plus `internal/app/`.
6. **Separate DTOs from models.** Wherever a model is serialized for transport, introduce a DTO and a converter.
7. **Add a dependency-direction lint** (`go-arch-lint`, `depguard`) to prevent regression.
8. **Wrap external SDKs in ACLs.** No third-party type may appear in a domain signature.
9. **Audit ceremony.** Remove abstractions with one implementation. Remove specifications and factories that have no real load. Collapse use-case interactors onto service methods where appropriate.

---

## 11. Compliance Checklist

A reviewer **MUST** verify each item before declaring a change MDCA-conformant:

- [ ] All new types and methods use ubiquitous language. (P2)
- [ ] Every package the domain imports passes the purity test of §7.1, transitively re-checked against that package's current import list. (P3, P8)
- [ ] No bounded context imports types from another bounded context. (P4)
- [ ] No interface introduced has only one implementation without a documented second use case on the horizon. (P5)
- [ ] All adapter wiring lives in the composition root. (P6)
- [ ] Cross-context communication uses events unless an exception is documented. (P7)
- [ ] Every external SDK is wrapped in an ACL package. (P8)
- [ ] DTOs carry only transport tags; models carry only persistence tags. (P9)
- [ ] Entity invariants are enforced at construction and via domain methods, not setters. (§8.1)
- [ ] Repository signatures use domain types, not driver or DTO types. (§8.4)
- [ ] Domain events are past-tense, persisted-then-published, with minimal payload. (§8.5)

---

## 12. Anti-Patterns (Informative)

The following patterns are non-conformant. Each row identifies the violated principle and the corrective action.

| Anti-Pattern | Violates | Fix |
|--------------|----------|-----|
| `domain/common/` package | P4 | Distribute types into the contexts that own them. |
| Domain importing `database/sql` | P3 | Define a port; move SQL to a repository. |
| Vendor SDK type in domain method signature | P8 | Wrap in ACL adapter; expose domain type. |
| `Manager`, `Processor`, `Util` type names in domain | P2 | Rename in business vocabulary. |
| Anemic entity (struct + getters/setters, no behavior) | P1 | Move behavior onto the entity; replace setters with operations. |
| Handler reading a domain model directly, with no DTO | P9 | Convert through a DTO at the transport boundary. |
| Repository returning DTOs | §8.4 | Repositories accept and return domain models. |
| Repository per database table | §8.4 | Repository per aggregate root. |
| Generic `Save(any)` repository | §8.4 | Type-specific repositories. |
| Event named `PlaceOrder` (imperative) | §8.5 | Rename to `OrderPlaced` (past tense). |
| Event published before persistence commits | §8.5 | Persist-then-publish; outbox pattern if guaranteed delivery is required. |
| `internal/services/` holding many files of mixed concerns | P1 | Split by domain feature, each as its own package. |
| Cross-service import of `services/<a>/internal/...` from `services/<b>` | P4 | Communicate via events; expose a stable contract when data is needed. |
| Composition root that constructs dozens of services in one file | P6 | Group by bounded context; introduce a sub-container per area. |
| Service locator / global container | P6 | Constructor injection from composition root. |
| Speculative interface with one implementation | P5 | Delete; reintroduce when a second implementation exists. |
| Model carrying both `db:""` and `json:""` tags | P9 | Split into model and DTO. |

---

## 13. Exceptions

A codebase **MAY** deviate from a **SHOULD** clause when a documented quality goal demands it. Deviations **MUST** be recorded in the initiative's Architecture Inception Canvas under *Architectural hypotheses* and reviewed at the next architecture review.

A codebase **MUST NOT** deviate from a **MUST** clause without explicit standard amendment.

---

## Appendix A — Glossary of Layer Imports (Go)

| From → To | Allowed? |
|-----------|----------|
| `cmd/<svc>` → `internal/app` | Yes |
| `cmd/<svc>` → `internal/repository` | Yes (composition only) |
| `internal/app` → `internal/domain/...` | Yes |
| `internal/handler` → `internal/domain/...` | Yes (read-only or via app) |
| `internal/repository` → `internal/domain/...` | Yes (to implement ports) |
| `internal/domain/...` → `internal/repository` | **No** |
| `internal/domain/...` → `internal/handler` | **No** |
| `internal/domain/<a>` → `internal/domain/<b>` | **No** (cross-context) |
| `internal/domain/...` → `database/sql`, `net/http`, vendor SDK | **No** |
| `internal/domain/...` → `pkg/...` | Allowed when the target package is pure (§7.1); forbidden otherwise |
| `internal/repository` → `internal/handler` | **No** |

---

## Appendix B — Server, Client, and Library Realizations (Informative)

MDCA is the same architecture in all three targets. Only the realization changes.

| Aspect | Go Server | Go Client | Go Library |
|--------|-----------|-----------|------------|
| **Composition root** | `cmd/<svc>/main.go` + `internal/app/app.go` | `cmd/<client>/main.go` + bootstrap | None — consumer composes |
| **Domain package** | `internal/domain/<feature>/` | `internal/<subsystem>/` (often ECS-shaped) | `pkg/<concept>/` |
| **Adapters** | `internal/repository/`, `internal/handler/`, NATS publishers | Plugins, drivers (graphics, audio, network) | None — library *is* an adapter to its consumer |
| **Communication** | NATS events internally; REST externally | Internal events; networking via the server | Pure function / type API |
| **Persistence** | Postgres, Redis | Local files, IndexedDB / SQLite | None |
| **Quality goals** | Reliability + maintainability dominate | Performance dominates (frame budget, input latency) | Stability of public API + minimal deps |
| **Reference doc** | `Go Server.md` | `Go Client.md` | `Go Library.md` |

---

## Appendix C — MDCA, SOLID, and Clean Code Alignment (Informative)

| Principle | MDCA realization |
|-----------|-------------------|
| **SRP** | One bounded context per package; one reason to change per type. |
| **OCP** | New behavior added by registering new adapters/subscribers, not editing existing services. |
| **LSP** | All adapters honor their port's contract; conformance test suites recommended. |
| **ISP** | Ports are small, role-based, defined by the consumer (the domain). |
| **DIP** | Domain defines interfaces; adapters depend on the domain. Imports flow inward. |
| **KISS / DRY** | One authoritative model per concept *within* a bounded context; mechanical layer-to-layer duplication is acceptable. |
| **Clean Code** | Stepdown ordering inside files; small functions; no dead code; intention-revealing names match ubiquitous language. |
| **Twelve-Factor** | Config at the composition root; logs to stdout; processes are stateless; backing services are attached resources via adapters. |
| **Trunk-Based Development** | Bounded contexts are small enough to ship in 1–2 day branches; new behaviors land behind feature flags inside their context. |
