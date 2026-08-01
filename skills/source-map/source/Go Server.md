# Go Application Architecture Instruction

Architectural reference for any Go service following Modular Domain-Centric Architecture (MDCA) with Domain-Driven Design, Event-Driven Design, and microservices inside a monorepository.

---

## Modular Domain-Centric Architecture (MDCA)

MDCA organizes a service around independent bounded-context modules and keeps infrastructure out of the domain: the domain declares the ports it needs, adapters implement them, and one composition root wires them together. It optimizes for performance efficiency, reliability, and maintainability, and introduces an abstraction only once a second concrete need for it exists.

The normative definition lives in **`mdca.md`** — principles `P1`–`P10` (cite clauses as `P3.2`), the layering rules and the purity test for domain imports in §7, the tactical rules in §8, and the compliance checklist in §11. This document does not restate them; it describes how a Go **server** realizes MDCA. General Go coding rules live in **`Go Best Practice.md`**.

---

## Folder Structure

The monorepository splits into three main sections: **services** (entry points + service logic), **internal** (private implementation), and **pkg** (shared infrastructure).

### Multi-Service Monorepository

```
root/
├── services/                      # Independent microservices
│   ├── serviceA/                  # One service per business capability
│   │   ├── cmd/
│   │   │   └── main.go            # Entry point: config, DI, server start
│   │   └── internal/              # Private to this service
│   │       ├── app/
│   │       │   └── app.go         # Dependency injection container
│   │       ├── domain/            # Core business logic
│   │       │   ├── featureA/      # One package per bounded context
│   │       │   │   ├── model.go         # DB entities (db:"" tags only)
│   │       │   │   ├── port.go          # Repository/dependency interfaces
│   │       │   │   ├── service.go       # Business logic + DTO↔Model converters
│   │       │   │   ├── service_test.go  # Unit tests with mocks
│   │       │   │   └── mocks/           # Auto-generated mock implementations
│   │       │   └── featureB/
│   │       │       └── ...
│   │       ├── dto/               # Data Transfer Objects (json:"" tags only)
│   │       │   ├── featureA/
│   │       │   │   ├── request.go       # Incoming request DTOs
│   │       │   │   └── response.go      # Outgoing response DTOs
│   │       │   └── featureB/
│   │       │       └── ...
│   │       ├── repository/        # Data access implementations
│   │       │   ├── featureA.go    # No _repo suffix — package provides context
│   │       │   └── featureB.go
│   │       ├── event/             # Messaging adapters implementing domain ports
│   │       │   └── featureA.go    # Wraps pkg/messaging; the domain never imports it
│   │       └── handler/           # Presentation layer (HTTP, gRPC, events)
│   │           ├── featureA.go    # Request handlers (uses DTOs, never models)
│   │           ├── meta.go        # Swagger/OpenAPI annotations
│   │           └── router.go      # Route registration
│   └── serviceB/
│       └── ...
├── pkg/                           # Shared reusable packages
│   ├── auth/                      # Authentication/authorization
│   ├── httpserver/                # HTTP framework wrapper, health checks
│   ├── postgres/                  # Database connection, DSN builder
│   ├── email/                     # SMTP client
│   ├── logger/                    # Structured logging
│   ├── messaging/                 # Event bus / message queue client
│   └── config/                    # Environment variable parsing
├── migrations/                    # Database migration files
├── docs/                          # Per-service documentation
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
└── README.md
```

### Single-Domain vs Multi-Domain Services

**Single-domain services** (e.g., accessapi with 1 domain, notifyapi with 1 domain) use a **flat** structure — no feature subfolders:

```
services/accessapi/internal/
├── domain/           # Flat — model.go, port.go, service.go directly here
│   ├── model.go
│   ├── port.go
│   ├── service.go
│   └── service_test.go
├── dto/              # Flat — request.go, response.go directly here
│   ├── request.go
│   └── response.go
├── repository/
│   └── repository.go # Single file, no suffix needed
└── handler/
    ├── handler.go
    └── router.go
```

**Multi-domain services** (e.g., scheduleapi with 11 domains) use **feature subfolders**:

```
services/scheduleapi/internal/
├── domain/
│   ├── shift/        # One package per bounded context
│   ├── trade/
│   └── timeaway/
├── dto/
│   ├── shift/
│   ├── trade/
│   └── timeaway/
├── repository/
│   ├── shift.go      # No _repo suffix — package name provides context
│   ├── trade.go
│   └── timeaway.go
└── handler/
    ├── shift.go
    ├── trade.go
    ├── timeaway.go
    └── router.go
```

