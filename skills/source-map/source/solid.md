# SOLID Principles in Go

**Source**: Robert C. Martin, "Agile Software Development: Principles, Patterns, and Practices" (2002); refined in "Clean Architecture" (2017).
**Adaptation**: Go-specific application — interfaces are structural, packages replace classes as the primary unit of cohesion. Aligned with Engineering Principles (POL-ENG-001), Clean Code (Go), Go Server / Go Library / Go Client architecture instructions, and MDCA / DDD.

Every rule in the five `### Rules` tables carries a stable identifier, numbered continuously across the whole
document (`SOLID-1` … `SOLID-21`). Cite it in a gap block (`SOLID-12`) rather than the principle name; identifiers
are never reused or renumbered, and new rules are appended.

---

## 1. Single Responsibility Principle (SRP)

> A module should have one, and only one, reason to change.

A "reason to change" is a **stakeholder** or **axis of variation** — pricing rules, persistence format, transport protocol, authentication scheme. If two responsibilities change for different reasons, they belong in different modules.

### Rules

| ID | Rule | Go Adaptation |
|----|------|---------------|
| `SOLID-1` | **One reason to change per type** | A `User` struct that knows its own DB schema AND its JSON shape AND its bcrypt hashing has three reasons to change. Split. |
| `SOLID-2` | **One reason to change per package** | `pkg/user` should not also handle email delivery. Move email to `pkg/notify`. |
| `SOLID-3` | **Cohesion over convenience** | Don't pile helpers into `util`/`common`/`shared`. These packages are SRP violations by definition. |
| `SOLID-4` | **Separate policy from mechanism** | Business rules (policy) live in the domain layer. I/O, serialization, frameworks (mechanism) live at the edges. See Go Server.md layering. |

### Smell List

- A type name ends in `Manager`, `Processor`, `Handler` (outside HTTP), `Service` with > 7 methods, or `Util`.
- A test for a type requires mocking 4+ collaborators — the type is doing too much.
- A change to one feature touches the same file as an unrelated feature ("shotgun surgery" in reverse).
- A package imports both `database/sql` and `net/http` and `crypto/...` — it's likely mixing layers.

### Example

```go
// BAD: three responsibilities in one type
type User struct {
    ID       int
    Email    string
    Password string
}

func (u *User) Save(db *sql.DB) error              { /* persistence */ }
func (u *User) HashPassword(plain string) error    { /* crypto */ }
func (u *User) SendWelcomeEmail(s *smtp.Client) error { /* notification */ }

// GOOD: each responsibility in its own type / package
// pkg/user
type User struct {
    ID       int
    Email    string
    PassHash []byte
}

// pkg/user — repository (persistence policy)
type Repository interface {
    Save(ctx context.Context, u *User) error
}

// pkg/auth — credential policy
type Hasher interface {
    Hash(plain string) ([]byte, error)
    Verify(hash []byte, plain string) error
}

// pkg/notify — notification policy
type Notifier interface {
    SendWelcome(ctx context.Context, to string) error
}
```

---

## 2. Open/Closed Principle (OCP)

> Software entities should be **open for extension, closed for modification**.

You should be able to add new behavior by **adding new code**, not by editing existing, working, tested code. The mechanism in Go is **interface-based polymorphism** plus **composition**.

### Rules

| ID | Rule | Go Adaptation |
|----|------|---------------|
| `SOLID-5` | **Depend on abstractions for varying behavior** | If a new payment provider is likely, define `PaymentProvider` interface. Add a new struct that implements it; don't edit the dispatcher. |
| `SOLID-6` | **Use the strategy pattern for swappable algorithms** | Pass a function or a small interface, not a `switch` over an enum that grows every release. |
| `SOLID-7` | **Prefer composition over modification** | Wrap an existing implementation with a decorator (`LoggingRepository`, `CachingRepository`) instead of editing the original. |
| `SOLID-8` | **Closed against what?** | OCP is always relative. Identify the **likely axis of change** and design around it. Don't speculatively abstract everything (YAGNI wins ties). |

### Anti-Patterns

