# X3: Documentation

## Status
- [ ] Not started

## Dependencies
- All other tasks

## Tasks
- [ ] README with setup instructions
- [ ] API documentation (endpoints, request/response formats)
- [ ] Local development guide

## Details

### Root README.md
```markdown
# Eurovision Vote Party

A browser-based Eurovision vote party app for hosting watch parties with friends.

## Features

- Create parties for Eurovision semifinals and grand final
- Share party code with guests to join
- Approve/reject join requests
- Cast votes using Eurovision scoring (12, 10, 8, 7, 6, 5, 4, 3, 2, 1)
- See aggregated results when voting closes

## Tech Stack

- **UI**: React + TypeScript + Vite + Tailwind CSS
- **Server**: Go (standard library + Firebase Admin SDK)
- **Database**: Firestore
- **Auth**: Firebase Authentication

## Quick Start

### Prerequisites

- Node.js 22+
- Go 1.22+
- Docker (for Firestore emulator)
- pnpm

### Local Development

1. Start Firestore emulator:
   ```bash
   cd server
   make emulator-start
   ```

2. Start the server:
   ```bash
   cd server
   make run
   ```

3. Start the UI:
   ```bash
   cd ui
   pnpm install
   pnpm dev
   ```

4. Open http://localhost:5173

### Running Tests

```bash
# Server tests
cd server && make test

# UI tests
cd ui && pnpm test
```

## Project Structure

```
eurovision-vote-party/
├── server/          # Go backend
├── ui/              # React frontend
├── tasks/           # Implementation task files
└── README.md
```

See individual project READMEs for more details:
- [Server README](./server/docs/README.md)
- [UI README](./ui/docs/README.md)
```

### Server API Documentation (server/docs/api.md)
```markdown
# API Documentation

Base URL: `http://localhost:8080/api`

## Authentication

Protected endpoints require a Firebase ID token in the Authorization header:
```
Authorization: Bearer <firebase-id-token>
```

## Endpoints

### Health
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /health | No | Health check |

### Parties
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /parties | Yes | Create party |
| GET | /parties | Yes | List admin's parties |
| GET | /parties/:code | No | Get party by code |
| GET | /parties/:id | Yes | Get party details |
| DELETE | /parties/:id | Yes | Delete party |

### Guests
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /parties/:code/join | No | Join party |
| GET | /parties/:code/guest-status | No | Get guest status |
| GET | /parties/:id/guests | Yes | List guests |
| GET | /parties/:id/join-requests | Yes | List pending requests |
| PUT | /parties/:id/guests/:guestId/approve | Yes | Approve guest |
| PUT | /parties/:id/guests/:guestId/reject | Yes | Reject guest |
| DELETE | /parties/:id/guests/:guestId | Yes | Remove guest |

### Acts
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /acts | No | List acts |
| GET | /acts?event=<type> | No | List acts by event |

### Votes
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /parties/:id/votes | No* | Submit vote |
| GET | /parties/:id/votes/:guestId | No* | Get guest's votes |
| PUT | /parties/:id/votes | No* | Update vote |
| PUT | /parties/:id/end-voting | Yes | End voting |
| GET | /parties/:id/results | No | Get results |

*Requires guest ID in request body

### Users
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /users/profile | Yes | Get profile |
| PUT | /users/profile | Yes | Update profile |

## Request/Response Examples

### Create Party
**Request:**
```json
POST /api/parties
Authorization: Bearer <token>
{
  "name": "My Eurovision Party",
  "eventType": "grandfinal"
}
```

**Response:**
```json
201 Created
{
  "id": "abc123",
  "name": "My Eurovision Party",
  "code": "XYZ789",
  "eventType": "grandfinal",
  "adminId": "user123",
  "status": "active",
  "createdAt": "2025-05-10T18:00:00Z"
}
```

### Submit Vote
**Request:**
```json
POST /api/parties/abc123/votes
{
  "guestId": "guest456",
  "votes": {
    "12": "se-2025",
    "10": "it-2025",
    "8": "ua-2025",
    "7": "fr-2025",
    "6": "de-2025",
    "5": "es-2025",
    "4": "uk-2025",
    "3": "nl-2025",
    "2": "ch-2025",
    "1": "no-2025"
  }
}
```

**Response:**
```json
201 Created
{
  "id": "vote789",
  "guestId": "guest456",
  "partyId": "abc123",
  "votes": {...},
  "createdAt": "2025-05-10T20:00:00Z"
}
```

## Error Responses

All errors follow this format:
```json
{
  "error": "Human-readable error message",
  "code": "ERROR_CODE"
}
```

| Status | Code | Description |
|--------|------|-------------|
| 400 | BAD_REQUEST | Invalid request data |
| 400 | VALIDATION_ERROR | Validation failed |
| 401 | UNAUTHORIZED | Missing/invalid auth |
| 403 | FORBIDDEN | Access denied |
| 404 | NOT_FOUND | Resource not found |
| 409 | CONFLICT | Resource conflict |
| 500 | INTERNAL_ERROR | Server error |
```

### UI Documentation (ui/docs/README.md)
```markdown
# Eurovision Vote Party - UI

React + TypeScript frontend for Eurovision Vote Party.

## Development

```bash
# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Run tests
pnpm test

# Build for production
pnpm build
```

## Project Structure

```
ui/
├── src/
│   ├── api/           # API client
│   ├── components/    # Reusable components
│   │   └── ui/        # Base UI components
│   ├── context/       # React contexts
│   ├── hooks/         # Custom hooks
│   ├── pages/         # Page components
│   │   ├── admin/     # Admin pages
│   │   └── guest/     # Guest pages
│   ├── routes/        # Route definitions
│   └── types/         # TypeScript types
└── tests/             # Test files
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_URL | API base URL | http://localhost:8080/api |
| VITE_FIREBASE_API_KEY | Firebase API key | - |
| VITE_FIREBASE_AUTH_DOMAIN | Firebase auth domain | - |
| VITE_FIREBASE_PROJECT_ID | Firebase project ID | - |

## Pages

### Guest Flow
1. `/` - Entry page (enter code + username)
2. `/waiting/:code` - Waiting for approval
3. `/party/:code` - Party overview
4. `/party/:code/vote` - Voting page
5. `/party/:code/results` - Results page

### Admin Flow
1. `/admin/login` - Login page
2. `/admin/profile` - Profile setup
3. `/admin` - Dashboard
4. `/admin/party/new` - Create party
5. `/admin/party/:id` - Party overview
6. `/admin/party/:id/requests` - Join requests
```

## TDD Approach
Documentation is created as the final task, summarizing all implemented features.

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- README provides clear setup instructions
- API docs match implemented endpoints
- Local dev guide is accurate
- All pages documented
