# S4: Data Models

## Status
- [ ] Not started

## Dependencies
- S1 (Project Setup)

## Tasks
- [ ] Define Party model (id, name, code, event_type, admin_id, status, created_at)
- [ ] Define Guest model (id, party_id, username, status: pending/approved/rejected)
- [ ] Define Act model (id, country, artist, song, running_order, event_type)
- [ ] Define Vote model (guest_id, party_id, votes: map[points]act_id)
- [ ] Define User model (id, username, email)
- [ ] Define VoteResult model for calculated results

## Details

### Party Model (`models/party.go`)
```go
type PartyStatus string
const (
    PartyStatusActive PartyStatus = "active"
    PartyStatusClosed PartyStatus = "closed"
)

type EventType string
const (
    EventSemifinal1 EventType = "semifinal1"
    EventSemifinal2 EventType = "semifinal2"
    EventGrandFinal EventType = "grandfinal"
)

type Party struct {
    ID        string      `firestore:"id" json:"id"`
    Name      string      `firestore:"name" json:"name"`
    Code      string      `firestore:"code" json:"code"`
    EventType EventType   `firestore:"eventType" json:"eventType"`
    AdminID   string      `firestore:"adminId" json:"adminId"`
    Status    PartyStatus `firestore:"status" json:"status"`
    CreatedAt time.Time   `firestore:"createdAt" json:"createdAt"`
}
```

### Guest Model (`models/guest.go`)
```go
type GuestStatus string
const (
    GuestStatusPending  GuestStatus = "pending"
    GuestStatusApproved GuestStatus = "approved"
    GuestStatusRejected GuestStatus = "rejected"
)

type Guest struct {
    ID        string      `firestore:"id" json:"id"`
    PartyID   string      `firestore:"partyId" json:"partyId"`
    Username  string      `firestore:"username" json:"username"`
    Status    GuestStatus `firestore:"status" json:"status"`
    CreatedAt time.Time   `firestore:"createdAt" json:"createdAt"`
}
```

### Act Model (`models/act.go`)
```go
type Act struct {
    ID           string    `json:"id"`
    Country      string    `json:"country"`
    Artist       string    `json:"artist"`
    Song         string    `json:"song"`
    RunningOrder int       `json:"runningOrder"`
    EventType    EventType `json:"eventType"`
}
```

### Vote Model (`models/vote.go`)
```go
type Vote struct {
    ID        string         `firestore:"id" json:"id"`
    GuestID   string         `firestore:"guestId" json:"guestId"`
    PartyID   string         `firestore:"partyId" json:"partyId"`
    Votes     map[int]string `firestore:"votes" json:"votes"` // points -> actID
    CreatedAt time.Time      `firestore:"createdAt" json:"createdAt"`
}
```

### User Model (`models/user.go`)
```go
type User struct {
    ID       string `firestore:"id" json:"id"`
    Username string `firestore:"username" json:"username"`
    Email    string `firestore:"email" json:"email"`
}
```

### VoteResult Model (`models/vote.go`)
```go
type VoteResult struct {
    ActID       string `json:"actId"`
    Country     string `json:"country"`
    Artist      string `json:"artist"`
    Song        string `json:"song"`
    TotalPoints int    `json:"totalPoints"`
    Rank        int    `json:"rank"`
}
```

## TDD Approach
1. Write tests for model validation
2. Implement models with validation methods
3. Verify with `go test ./...`

## Verification
- All models compile
- Validation methods work correctly
- JSON/Firestore tags are correct
