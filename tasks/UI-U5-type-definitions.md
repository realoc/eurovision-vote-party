# U5: Type Definitions

## Status
- [x] Complete

## Dependencies
- U1 (Project Setup)

## Tasks
- [x] Define TypeScript types matching server models
- [x] Party, Guest, Act, Vote, VoteResult types
- [x] API response types

## Details

### Type Definitions
```typescript
// src/types/index.ts

// Enums
export type PartyStatus = 'active' | 'closed'
export type GuestStatus = 'pending' | 'approved' | 'rejected'
export type EventType = 'semifinal1' | 'semifinal2' | 'grandfinal'

// Eurovision points
export const EUROVISION_POINTS = [12, 10, 8, 7, 6, 5, 4, 3, 2, 1] as const
export type EurovisionPoints = typeof EUROVISION_POINTS[number]

// Models
export interface Party {
  id: string
  name: string
  code: string
  eventType: EventType
  adminId: string
  status: PartyStatus
  createdAt: string
}

export interface Guest {
  id: string
  partyId: string
  username: string
  status: GuestStatus
  createdAt: string
}

export interface Act {
  id: string
  country: string
  artist: string
  song: string
  runningOrder: number
  eventType: EventType
}

export interface Vote {
  id: string
  guestId: string
  partyId: string
  votes: Record<number, string> // points -> actId
  createdAt: string
}

export interface User {
  id: string
  username: string
  email: string
}

export interface VoteResult {
  rank: number
  actId: string
  country: string
  artist: string
  song: string
  totalPoints: number
}

export interface VoteResults {
  partyId: string
  partyName: string
  totalVoters: number
  results: VoteResult[]
}

// Request types
export interface CreatePartyRequest {
  name: string
  eventType: EventType
}

export interface JoinPartyRequest {
  username: string
}

export interface SubmitVoteRequest {
  guestId: string
  votes: Record<number, string>
}

export interface UpdateProfileRequest {
  username: string
}

// Response types
export interface PartyListResponse {
  parties: Party[]
}

export interface GuestListResponse {
  guests: Guest[]
}

export interface ActListResponse {
  acts: Act[]
}

export interface GuestStatusResponse {
  status: GuestStatus
}

// API Error
export interface ApiErrorResponse {
  error: string
  code?: string
}
```

### Type Guards
```typescript
// src/types/guards.ts
export function isPartyActive(party: Party): boolean {
  return party.status === 'active'
}

export function isGuestApproved(guest: Guest): boolean {
  return guest.status === 'approved'
}

export function isValidEventType(value: string): value is EventType {
  return ['semifinal1', 'semifinal2', 'grandfinal'].includes(value)
}
```

### Utility Types
```typescript
// For vote form state
export type VoteFormData = Partial<Record<EurovisionPoints, string>>

// For checking if vote is complete
export function isVoteComplete(votes: VoteFormData): votes is Record<EurovisionPoints, string> {
  return EUROVISION_POINTS.every(points => votes[points] !== undefined)
}
```

## TDD Approach
1. Write tests for type guards
2. Write tests for utility functions
3. Implement types and guards
4. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Types compile without errors
- Types match server response structure
- Type guards work correctly
- Utility functions work correctly