### Domain Internal Structure

Each domain feature directory follows a consistent pattern:

| File | Purpose |
|------|---------|
| `model.go` | DB entities with `db:""` struct tags only. **Never** `json:""` tags — those belong in DTOs. Domain-specific value objects and enums. |
| `port.go` | Interfaces (ports) that define what the domain needs from infrastructure — typically a Repository, plus any outbound port such as a Publisher. Includes `//go:generate` directives for mock generation. |
| `service.go` | Business logic + DTO↔Model converter functions (`NewXResponse`, `NewXFromRequest`). Orchestrates calls to ports and applies domain rules. Each public method represents a use case. |
| `service_test.go` | Unit tests using generated mocks. Tests business logic in isolation from infrastructure. |
| `mocks/` | Auto-generated mock implementations (mockgen, moq, etc.). |

### DTO Structure

DTOs live in a **separate** `dto/` directory (not inside domain):

| File | Purpose |
|------|---------|
| `request.go` | Incoming request types with `json:""` struct tags only. **Never** `db:""` tags. |
| `response.go` | Outgoing response types with `json:""` struct tags only. **Never** `db:""` tags. |

**Flow:** Handler → DTO → Service (converts to Model) → Repository. Handler never sees models. Repository never sees DTOs.

---

## Section Responsibilities

### cmd/main.go -- Entry Point

The entry point is **minimal**. It wires everything together and starts the server. It must **NOT** contain business logic.

Responsibilities:
- Parse configuration from environment variables
- Initialize infrastructure (database, message broker, logger)
- Create the DI container (`app.NewServices`)
- Register HTTP/gRPC routes or event subscriptions
- Start the server and handle graceful shutdown

### internal/app/app.go -- Dependency Injection Container

A plain struct that holds all domain services. A constructor function creates repositories, injects them into services, and returns the container. No framework needed -- Go structs and constructors are sufficient.

```go
type Services struct {
    Order   *order.Service
    Payment *payment.Service
}

func NewServices(db *sqlx.DB, bus messaging.Publisher) *Services {
    orderRepo := repository.NewOrderRepo(db)
    paymentRepo := repository.NewPaymentRepo(db)
    return &Services{
        // The broker is wrapped in an adapter satisfying order.Publisher,
        // so the domain package never sees pkg/messaging.
        Order:   order.NewService(orderRepo, event.NewOrderPublisher(bus)),
        Payment: payment.NewService(paymentRepo),
    }
}
```

The composition root is the one place allowed to name concrete infrastructure. Everything it hands to a domain service is a domain-declared port.

### internal/domain/{feature}/ -- Domain Layer

The heart of the application. Each feature (bounded context) gets its own package containing model, port, service, and tests.

**Domain packages must have ZERO infrastructure imports** -- no database drivers, no HTTP libraries, no message queue clients, no vendor SDKs. They depend on the Go standard library, their own port interfaces, and only such other packages as pass the purity test in `mdca.md` §7.1: a package is importable by the domain only when its own imports are standard-library-only (transitively) and it performs no I/O and holds no mutable global state. The folder a package lives in is not the test — a `pkg/` package may qualify, and a package anywhere may fail.

### internal/repository/ -- Data Access Layer

Implementations of `port.Repository` interfaces. Each file corresponds to one domain feature. Repositories own SQL queries, data mapping, and database-specific concerns. They accept `context.Context` for cancellation and return domain models (not ORM objects).

### internal/event/ -- Messaging Adapters

Implementations of the domain's outbound messaging ports. Each adapter wraps the broker client from `pkg/messaging` (infrastructure: it opens network connections, so the domain must not import it) and translates domain events into broker topics and payloads. This is the anticorruption boundary for the message bus.

### internal/handler/ -- Presentation Layer

Request handlers that translate HTTP/gRPC/event inputs into domain service calls and format responses.

We do **NOT** split handler logic from validation or response formatting because of the KISS principle -- changing the communication layer usually changes how events are processed. We should not add flexibility to systems that do not change frequently.

| File | Purpose |
|------|---------|
| `router.go` | Route registration, middleware wiring |
| `meta.go` | Swagger/OpenAPI general info annotations |
| `{feature}.go` | Handler methods grouped by domain feature |

