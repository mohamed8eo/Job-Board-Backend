# Development Guide

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.25+ | Programming language |
| PostgreSQL | 14+ | Database |
| Goose | Latest | Database migrations |
| sqlc | Latest | Query compiler |
| gqlgen | Latest | GraphQL code generation |
| Air | Latest | Hot reload (development) |

---

## Installation

### Install Dependencies

```bash
# Clone the repository
git clone https://github.com/mohamed8eo/jobBoard.git
cd jobBoard

# Install Go modules
go mod download
```

### Install Tools

```bash
# Goose (migrations)
go install github.com/pressly/goose/v3/cmd/goose@latest

# sqlc (query generation)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# gqlgen (GraphQL generation)
go install github.com/99designs/gqlgen@latest

# Air (hot reload)
go install github.com/cosmtrek/air@latest
```

---

## Environment Setup

### 1. Create Environment File

```bash
cp .env.example .env
```

### 2. Configure Variables

Edit `.env`:

```env
# Server
PORT=8080
ENV=dev

# Database
DB_URL=postgres://postgres:password@localhost:5432/jobboard?sslmode=disable

# JWT
SECRET_JWT=your-super-secret-jwt-key-at-least-32-characters-long
```

### Generate Secure JWT Secret

```bash
openssl rand -base64 32
```

---

## Database Setup

### 1. Create PostgreSQL Database

```bash
# Using psql
psql -U postgres
CREATE DATABASE jobboard;
\q

# Or using createdb
createdb -U postgres jobboard
```

### 2. Run Migrations

```bash
# Set DB_URL environment
export DB_URL="postgres://postgres:password@localhost:5432/jobboard?sslmode=disable"

# Run all migrations
goose -dir internal/db/migrations postgres "$DB_URL" up
```

### Migration Commands

| Command | Action |
|---------|--------|
| `goose up` | Apply all pending migrations |
| `goose down` | Rollback last migration |
| `goose status` | Show migration status |
| `goose reset` | Rollback ALL migrations |
| `goose redo` | Rollback and reapply last migration |

### Create New Migration

```bash
goose -dir internal/db/migrations create add_something_table sql
```

---

## Code Generation

This project uses **gqlgen** for GraphQL and **sqlc** for database queries.

### gqlgen (GraphQL)

