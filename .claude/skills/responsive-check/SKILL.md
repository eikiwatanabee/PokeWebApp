---
name: responsive-check
description: Check and fix responsive design issues across breakpoints. Use when reviewing page layouts.
argument-hint: [page-or-component-path]
---

# Responsive Design Check

Check responsive design for: $ARGUMENTS

## Breakpoints (Tailwind CSS v4)

| Breakpoint | Width | Target |
|---|---|---|
| default | < 640px | Mobile |
| sm | ≥ 640px | Large phone |
| md | ≥ 768px | Tablet |
| lg | ≥ 1024px | Desktop |
| xl | ≥ 1280px | Large desktop |

## Checklist Per Breakpoint

### Mobile (< 640px)
- [ ] Single column layout
- [ ] Touch-friendly tap targets (44px+)
- [ ] No horizontal scroll
- [ ] Navigation: hamburger menu or bottom nav
- [ ] Book cards: full width, stacked
- [ ] Pokemon grid: 2-3 columns
- [ ] Readable text without zooming

### Tablet (768px+)
- [ ] 2-column layouts where appropriate
- [ ] Side-by-side book detail + memos
- [ ] Pokemon grid: 4-5 columns
- [ ] Sidebar navigation option

### Desktop (1024px+)
- [ ] Full layout with sidebars
- [ ] Multi-column book list
- [ ] Pokemon grid: 6+ columns
- [ ] Hover effects active

## Common Issues to Flag
- Fixed widths instead of responsive
- Images without max-width
- Overflow hidden cutting content
- Font sizes too small on mobile
- Modals not fitting mobile screens
