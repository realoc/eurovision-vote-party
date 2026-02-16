import type { CreatePartyRequest, Party, PublicParty } from '../types/api';
import { apiFetch } from './client';

export function createParty(req: CreatePartyRequest): Promise<Party> {
	return apiFetch('/api/parties', {
		method: 'POST',
		body: req,
		authenticated: true,
	});
}

export function listParties(): Promise<Party[]> {
	return apiFetch('/api/parties', {
		authenticated: true,
	});
}

export function getPartyByCode(code: string): Promise<PublicParty> {
	return apiFetch(`/api/parties/${code}`);
}

export function getPartyById(id: string): Promise<Party> {
	return apiFetch(`/api/parties/${id}`, {
		authenticated: true,
	});
}

export function deleteParty(id: string): Promise<void> {
	return apiFetch(`/api/parties/${id}`, {
		method: 'DELETE',
		authenticated: true,
	});
}
