---
name: ux-review
description: UI/UX heuristic evaluation based on Nielsen's 10 usability heuristics. Use when reviewing frontend components or pages.
argument-hint: [page-or-component-path]
---

# UX Heuristic Review

Review UI/UX for: $ARGUMENTS

## Nielsen's 10 Usability Heuristics

### 1. Visibility of System Status
- Loading states shown during API calls?
- Book status changes reflected immediately?
- Pokemon catch animation/feedback?

### 2. Match Between System and Real World
- Pokemon terminology used correctly?
- Book status labels intuitive? (未読/読書中/読了)
- Icons match expectations?

### 3. User Control and Freedom
- Undo available for status changes?
- Easy navigation back?
- Can revert from 読書中 to 未読?

### 4. Consistency and Standards
- Consistent button styles and placements?
- Same patterns for similar actions?
- Design system followed?

### 5. Error Prevention
- Confirmation before destructive actions (delete book)?
- Prevent duplicate book registration?
- Validate form inputs before submission?

### 6. Recognition Rather Than Recall
- Tags shown as suggestions during book registration?
- Recent books easily accessible?
- Pokemon sprites shown (not just names)?

### 7. Flexibility and Efficiency of Use
- Keyboard shortcuts for power users?
- Quick actions on book list?
- Search/filter easily accessible?

### 8. Aesthetic and Minimalist Design
- Only essential information shown?
- Visual hierarchy clear?
- Pokemon theme enhances, not distracts?

### 9. Help Users Recognize, Diagnose, and Recover from Errors
- Clear error messages?
- Suggested actions after errors?
- Form validation messages next to fields?

### 10. Help and Documentation
- Onboarding flow for new users?
- Tooltips for non-obvious features?

## Output
Rate each heuristic: ✅ Good / ⚠️ Needs Work / ❌ Violation
Provide specific fix recommendations.
