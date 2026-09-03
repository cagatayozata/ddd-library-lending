# DDD Library Lending

A small Go sample of a **library catalog + book lending** domain, designed as a **learning project** to practice Domain-Driven Design (DDD) together with Test-Driven Development (TDD).

It is intentionally simple: enough structure to see the building blocks in real code, without HTTP, auth, or a real database.

## What this project is for

Learning and applying these topics:

- DDD & TDD
- Anemic Domain Model vs Rich Domain
- Ubiquitous Language
- Bounded Context
- Tactical DDD concepts
- Entity
- Aggregate & Aggregate Root
- Value Object

## Domain

Members **borrow**, **renew**, and **return** books. The catalog owns shelf copies; lending owns the loan lifecycle. Contexts reference each other by **ISBN**, not by sharing a `Book` object.

## Architecture

```text
(Controller or thin app proxy)
        →  Aggregate Root (Book / Loan)
        →  Repository port
                ↑
         in-memory adapter
```

| Bounded context | Aggregate root | Notes |
|-----------------|----------------|--------|
| `catalog` | `Book` | Owns available copies; raises its own events |
| `lending` | `Loan` | References `ISBN` only; raises its own events |
| `shared` | `DomainEvent` interface only | No shared `AggregateRoot` base type |

## Layout

```text
internal/
  shared/domain/           # DomainEvent interface only
  catalog/
    domain/                # Book (AR), ISBN, Copies, events, repo port
    application/           # thin RegisterBook proxy
    infrastructure/memory/
  lending/
    domain/                # Loan (AR), LoanID/MemberID/ISBN/DueDate VOs, events
    application/           # thin Borrow / Renew / Return proxies
    infrastructure/memory/
cmd/demo/
```

## Concept → code

| Concept | Where |
|---------|--------|
| Rich domain | `Loan.Renew`, `Loan.Return`, `Book.ReserveCopy` |
| Anemic vs rich | rules on the root — not in application if-chains |
| Value Object | `ISBN`, `Copies`, `DueDate`, `LoanID`, `MemberID` |
| Entity / Aggregate Root | `Book`, `Loan` themselves (no shared AR base) |
| Aggregate | Book owns copies; Loan owns due date / status |
| Domain Event | collected on the root (`PullEvents`) |
| Bounded Context | `catalog` vs `lending` |
| Ubiquitous Language | `Borrow`, `Renew`, `Return`, `RegisterBook` |
| Application proxy | load → root method → save (optional; not CQRS) |
| Repository port/adapter | domain ports → `infrastructure/memory` |
| TDD | domain + application `*_test.go` files |

## Run

```bash
cd ddd-library-lending
go test ./... -count=1
go run ./cmd/demo
```
