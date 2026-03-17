# PokeWebApp

PokeWebApp combines a Go API, a Next.js frontend, and PostgreSQL.

## Prerequisites

- Docker Desktop or Docker Engine with Compose
- Go 1.22 or later for local API development
- Node.js 22 and npm for local frontend development

## Startup Commands

All startup commands are collected in `Makefile`.

Run `make help` to list the available targets.

### Full Stack with Docker

```bash
make up
```

Detached mode:

```bash
make up-d
```

Useful follow-up commands:

```bash
make logs
make ps
make down
```

Services:

- Frontend: `http://localhost:3000`
- API: `http://localhost:8080`
- API health check: `http://localhost:8080/api/health`
- PostgreSQL: `localhost:5432`

### Local Development

Install dependencies once:

```bash
make install
```

Start only PostgreSQL in Docker:

```bash
make dev
```

Then run the API and frontend in separate terminals:

```bash
make dev-api
make dev-frontend
```

If a root `.env` file exists, `make dev-api` loads it automatically before starting the Go server. You can start from `.env.example` when you need custom OAuth or JWT values.
