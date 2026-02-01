# Eurovision Vote Party - Implementation Tasks

This directory contains all implementation tasks for the Eurovision Vote Party application.

## Task Organization

Tasks are prefixed by their domain:
- `Server-S*` - Go backend tasks
- `UI-U*` - React frontend tasks
- `X*` - Cross-cutting tasks

## Implementation Order

### Phase 1: Foundation (Parallel)
| Server | UI |
|--------|-----|
| S1 - Project Setup | U1 - Project Setup |
| S2 - Docker/Firestore | U17 - Component Library |
| S4 - Data Models | U5 - Type Definitions |

### Phase 2: Auth
| Server | UI |
|--------|-----|
| S3 - Firebase Auth | U2 - Routing |
| S10 - User Profile | U3 - Firebase Auth |

### Phase 3: Core Features
| Server | UI |
|--------|-----|
| S5 - Party Endpoints | U4 - API Client |
| S6 - Guest Endpoints | U6 - Guest Entry |
| S7 - Acts Endpoints | U7 - Waiting Page |
| S8 - Voting Endpoints | U8 - Party Overview |
| | U9 - Voting Page |

### Phase 4: Admin Features
| Server | UI |
|--------|-----|
| S9 - Results | U11 - Login Page |
| | U12 - Profile Setup |
| | U13 - Dashboard |
| | U14 - Create Party |
| | U15 - Admin Overview |
| | U16 - Join Requests |
| | U10 - Results Page |

### Phase 5: Polish & Testing
| Server | UI | Cross-cutting |
|--------|-----|---------------|
| S11 - Integration Tests | U18 - Integration Tests | X1 - CORS Config |
| | | X2 - Error Handling |
| | | X3 - Documentation |

## Task Dependencies

```
Server: S1 → S2 → S3 → S5 → S6 → S8 → S9
                ↘ S4 ↗      ↘ S7 ↗
                  ↘ S10

UI: U1 → U2 → U3 → U4 → U6 → U7 → U8 → U9
      ↘ U17        ↘ U5      ↘ U11 → U12 → U13 → U14 → U15 → U16
                                                          ↘ U10
```

## Task File Format

Each task file contains:
1. **Status** - Checkbox list of sub-tasks
2. **Dependencies** - Required completed tasks
3. **Tasks** - Implementation checklist
4. **Details** - Code snippets, specifications
5. **TDD Approach** - Test-first workflow
6. **Verification** - Acceptance criteria

## TDD Workflow

Every task follows this workflow:
1. **Write tests first** - Define expected behavior
2. **Implement & verify** - Make tests pass
3. **Write documentation** - For coding agents

## All Tasks

### Server Tasks (11)
- [Server-S1-project-setup.md](./Server-S1-project-setup.md)
- [Server-S2-docker-firestore-setup.md](./Server-S2-docker-firestore-setup.md)
- [Server-S3-firebase-auth-integration.md](./Server-S3-firebase-auth-integration.md)
- [Server-S4-data-models.md](./Server-S4-data-models.md)
- [Server-S5-party-management-endpoints.md](./Server-S5-party-management-endpoints.md)
- [Server-S6-guest-management-endpoints.md](./Server-S6-guest-management-endpoints.md)
- [Server-S7-acts-data-endpoints.md](./Server-S7-acts-data-endpoints.md)
- [Server-S8-voting-endpoints.md](./Server-S8-voting-endpoints.md)
- [Server-S9-voting-end-results.md](./Server-S9-voting-end-results.md)
- [Server-S10-admin-user-profile.md](./Server-S10-admin-user-profile.md)
- [Server-S11-integration-tests.md](./Server-S11-integration-tests.md)

### UI Tasks (18)
- [UI-U1-project-setup.md](./UI-U1-project-setup.md)
- [UI-U2-routing-layout.md](./UI-U2-routing-layout.md)
- [UI-U3-firebase-auth-setup.md](./UI-U3-firebase-auth-setup.md)
- [UI-U4-api-client.md](./UI-U4-api-client.md)
- [UI-U5-type-definitions.md](./UI-U5-type-definitions.md)
- [UI-U6-guest-entry-page.md](./UI-U6-guest-entry-page.md)
- [UI-U7-guest-waiting-page.md](./UI-U7-guest-waiting-page.md)
- [UI-U8-party-overview-page.md](./UI-U8-party-overview-page.md)
- [UI-U9-voting-page.md](./UI-U9-voting-page.md)
- [UI-U10-results-page.md](./UI-U10-results-page.md)
- [UI-U11-admin-login-page.md](./UI-U11-admin-login-page.md)
- [UI-U12-admin-profile-setup.md](./UI-U12-admin-profile-setup.md)
- [UI-U13-admin-dashboard.md](./UI-U13-admin-dashboard.md)
- [UI-U14-party-creation-page.md](./UI-U14-party-creation-page.md)
- [UI-U15-admin-party-overview-page.md](./UI-U15-admin-party-overview-page.md)
- [UI-U16-join-requests-page.md](./UI-U16-join-requests-page.md)
- [UI-U17-component-library.md](./UI-U17-component-library.md)
- [UI-U18-integration-tests.md](./UI-U18-integration-tests.md)

### Cross-cutting Tasks (3)
- [X1-cors-api-configuration.md](./X1-cors-api-configuration.md)
- [X2-error-handling-loading-states.md](./X2-error-handling-loading-states.md)
- [X3-documentation.md](./X3-documentation.md)

**Total: 32 tasks**
