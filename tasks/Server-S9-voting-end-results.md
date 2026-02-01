# S9: Voting End & Results

## Status
- [ ] Not started

## Dependencies
- S8 (Voting Endpoints)

## Tasks
- [ ] `PUT /api/parties/:id/end-voting` - Close voting (admin only)
- [ ] Calculate total points per act across all guests
- [ ] `GET /api/parties/:id/results` - Get calculated results (sorted by points)

## Details

### Endpoint Specifications

#### PUT /api/parties/:id/end-voting
**Auth**: Required (admin of party)
**Response** (200):
```json
{
  "id": "party-uuid",
  "status": "closed"
}
```

#### GET /api/parties/:id/results
**Auth**: Required (admin or approved guest)
**Response** (200):
```json
{
  "partyId": "party-uuid",
  "partyName": "My Eurovision Party",
  "totalVoters": 10,
  "results": [
    {
      "rank": 1,
      "actId": "se-2025",
      "country": "Sweden",
      "artist": "Artist Name",
      "song": "Song Title",
      "totalPoints": 120
    },
    {
      "rank": 2,
      "actId": "it-2025",
      "country": "Italy",
      "artist": "Another Artist",
      "song": "Another Song",
      "totalPoints": 95
    }
  ]
}
```

### Results Calculation
1. Get all votes for party
2. Sum points per act across all votes
3. Sort by total points (descending)
4. Assign ranks (handle ties with same rank)

### Business Rules
- Results only available after voting ends
- Results include all acts (even those with 0 points)
- Ties get the same rank

### Files to Create/Modify
- `handlers/votes.go` (add end-voting, results)
- `services/vote_service.go` (add calculation logic)

## TDD Approach
1. Write tests for results calculation logic
2. Write handler tests for end-voting
3. Write handler tests for results endpoint
4. Implement calculation and handlers
5. Verify with `go test ./...`

## Verification
- Admin can end voting
- Party status changes to "closed"
- Cannot submit/update votes after closing
- Results calculated correctly
- Results sorted by points
- Ties handled correctly
