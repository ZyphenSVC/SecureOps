# SecureOps

SecureOps is a Dockerized Go/PostgreSQL backend for a secure internal operations platform. It implements authentication, JWT-protected routes, role-based access control, permission-guarded admin endpoints, and audit logging.

## Tech Stack

- Go
- PostgreSQL
- Docker Compose
- JWT
- bcrypt
- RBAC

## Features

- User registration and login
- Password hashing with bcrypt
- JWT authentication
- Protected `/auth/me` route
- Role-based permission checks
- Admin user listing protected by `users:read`
- Audit log listing protected by `audit_logs:read`
- PostgreSQL migrations
- Dockerized local development

## Run Locally

```bash
docker compose up --build

curl -i http://localhost:8080/health
curl -i http://localhost:8080/db-health

TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Password123!"
  }' | jq -r .token)

curl -i http://localhost:8080/auth/me \
  -H "Authorization: Bearer $TOKEN"

curl -i http://localhost:8080/admin/users \
  -H "Authorization: Bearer $TOKEN"

curl -i http://localhost:8080/admin/audit-logs \
  -H "Authorization: Bearer $TOKEN"
```