- A `switch type` or `switch enum` that grows every time a new feature ships. Replace with a registry of implementations.
- `if customerType == "premium" { ... } else if customerType == "trial" { ... }` repeated in 10 places.
- A package that requires a new field on a shared struct every time a new use case is added.

### Example

```go
// Money is the value object from ddd.md — a struct (Amount decimal.Decimal, Currency string),
// not a number — so a rate is built with a constructor, never by multiplying a weight.
// Order.Weight is a decimal.Decimal (kg).
func usd(amount decimal.Decimal) Money { return Money{Amount: amount, Currency: "USD"} }

// BAD: every new shipper requires editing Calculate
func Calculate(order Order, shipper string) Money {
    switch shipper {
    case "ups":   return usd(order.Weight.Mul(decimal.NewFromInt(2)))
    case "fedex": return usd(order.Weight.Mul(decimal.NewFromInt(2)).Add(decimal.NewFromInt(5)))
    case "dhl":   return usd(order.Weight.Mul(decimal.NewFromInt(3))) // added last week
    // case "...": added next week, and the week after...
    }
}

// GOOD: extension point is an interface
type RateCalculator interface {
    Rate(order Order) Money
}

type UPS struct{}
func (UPS) Rate(o Order) Money { return usd(o.Weight.Mul(decimal.NewFromInt(2))) }

type FedEx struct{}
func (FedEx) Rate(o Order) Money {
    return usd(o.Weight.Mul(decimal.NewFromInt(2)).Add(decimal.NewFromInt(5)))
}

// New shipper = new file, new struct. Calculate never changes.
func Calculate(order Order, c RateCalculator) Money { return c.Rate(order) }
```

### Go-Specific Notes

- **Don't pre-abstract.** Wait for the second concrete need before extracting an interface. The first interface designed without a real second implementation is almost always wrong.
- **Decorators are idiomatic in Go** for cross-cutting concerns (logging, tracing, caching, retries) — they preserve OCP by leaving the base implementation untouched. See Go Server.md middleware patterns.

---

## 3. Liskov Substitution Principle (LSP)

> Subtypes must be substitutable for their base types **without altering the correctness of the program**.

Go has no subclassing, so LSP applies to **interface implementations**: every type that satisfies an interface must obey the interface's behavioral contract — not just its method signatures.

### Rules

| ID | Rule | Go Adaptation |
|----|------|---------------|
| `SOLID-9` | **Honor the contract, not just the signature** | A `Repository.Save` that silently drops writes when the cache is full violates LSP even though it compiles. |
| `SOLID-10` | **Don't strengthen preconditions** | If `io.Writer.Write` accepts any byte slice, your implementation cannot reject empty slices with an error. |
| `SOLID-11` | **Don't weaken postconditions** | If `Reader.Read` guarantees `n` bytes were read on success, your implementation cannot return success with `n=0` unless `len(p)==0`. |
| `SOLID-12` | **Don't surprise the caller** | A `Logger.Info` implementation that blocks on network I/O for 30 seconds violates the implicit "logging is fast" contract. |
| `SOLID-13` | **Errors are part of the contract** | Document which sentinel errors a method can return (`io.EOF`, `sql.ErrNoRows`). Implementations must respect this. |

### Smell List

- An interface implementation panics, returns "not supported", or no-ops where callers expect work to happen.
- A wrapper type that implements an interface but ignores some methods.
- Tests for an implementation differ in semantics from tests for the interface (they should be the same conformance suite).
- The interface's godoc says "may return X"; one implementation always returns X, another never does — callers can't write generic code.

### Example

```go
// Contract: Storage.Get returns (value, true) if key exists, ("", false) otherwise.
// MUST NOT return an error for a missing key.
type Storage interface {
    Get(key string) (string, bool)
}

// BAD: violates LSP — caller code that uses InMemory will break with Redis
type Redis struct{ /* ... */ }
func (r *Redis) Get(key string) (string, bool) {
    v, err := r.client.Get(key)
    if err != nil {
        panic(err) // NO — contract says return ("", false), not panic
    }
    return v, true
}

// GOOD: respects the contract
func (r *Redis) Get(key string) (string, bool) {
    v, err := r.client.Get(key)
    if errors.Is(err, redis.ErrNil) || err != nil {
        return "", false
    }
    return v, true
}
```

