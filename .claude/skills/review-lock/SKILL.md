---
name: review-lock
description: Check for lock ordering violations and potential deadlocks. Use during code review for write operations.
---

# Lock Order Review

Analyze recent changes for lock ordering violations and deadlock risks.

## Lock Order Definition

```
Tenant (1) → User (2) → Book (3) → Memo (4) → Tag (5) → UserPokemon (6)
```

## Review Steps

1. Find all `FOR UPDATE` queries and `Save`/`Update`/`Delete` operations
2. For each command handler, trace the order of locked resources
3. Verify the order is strictly ascending (lower number → higher number)
4. Flag any violation

## What to Look For

### Violations
- `BookRepo.FindByIDForUpdate()` followed by `UserRepo.FindByIDForUpdate()` (3 → 2, WRONG)
- Multiple aggregates locked in wrong order
- Nested transactions

### Safe Patterns
- Single aggregate operations (no ordering concern)
- Read-only queries (no locks)
- Ascending order: `Book(3) → Memo(4) → UserPokemon(6)` ✓

## Output

For each command handler:
```
Handler: FinishReadingHandler
Locks: Book(3) → UserPokemon(6)
Order: ✅ VALID (ascending)
```

Or if violation:
```
Handler: SomeHandler
Locks: UserPokemon(6) → Book(3)
Order: ❌ VIOLATION (descending: 6 → 3)
Fix: Reorder to lock Book first, then UserPokemon
```