**Files involved:**
- `graph/schema.graphqls` - Your schema definition
- `graph/schema.resolvers.go` - Your resolver implementations
- `graph/generated.go` - **Generated** code (don't edit)
- `graph/model/models_gen.go` - **Generated** models (don't edit)

**When to regenerate:**
- After modifying `schema.graphqls`
- After creating new queries/mutations

```bash
# Regenerate GraphQL code
go generate ./...

# Or directly
gqlgen generate
```

### sqlc (Database)

**Files involved:**
- `internal/db/queries/*.sql` - Your SQL queries
- `internal/db/sqlc/*.sql.go` - **Generated** code (don't edit)
- `internal/db/sqlc/models.go` - **Generated** models (don't edit)
- `sqlc.yaml` - sqlc configuration

**When to regenerate:**
- After adding/modifying queries in `internal/db/queries/`
- After adding new database tables

```bash
# Regenerate database code
sqlc generate
```

### sqlc.yaml Configuration

```yaml
version: "2"
sql:
  - schema: "internal/db/migrations"
    queries: "internal/db/queries"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "internal/db/sqlc"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
```

---

## Running the Server

### Development (Hot Reload)

```bash
# Using Air (recommended)
air

# Or using go run
go run cmd/api/main.go
```

### Production

```bash
# Build
go build -o server cmd/api/main.go

# Run
./server
```

---

## GraphQL Playground

When `ENV=dev`, the GraphQL playground is available:

- **Playground URL:** `http://localhost:8080/`
- **API Endpoint:** `http://localhost:8080/query`

---

## Project Structure Deep Dive

```
jobBoard/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
│
├── graph/                    # GraphQL layer
│   ├── schema.graphqls       # Schema definition (EDIT THIS)
│   ├── schema.resolvers.go   # Resolver logic (EDIT THIS)
│   ├── resolver.go           # Resolver base
│   ├── model/
│   │   └── models_gen.go     # Generated models (DON'T EDIT)
│   └── generated.go          # Generated gqlgen code (DON'T EDIT)
│
├── internal/                 # Private application code
│   ├── auth/
│   │   └── jwt.go            # JWT and password utilities
│   │
│   ├── db/
│   │   ├── migrations/       # Database migrations (EDIT HERE)
│   │   ├── queries/          # SQL queries (EDIT HERE)
│   │   └── sqlc/             # Generated code (DON'T EDIT)
│   │
│   ├── middleware/
│   │   ├── auth.go           # Authentication middleware
│   │   └── logger.go         # Request logging
│   │
│   └── validator/
│       └── validator.go      # Input validation
│
├── go.mod
├── go.sum
└── sqlc.yaml                 # sqlc configuration
```

---

## Development Workflow

### Adding a New GraphQL Query/Mutation

1. **Edit** `graph/schema.graphqls`:
   ```graphql
   type Query {
       # Existing...
       myNewQuery: String!
   }
   ```

2. **Regenerate** code:
   ```bash
   go generate ./...
   ```

3. **Implement** the resolver in `graph/schema.resolvers.go`:
   ```go
   func (r *queryResolver) MyNewQuery(ctx context.Context) (string, error) {
       // Your logic here
       return "Hello!", nil
   }
   ```

---

### Adding a New Database Query

1. **Edit** the appropriate `.sql` file in `internal/db/queries/`:
   ```sql
   -- name: GetUserByEmail :one
   SELECT * FROM users WHERE email = $1;
   ```

2. **Regenerate** code:
   ```bash
   sqlc generate
   ```

3. **Use** the generated function in your resolver:
   ```go
   user, err := r.Queries.GetUserByEmail(ctx, "email@example.com")
   ```

---

### Adding a New Database Table

1. **Create** a new migration:
   ```bash
   goose -dir internal/db/migrations create new_table sql
   ```

2. **Edit** the migration file:
   ```sql
   -- +goose Up
   CREATE TABLE new_table (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name TEXT NOT NULL
   );
   
   -- +goose Down
   DROP TABLE new_table;
   ```

3. **Run** the migration:
   ```bash
   goose -dir internal/db/migrations postgres "$DB_URL" up
   ```

4. **Add** queries in `internal/db/queries/`

5. **Regenerate** with `sqlc generate`

---

## Configuration Files

### gqlgen.yml

Located in project root. Controls GraphQL generation.

```yaml
schema:
  - graph/schema.graphqls

exec:
  filename: graph/generated.go
  package: graph

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph
```

### sqlc.yaml

Controls database code generation.

```yaml
version: "2"
sql:
  - schema: "internal/db/migrations"
    queries: "internal/db/queries"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "internal/db/sqlc"
        emit_json_tags: true
        emit_interface: true
```

### .air.toml

Hot reload configuration for Air.

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"
  exclude_dir = ["tmp", "vendor"]
  include_ext = ["go", "graphql", "graphqls"]
```

---

## Common Commands

| Task | Command |
|------|---------|
| Start dev server | `air` |
| Build server | `go build -o server cmd/api/main.go` |
| Regenerate GraphQL | `go generate ./...` |
| Regenerate DB code | `sqlc generate` |
| Run migrations | `goose -dir internal/db/migrations postgres "$DB_URL" up` |
| Check dependencies | `go mod tidy` |
| Format code | `gofmt -w .` |

---

## Troubleshooting

### Migration Issues

**Problem:** Goose can't connect to database.

**Solution:** Ensure `DB_URL` is correct and PostgreSQL is running.

```bash
# Test connection
psql "postgres://postgres:password@localhost:5432/jobboard?sslmode=disable"
```

**Problem:** Migration failed and DB is in inconsistent state.

**Solution:**
```bash
# Check status
goose status

# Force version
goose fix <version>
```

### Code Generation Issues

**Problem:** gqlgen not generating resolvers.

**Solution:** Check `gqlgen.yml` configuration and run:
```bash
go generate ./...
```

**Problem:** sqlc complaining about unknown types.

**Solution:** Regenerate after running migrations:
```bash
goose up
sqlc generate
```

### Connection Issues

**Problem:** Server starts but can't connect to DB.

**Solution:**
1. Ensure PostgreSQL is running
2. Check `DB_URL` in `.env`
3. Verify database exists
4. Check SSL mode (`sslmode=disable` for local dev)

---

## Useful Links

- [gqlgen Documentation](https://gqlgen.com/)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [Goose Documentation](https://pressly.github.io/goose/)
- [GraphQL Spec](https://spec.graphql.org/)