### Go-Specific Notes

- Define a **conformance test suite** for any interface with multiple implementations. Run it against every implementation. This is the most reliable way to enforce LSP.
- The standard library does this — see `testing/iotest` and `net/http/httptest` patterns.
- If different implementations need different contracts, they need different interfaces. Don't fake-satisfy an interface.

---

## 4. Interface Segregation Principle (ISP)

> Clients should not be forced to depend on methods they do not use.

Many small, role-based interfaces are better than one large interface with everything. In Go, this is the **single most important** SOLID principle, and it's the one Go is designed around.

### Rules

| ID | Rule | Go Adaptation |
|----|------|---------------|
| `SOLID-14` | **Define interfaces at the point of use** | "Accept interfaces, return structs." A consumer declares the minimum interface it needs, not what the producer offers. |
| `SOLID-15` | **Prefer 1–3 method interfaces** | `io.Reader`, `io.Writer`, `io.Closer`, `fmt.Stringer`. Compose with embedding (`io.ReadWriter`) when needed. |
| `SOLID-16` | **No `XxxService` god interfaces** | A 12-method `UserService` interface forces every test mock to stub all 12. Split by use case. |
| `SOLID-17` | **Don't expose more than callers need** | If your handler needs only `GetUser`, it should depend on `interface{ GetUser(ctx, id) (User, error) }`, not the full repository. |

### Smell List

- A test file with 80 lines of mock methods that all `panic("not implemented")` — the interface is too big.
- An interface where consumers say "I only call one of these five methods."
- An interface defined in the same package as its only implementation, used by external callers — invert it: define it where it's consumed.
- Method names on an interface span unrelated concerns (`SaveUser`, `SendEmail`, `LogAudit`) — that's three interfaces.

### Example

```go
// BAD: one fat interface — every consumer takes a hard dependency on all 6 methods
type UserStore interface {
    Create(ctx context.Context, u User) error
    Get(ctx context.Context, id int) (User, error)
    List(ctx context.Context) ([]User, error)
    Update(ctx context.Context, u User) error
    Delete(ctx context.Context, id int) error
    Search(ctx context.Context, q string) ([]User, error)
}

// A handler that only reads one user is forced to mock all six.

// GOOD: small interfaces declared by consumers
// internal/api/profile_handler.go
type userReader interface {
    Get(ctx context.Context, id int) (User, error)
}

type ProfileHandler struct{ users userReader }

// internal/api/admin_handler.go
type userAdmin interface {
    List(ctx context.Context) ([]User, error)
    Delete(ctx context.Context, id int) error
}

// The same concrete *postgres.UserStore satisfies both. Tests stub only what they use.
```

### Go-Specific Notes

- **Define interfaces in the consumer package**, not the producer package. Producers return concrete structs. Consumers define the slice of behavior they need. This is the opposite of Java/C# convention and is idiomatic Go.
- **Embed interfaces** to compose larger contracts: `type ReadWriter interface { Reader; Writer }`. Keep the embedded interfaces small.
- **Empty interface** (`any`) is the ultimate ISP violation — it asserts the consumer needs nothing, then forces type assertions. Avoid except in genuinely generic containers (and prefer generics now: `[T any]`).

---

## 5. Dependency Inversion Principle (DIP)

> High-level modules should not depend on low-level modules. Both should depend on abstractions.
> Abstractions should not depend on details. Details should depend on abstractions.

The direction of the **source-code dependency** should be the opposite of the **flow of control**. Domain logic must not import infrastructure. Infrastructure imports domain.

### Rules

| ID | Rule | Go Adaptation |
|----|------|---------------|
| `SOLID-18` | **Domain defines interfaces; infra implements them** | The `domain/order` package defines `Repository`. The `infra/postgres` package imports `domain/order` and implements it. Never the reverse. |
| `SOLID-19` | **Inject dependencies at construction** | `func NewService(repo Repository, log *slog.Logger) *Service`. No package-level singletons. No global `db`. |
| `SOLID-20` | **Wire at the edges** | `main.go` (or a `cmd/...` entry point) constructs concrete implementations and injects them downward. The center of the app sees only interfaces. |
| `SOLID-21` | **No imports from inner layers to outer layers** | `domain → application → infra → cmd`. Imports point inward. Enforce with a linter (e.g., `go-arch-lint`, `depguard`). |

