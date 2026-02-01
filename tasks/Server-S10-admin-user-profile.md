# S10: Admin User Profile

## Status
- [ ] Not started

## Dependencies
- S3 (Firebase Auth Integration)

## Tasks
- [ ] `PUT /api/users/profile` - Set/update admin username
- [ ] `GET /api/users/profile` - Get admin profile
- [ ] Store admin profiles in Firestore

## Details

### Endpoint Specifications

#### PUT /api/users/profile
**Auth**: Required
**Request**:
```json
{
  "username": "AdminName"
}
```
**Response** (200):
```json
{
  "id": "firebase-uid",
  "username": "AdminName",
  "email": "admin@example.com"
}
```

#### GET /api/users/profile
**Auth**: Required
**Response** (200):
```json
{
  "id": "firebase-uid",
  "username": "AdminName",
  "email": "admin@example.com"
}
```
**Response** (404) - if no profile exists:
```json
{
  "error": "Profile not found"
}
```

### Profile Creation Flow
1. User logs in with Firebase Auth
2. UI calls GET /api/users/profile
3. If 404, redirect to profile setup page
4. User enters username
5. UI calls PUT /api/users/profile
6. Redirect to dashboard

### Validation
- Username required
- Username 3-30 characters
- Username alphanumeric + underscores only

### Files to Create/Modify
- `handlers/user.go`
- `services/user_service.go`
- `persistence/user_dao.go`

## TDD Approach
1. Write handler tests with mocked service
2. Write service tests with mocked DAO
3. Write DAO tests with Firestore emulator
4. Implement each layer
5. Verify with `go test ./...`

## Verification
- Can create profile
- Can update profile
- Returns 404 if no profile
- Username validation works
- Email populated from Firebase token
