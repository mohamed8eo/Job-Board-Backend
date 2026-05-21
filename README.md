# Job Board API (GraphQL)

A modern, production-ready job board API built with Go, featuring role-based access control, JWT authentication, and a GraphQL interface.

## Features

- 🔐 **JWT Authentication** - Secure user and company authentication
- 🎭 **Role-Based Access Control** - Separate permissions for users and companies
- 📊 **PostgreSQL Backend** - Robust data persistence
- 🎯 **GraphQL Interface** - Flexible querying with gqlgen
- 🔒 **Secure Password Hashing** - Argon2id for password protection
- 📝 **Database Migrations** - Versioned schema management
- 🚀 **Hot Reload** - Air for development

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| GraphQL | gqlgen |
| Database | PostgreSQL |
| Query Builder | sqlc |
| Migrations | Goose |
| Auth | JWT + Argon2id |
| Dev Server | Air |

## Project Structure

```
jobBoard/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── graph/
│   ├── schema.graphqls       # GraphQL schema definition
│   ├── schema.resolvers.go   # Resolver implementations
│   ├── model/
│   │   └── models_gen.go     # Generated GraphQL models
│   └── generated.go          # Generated gqlgen code
├── internal/
│   ├── auth/
│   │   └── jwt.go            # JWT generation and validation
│   ├── db/
│   │   ├── migrations/       # Database migrations
│   │   ├── queries/          # SQL queries for sqlc
│   │   └── sqlc/             # Generated database code
│   ├── middleware/
│   │   ├── auth.go           # Authentication middleware
│   │   └── logger.go         # Request logging middleware
│   └── validator/
│       └── validator.go      # Input validation
├── go.mod
├── go.sum
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL
- Air (for hot reload)
- Goose (for migrations)

### Installation

```bash
# Clone the repository
git clone https://github.com/mohamed8eo/jobBoard.git
cd jobBoard

# Install dependencies
go mod download

# Create .env file (see DEVELOPMENT.md for details)
cp .env.example .env

# Run database migrations
just up 

#Run  sqlc queries 
just sqlc

# Start development server
just air 
```

### GraphQL Playground

The GraphQL playground is available at:
```
http://localhost:8080/
```

API endpoint:
```
http://localhost:8080/query
```

## Documentation

- [API Reference](API.md) - Complete GraphQL schema and examples
- [Database Schema](DATABASE.md) - Data model and migrations
- [Security & Authentication](SECURITY.md) - JWT structure and RBAC
- [Development Guide](DEVELOPMENT.md) - Setup and code generation
- [Contributing](CONTRIBUTING.md) - How to contribute

## API Examples

### Public Queries

```graphql
# Get all jobs
query {
  jobs {
    id
    title
    description
    location
    remote
    salary
  }
}

# Get available skills
query {
  skills {
    id
    name
  }
}
```

### Protected Mutations (Company)

```graphql
# Create a job posting (company only)
mutation {
  createJob(input: {
    title: "Senior Go Developer"
    description: "Looking for experienced Go developers"
    location: "Remote"
    remote: true
    salary: 150000
  }) {
    id
    title
  }
}
```

### Protected Mutations (User)

```graphql
# Apply to a job (user only)
mutation {
  applyToJob(jobID: "uuid-goes-here") {
    id
    status
    createdAt
  }
}
```

See [API Reference](API.md) for complete documentation.

## License

MIT License - feel free to use this in your projects or as a reference.

## Author

Created by [mohamed8eo](https://github.com/mohamed8eo)
