---
name: a11y-check
description: Accessibility audit based on WCAG 2.1 guidelines. Use when reviewing frontend components for accessibility compliance.
argument-hint: [component-or-page-path]
---

# Accessibility Audit

Audit accessibility for: $ARGUMENTS

## WCAG 2.1 Checklist

### Perceivable
- [ ] Images have alt text (Pokemon sprites, book covers)
- [ ] Color is not the only way to convey information
- [ ] Sufficient color contrast (4.5:1 for text, 3:1 for large text)
- [ ] Text resizable up to 200% without loss
- [ ] Captions/labels for form inputs

### Operable
- [ ] All functionality available via keyboard
- [ ] No keyboard traps
- [ ] Skip navigation link available
- [ ] Focus indicators visible
- [ ] Touch targets at least 44x44px

### Understandable
- [ ] Language attribute set (`lang="ja"`)
- [ ] Form labels and instructions clear
- [ ] Error messages identify the field and suggest fix
- [ ] Consistent navigation across pages

### Robust
- [ ] Valid HTML semantics (headings hierarchy, landmarks)
- [ ] ARIA roles used correctly
- [ ] Works with screen readers
- [ ] Compatible with assistive technologies

## Pokemon-Specific
- [ ] Pokemon type colors have text labels too
- [ ] Catch animation has non-visual alternative
- [ ] Pokedex grid navigable by keyboard

## Tools to Run
```bash
# In frontend directory
npx axe-core-cli http://localhost:3000
npx pa11y http://localhost:3000
```