### pkg/ -- Shared Packages

Reusable infrastructure packages shared across services. Each package owns its configuration prefix and exposes a `GetEnvConfig()` function. Prefixes are by **FUNCTIONALITY** (e.g., `POSTGRES_`, `APP_`, `AUTH_`), not by service name.

**Rules for pkg/:**
- Must be genuinely reusable across 2+ services
- Must not import from any service's `internal/` package
- Must not contain business logic -- only infrastructure concerns
- Each package should be independently testable

---

## Domain-Driven Design (DDD) in Go

DDD organizes code around the business domain rather than technical layers. In MDCA, we apply DDD **pragmatically** -- use the patterns that help, skip those that add complexity without value.

### Bounded Contexts

Each domain feature package is a bounded context. It has its own models, its own language (ubiquitous language), and its own rules. **Do not share models across bounded contexts** -- if two domains need similar data, each defines its own struct.

```go
// domain/order/model.go -- Order's view of a product
type OrderItem struct {
    ProductID string
    Quantity  int
    UnitPrice Money // order's own Money value object
}

// domain/catalog/model.go -- Catalog's view of a product
type Product struct {
    ID          string
    Name        string
    Description string
    Price       Money // catalog's own Money value object
    Stock       int
}
```

Each context owns its `Money` type too (see below). Reaching for a shared third-party money or decimal package would both couple the contexts and pull a vendor type into the domain.

### Entities and Value Objects

**Entities** have identity (an ID field) and a lifecycle. **Value objects** are defined by their attributes and are immutable. In Go, both are plain structs. Use pointer receivers for entity methods that mutate state; use value receivers for value object methods.

```go
// Entity -- has identity
type Order struct {
    ID        string
    Status    OrderStatus
    Items     []OrderItem
    CreatedAt time.Time
}

func (o *Order) Cancel() error {
    if o.Status != StatusPending {
        return ErrCannotCancel
    }
    o.Status = StatusCancelled
    return nil
}

// Value Object -- defined by attributes, immutable, domain-owned.
// Amounts are held in the currency's minor unit (cents) as an int64: exact
// arithmetic without a third-party decimal type, so the package keeps
// standard-library-only imports.
type Money struct {
    Minor    int64  // e.g. 1099 == 10.99 USD
    Currency string // ISO 4217 code
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, ErrCurrencyMismatch
    }
    return Money{Minor: m.Minor + other.Minor, Currency: m.Currency}, nil
}
```

### Ports and Adapters (Hexagonal Architecture)

**Ports** are interfaces defined in the domain layer. **Adapters** are implementations in the repository, event, or handler layer. The domain depends on abstractions it declares itself, not on concrete infrastructure.

```go
// domain/order/port.go
//go:generate go tool mockgen -source=port.go -destination=mocks/mock_repository.go
type Repository interface {
    GetByID(ctx context.Context, id string) (Order, error)
    Save(ctx context.Context, order Order) error
    ListByUser(ctx context.Context, userID string, filter Filter) ([]Order, error)
}

// Publisher is the domain's outbound port for announcing what happened.
// It speaks domain events declared in this package, not broker topics, so
// the domain never imports pkg/messaging -- which fails the purity test.
type Publisher interface {
    PublishOrderPlaced(ctx context.Context, evt OrderPlaced) error
}
```

```go
// repository/order.go — no _repo suffix; package name provides context
type OrderRepo struct {
    db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
    return &OrderRepo{db: db}
}

func (r *OrderRepo) GetByID(ctx context.Context, id string) (order.Order, error) {
    const query = `SELECT id, status, created_at FROM orders WHERE id = $1`
    var o order.Order
    if err := r.db.GetContext(ctx, &o, query, id); err != nil {
        return order.Order{}, fmt.Errorf("get order by id: %w", err)
    }
    return o, nil
}
```

The same shape applies to non-database ports. `order.Publisher` is implemented by an adapter over `pkg/messaging`; the dependency points inward, from adapter to domain:

```go
// event/order.go — adapter implementing order.Publisher
var _ order.Publisher = (*OrderPublisher)(nil)

type OrderPublisher struct {
    bus messaging.Publisher
}

func NewOrderPublisher(bus messaging.Publisher) *OrderPublisher {
    return &OrderPublisher{bus: bus}
}

func (p *OrderPublisher) PublishOrderPlaced(ctx context.Context, evt order.OrderPlaced) error {
    if err := p.bus.Publish(ctx, "order.placed", newOrderPlacedPayload(evt)); err != nil {
        return fmt.Errorf("publish order.placed: %w", err)
    }
    return nil
}
```

