import type { Guest, JoinPartyRequest, StatusOkResponse } from '../types/api';
import { apiFetch } from './client';

export function joinParty(code: string, req: JoinPartyRequest): Promise<Guest> {
	return apiFetch(`/api/parties/${code}/join`, {
		method: 'POST',
		body: req,
	});
}

export function getGuestStatus(code: string, guestId: string): Promise<Guest> {
	return apiFetch(`/api/parties/${code}/guest-status`, {
		params: { guestId },
	});
}

export function listGuests(partyId: string): Promise<Guest[]> {
	return apiFetch(`/api/parties/${partyId}/guests`, {
		authenticated: true,
	});
}

export function listApprovedGuests(partyId: string): Promise<Guest[]> {
	return apiFetch(`/api/parties/${partyId}/guests`);
}

export function listJoinRequests(partyId: string): Promise<Guest[]> {
	return apiFetch(`/api/parties/${partyId}/join-requests`, {
		authenticated: true,
	});
}

export function approveGuest(
	partyId: string,
	guestId: string,
): Promise<StatusOkResponse> {
	return apiFetch(`/api/parties/${partyId}/guests/${guestId}/approve`, {
		method: 'PUT',
		authenticated: true,
	});
}

export function rejectGuest(
	partyId: string,
	guestId: string,
): Promise<StatusOkResponse> {
	return apiFetch(`/api/parties/${partyId}/guests/${guestId}/reject`, {
		method: 'PUT',
		authenticated: true,
	});
}

export function removeGuest(partyId: string, guestId: string): Promise<void> {
	return apiFetch(`/api/parties/${partyId}/guests/${guestId}`, {
		method: 'DELETE',
		authenticated: true,
	});
}
