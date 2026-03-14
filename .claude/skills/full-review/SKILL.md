---
name: full-review
description: Comprehensive multi-perspective code review covering architecture, security, performance, and quality. Use before major merges.
---

# Full Code Review

Run a comprehensive review on recent changes.

## Review Dimensions

### 1. Architecture (invoke /review-ddd mentally)
- Clean Architecture layer violations
- DDD pattern compliance
- CQRS separation
- Lock ordering

### 2. Security (invoke /security-review mentally)
- OWASP Top 10
- Auth/AuthZ
- Data leakage
- Input validation

### 3. Performance
- N+1 query detection
- Missing database indexes
- Unnecessary allocations
- PokeAPI call optimization (caching)
- Large payload responses

### 4. Code Quality
- Error handling consistency
- Naming conventions (Go idioms)
- Test coverage
- Dead code
- Duplicated logic

### 5. Go Best Practices
- Context propagation
- Goroutine leak risks
- Interface design (accept interfaces, return structs)
- Error wrapping with `fmt.Errorf("%w", err)`

## Output Format

```
## Architecture: ✅ / ⚠️ / ❌
[Findings...]

## Security: ✅ / ⚠️ / ❌
[Findings...]

## Performance: ✅ / ⚠️ / ❌
[Findings...]

## Code Quality: ✅ / ⚠️ / ❌
[Findings...]

## Summary
- Critical: N issues
- Warnings: N issues
- Suggestions: N improvements
```
