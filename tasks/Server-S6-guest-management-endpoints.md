# S6: Guest Management Endpoints

## Status
- [ ] Not started

## Dependencies
- S5 (Party Management Endpoints)

## Tasks
- [ ] `POST /api/parties/:code/join` - Request to join party (guest submits username)
- [ ] `GET /api/parties/:id/guests` - List guests in party
- [ ] `GET /api/parties/:id/join-requests` - List pending join requests (admin only)
- [ ] `PUT /api/parties/:id/guests/:guestId/approve` - Approve guest (admin only)
- [ ] `PUT /api/parties/:id/guests/:guestId/reject` - Reject guest (admin only)
- [ ] `DELETE /api/parties/:id/guests/:guestId` - Remove guest (admin only)

## Details

### Endpoint Specifications

#### POST /api/parties/:code/join
**Auth**: None (public)
**Request**:
```json
{
  "username": "JohnDoe"
}
```
**Response** (201):
```json
{
  "id": "guest-uuid",
  "partyId": "party-uuid",
  "username": "JohnDoe",
  "status": "pending",
  "createdAt": "2025-05-10T18:30:00Z"
}
```

#### GET /api/parties/:id/guests
**Auth**: Required (admin of party) or Guest token
**Response** (200):
```json
{
  "guests": [
    {
      "id": "guest-uuid",
      "username": "JohnDoe",
      "status": "approved"
    }
  ]
}
```

#### GET /api/parties/:id/join-requests
**Auth**: Required (admin of party)
**Response** (200):
```json
{
  "requests": [
    {
      "id": "guest-uuid",
      "username": "JaneDoe",
      "status": "pending",
      "createdAt": "2025-05-10T18:35:00Z"
    }
  ]
}
```

#### PUT /api/parties/:id/guests/:guestId/approve
**Auth**: Required (admin of party)
**Response** (200):
```json
{
  "id": "guest-uuid",
  "status": "approved"
}
```

#### PUT /api/parties/:id/guests/:guestId/reject
**Auth**: Required (admin of party)
**Response** (200):
```json
{
  "id": "guest-uuid",
  "status": "rejected"
}
```

#### DELETE /api/parties/:id/guests/:guestId
**Auth**: Required (admin of party)
**Response** (204): No content

### Guest Status Flow
```
pending -> approved
pending -> rejected
approved -> (removed via DELETE)
```

### Guest Status Polling (for waiting page)
Add endpoint for guests to check their status:
#### GET /api/parties/:code/guest-status?guestId=xxx
**Auth**: None (guest ID acts as simple token)
**Response** (200):
```json
{
  "status": "pending" | "approved" | "rejected"
}
```

### Files to Create/Modify
- `handlers/guest.go`
- `services/guest_service.go`
- `persistence/guest_dao.go`

## TDD Approach
1. Write handler tests with mocked service
2. Write service tests with mocked DAO
3. Write DAO tests with Firestore emulator
4. Implement each layer
5. Verify with `go test ./...`

## Verification
- Guest can join party with username
- Admin can see pending requests
- Admin can approve/reject guests
- Guest status endpoint works for polling
