import { vi } from 'vitest';

vi.mock('../../src/api/client', () => ({
	apiFetch: vi.fn(),
}));

const { apiFetch } = await import('../../src/api/client');
const {
	joinParty,
	getGuestStatus,
	listGuests,
	listJoinRequests,
	approveGuest,
	rejectGuest,
	removeGuest,
} = await import('../../src/api/guests');

describe('guests API', () => {
	beforeEach(() => {
		vi.mocked(apiFetch).mockReset();
	});

	it('joinParty calls POST /api/parties/:code/join without auth', async () => {
		const guest = { id: 'g-1', status: 'pending' };
		vi.mocked(apiFetch).mockResolvedValue(guest);

		const result = await joinParty('ABC123', { username: 'alice' });

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/ABC123/join', {
			method: 'POST',
			body: { username: 'alice' },
		});
		expect(result).toEqual(guest);
	});

	it('getGuestStatus calls GET /api/parties/:code/guest-status with guestId param', async () => {
		const guest = { id: 'g-1', status: 'approved' };
		vi.mocked(apiFetch).mockResolvedValue(guest);

		const result = await getGuestStatus('ABC123', 'g-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/ABC123/guest-status', {
			params: { guestId: 'g-1' },
		});
		expect(result).toEqual(guest);
	});

	it('listGuests calls GET /api/parties/:id/guests with auth', async () => {
		const guests = [{ id: 'g-1' }];
		vi.mocked(apiFetch).mockResolvedValue(guests);

		const result = await listGuests('party-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/guests', {
			authenticated: true,
		});
		expect(result).toEqual(guests);
	});

	it('listJoinRequests calls GET /api/parties/:id/join-requests with auth', async () => {
		const guests = [{ id: 'g-1', status: 'pending' }];
		vi.mocked(apiFetch).mockResolvedValue(guests);

		const result = await listJoinRequests('party-1');

		expect(apiFetch).toHaveBeenCalledWith(
			'/api/parties/party-1/join-requests',
			{
				authenticated: true,
			},
		);
		expect(result).toEqual(guests);
	});

	it('approveGuest calls PUT /api/parties/:id/guests/:guestId/approve with auth', async () => {
		vi.mocked(apiFetch).mockResolvedValue({ status: 'ok' });

		const result = await approveGuest('party-1', 'g-1');

		expect(apiFetch).toHaveBeenCalledWith(
			'/api/parties/party-1/guests/g-1/approve',
			{
				method: 'PUT',
				authenticated: true,
			},
		);
		expect(result).toEqual({ status: 'ok' });
	});

	it('rejectGuest calls PUT /api/parties/:id/guests/:guestId/reject with auth', async () => {
		vi.mocked(apiFetch).mockResolvedValue({ status: 'ok' });

		const result = await rejectGuest('party-1', 'g-1');

		expect(apiFetch).toHaveBeenCalledWith(
			'/api/parties/party-1/guests/g-1/reject',
			{
				method: 'PUT',
				authenticated: true,
			},
		);
		expect(result).toEqual({ status: 'ok' });
	});

	it('removeGuest calls DELETE /api/parties/:id/guests/:guestId with auth', async () => {
		vi.mocked(apiFetch).mockResolvedValue(null);

		await removeGuest('party-1', 'g-1');

		expect(apiFetch).toHaveBeenCalledWith('/api/parties/party-1/guests/g-1', {
			method: 'DELETE',
			authenticated: true,
		});
	});
});
