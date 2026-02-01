# S5: Party Management Endpoints

## Status
- [ ] Not started

## Dependencies
- S3 (Firebase Auth Integration)
- S4 (Data Models)

## Tasks
- [ ] `POST /api/parties` - Create party (admin only, returns party code)
- [ ] `GET /api/parties` - List admin's parties
- [ ] `GET /api/parties/:code` - Get party by code (for guests)
- [ ] `GET /api/parties/:id` - Get party details (for admin)
- [ ] `DELETE /api/parties/:id` - Delete party (admin only)

## Details

### Endpoint Specifications

#### POST /api/parties
**Auth**: Required (admin)
**Request**:
```json
{
  "name": "My Eurovision Party",
  "eventType": "grandfinal"
}
```
**Response** (201):
```json
{
  "id": "uuid",
  "name": "My Eurovision Party",
  "code": "ABC123",
  "eventType": "grandfinal",
  "adminId": "firebase-uid",
  "status": "active",
  "createdAt": "2025-05-10T18:00:00Z"
}
```

#### GET /api/parties
**Auth**: Required (admin)
**Response** (200):
```json
{
  "parties": [...]
}
```

#### GET /api/parties/:code
**Auth**: None (public for guests)
**Response** (200):
```json
{
  "id": "uuid",
  "name": "My Eurovision Party",
  "eventType": "grandfinal",
  "status": "active"
}
```

#### GET /api/parties/:id
**Auth**: Required (admin, must own party)
**Response** (200): Full party object

#### DELETE /api/parties/:id
**Auth**: Required (admin, must own party)
**Response** (204): No content

### Party Code Generation
- 6 alphanumeric characters (uppercase)
- Must be unique
- Easy to read/type (avoid ambiguous chars: 0/O, 1/I/L)

### Files to Create/Modify
- `handlers/party.go`
- `services/party_service.go`
- `persistence/party_dao.go`

## TDD Approach
1. Write handler tests with mocked service
2. Write service tests with mocked DAO
3. Write DAO tests with Firestore emulator
4. Implement each layer
5. Verify with `go test ./...`

## Verification
- All endpoints return correct status codes
- Party code is generated and unique
- Only admin can access protected endpoints
- Admin can only access their own parties
