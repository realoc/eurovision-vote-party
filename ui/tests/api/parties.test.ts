import { vi } from 'vitest';

vi.mock('../../src/api/client', () => ({
	apiFetch: vi.fn(),
}));

const { apiFetch } = await import('../../src/api/client');
const { createParty, listParties, getPartyByCode, getPartyById, deleteParty } =
	await import('../../src/api/parties');

describe('parties API', () => {
	beforeEach(() => {
		vi.mocked(apiFetch).mockReset();
	});

	it('createParty calls POST /api/parties with auth', async () => {
		const party = { id: '1', name: 'Test' };
		vi.mocked(apiFetch).mockResolvedValue(party);

		const result = await createParty({ name: 'Test', eventType: 'grandfinal' });

		expect(apiFetch).toHaveBeenCalledWith('/api/parties', {
			method: 'POST',
			body: { name: 'Test', eventType: 'grandfinal' },
			authenticated: true,
		});
		expect(result).toEqual(party);
	});

	it('listParties calls GET /api/parties with auth', async () => {
		const parties = [{ id: '1' }];
		vi.mocked(apiFetch).mockResolvedValue(parties);

		const result = await listParties();

		expect(apiFetch).toHaveBeenCalledWith('/api/parties', {
			authenticated: true,
		});
		expect(result).toEqual(parties);
	});

	it('getPartyByCode calls GET /api/parties/:code without auth', async () => {
		const party = { id: '1', code: 'ABC123' };
		vi.mocked(apiFetch).mockResolvedValue(party);

		const result = await getPartyByCode('ABC123');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/ABC123');
		expect(result).toEqual(party);
	});

	it('getPartyById calls GET /api/parties/:id with auth', async () => {
		const party = { id: 'uuid-1' };
		vi.mocked(apiFetch).mockResolvedValue(party);

		const result = await getPartyById('uuid-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/uuid-1', {
			authenticated: true,
		});
		expect(result).toEqual(party);
	});

	it('deleteParty calls DELETE /api/parties/:id with auth', async () => {
		vi.mocked(apiFetch).mockResolvedValue(null);

		await deleteParty('uuid-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/uuid-1', {
			method: 'DELETE',
			authenticated: true,
		});
	});
});
