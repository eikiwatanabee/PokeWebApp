---
name: tdd-cycle
description: Run a Test-Driven Development cycle (Red → Green → Refactor). Use when implementing new features with TDD.
argument-hint: [feature-description]
---

# TDD Cycle

Implement a feature using strict TDD for: $ARGUMENTS

## Process

### Phase 1: RED (Write failing test)
1. Create test file in appropriate `_test.go` location
2. Write the minimal test that describes the expected behavior
3. Run `go test ./...` to confirm it FAILS
4. Do NOT write implementation yet

### Phase 2: GREEN (Make it pass)
1. Write the minimal code to make the test pass
2. Run `go test ./...` to confirm it PASSES
3. Do NOT optimize or refactor yet

### Phase 3: REFACTOR
1. Clean up the implementation
2. Remove duplication
3. Ensure all tests still pass
4. Run `/simplify` on changed files

## Testing Conventions
- Table-driven tests for multiple cases
- Use `testify/assert` and `testify/require`
- Mock interfaces with `testify/mock` or `gomock`
- Test file location mirrors source: `entity/book.go` → `entity/book_test.go`
- Name tests: `Test{Function}_{Scenario}_{ExpectedResult}`

## Example

```go
func TestBook_Finish_WhenReading_ShouldTransitionToFinished(t *testing.T) {
    // Arrange
    book := &entity.Book{Status: valueobject.Reading}

    // Act
    err := book.Finish()

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, valueobject.Finished, book.Status)
    assert.NotNil(t, book.FinishedAt)
}

func TestBook_Finish_WhenAlreadyFinished_ShouldReturnError(t *testing.T) {
    book := &entity.Book{Status: valueobject.Finished}
    err := book.Finish()
    assert.ErrorIs(t, err, entity.ErrAlreadyFinished)
}
```
