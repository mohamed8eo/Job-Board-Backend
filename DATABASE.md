# Database Schema

## Overview

This project uses PostgreSQL as the database backend. Database queries are type-safe thanks to [sqlc](https://sqlc.dev/). Migrations are managed with [Goose](https://github.com/pressly/goose).

## Schema Diagram

```
┌─────────────┐       ┌─────────────┐
│  companies  │       │    users    │
├─────────────┤       ├─────────────┤
│ id (UUID)   │       │ id (UUID)   │
│ name        │       │ name        │
│ email       │       │ email       │
│ password    │       │ password    │
│ created_at  │       │ created_at  │
└──────┬──────┘       └──────┬──────┘
       │                     │
       │ ┌───────────────────┘
       │ │
       ▼ ▼
┌─────────────────┐
│   job_apps      │
├─────────────────┤
│ id (UUID)       │
│ company_id (FK) │◄─────────────────────
│ title           │                      │
│ description     │                      │
│ location        │                      │
│ remote          │                      │
│ salary          │                      │
│ created_at      │                      │
└──────┬──────────┘                      │
       │                                 │
       │           ┌─────────────────────┤
       │           │                     │
       ▼           ▼                     │
┌───────────────┐ ┌─────────────────┐    │
│  job_skills   │ │  applications   │    │
├───────────────┤ ├─────────────────┤    │
│ id (UUID)     │ │ id (UUID)       │    │
│ job_app_id(FK)│ │ user_id (FK)    │────┘
│ skill_id (FK) │ │ job_app_id (FK) │
└──────┬────────┘ │ status          │
       │          │ created_at      │
       │          └─────────────────┘
       │
       ▼
┌─────────────┐       ┌───────────────┐
│   skills    │       │  user_skills  │
├─────────────┤       ├───────────────┤
│ id (UUID)   │◄──────│ id (UUID)     │
│ name        │       │ user_id (FK)  │
└─────────────┘       │ skill_id (FK) │
                      └───────────────┘
```

## Tables

### companies

Companies/employers who post jobs.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique company identifier |
| `name` | TEXT | NOT NULL | Company name |
| `email` | TEXT | NOT NULL, UNIQUE | Company email |
| `password` | TEXT | NOT NULL | Hashed password |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT now() | Creation timestamp |

---

### users

Job seekers who apply for jobs.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique user identifier |
| `name` | TEXT | NOT NULL | User's full name |
| `email` | TEXT | NOT NULL, UNIQUE | User's email |
| `password` | TEXT | NOT NULL | Hashed password |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT now() | Creation timestamp |

---

### skills

Master list of skills that can be associated with jobs or users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique skill identifier |
| `name` | TEXT | NOT NULL | Skill name (e.g., "Go", "Python") |

---

### job_apps

Job postings created by companies.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique job identifier |
| `company_id` | UUID | NOT NULL, FOREIGN KEY to companies | Company that posted the job |
| `title` | TEXT | NOT NULL | Job title |
| `description` | TEXT | NOT NULL | Job description |
| `location` | TEXT | NOT NULL | Job location |
| `remote` | BOOLEAN | NOT NULL | Whether job is remote |
| `salary` | INT4 | NULLABLE | Salary amount |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT now() | Creation timestamp |

---

### user_skills

Junction table connecting users to skills (many-to-many).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique identifier |
| `user_id` | UUID | NOT NULL, FOREIGN KEY to users | User who has the skill |
| `skill_id` | UUID | NOT NULL, FOREIGN KEY to skills | The skill |

---

### job_skills

Junction table connecting jobs to skills (many-to-many).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique identifier |
| `job_app_id` | UUID | NOT NULL, FOREIGN KEY to job_apps | Job requiring the skill |
| `skill_id` | UUID | NOT NULL, FOREIGN KEY to skills | The skill |

---

### applications

Applications from users to jobs.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique application ID |
| `user_id` | UUID | NOT NULL, FOREIGN KEY to users | User who applied |
| `job_app_id` | UUID | NOT NULL, FOREIGN KEY to job_apps | Job applied to |
| `status` | TEXT | DEFAULT 'PENDING' | Application status |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT now() | Creation timestamp |

---

## Relationships

| Table | Related Table | Relationship |
|-------|---------------|--------------|
| `companies` | `job_apps` | One-to-Many (1 company → many jobs) |
| `users` | `applications` | One-to-Many (1 user → many apps) |
| `job_apps` | `applications` | One-to-Many (1 job → many apps) |
| `users` | `skills` | Many-to-Many (via `user_skills`) |
| `job_apps` | `skills` | Many-to-Many (via `job_skills`) |

---

## Migration Files

Migrations are located in `internal/db/migrations/`:

| File | Description |
|------|-------------|
| `00001_companies.sql` | Create companies table |
| `00002_users.sql` | Create users table |
| `00003_skills.sql` | Create skills table |
| `00004_job_apps.sql` | Create job_apps table with foreign key to companies |
| `00005_job_skills.sql` | Create job_skills junction table |
| `00006_user_skills.sql` | Create user_skills junction table |
| `00007_applications.sql` | Create applications table |

---

## Running Migrations

### Using Goose

```bash
# Set database URL
export DB_URL="postgres://user:password@localhost:5432/jobboard?sslmode=disable"

# Run all migrations
goose -dir internal/db/migrations postgres "$DB_URL" up

# Rollback last migration
goose -dir internal/db/migrations postgres "$DB_URL" down

# Check status
goose -dir internal/db/migrations postgres "$DB_URL" status
```

### From Go Code

See `internal/db/queries/` for SQL queries and `internal/db/sqlc/` for generated Go code.

---

## Query Organization

### SQL Queries (`internal/db/queries/`)

| File | Purpose |
|------|---------|
| `company.sql` | Company CRUD operations |
| `user.sql` | User CRUD operations |
| `skill.sql` | Skill CRUD operations |
| `job_app.sql` | Job posting CRUD |
| `job_skill.sql` | Job-skill associations |
| `user_skill.sql` | User-skill associations |
| `application.sql` | Application CRUD + status updates |

### Generated Code (`internal/db/sqlc/`)

| File | Generated From |
|------|----------------|
| `company.sql.go` | `company.sql` + sqlc config |
| `user.sql.go` | `user.sql` + sqlc config |
| ... | ... |

---

## sqlc Configuration

See `sqlc.yaml` for sqlc configuration. To regenerate after modifying queries:

```bash
sqlc generate
```