### Smell List

- A domain package imports `database/sql`, `net/http`, `github.com/aws/...`, or any framework. The dependency points the wrong way.
- A function takes `*sql.DB` directly when it could take a small `Querier` interface.
- Constructors call `os.Getenv` or read config — config belongs at composition root (`main`), not in business logic. See Twelve-Factor App.md.
- A test requires spinning up Postgres because the code calls `sql.Open` deep inside a service method.

### Example

```go
// BAD: domain depends on infrastructure
// internal/domain/order/service.go
package order

import "database/sql" // <-- domain importing infra. Wrong direction.

type Service struct{ db *sql.DB }

func (s *Service) Place(ctx context.Context, o Order) error {
    _, err := s.db.ExecContext(ctx, "INSERT INTO orders ...")
    return err
}

// GOOD: domain defines abstraction; infra depends on domain
// internal/domain/order/service.go
package order

type Repository interface {
    Save(ctx context.Context, o Order) error
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Place(ctx context.Context, o Order) error {
    return s.repo.Save(ctx, o)
}

// internal/infra/postgres/order_repo.go
package postgres

import (
    "database/sql"
    "myapp/internal/domain/order" // <-- infra imports domain. Correct direction.
)

type OrderRepo struct{ db *sql.DB }

func (r *OrderRepo) Save(ctx context.Context, o order.Order) error { /* SQL */ }

// cmd/server/main.go — composition root
func main() {
    db := openDB()
    repo := &postgres.OrderRepo{DB: db}
    svc := order.NewService(repo) // wires concrete into abstract
    // ...
}
```

### Go-Specific Notes

- Go's structural interfaces make DIP nearly free — the implementer doesn't even need to import the interface package. This is unique to Go and Rust-style structural typing.
- **Composition root** is `main` (or a `wire`/`fx` setup function). All `new*` calls and concrete bindings happen there. Nowhere else.
- **No service locators**. No `container.Get("UserService")`. Pass dependencies explicitly through constructors.
- For test seams: production code receives interfaces. Tests pass fakes/stubs/mocks that satisfy them. No build tags, no `if testing { ... }`.

---

## Quick Self-Check

Before merging a change, ask:

1. **SRP** — If I describe what this type/package does, do I need the word "and"? If yes, split.
2. **OCP** — To support the next likely variation, will I edit existing code or add a new file? If edit, reconsider the abstraction.
3. **LSP** — Can I swap any implementation of this interface for any other and have all callers still work correctly? If no, the contract is wrong.
4. **ISP** — Does any consumer of this interface use less than half its methods? If yes, segregate.
5. **DIP** — Does the import graph point from concrete edges toward the abstract center? If imports flow outward, invert them.

---

## Anti-Patterns Specific to Go

| Anti-Pattern | Principle Violated | Fix |
|--------------|--------------------|-----|
| `util` / `common` / `helpers` package | SRP (`SOLID-3`) | Move each helper to the package that owns its concept. |
| Interface defined in producer package, only used externally | ISP (`SOLID-14`), DIP (`SOLID-18`) | Move interface to consumer; producer returns concrete struct. |
| `switch x.(type)` over a closed set of types in business code | OCP (`SOLID-6`) | Replace with method dispatch on an interface. |
| Global `var DB *sql.DB` | DIP (`SOLID-19`) | Inject through constructors. |
| 15-method `*Service` struct | SRP (`SOLID-1`) | Split by use case (one struct per command/query group). |
| Mocks that `panic("not implemented")` for half their methods | ISP (`SOLID-16`) | The interface is too big — split it. |
| `init()` that opens DB connections, reads env, starts goroutines | DIP (`SOLID-20`), SRP (`SOLID-4`) | Move to `main`. `init` is for registering codecs and pure setup only. |
