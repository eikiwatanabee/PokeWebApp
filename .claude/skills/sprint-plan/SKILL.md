---
name: sprint-plan
description: Create a sprint plan with prioritized tasks. Use when planning a development iteration.
argument-hint: [sprint-goal]
---

# Sprint Planning

Sprint goal: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for backlog context
2. Review existing code to assess current state
3. Prioritize tasks by dependency order and value

## Output Template

### Sprint Goal
One sentence describing what this sprint delivers.

### Sprint Backlog

| Priority | Task | Type | Layer | Status |
|---|---|---|---|---|
| P0 | Must complete | Feature/Bug/Tech | Domain/App/Infra/FE | Todo |
| P1 | Should complete | ... | ... | Todo |
| P2 | Could complete | ... | ... | Todo |

### Definition of Done
- [ ] Code implemented
- [ ] Tests written and passing
- [ ] /review-ddd passed
- [ ] /review-lock passed (for commands)
- [ ] /security-review passed
- [ ] API documentation updated
- [ ] Frontend responsive on mobile/desktop

### Dependencies & Blockers
- External: PokeAPI availability
- Internal: task dependency chains

### Daily Focus Suggestion
- Day 1-2: Domain + Infrastructure
- Day 3-4: Application Layer + API
- Day 5: Frontend + Integration
- Day 6: Testing + Review
