---
name: ui-component
description: Generate a React component following the project's design system with Tailwind CSS v4. Use when creating new UI components.
argument-hint: [ComponentName]
---

# UI Component Generator

Generate component: $ARGUMENTS

## Instructions

1. Create in `frontend/components/`
2. Use TypeScript + React 19 + Tailwind CSS v4
3. Follow project conventions

## Component Template

```tsx
import { type FC } from 'react'

interface ${ComponentName}Props {
  // Define props
}

export const ${ComponentName}: FC<${ComponentName}Props> = ({ ...props }) => {
  return (
    <div>
      {/* Implementation */}
    </div>
  )
}
```

## Design System Rules

### Colors (Pokemon Theme)
- Primary: Pokemon Red (`#DC0A2D`)
- Secondary: Pokemon Blue (`#3B4CCA`)
- Accent: Pokemon Yellow (`#FFDE00`)
- Background: Neutral grays
- Status: Green (読了), Blue (読書中), Gray (未読)

### Typography
- Headings: Bold, clear hierarchy
- Body: Readable size (16px base)
- Japanese text: system font stack with good JP support

### Spacing
- Consistent 4px grid (p-1, p-2, p-4, p-6, p-8)
- Card padding: p-4 or p-6
- Section gaps: space-y-4 or space-y-6

### Components
- Rounded corners: rounded-lg
- Shadows: shadow-sm for cards
- Transitions: transition-colors duration-200
- Interactive: hover/focus states on all clickable elements

## Accessibility
- Semantic HTML elements
- ARIA labels where needed
- Keyboard navigable
- Focus visible styles
