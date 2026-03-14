---
name: review-ddd
description: Review code for DDD and Clean Architecture violations. Use after implementing features to ensure architectural integrity.
---

# DDD & Clean Architecture Review

Review recent changes for architectural violations.

## Checklist

### Layer Dependency Rules
- [ ] **Domain layer** has NO imports from application, infrastructure, or presentation
- [ ] **Application layer** imports only from domain
- [ ] **Infrastructure layer** implements domain interfaces
- [ ] **Presentation layer** calls only application layer (commands/queries)

### DDD Rules
- [ ] Entities have proper constructors with validation
- [ ] Value Objects are immutable
- [ ] Aggregate boundaries are respected (no cross-aggregate direct references)
- [ ] Repository interfaces are in domain layer
- [ ] Repository implementations are in infrastructure layer
- [ ] Domain events are used for cross-aggregate side effects

### CQRS Rules
- [ ] Commands (writes) use UoW for transactions
- [ ] Queries (reads) do NOT use UoW or locks
- [ ] Commands return minimal results
- [ ] Queries return DTOs, not entities

### Lock Order Rules
- [ ] Locks follow order: Tenant(1) → User(2) → Book(3) → Memo(4) → Tag(5) → UserPokemon(6)
- [ ] No reverse-order lock acquisition
- [ ] SELECT ... FOR UPDATE only in commands, never in queries

### Anti-Patterns to Flag
- Business logic in handlers (should be in domain)
- GORM tags in domain entities (should be in infrastructure models)
- Direct DB calls from handlers (should go through application layer)
- Shared mutable state without proper synchronization

## Output Format
Report findings as:
- **VIOLATION**: Must fix before merge
- **WARNING**: Should fix, lower priority
- **SUGGESTION**: Nice to have improvement
