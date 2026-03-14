---
name: plan-feature
description: Break down a feature into requirements, tasks, and estimates. Use when starting a new feature implementation.
argument-hint: [feature-name]
---

# Feature Planning

Plan feature: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for project context
2. Break down the feature into deliverables

## Output Template

### Feature Overview
- **Name**: Feature name
- **Goal**: What problem does this solve?
- **User Story**: As a [role], I want [action], so that [benefit]

### Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

### Technical Design
- **Domain Changes**: New entities, VOs, domain events
- **Application Layer**: New commands and queries
- **Infrastructure**: New repos, external APIs
- **Presentation**: New endpoints, handlers
- **Frontend**: New pages, components

### Task Breakdown

| # | Task | Layer | Dependencies |
|---|---|---|---|
| 1 | Create Entity | Domain | - |
| 2 | Add Repository Interface | Domain | 1 |
| 3 | Implement Repository | Infra | 2 |
| 4 | Create Command Handler | App | 2, 3 |
| 5 | Create Query Handler | App | 2, 3 |
| 6 | Add API Endpoint | Presentation | 4, 5 |
| 7 | Create Frontend Page | Frontend | 6 |
| 8 | Write Tests | All | 1-7 |

### Risks & Considerations
- Lock ordering implications
- Performance concerns
- Security considerations
- PokeAPI rate limits
