---
name: user-story
description: Generate user stories with acceptance criteria in standard format. Use when defining requirements for new features.
argument-hint: [feature-description]
---

# User Story Generator

Generate user stories for: $ARGUMENTS

## Format

### User Story
```
As a [ユーザータイプ],
I want to [アクション],
so that [目的/価値].
```

### Acceptance Criteria (Given-When-Then)
```
Given [前提条件],
When [アクション],
Then [期待結果].
```

## Instructions

1. Read `docs/REQUIREMENTS.md` for context
2. Identify user types: 一般ユーザー, 管理者
3. Generate stories covering:
   - Happy path
   - Edge cases
   - Error scenarios
4. Prioritize: MUST / SHOULD / COULD

## Example

### US-001: 本を読了してポケモンをゲット
```
As a 一般ユーザー,
I want to mark a book as finished,
so that I can catch a random Pokemon.
```

**Acceptance Criteria:**
```
Given I am reading a book (status: 読書中),
When I click the "読了" button,
Then the book status changes to "読了",
And I receive a random Pokemon from PokeAPI,
And the Pokemon appears in my Pokedex,
And I see a catch animation with the Pokemon's sprite.
```

```
Given I have a book with status "未読",
When I try to mark it as finished,
Then I see an error message "先に読書を開始してください".
```

## Output
Generate 3-5 user stories with full acceptance criteria for the given feature.
