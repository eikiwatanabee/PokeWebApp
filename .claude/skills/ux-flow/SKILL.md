---
name: ux-flow
description: Design user flows and screen transition diagrams. Use when planning new features or pages.
argument-hint: [feature-name]
---

# User Flow Designer

Design user flow for: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for feature context
2. Design the user flow as a Mermaid diagram
3. Identify screens, actions, and decision points
4. Consider error states and edge cases

## Output Format

### Flow Diagram (Mermaid)
```mermaid
graph TD
    A[Start] --> B{User logged in?}
    B -->|No| C[Login Page]
    B -->|Yes| D[Dashboard]
    C --> D
    D --> E[Book List]
    E --> F[Book Detail]
    F --> G{Action?}
    G -->|Start Reading| H[Status: 読書中]
    G -->|Finish| I[Pokemon Catch!]
    I --> J[Show Pokemon]
    J --> K[Back to Book Detail]
```

### Screen List
For each screen:
- **Name**: Screen name
- **Route**: URL path
- **Components**: Key UI components
- **Data**: Required API calls
- **Actions**: User interactions available
- **Next**: Possible navigation targets

### Edge Cases
- What happens if PokeAPI is down?
- What if user has no books yet?
- Empty states for each screen
