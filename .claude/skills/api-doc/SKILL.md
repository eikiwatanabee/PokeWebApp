---
name: api-doc
description: Generate or update API documentation from handler code. Use after adding or modifying endpoints.
argument-hint: [handler-file-path]
---

# API Documentation Generator

Generate docs for: $ARGUMENTS

## Instructions

1. Read the handler file(s)
2. Extract endpoint definitions
3. Generate documentation

## Output Format (per endpoint)

### `METHOD /api/path`

**Description**: What this endpoint does

**Authentication**: Required / Not required

**Request**:
```json
{
  "field": "type - description"
}
```

**Response** (200):
```json
{
  "field": "type - description"
}
```

**Errors**:
| Status | Code | Description |
|---|---|---|
| 400 | INVALID_INPUT | Validation failed |
| 401 | UNAUTHORIZED | Missing/invalid token |
| 404 | NOT_FOUND | Resource not found |

**Example**:
```bash
curl -X POST http://localhost:8080/api/books/123/finish \
  -H "Authorization: Bearer <token>"
```

## Rules
- Document all request/response fields
- Include all possible error responses
- Show curl examples
- Note CQRS type (Command/Query)
- Reference lock order for Commands
