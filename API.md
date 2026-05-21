# GraphQL API Reference

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Role-Based Permissions](#role-based-permissions)
4. [Queries](#queries)
5. [Mutations](#mutations)
6. [Examples](#examples)

---

## Overview

Base URL: `http://localhost:8080/query`

This is a GraphQL API. All requests are POST requests to the `/query` endpoint.

## Authentication

Most endpoints require authentication using JWT tokens. Send the token in the Authorization header:

```
Authorization: Bearer <your-token>
```

Tokens are issued when you register or login.

## Role-Based Permissions

| Endpoint | User Role | Company Role | Public |
|----------|:----------:|:-----------:|:------:|
| registerUser | - | - | ✅ |
| registerCompany | - | - | ✅ |
| loginUser | - | - | ✅ |
| loginCompany | - | - | ✅ |
| jobs | - | - | ✅ |
| job(id) | - | - | ✅ |
| jobsByCompany | - | - | ✅ |
| skills | - | - | ✅ |
| createJob | ❌ | ✅ | ❌ |
| deleteJob | ❌ | ✅ | ❌ |
| createSkill | ❌ | ✅ | ❌ |
| addSkillToUser | ✅ | ❌ | ❌ |
| removeSkillFromUser | ✅ | ❌ | ❌ |
| applyToJob | ✅ | ❌ | ❌ |
| updateApplicationStatus | ❌ | ✅ | ❌ |
| me | ✅ | ✅ | ❌ |
| myApplications | ✅ | ❌ | ❌ |
| jobApplications | ❌ | ✅ | ❌ |

---

## Queries

### jobs

Get all job listings.

```graphql
query {
  jobs {
    id
    title
    description
    location
    remote
    salary
    createdAt
  }
}
```

**Response:**
```json
{
  "data": {
    "jobs": [
      {
        "id": "uuid-1",
        "title": "Senior Go Developer",
        "description": "Join our growing team...",
        "location": "Remote",
        "remote": true,
        "salary": 150000,
        "createdAt": "2024-..."
      }
    ]
  }
}
```

---

### job(id)

Get a specific job by ID.

```graphql
query {
  job(id: "job-uuid-here") {
    id
    title
    description
    location
    remote
    salary
  }
}
```

**Arguments:
- `id` (ID!, required) - Job UUID

---

### jobsByCompany(companyID)

Get all jobs posted by a specific company.

```graphql
query {
  jobsByCompany(companyID: "company-uuid-here") {
    id
    title
    description
  }
}
```

**Arguments:**
- `companyID` (ID!, required) - Company UUID

---

### skills

Get all available skills.

```graphql
query {
  skills {
    id
    name
  }
}
```

**Response:**
```json
{
  "data": {
    "skills": [
      { "id": "uuid-1", "name": "Go" },
      { "id": "uuid-2", "name": "Python" },
      { "id": "uuid-3", "name": "PostgreSQL" }
    ]
  }
}
```

---

### me

Get the currently authenticated user or company.

**Requires Authentication**

```graphql
query {
  me {
    id
    name
    email
    createdAt
  }
}
```

**Response (User):**
```json
{
  "data": {
    "me": {
      "id": "user-uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2024-..."
    }
  }
}
```

---

### myApplications

Get all applications submitted by the current user.

**Requires Authentication (User role only)**

```graphql
query {
  myApplications {
    id
    status
    createdAt
    job {
      id
      title
    }
  }
}
```

**Response:**
```json
{
  "data": {
    "myApplications": [
      {
        "id": "app-uuid",
        "status": "PENDING",
        "createdAt": "2024-...",
        "job": {
          "id": "job-uuid",
          "title": "Senior Go Developer"
        }
      }
    ]
  }
}
```

---

### jobApplications(jobID)

Get all applications for a specific job.

**Requires Authentication (Company role only)

```graphql
query {
  jobApplications(jobID: "job-uuid-here") {
    id
    status
    createdAt
    user {
      id
      name
      email
    }
  }
}
```

**Arguments:**
- `jobID` (ID!, required) - Job UUID

---

## Mutations

### registerUser

Register as a job seeker.

```graphql
mutation {
  registerUser(input: {
    name: "John Doe"
    email: "john@example.com"
    password: "SecurePass123!"
  }) {
    token
  }
}
```

**Input:**
- `name` (String!, required) - User's name
- `email` (String!, required) - User's email
- `password` (String!, required) - User's password

**Response:**
```json
{
  "data": {
    "registerUser": {
      "token": "eyJhbGciOiJIUzI1NiIs..."
    }
  }
}
```

---

### registerCompany

Register as a company/employer.

```graphql
mutation {
  registerCompany(input: {
    name: "Tech Company Inc."
    email: "careers@techcompany.com"
    password: "SecurePass123!"
  }) {
    token
  }
}
```

**Input:**
- `name` (String!, required) - Company name
- `email` (String!, required) - Company email
- `password` (String!, required) - Company password

---

### loginUser

Login as an existing user.

```graphql
mutation {
  loginUser(input: {
    email: "john@example.com"
    password: "SecurePass123!"
  }) {
    token
  }
}
```

---

### loginCompany

Login as an existing company.

```graphql
mutation {
  loginCompany(input: {
    email: "careers@techcompany.com"
    password: "SecurePass123!"
  }) {
    token
  }
}
```

---

### createJob

Create a new job posting.

**Requires Authentication (Company role only)**

```graphql
mutation {
  createJob(input: {
    title: "Senior Go Developer"
    description: "We are looking for..."
    location: "Remote"
    remote: true
    salary: 150000
  }) {
    id
    title
    description
  }
}
```

**Input:**
- `title` (String!, required) - Job title
- `description` (String!, required) - Job description
- `location` (String!, required) - Job location
- `remote` (Boolean!, required) - Whether job is remote
- `salary` (Int, optional) - Salary amount

---

### deleteJob

Delete a job posting.

**Requires Authentication (Company role only)**

```graphql
mutation {
  deleteJob(id: "job-uuid-here")
}
```

**Arguments:**
- `id` (ID!, required) - Job UUID to delete

**Response:**
```json
{
  "data": {
    "deleteJob": true
  }
}
```

---

### createSkill

Create a new skill.

**Requires Authentication (Company role only)**

```graphql
mutation {
  createSkill(input: {
    name: "Rust"
  }) {
    id
    name
  }
}
```

---

### addSkillToUser

Add a skill to the current user.

**Requires Authentication (User role only)**

```graphql
mutation {
  addSkillToUser(skillID: "skill-uuid-here") {
    id
    name
  }
}
```

---

### removeSkillFromUser

Remove a skill from the current user.

**Requires Authentication (User role only)**

```graphql
mutation {
  removeSkillFromUser(skillID: "skill-uuid-here") {
    id
    name
  }
}
```

---

### applyToJob

Apply to a job posting.

**Requires Authentication (User role only)**

```graphql
mutation {
  applyToJob(jobID: "job-uuid-here") {
    id
    status
    createdAt
  }
}
```

**Response:**
```json
{
  "data": {
    "applyToJob": {
      "id": "app-uuid",
      "status": "PENDING",
      "createdAt": "2024-..."
    }
  }
}
```

---

### updateApplicationStatus

Update the status of an application.

**Requires Authentication (Company role only)**

```graphql
mutation {
  updateApplicationStatus(input: {
    applicationID: "app-uuid-here"
    status: "ACCEPTED"
  }) {
    id
    status
  }
}
```

**Input:**
- `applicationID` (ID!, required) - Application UUID
- `status` (String!, required) - New status

Common status values: `"PENDING"`, `"ACCEPTED"`, `"REJECTED"`, `"INTERVIEW"`, `"OFFER"`

---

## Examples

### cURL Examples

**Register User:**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { registerUser(input: { name: \"John\", email: \"john@test.com\", password: \"Test123!\" }) { token } }"
  }'
```

**Get All Jobs:**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { jobs { id title } }"
  }'
```

**Protected Request (with token):**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR..." \
  -d '{
    "query": "query { me { id name email } }"
  }'
```
