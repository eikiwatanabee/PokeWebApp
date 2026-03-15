---
name: docker-optimize
description: Optimize Dockerfiles and docker-compose configuration. Use when improving build times or image sizes.
---

# Docker Optimization

Analyze and optimize Docker configuration.

## Checklist

### Dockerfile (Go API)
- [ ] Multi-stage build (build stage + runtime stage)
- [ ] Alpine-based final image
- [ ] Non-root user in final image
- [ ] `.dockerignore` excludes unnecessary files
- [ ] Go binary statically compiled (`CGO_ENABLED=0`)
- [ ] Layer caching optimized (COPY go.mod/go.sum first)

### Dockerfile (Next.js)
- [ ] Multi-stage build (deps → build → runtime)
- [ ] Standalone output mode
- [ ] Non-root user
- [ ] Only production dependencies in final image

### docker-compose.yml
- [ ] Health checks for all services
- [ ] Proper depends_on with condition
- [ ] Volume mounts for development (hot reload)
- [ ] Environment variables via .env file
- [ ] Network isolation

### Performance
- [ ] Build cache utilized
- [ ] Parallel builds enabled
- [ ] Image sizes minimized
- [ ] Startup order correct (DB → API → Frontend)

## Ideal Image Sizes
- Go API: < 30MB
- Next.js: < 200MB
- PostgreSQL: standard image
