# S1: Server Project Setup

## Status
- [ ] Not started

## Dependencies
None

## Tasks
- [ ] Initialize Go module (`go.mod`)
- [ ] Create main.go with HTTP server using standard lib
- [ ] Setup project structure: `handlers/`, `services/`, `persistence/`, `models/`
- [ ] Add health check endpoint

## Details

### Project Structure
```
server/
├── main.go
├── go.mod
├── handlers/
├── services/
├── persistence/
├── models/
├── middleware/
└── data/
```

### Health Check Endpoint
- `GET /api/health` - Returns 200 OK with `{"status": "ok"}`

## TDD Approach
1. Write test for health endpoint first
2. Implement health handler
3. Verify with `go test ./...`

## Verification
- Server starts on port 8080
- Health endpoint returns 200 OK
