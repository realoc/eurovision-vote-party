import { vi } from 'vitest';

vi.mock('../../src/api/client', () => ({
	apiFetch: vi.fn(),
}));

const { apiFetch } = await import('../../src/api/client');
const { getProfile, updateProfile } = await import('../../src/api/users');

describe('users API', () => {
	beforeEach(() => {
		vi.mocked(apiFetch).mockReset();
	});

	it('getProfile calls GET /api/users/profile with auth', async () => {
		const user = { id: 'u-1', username: 'alice', email: 'a@b.com' };
		vi.mocked(apiFetch).mockResolvedValue(user);

		const result = await getProfile();

		expect(apiFetch).toHaveBeenCalledWith('/api/users/profile', {
			authenticated: true,
		});
		expect(result).toEqual(user);
	});

	it('updateProfile calls PUT /api/users/profile with auth', async () => {
		const user = { id: 'u-1', username: 'bob', email: 'a@b.com' };
		vi.mocked(apiFetch).mockResolvedValue(user);

		const result = await updateProfile({ username: 'bob' });

		expect(apiFetch).toHaveBeenCalledWith('/api/users/profile', {
			method: 'PUT',
			body: { username: 'bob' },
			authenticated: true,
		});
		expect(result).toEqual(user);
	});
});
