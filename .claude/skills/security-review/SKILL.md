---
name: security-review
description: Review code for security vulnerabilities (OWASP Top 10, auth issues, data leakage). Use before merging features.
---

# Security Review

Analyze recent changes for security vulnerabilities.

## Checklist

### Authentication & Authorization
- [ ] All API endpoints require authentication (except /auth/*)
- [ ] JWT tokens are validated on every request
- [ ] Tenant isolation: users can only access their tenant's data
- [ ] User can only access/modify their own books and memos
- [ ] Admin-only endpoints check user role

### Injection
- [ ] SQL: All queries use GORM parameter binding (no string concatenation)
- [ ] XSS: User input is sanitized before rendering
- [ ] Command Injection: No shell execution with user input

### Data Protection
- [ ] No sensitive data in logs (tokens, passwords, emails in bulk)
- [ ] No sensitive data in error responses
- [ ] CORS configured to allow only expected origins
- [ ] JWT secret is not hardcoded

### API Security
- [ ] Rate limiting on auth endpoints
- [ ] Input validation on all request bodies
- [ ] Proper HTTP status codes (don't leak info with 404 vs 403)
- [ ] No mass assignment (explicit field binding, not blind bind)

### PokeAPI Specific
- [ ] PokeAPI responses are validated before storing
- [ ] No SSRF via user-controlled URLs

## Output
- **CRITICAL**: Immediate security risk, must fix
- **HIGH**: Significant risk, fix before production
- **MEDIUM**: Should be addressed
- **LOW**: Best practice recommendation
