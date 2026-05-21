# Security & Authentication

## Overview

This API uses JWT (JSON Web Tokens) for authentication and Role-Based Access Control (RBAC) for authorization.

## Authentication Flow

### 1. Registration

```
User/Company ──► POST /query ──► Register mutation
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Validate Input  │
                                    │ (email, password)│
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Hash Password   │
                                    │ (Argon2id)      │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Insert into DB  │
                                    │ (users/companies)│
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Generate JWT    │
                                    │ with role claim │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Return Token    │
                                    └─────────────────┘
```

### 2. Login

```
User/Company ──► POST /query ──► Login mutation
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ Get User by     │
                                    │ Email from DB   │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Compare Password │
                                    │ (Argon2id)      │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Generate JWT    │
                                    │ with role claim │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │ Return Token    │
                                    └─────────────────┘
```

### 3. Authenticated Request

```
Client ──► POST /query
              │
              ├─ Headers:
              │   Authorization: Bearer <token>
              │
              ▼
    ┌───────────────────┐
    │  Middleware       │
    │  - Extract token  │
    │  - Validate JWT   │
    │  - Set context    │
    │    (userID, role) │
    └─────────┬─────────┘
              │
              ▼
    ┌───────────────────┐
    │  Resolver         │
    │  - Check role     │
    │    (RequireAuth)  │
    │  - Execute query  │
    └─────────┬─────────┘
              │
              ▼
    ┌───────────────────┐
    │  Return Response  │
    └───────────────────┘
```

---

## JWT Structure

### Token Claims

```json
{
  "role": "user",
  "iss": "token",
  "iat": 1700000000,
  "exp": 1700086400,
  "sub": "user-uuid-here"
}
```

### Claim Descriptions

| Claim | Type | Description |
|-------|------|-------------|
| `role` | String | `"user"` or `"company"` |
| `iss` | String | Issuer (always `"token"`) |
| `iat` | Number | Issued at timestamp |
| `exp` | Number | Expiration timestamp (24 hours) |
| `sub` | String | Subject - User or Company UUID |

### Token Generation

```go
// In internal/auth/jwt.go
func Jwt(userID string, role string) (string, error) {
    claims := Claims{
        Role: role,  // "user" or "company"
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "token",
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            Subject:   userID,  // UUID
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
```

### Token Validation

```go
func ValidateJwt(tokenString string, secret string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(
        tokenString,
        claims,
        func(token *jwt.Token) (any, error) {
            return []byte(secret), nil
        },
    )
    
    if err != nil || !token.Valid {
        return nil, err
    }
    
    return claims, nil
}
```

---

## Role-Based Access Control (RBAC)

### Available Roles

| Role | Description |
|------|-------------|
| `user` | Job seekers who can apply to jobs |
| `company` | Employers who can post jobs |

### Permission Matrix

| Operation | user | company | Public |
|-----------|:----:|:-------:|:------:|
| **Authentication** |
| registerUser | ❌ | ❌ | ✅ |
| registerCompany | ❌ | ❌ | ✅ |
| loginUser | ❌ | ❌ | ✅ |
| loginCompany | ❌ | ❌ | ✅ |
| **Jobs (Company)** |
| createJob | ❌ | ✅ | ❌ |
| deleteJob | ❌ | ✅ | ❌ |
| **Jobs (Public)** |
| jobs | ✅ | ✅ | ✅ |
| job(id) | ✅ | ✅ | ✅ |
| jobsByCompany | ✅ | ✅ | ✅ |
| **Skills** |
| skills | ✅ | ✅ | ✅ |
| createSkill | ❌ | ✅ | ❌ |
| addSkillToUser | ✅ | ❌ | ❌ |
| removeSkillFromUser | ✅ | ❌ | ❌ |
| **Applications** |
| applyToJob | ✅ | ❌ | ❌ |
| updateApplicationStatus | ❌ | ✅ | ❌ |
| myApplications | ✅ | ❌ | ❌ |
| jobApplications | ❌ | ✅ | ❌ |
| **Profile** |
| me | ✅ | ✅ | ❌ |

---

## Middleware Architecture

### Middleware Chain

```
Request
    │
    ▼
┌───────────────┐
│   Logger      │ ← Logs request/response
│   Middleware  │   with status & duration
└───────┬───────┘
        │
        ▼
┌───────────────┐
│    Auth       │ ← Extracts & validates JWT
│   Middleware  │   Sets userID & role in context
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ GraphQL       │ ← Resolver checks role
│ Resolver      │   with RequireAuth helpers
└───────┬───────┘
        │
        ▼
┌───────────────┐
│  Database     │ ← Executes query
└───────┬───────┘
        │
        ▼
    Response
```

