import type { EurovisionPoints, EventType, Guest, Party } from './api';
import { EUROVISION_POINTS } from './api';

export function isPartyActive(party: Party): boolean {
	return party.status === 'active';
}

export function isGuestApproved(guest: Guest): boolean {
	return guest.status === 'approved';
}

export function isValidEventType(value: string): value is EventType {
	return ['semifinal1', 'semifinal2', 'grandfinal'].includes(value);
}

export type VoteFormData = Partial<Record<EurovisionPoints, string>>;

export function isVoteComplete(
	votes: VoteFormData,
): votes is Record<EurovisionPoints, string> {
	return EUROVISION_POINTS.every((points) => votes[points] !== undefined);
}