### Domain Services

Domain services contain business logic that does not naturally belong to a single entity. They orchestrate repositories and apply domain rules. Each public method is a use case.

```go
// domain/order/service.go -- depends only on ports declared in this package
type Service struct {
    repo      Repository
    publisher Publisher
}

func NewService(repo Repository, pub Publisher) *Service {
    return &Service{repo: repo, publisher: pub}
}

func (s *Service) PlaceOrder(ctx context.Context, req PlaceOrderRequest) (Order, error) {
    o := NewOrder(req.UserID, req.Items)
    if err := o.Validate(); err != nil {
        return Order{}, fmt.Errorf("validate order: %w", err)
    }
    if err := s.repo.Save(ctx, o); err != nil {
        return Order{}, fmt.Errorf("save order: %w", err)
    }
    if err := s.publisher.PublishOrderPlaced(ctx, OrderPlaced{OrderID: o.ID}); err != nil {
        return Order{}, fmt.Errorf("publish order placed: %w", err)
    }
    return o, nil
}
```

The publish error is handled, never discarded -- a swallowed publish leaves the rest of the system unaware of a state change that already happened. When the order must survive a broker outage, persist the event with the order instead and let a relay publish it (see the Outbox Pattern below).

---

## Event-Driven Design

Event-driven design decouples services by communicating through events rather than direct calls. In MDCA, events are the primary mechanism for cross-service communication within the monorepository.

### Event Types

There are three kinds of events in an event-driven system:

| Type | Description | Example |
|------|-------------|---------|
| **Domain Events** | Something happened within a bounded context. Published by domain services after a state change succeeds. | `OrderPlaced`, `PaymentCompleted` |
| **Integration Events** | Domain events translated for cross-service consumption. Carry only the data other services need, not the full domain model. | `OrderPlacedIntegration` |
| **Command Events** | Requests for another service to perform an action. Used for async task delegation. | `SendNotification`, `ProcessPayment` |

### Event Structure

Events are plain Go structs. They carry enough context for consumers to act without calling back to the producer. Include correlation IDs for tracing.

```go
type Event struct {
    ID            string    // Unique event identifier
    Type          string    // Event type (e.g., "order.placed")
    Source        string    // Producing service name
    CorrelationID string    // Trace correlation
    OccurredAt    time.Time // When the event happened
    Payload       any       // Event-specific data
}

// Wire payload -- built by the adapter from the domain event, so vendor and
// transport concerns stay outside the domain package.
type OrderPlacedPayload struct {
    OrderID    string
    UserID     string
    TotalMinor int64  // total in the currency's minor unit
    Currency   string // ISO 4217 code
    ItemCount  int
}
```

### Publisher/Subscriber Pattern

`pkg/messaging` defines the broker-facing interfaces used by the application and infrastructure layers. Implementations can use NATS, Kafka, RabbitMQ, or even in-process channels -- the rest of the system does not care.

Domain services do **not** depend on these interfaces: `pkg/messaging` ships a broker client, so it fails the purity test and is off-limits to the domain. The domain declares its own port (`order.Publisher` above) and an adapter in `internal/event/` bridges the two.

```go
// pkg/messaging/publisher.go
type Publisher interface {
    Publish(ctx context.Context, topic string, event any) error
}

// pkg/messaging/subscriber.go
type Handler func(ctx context.Context, event []byte) error

type Subscriber interface {
    Subscribe(ctx context.Context, topic string, handler Handler) error
    Close() error
}
```

### Event Handling in Handlers

Event subscribers are registered alongside HTTP routes in the handler layer. Each handler translates the raw event into a domain service call.

```go
func RegisterEventHandlers(ctx context.Context, sub messaging.Subscriber, svc *app.Services) error {
    return sub.Subscribe(ctx, "order.placed", func(ctx context.Context, data []byte) error {
        var evt OrderPlacedPayload
        if err := json.Unmarshal(data, &evt); err != nil {
            return fmt.Errorf("unmarshal order.placed: %w", err)
        }
        return svc.Notification.SendOrderConfirmation(ctx, evt.UserID, evt.OrderID)
    })
}
```

### Eventual Consistency

