# U4: API Client

## Status
- [ ] Not started

## Dependencies
- U3 (Firebase Auth Setup)

## Tasks
- [ ] Create base API client with fetch
- [ ] Add auth token injection for admin requests
- [ ] Create typed API functions for all endpoints
- [ ] Add error handling

## Details

### Base API Client
```typescript
// src/api/client.ts
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

interface ApiOptions {
  authenticated?: boolean
}

async function apiRequest<T>(
  path: string,
  options: RequestInit & ApiOptions = {}
): Promise<T> {
  const { authenticated = false, ...fetchOptions } = options

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...fetchOptions.headers,
  }

  if (authenticated) {
    const token = await getIdToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...fetchOptions,
    headers,
  })

  if (!response.ok) {
    throw new ApiError(response.status, await response.text())
  }

  return response.json()
}
```

### API Functions

#### Party API
```typescript
// src/api/parties.ts
export const partyApi = {
  create: (data: CreatePartyRequest) =>
    apiRequest<Party>('/parties', {
      method: 'POST',
      body: JSON.stringify(data),
      authenticated: true
    }),

  list: () =>
    apiRequest<{ parties: Party[] }>('/parties', { authenticated: true }),

  getByCode: (code: string) =>
    apiRequest<Party>(`/parties/${code}`),

  getById: (id: string) =>
    apiRequest<Party>(`/parties/${id}`, { authenticated: true }),

  delete: (id: string) =>
    apiRequest<void>(`/parties/${id}`, { method: 'DELETE', authenticated: true }),
}
```

#### Guest API
```typescript
// src/api/guests.ts
export const guestApi = {
  join: (code: string, username: string) =>
    apiRequest<Guest>(`/parties/${code}/join`, {
      method: 'POST',
      body: JSON.stringify({ username })
    }),

  getStatus: (code: string, guestId: string) =>
    apiRequest<{ status: GuestStatus }>(`/parties/${code}/guest-status?guestId=${guestId}`),

  list: (partyId: string) =>
    apiRequest<{ guests: Guest[] }>(`/parties/${partyId}/guests`, { authenticated: true }),

  approve: (partyId: string, guestId: string) =>
    apiRequest<Guest>(`/parties/${partyId}/guests/${guestId}/approve`, {
      method: 'PUT',
      authenticated: true
    }),

  reject: (partyId: string, guestId: string) =>
    apiRequest<Guest>(`/parties/${partyId}/guests/${guestId}/reject`, {
      method: 'PUT',
      authenticated: true
    }),
}
```

#### Acts API
```typescript
// src/api/acts.ts
export const actsApi = {
  list: (eventType?: EventType) =>
    apiRequest<{ acts: Act[] }>(`/acts${eventType ? `?event=${eventType}` : ''}`),
}
```

#### Votes API
```typescript
// src/api/votes.ts
export const votesApi = {
  submit: (partyId: string, guestId: string, votes: VoteData) =>
    apiRequest<Vote>(`/parties/${partyId}/votes`, {
      method: 'POST',
      body: JSON.stringify({ guestId, votes })
    }),

  get: (partyId: string, guestId: string) =>
    apiRequest<Vote>(`/parties/${partyId}/votes/${guestId}`),

  update: (partyId: string, guestId: string, votes: VoteData) =>
    apiRequest<Vote>(`/parties/${partyId}/votes`, {
      method: 'PUT',
      body: JSON.stringify({ guestId, votes })
    }),

  endVoting: (partyId: string) =>
    apiRequest<Party>(`/parties/${partyId}/end-voting`, {
      method: 'PUT',
      authenticated: true
    }),

  getResults: (partyId: string) =>
    apiRequest<VoteResults>(`/parties/${partyId}/results`),
}
```

#### User API
```typescript
// src/api/users.ts
export const userApi = {
  getProfile: () =>
    apiRequest<User>('/users/profile', { authenticated: true }),

  updateProfile: (username: string) =>
    apiRequest<User>('/users/profile', {
      method: 'PUT',
      body: JSON.stringify({ username }),
      authenticated: true
    }),
}
```

### Error Handling
```typescript
export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}
```

## TDD Approach
1. Write tests for API client with mocked fetch
2. Write tests for each API function
3. Implement API client and functions
4. Verify with `pnpm test`

## Verification
- API calls include correct headers
- Auth token injected for authenticated requests
- Error handling works
- Types are correct
