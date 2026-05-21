# Contributing to Job Board API

Thank you for your interest in contributing to this project! This document outlines the guidelines and best practices for contributing.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [How Can I Contribute?](#how-can-i-contribute)
3. [Development Setup](#development-setup)
4. [Pull Request Guidelines](#pull-request-guidelines)
5. [Coding Standards](#coding-standards)
6. [Commit Messages](#commit-messages)

---

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone.

### Our Standards

Examples of behavior that contributes to creating a positive environment:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

---

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the [existing issues](../../issues) to avoid duplicates.

**When submitting a bug report, please include:**

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version
- PostgreSQL version
- Any relevant logs or error messages

### Suggesting Enhancements

We welcome suggestions for improvements! When suggesting an enhancement:

- Use a clear, descriptive title
- Provide a detailed description of the proposed feature
- Explain why this enhancement would be useful
- Include code examples if applicable

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests (if applicable)
5. Ensure the code compiles and works
6. Submit a pull request

---

## Development Setup

### Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Goose (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

### Setup Steps

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/mohamed8eo/jobBoard.git
   cd jobBoard
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   createdb jobboard
   ```

4. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. **Run migrations**
   ```bash
   export DB_URL="postgres://user:pass@localhost:5432/jobboard?sslmode=disable"
   goose -dir internal/db/migrations postgres "$DB_URL" up
   ```

6. **Start the development server**
   ```bash
   air
   # Or: go run cmd/api/main.go
   ```

---

## Pull Request Guidelines

### Branch Naming

Use the following naming convention for branches:

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/<description>` | `feature/add-email-verification` |
| Bug Fix | `fix/<description>` | `fix/login-timeout-bug` |
| Documentation | `docs/<description>` | `docs/update-readme` |
| Refactor | `refactor/<description>` | `refactor/auth-middleware` |

### PR Checklist

Before submitting a PR, please ensure:

- [ ] Code compiles successfully (`go build ./...`)
- [ ] Tests pass (if applicable)
- [ ] Code is formatted (`gofmt -w .`)
- [ ] Documentation is updated (if needed)
- [ ] Commit messages follow the [conventional commits](#commit-messages) format
- [ ] No breaking changes (or clearly documented)

### What to Expect

After submitting your PR:

1. The code will be reviewed
2. Feedback may be provided
3. Changes may be requested
4. Once approved, your PR will be merged

---

## Coding Standards

### General Guidelines

1. **Follow Go conventions**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://go.dev/doc/effective_go)

2. **Keep functions small and focused**
   - Ideally under 50 lines
   - One function, one responsibility

3. **Handle errors properly**
   - Never ignore errors
   - Return meaningful error messages
   - Wrap errors with context

4. **Write comments**
   - Public functions must have docstrings
   - Complex logic needs explanation
   - No unnecessary comments

### Project-Specific Standards

#### GraphQL Schema

- All queries/mutations must have clear return types
- Use meaningful names (no abbreviations)
- Document with comments in schema

#### Resolvers

- Keep resolvers thin (delegate logic)
- Always check context for auth
- Return appropriate errors

#### Database Queries

- Use parameterized queries (sqlc handles this)
- Keep queries simple
- One query file per table

### Directory Structure

- **`cmd/`**: Application entry points
- **`graph/`**: GraphQL schema and resolvers
- **`internal/`**: Private application code
  - **`auth/`**: JWT and password handling
  - **`db/`**: Database layer (migrations, queries, generated code)
  - **`middleware/`**: HTTP middleware
  - **`validator/`**: Input validation

---

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, whitespace, etc. |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | Performance improvement |
| `test` | Adding tests |
| `chore` | Build, CI, or maintenance |

### Examples

```
feat(graphql): add job search query
```

```
fix(auth): correct password validation
```

```
docs(readme): update setup instructions
```

```
refactor(db): simplify query logic
```

---

## Testing

### Manual Testing

1. **GraphQL Playground**
   - Start server: `air`
   - Visit: `http://localhost:8080/`
   - Test queries and mutations

2. **API Testing**
   ```bash
   # Example: Register a user
   curl -X POST http://localhost:8080/query \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { registerUser(input: { name:\"Test\", email:\"test@test.com\", password:\"Test123!\" }) { token } }"}'
   ```

### Adding Tests

When adding new features, consider:

1. **Unit tests** for utility functions
2. **Integration tests** for database operations
3. **GraphQL tests** for resolvers

---

## Database Changes

### Adding a New Table

1. **Create migration**
   ```bash
   goose -dir internal/db/migrations create my_new_table sql
   ```

2. **Write migration SQL**
   - Always include `+goose Up` and `+goose Down`
   - Make sure down migration is correct

3. **Add queries**
   - Add to appropriate file in `internal/db/queries/`
   - Follow sqlc comment format: `-- name: QueryName :one` or `:many`

4. **Regenerate**
   ```bash
   sqlc generate
   ```

### Modifying Existing Tables

**Always create a new migration.** Never edit existing migrations.

```bash
goose -dir internal/db/migrations create add_column_to_table sql
```

---

## Getting Help

If you need help with contributing:

1. Check the [README.md](README.md) for setup instructions
2. Look at [DEVELOPMENT.md](DEVELOPMENT.md) for detailed setup
3. Search existing [issues](../../issues)
4. Create a new issue if your question isn't addressed

---

## License

By contributing, you agree that your contributions will be licensed under the project's [LICENSE](LICENSE) (MIT License).

---

Thank you for contributing! 🎉
