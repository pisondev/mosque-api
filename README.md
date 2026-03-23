# Mosque API (Core SaaS Engine)

Current release: **v1.0.0**

Backend service for Mosque SaaS platform, built with a Modular Monolith architecture.

## Tech Stack
- Go (Golang)
- Go Fiber v2
- PostgreSQL (`pgx/v5`)
- Migrations via `golang-migrate` (invoked from Go)
- Logrus

## Prerequisites
- Go 1.21+

## Local Setup (Windows, non-Docker)

1. Clone repository and open project root.
2. Setup local PostgreSQL portable:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/setup_local_postgres.ps1
```

3. Set environment for API:

```powershell
$env:APP_PORT='8080'
$env:DB_URL='postgres://root:secretpassword@localhost:5432/mosque_saas?sslmode=disable'
$env:JWT_SECRET='dev-secret'
```

4. Run migration:

```powershell
make migrate-up
```

5. Start API server:

```powershell
make run
```

## Local Setup (manual PostgreSQL)

If you use your own PostgreSQL service, just set:

```env
APP_PORT=8080
DB_URL=postgres://root:secretpassword@localhost:5432/mosque_saas?sslmode=disable
JWT_SECRET=dev-secret
```
