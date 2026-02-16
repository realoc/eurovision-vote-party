import { vi } from 'vitest';

vi.mock('../../src/api/client', () => ({
	apiFetch: vi.fn(),
}));

const { apiFetch } = await import('../../src/api/client');
const { submitVote, updateVote, getGuestVotes, endVoting, getResults } =
	await import('../../src/api/votes');

describe('votes API', () => {
	beforeEach(() => {
		vi.mocked(apiFetch).mockReset();
	});

	it('submitVote calls POST /api/parties/:partyId/votes without auth', async () => {
		const vote = { id: 'v-1', guestId: 'g-1' };
		const req = { guestId: 'g-1', votes: { '12': 'act-1', '10': 'act-2' } };
		vi.mocked(apiFetch).mockResolvedValue(vote);

		const result = await submitVote('party-1', req);

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/votes', {
			method: 'POST',
			body: req,
		});
		expect(result).toEqual(vote);
	});

	it('updateVote calls PUT /api/parties/:partyId/votes without auth', async () => {
		const vote = { id: 'v-1', guestId: 'g-1' };
		const req = { guestId: 'g-1', votes: { '12': 'act-3' } };
		vi.mocked(apiFetch).mockResolvedValue(vote);

		const result = await updateVote('party-1', req);

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/votes', {
			method: 'PUT',
			body: req,
		});
		expect(result).toEqual(vote);
	});

	it('getGuestVotes calls GET /api/parties/:partyId/votes/:guestId without auth', async () => {
		const vote = { id: 'v-1' };
		vi.mocked(apiFetch).mockResolvedValue(vote);

		const result = await getGuestVotes('party-1', 'g-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/votes/g-1');
		expect(result).toEqual(vote);
	});

	it('endVoting calls POST /api/parties/:partyId/end-voting with auth', async () => {
		const response = { id: 'party-1', status: 'closed' };
		vi.mocked(apiFetch).mockResolvedValue(response);

		const result = await endVoting('party-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/end-voting', {
			method: 'POST',
			authenticated: true,
		});
		expect(result).toEqual(response);
	});

	it('getResults calls GET /api/parties/:partyId/results without auth', async () => {
		const results = { partyId: 'party-1', totalVoters: 5, results: [] };
		vi.mocked(apiFetch).mockResolvedValue(results);

		const result = await getResults('party-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/results');
		expect(result).toEqual(results);
	});
});
