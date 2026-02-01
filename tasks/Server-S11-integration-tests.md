# S11: Integration Tests & Test Infrastructure

## Status
- [ ] Not started

## Dependencies
- S5 (Party Management Endpoints)
- S6 (Guest Management Endpoints)
- S7 (Acts Data Endpoints)
- S8 (Voting Endpoints)
- S9 (Voting End & Results)
- S10 (Admin User Profile)

## Note
Unit tests are written as part of each task (TDD). This task is for integration tests.

## Tasks
- [ ] End-to-end integration tests with Firestore emulator
- [ ] Test full user flows (create party -> join -> vote -> results)
- [ ] Add test scripts to Makefile (`make test`, `make test-integration`)
- [ ] CI configuration for running tests

## Details

### Integration Test Scenarios

#### Scenario 1: Complete Party Flow
1. Admin creates profile
2. Admin creates party
3. Guest joins party (pending)
4. Admin approves guest
5. Guest submits vote
6. Admin ends voting
7. Results calculated correctly

#### Scenario 2: Multiple Guests
1. Create party with 3 guests
2. All guests approved
3. All guests vote differently
4. Verify aggregated results

#### Scenario 3: Edge Cases
1. Guest rejected (cannot vote)
2. Vote after party closed (rejected)
3. Invalid vote data (rejected)
4. Duplicate party code handling

### Test Infrastructure

#### Makefile Targets
```makefile
test:
	go test ./... -v

test-integration:
	docker-compose up -d firestore
	sleep 5
	FIRESTORE_EMULATOR_HOST=localhost:8081 go test ./... -tags=integration -v
	docker-compose down

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
```

#### Test Tags
Use build tags to separate unit and integration tests:
```go
//go:build integration
```

### CI Configuration
Create `.github/workflows/test.yml`:
- Run unit tests on every push
- Run integration tests with Firestore emulator
- Report coverage

## TDD Approach
1. Write integration test scenarios
2. Ensure all pass with current implementation
3. Add CI configuration
4. Verify with `make test-integration`

## Verification
- All integration tests pass
- Tests run against Firestore emulator
- CI pipeline works
- Coverage report generated