### Auth Middleware

```go
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        
        if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            secret := os.Getenv("SECRET_JWT")
            
            claims, err := auth.ValidateJwt(tokenString, secret)
            if err == nil {
                // Set context for resolvers
                ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
                ctx = context.WithValue(ctx, RoleKey, claims.Role)
                r = r.WithContext(ctx)
            }
        }
        
        // Always continue - public endpoints work without token
        next.ServeHTTP(w, r)
    })
}
```

**Key Point:** The middleware is **permissive**. It doesn't block requests. It only extracts user info if a valid token exists. Resolvers decide whether to allow/deny.

### RequireAuth Helpers

Resolvers use these helpers to enforce authentication and roles:

```go
// Any authenticated user
func RequireAuth(ctx context.Context) (string, error) {
    userID, ok := ctx.Value(UserIDKey).(string)
    if !ok || userID == "" {
        return "", errors.New("authentication required")
    }
    return userID, nil
}

// Only company role
func RequireCompany(ctx context.Context) (string, error) {
    userID, err := RequireAuth(ctx)
    if err != nil {
        return "", err
    }
    
    role := ctx.Value(RoleKey).(string)
    if role != "company" {
        return "", errors.New("company role required")
    }
    
    return userID, nil
}

// Only user role
func RequireUser(ctx context.Context) (string, error) {
    userID, err := RequireAuth(ctx)
    if err != nil {
        return "", err
    }
    
    role := ctx.Value(RoleKey).(string)
    if role != "user" {
        return "", errors.New("user role required")
    }
    
    return userID, nil
}
```

### Resolver Usage Example

```go
func (r *mutationResolver) CreateJob(ctx context.Context, input model.NewJob) (*model.Job, error) {
    // Enforce company role
    companyID, err := middleware.RequireCompany(ctx)
    if err != nil {
        return nil, err  // Returns error if not authenticated or wrong role
    }
    
    // Continue with business logic...
    pgUUID, _ := auth.ParseUUID(companyID)
    job, _ := r.Queries.CreateJobApp(ctx, ...)
    
    return &model.Job{...}, nil
}
```

---

## Password Security

### Hashing Algorithm

Passwords are hashed using **Argon2id**, a modern memory-hard password hashing algorithm.

```go
// Hash password
func HashPassword(password string) (string, error) {
    hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
    return hash, err
}

// Verify password
func CheckPassword(password, hash string) (bool, error) {
    match, err := argon2id.ComparePasswordAndHash(password, hash)
    return match, err
}
```

### Argon2id Default Parameters

| Parameter | Value |
|-----------|-------|
| Time | 3 |
| Memory | 64MB |
| Threads | 4 |
| Key Length | 32 bytes |
| Salt Length | 16 bytes |

---

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `SECRET_JWT` | JWT signing secret | ✅ Yes |
| `DB_URL` | PostgreSQL connection string | ✅ Yes |
| `PORT` | API server port (default: 8080) | No |
| `ENV` | Environment (`dev` enables introspection) | No |

**Never commit secrets to version control!**

---

## Security Best Practices

### 1. JWT Secret

- Use a **cryptographically secure random string** (32+ characters)
- Never hardcode in source code
- Rotate secrets periodically

**Generate secure secret:**
```bash
openssl rand -base64 32
```

### 2. Database Connection

- Use `sslmode=require` in production
- Use environment variables for credentials
- Never commit connection strings

### 3. HTTPS

- Always use HTTPS in production
- JWTs can be intercepted on unencrypted connections

### 4. Token Expiry

- Tokens expire after **24 hours**
- Consider implementing refresh tokens for longer sessions

---

## Attack Vectors & Mitigations

| Attack | Mitigation |
|--------|------------|
| **SQL Injection** | Uses parameterized queries (sqlc) |
| **XSS** | API-only, no HTML rendering |
| **CSRF** | GraphQL POST-only, token in header |
| **Brute Force** | Consider rate limiting (not implemented) |
| **Password Leaks** | Argon2id hashing, never store plaintext |
| **JWT Forgery** | HS256 signing, secret validation |

---

## Error Messages

For security, error messages are intentionally generic:

| Scenario | Error Message |
|----------|---------------|
| Invalid credentials | `authentication required` |
| Wrong role | `access denied: this operation requires X role` |
| Invalid token | (Middleware silently ignores, resolvers fail auth) |

**Never expose:**
- Database errors
- Whether an email exists in the system
- Internal implementation details
