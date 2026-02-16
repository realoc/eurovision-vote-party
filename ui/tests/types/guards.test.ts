/// <reference types="vitest" />

import type { Guest, Party } from '../../src/types/api';
import { EUROVISION_POINTS } from '../../src/types/api';
import {
	isGuestApproved,
	isPartyActive,
	isValidEventType,
	isVoteComplete,
} from '../../src/types/guards';

describe('EUROVISION_POINTS', () => {
	it('has the exact expected values', () => {
		expect(EUROVISION_POINTS).toEqual([12, 10, 8, 7, 6, 5, 4, 3, 2, 1]);
	});

	it('has length 10', () => {
		expect(EUROVISION_POINTS).toHaveLength(10);
	});
});

describe('isPartyActive', () => {
	it('returns true for a party with status active', () => {
		const party: Party = {
			id: 'p1',
			name: 'Test Party',
			code: 'ABC123',
			eventType: 'grandfinal',
			adminId: 'admin1',
			status: 'active',
			createdAt: '2026-01-01T00:00:00Z',
		};
		expect(isPartyActive(party)).toBe(true);
	});

	it('returns false for a party with status closed', () => {
		const party: Party = {
			id: 'p1',
			name: 'Test Party',
			code: 'ABC123',
			eventType: 'grandfinal',
			adminId: 'admin1',
			status: 'closed',
			createdAt: '2026-01-01T00:00:00Z',
		};
		expect(isPartyActive(party)).toBe(false);
	});
});

describe('isGuestApproved', () => {
	it('returns true for a guest with status approved', () => {
		const guest: Guest = {
			id: 'g1',
			partyId: 'p1',
			username: 'alice',
			status: 'approved',
			createdAt: '2026-01-01T00:00:00Z',
		};
		expect(isGuestApproved(guest)).toBe(true);
	});

	it('returns false for a guest with status pending', () => {
		const guest: Guest = {
			id: 'g1',
			partyId: 'p1',
			username: 'bob',
			status: 'pending',
			createdAt: '2026-01-01T00:00:00Z',
		};
		expect(isGuestApproved(guest)).toBe(false);
	});

	it('returns false for a guest with status rejected', () => {
		const guest: Guest = {
			id: 'g1',
			partyId: 'p1',
			username: 'charlie',
			status: 'rejected',
			createdAt: '2026-01-01T00:00:00Z',
		};
		expect(isGuestApproved(guest)).toBe(false);
	});
});

describe('isValidEventType', () => {
	it('returns true for semifinal1', () => {
		expect(isValidEventType('semifinal1')).toBe(true);
	});

	it('returns true for semifinal2', () => {
		expect(isValidEventType('semifinal2')).toBe(true);
	});

	it('returns true for grandfinal', () => {
		expect(isValidEventType('grandfinal')).toBe(true);
	});

	it('returns false for final', () => {
		expect(isValidEventType('final')).toBe(false);
	});

	it('returns false for empty string', () => {
		expect(isValidEventType('')).toBe(false);
	});
});

describe('isVoteComplete', () => {
	it('returns true when all 10 points are assigned', () => {
		const votes = {
			12: 'act-a',
			10: 'act-b',
			8: 'act-c',
			7: 'act-d',
			6: 'act-e',
			5: 'act-f',
			4: 'act-g',
			3: 'act-h',
			2: 'act-i',
			1: 'act-j',
		};
		expect(isVoteComplete(votes)).toBe(true);
	});

	it('returns false when only some points are assigned', () => {
		const votes = {
			12: 'act-a',
			10: 'act-b',
		};
		expect(isVoteComplete(votes)).toBe(false);
	});

	it('returns false for an empty object', () => {
		expect(isVoteComplete({})).toBe(false);
	});
});
