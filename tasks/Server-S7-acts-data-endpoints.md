# S7: Acts Data Endpoints

## Status
- [ ] Not started

## Dependencies
- S4 (Data Models)

## Tasks
- [ ] `GET /api/acts` - List all acts
- [ ] `GET /api/acts?event=semifinal1|semifinal2|grandfinal` - Filter by event
- [ ] Create `data/acts.json` with Eurovision 2025 acts (hardcoded)
- [ ] Load acts from JSON file on server start
- [ ] Include: country, artist, song, running_order, event_type

## Details

### Endpoint Specifications

#### GET /api/acts
**Auth**: None (public)
**Query Params**:
- `event` (optional): `semifinal1`, `semifinal2`, `grandfinal`

**Response** (200):
```json
{
  "acts": [
    {
      "id": "se-2025",
      "country": "Sweden",
      "artist": "Artist Name",
      "song": "Song Title",
      "runningOrder": 1,
      "eventType": "grandfinal"
    }
  ]
}
```

### Acts Data File (`data/acts.json`)
Create JSON file with Eurovision 2025 participating countries.

Structure:
```json
{
  "acts": [
    {
      "id": "country-code-year",
      "country": "Country Name",
      "artist": "Artist Name",
      "song": "Song Title",
      "runningOrder": 1,
      "eventType": "grandfinal"
    }
  ]
}
```

### Eurovision 2025 Events
- **Semifinal 1**: ~18 countries
- **Semifinal 2**: ~18 countries
- **Grand Final**: 26 countries (Big 5 + host + qualifiers)

### Files to Create/Modify
- `handlers/acts.go`
- `services/acts_service.go`
- `data/acts.json`

### Data Loading
- Load acts from JSON on server startup
- Store in memory (no database needed for acts)
- Acts service provides filtered access

## TDD Approach
1. Write service tests for filtering logic
2. Write handler tests with mocked service
3. Implement service and handler
4. Create acts.json data file
5. Verify with `go test ./...`

## Verification
- Acts endpoint returns all acts
- Filtering by event type works
- Running order is correct per event
- All Eurovision 2025 countries included
