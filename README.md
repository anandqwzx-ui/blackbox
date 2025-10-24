# Crudbox

Crudbox is a two-part application for designing and serving mocked HTTP APIs. The backend is a Golang microservice powered by Gin and PostgreSQL, and the frontend is a minimalist black-and-white Next.js dashboard for managing organisations, projects, and endpoints.

## Project Structure

```
backend/   # Go service with Gin
frontend/  # Next.js application (App Router)
```

## Backend

### Features

- Email + password authentication with JWT sessions.
- Organisation, project, and endpoint management with soft-delete audit fields.
- Per-project five-character alphanumeric codes for publicly accessible mock endpoints.
- Raw SQL data access powered by `sqlx` on top of the pgx driver (no ORM).
- Built-in migration runner that executes embedded SQL files on startup.
- Graceful shutdown and signal handling for safe restarts.
- IDOR protections: every data interaction checks ownership before proceeding.

### Getting Started

1. Copy the environment template and adjust values if needed:
   ```bash
   cp backend/.env.example backend/.env
   ```
2. Ensure PostgreSQL is running and accessible via the connection string from the `.env` file.
3. Run the API service:
   ```bash
   cd backend
   go run ./cmd/server
   ```

### Database migrations

SQL migrations live in `backend/internal/database/migrations` and are executed automatically whenever the API starts. If you need to apply them without keeping the server running (for example during CI deploy steps), you can run:

```
cd backend
go run ./cmd/server --migrate-only
```

The server listens on `PORT` (default `8080`). API routes are namespaced under `/api/v1`, while mock endpoints resolve directly from the root using the project code (e.g. `/{code}/path`).

## Frontend

### Features

- Dark monochrome aesthetic with glassmorphism-inspired panels.
- Signup and login flows that persist JWTs in `localStorage`.
- Dashboard to manage organisations, projects, and mock endpoints end-to-end.
- Inline editing for projects and endpoints, including toggleable enable/disable switches.

### Getting Started

1. Copy the environment template and adjust the API base URL if the backend runs elsewhere:
   ```bash
   cp frontend/.env.example frontend/.env.local
   ```
2. Install dependencies and run the development server:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

Visit `http://localhost:3000` for the UI. The dashboard expects the backend to be reachable at the URL defined by `NEXT_PUBLIC_API_BASE_URL` (defaults to `http://localhost:8080/api/v1`).

## Mock Endpoints

Once you have created a project and defined endpoints, mock responses are served directly from the backend. For example, if a project code is `aB12C` and an endpoint is registered at `/users` for `GET`, requests to `http://localhost:8080/aB12C/users` will respond with the configured payload, status code, and headers.

## Tooling

- **Backend:** Go 1.24, Gin, sqlx, pgx, JWT, bcrypt.
- **Frontend:** Next.js 14 (App Router), React 18, TypeScript.

Feel free to extend the platform with additional features such as shared projects, request logging, or versioned endpoint definitions.
