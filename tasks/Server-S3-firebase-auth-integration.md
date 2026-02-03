# S3: Firebase Auth Integration

## Status
- [x] Completed

## Dependencies
- S2 (Docker & Firestore Setup)

## Tasks
- [x] Add Firebase Admin SDK for token verification
- [x] Create auth middleware for protected endpoints
- [x] Implement token validation logic

## Details

### Auth Middleware
Located in `middleware/auth.go`

```go
// AuthMiddleware verifies Firebase ID tokens
func AuthMiddleware(next http.Handler) http.Handler
```

### Token Validation
- Extract `Authorization: Bearer <token>` header
- Verify token with Firebase Admin SDK
- Set user ID in request context
- Return 401 if invalid/missing token

### Protected vs Public Endpoints
**Protected (require auth)**:
- All `/api/parties` endpoints (except GET by code)
- All `/api/users` endpoints

**Public (no auth)**:
- `GET /api/health`
- `GET /api/parties/:code` (for guests)
- `POST /api/parties/:code/join`
- `GET /api/acts`

## TDD Approach
1. Write tests for middleware with mock tokens
2. Implement middleware
3. Verify with `go test ./...`

## Verification
- Requests without token return 401
- Requests with invalid token return 401
- Requests with valid token pass through with user ID in context
