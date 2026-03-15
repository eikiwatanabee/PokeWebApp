---
name: retrospective
description: Run an implementation retrospective analyzing what went well, what didn't, and improvements. Use after completing features.
---

# Retrospective

Analyze the recent implementation session.

## Instructions

1. Review git log for recent commits
2. Review code changes (git diff)
3. Check test coverage
4. Analyze architectural decisions made

## Analysis Template

### What Went Well
- Patterns that worked
- Efficient implementations
- Good test coverage areas

### What Could Be Better
- Architectural concerns
- Code smells detected
- Missing test coverage
- Performance issues

### Action Items

| # | Action | Priority | Category |
|---|---|---|---|
| 1 | ... | High/Med/Low | Code/Arch/Process |

### Technical Debt Identified
- Shortcuts taken that need revisiting
- TODOs left in code
- Missing error handling
- Hardcoded values to externalize

### Metrics
- Files changed: N
- Lines added: N
- Lines removed: N
- Test coverage: N%
- New commands: N
- New queries: N
