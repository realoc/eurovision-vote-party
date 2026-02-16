// --- Enums as union types ---

export type EventType = 'semifinal1' | 'semifinal2' | 'grandfinal';

export type PartyStatus = 'active' | 'closed';

export type GuestStatus = 'pending' | 'approved' | 'rejected';

// --- Domain Models ---

export type Party = {
	id: string;
	name: string;
	code: string;
	eventType: EventType;
	adminId: string;
	status: PartyStatus;
	createdAt: string;
};

// Public party info (returned by getPartyByCode - same fields but conceptually public)
export type PublicParty = {
	id: string;
	name: string;
	code: string;
	eventType: EventType;
	status: PartyStatus;
};

export type Guest = {
	id: string;
	partyId: string;
	username: string;
	status: GuestStatus;
	createdAt: string;
};

export type Act = {
	id: string;
	country: string;
	artist: string;
	song: string;
	runningOrder: number;
	eventType: EventType;
};

export type Vote = {
	id: string;
	guestId: string;
	partyId: string;
	votes: Record<string, string>; // Go marshals map[int]string with string keys
	createdAt: string;
};

export type VoteResult = {
	actId: string;
	country: string;
	artist: string;
	song: string;
	totalPoints: number;
	rank: number;
};

export type PartyResults = {
	partyId: string;
	partyName: string;
	totalVoters: number;
	results: VoteResult[];
};

export type User = {
	id: string;
	username: string;
	email: string;
};

// --- Request Types ---

export type CreatePartyRequest = {
	name: string;
	eventType: EventType;
};

export type JoinPartyRequest = {
	username: string;
};

export type SubmitVoteRequest = {
	guestId: string;
	votes: Record<string, string>;
};

export type UpdateProfileRequest = {
	username: string;
};

// --- Response Types ---

export type ActsResponse = {
	acts: Act[];
};

export type EndVotingResponse = {
	id: string;
	status: PartyStatus;
};

export type StatusOkResponse = {
	status: string;
};

// --- Eurovision Scoring ---
export const EUROVISION_POINTS = [12, 10, 8, 7, 6, 5, 4, 3, 2, 1] as const;
export type EurovisionPoints = (typeof EUROVISION_POINTS)[number];

// --- Error Response ---
export type ApiErrorResponse = {
	error: string;
};
