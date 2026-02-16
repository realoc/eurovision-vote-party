import type {
	EndVotingResponse,
	PartyResults,
	SubmitVoteRequest,
	Vote,
} from '../types/api';
import { apiFetch } from './client';

export function submitVote(
	partyId: string,
	req: SubmitVoteRequest,
): Promise<Vote> {
	return apiFetch(`/api/parties/${partyId}/votes`, {
		method: 'POST',
		body: req,
	});
}

export function updateVote(
	partyId: string,
	req: SubmitVoteRequest,
): Promise<Vote> {
	return apiFetch(`/api/parties/${partyId}/votes`, {
		method: 'PUT',
		body: req,
	});
}

export function getGuestVotes(partyId: string, guestId: string): Promise<Vote> {
	return apiFetch(`/api/parties/${partyId}/votes/${guestId}`);
}

export function endVoting(partyId: string): Promise<EndVotingResponse> {
	return apiFetch(`/api/parties/${partyId}/end-voting`, {
		method: 'POST',
		authenticated: true,
	});
}

export function getResults(partyId: string): Promise<PartyResults> {
	return apiFetch(`/api/parties/${partyId}/results`);
}
