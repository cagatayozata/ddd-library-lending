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
Command → Command Handler → Domain (aggregate) → Repository port
                                    ↑
                         infrastructure adapter (in-memory)
```

| Bounded context | Aggregate root | Notes |
|-----------------|----------------|--------|
| `catalog` | `Book` | Owns available copies |
| `lending` | `Loan` | References `ISBN` only |
| `shared` | `AggregateRoot`, `DomainEvent`, `Invariants` | Shared tactical pieces |

## Layout

```text
internal/
  shared/domain/           # AggregateRoot, DomainEvent, Invariants
  catalog/
    domain/                # Book, ISBN, Copies, domain events, repo port
    application/command/   # RegisterBook
    infrastructure/memory/
  lending/
    domain/                # Loan, DueDate, status, domain events, repo port
    application/command/   # BorrowBook, RenewLoan, ReturnBook
    infrastructure/memory/
cmd/demo/                  # end-to-end walkthrough
```

## Concept → code

| Concept | Where |
|---------|--------|
| Rich domain | `Loan.Renew`, `Loan.Return`, `Book.ReserveCopy` |
| Anemic vs rich | rich behaviour lives on aggregates/VOs (not anemic setters + service ifs) |
| Value Object | `ISBN`, `Copies`, `DueDate`, typed IDs |
| Entity / Aggregate Root | `Book`, `Loan` (+ `shared.AggregateRoot`) |
| Aggregate | Book owns copy count; Loan owns due date / status |
| Domain Event | `BookRegistered`, `LoanBorrowed`, `LoanReturned`, … |
| Bounded Context | `catalog` vs `lending` |
| Ubiquitous Language | `Borrow`, `Renew`, `Return`, `RegisterBook` |
| CQRS L1 (commands) | `application/command` handlers |
| Repository port/adapter | domain ports → `infrastructure/memory` |
| TDD | domain + handler `*_test.go` files |

## Run

```bash
cd ddd-library-lending
go test ./... -count=1
go run ./cmd/demo
```
