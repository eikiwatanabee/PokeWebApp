---
name: test-gen
description: Automatically generate comprehensive tests for existing code. Use when you need to add test coverage.
argument-hint: [file-or-package-path]
---

# Test Generator

Generate tests for: $ARGUMENTS

## Instructions

1. Read the target file/package
2. Identify all public functions and methods
3. Generate table-driven tests covering:
   - Happy path
   - Edge cases (empty input, nil, zero values)
   - Error cases
   - Boundary conditions

## Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr error
    }{
        {name: "valid input", input: ..., want: ..., wantErr: nil},
        {name: "empty input", input: ..., want: ..., wantErr: ErrEmpty},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Layer-Specific Guidelines

### Domain Layer Tests
- No mocks needed (pure logic)
- Test entity constructors, methods, value object validation

### Application Layer Tests
- Mock repositories and external services
- Test command/query handlers
- Verify UoW commit/rollback behavior

### Presentation Layer Tests
- Use `httptest` for HTTP testing
- Test request parsing, response format, error handling
- Mock application layer handlers

## Run after generation
```bash
go test ./... -v -cover
```
