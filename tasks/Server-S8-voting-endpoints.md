# S8: Voting Endpoints

## Status
- [ ] Not started

## Dependencies
- S6 (Guest Management Endpoints)
- S7 (Acts Data Endpoints)

## Tasks
- [ ] `POST /api/parties/:id/votes` - Submit vote (guest/admin)
- [ ] `GET /api/parties/:id/votes/:guestId` - Get guest's votes
- [ ] `PUT /api/parties/:id/votes` - Update vote (guest/admin)
- [ ] Validate: each act selected exactly once, all 10 point values used

## Details

### Eurovision Scoring System
Points: 12, 10, 8, 7, 6, 5, 4, 3, 2, 1 (10 acts receive points)

### Endpoint Specifications

#### POST /api/parties/:id/votes
**Auth**: Guest ID or Admin token
**Request**:
```json
{
  "guestId": "guest-uuid",
  "votes": {
    "12": "act-id-1",
    "10": "act-id-2",
    "8": "act-id-3",
    "7": "act-id-4",
    "6": "act-id-5",
    "5": "act-id-6",
    "4": "act-id-7",
    "3": "act-id-8",
    "2": "act-id-9",
    "1": "act-id-10"
  }
}
```
**Response** (201):
```json
{
  "id": "vote-uuid",
  "guestId": "guest-uuid",
  "partyId": "party-uuid",
  "votes": {...},
  "createdAt": "2025-05-10T19:00:00Z"
}
```

#### GET /api/parties/:id/votes/:guestId
**Auth**: Guest ID matches or Admin of party
**Response** (200):
```json
{
  "id": "vote-uuid",
  "guestId": "guest-uuid",
  "votes": {...}
}
```

#### PUT /api/parties/:id/votes
**Auth**: Guest ID or Admin token
**Request**: Same as POST
**Response** (200): Updated vote object

### Validation Rules
1. Exactly 10 point values must be provided: 12, 10, 8, 7, 6, 5, 4, 3, 2, 1
2. Each act ID can only appear once
3. All act IDs must be valid acts for the party's event type
4. Guest must be approved to vote
5. Party must be in "active" status (voting not ended)

### Error Responses
- 400: Invalid vote structure
- 403: Guest not approved / Party voting closed
- 404: Party or guest not found

### Files to Create/Modify
- `handlers/votes.go`
- `services/vote_service.go`
- `persistence/vote_dao.go`

## TDD Approach
1. Write validation tests first
2. Write handler tests with mocked service
3. Write service tests with mocked DAO
4. Write DAO tests with Firestore emulator
5. Implement each layer
6. Verify with `go test ./...`

## Verification
- Can submit complete vote
- Validation rejects incomplete votes
- Validation rejects duplicate acts
- Cannot vote if not approved
- Cannot vote if party closed
- Can update existing vote