Event-driven systems are eventually consistent. Each service owns its data and updates it based on events it receives. Do not rely on synchronous reads across service boundaries. If a service needs data from another, it maintains its own **local projection** updated via events.

### Idempotency

Event handlers **MUST** be idempotent -- processing the same event twice must produce the same result. Use the event ID to deduplicate. Store processed event IDs in a table or use database constraints to prevent double processing.

```go
func (s *Service) HandlePaymentCompleted(ctx context.Context, evt PaymentCompletedEvent) error {
    if s.repo.EventAlreadyProcessed(ctx, evt.ID) {
        return nil // idempotent -- already handled
    }
    // ... process event ...
    return s.repo.MarkEventProcessed(ctx, evt.ID)
}
```

### Outbox Pattern

To guarantee that database changes and events are published atomically, use the **transactional outbox pattern**: write the event to an outbox table in the same database transaction as the state change. A separate process polls the outbox and publishes events to the message broker.

```go
func (r *OrderRepo) SaveWithEvent(ctx context.Context, order Order, event Event) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    if _, err := tx.ExecContext(ctx, `INSERT INTO orders ...`, order); err != nil {
        return fmt.Errorf("insert order: %w", err)
    }
    if _, err := tx.ExecContext(ctx, `INSERT INTO outbox ...`, event); err != nil {
        return fmt.Errorf("insert outbox: %w", err)
    }
    return tx.Commit()
}
```

---

## Microservices in a Monorepository

A monorepository hosts multiple independently deployable services that share infrastructure code via `pkg/`. Each service is a self-contained unit with its own `main.go`, `internal/` tree, and deployment configuration.

### Service Boundaries

- Each service owns exactly one business capability
- Services **MUST NOT** import another service's `internal/` packages
- Services communicate via events, HTTP APIs, or gRPC -- never via direct function calls
- Shared code lives in `pkg/` and must be infrastructure-only (no business logic)

### Independent Deployability

- Each service compiles to its own binary
- Services can be deployed, scaled, and versioned independently
- A change in serviceA should not require redeploying serviceB
- Database migrations are per-service or shared via a controlled migration directory

### Data Ownership

- Each service owns its database tables
- No service reads or writes another service's tables directly
- Cross-service data access happens through APIs or event-driven projections

---

## Go Coding Rules

General Go rules — formatting, naming, errors, control structures, `defer`, interfaces, concurrency, data structures, comments, constants, `init`, embedding, and testing — live in **`Go Best Practice.md`**; this document covers server architecture only.

---

## How to Start a New Service

1. **Domain Analysis** -- identify the bounded context and its core use cases.
2. **Define Module Boundaries** -- one package per feature under `domain/`.
3. **Establish Domain Models** -- define DB entities in `domain/{feature}/model.go` with `db:""` tags only.
4. **Define DTOs** -- request/response types in `dto/{feature}/request.go` and `response.go` with `json:""` tags only.
5. **Define Ports** -- write the Repository and any outbound interfaces (e.g. Publisher) in `domain/{feature}/port.go`.
6. **Implement Domain Service** -- business logic + DTO↔Model converters in `service.go`, tests in `service_test.go`.
7. **Implement Adapters** -- data access in `repository/{feature}.go` (no `_repo` suffix — package name provides context), messaging in `event/{feature}.go`.
8. **Implement Handlers** -- HTTP/event handlers in `handler/{feature}.go` (bind DTOs, call service, return DTOs).
9. **Wire Everything** -- DI in `app/app.go`, routes in `handler/router.go`.
10. **Create Entry Point** -- minimal `main.go`: config -> infra -> DI -> serve.

### Real-World Implementation Tips

- **Start small** -- apply the architecture to a bounded area first
- **Evolve gradually** -- refactor existing systems incrementally
- **Focus on business value** -- prioritize domains most critical to the business
- **Continuous refactoring** -- revisit boundaries as domain understanding evolves
- **Do not add abstraction before you need it**
- **Keep the architecture solution behind the scenes** -- we should not always use an interface or follow SOLID if it will not improve the specific situation
- **Do not use an interface when the cost of abstraction exceeds the cost of rewriting**

### Example of High Interface Price

Consider a dynamic subscriber that entirely relies on a specific message queue implementation (e.g., NATS). Separating the subscriber from NATS to make it generic costs ~8 hours. But if we decide to use Kafka instead of NATS, rewriting the subscriber directly takes ~3 hours. The abstraction costs more than the rewrite -- avoid it.
