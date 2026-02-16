import { vi } from 'vitest';

vi.mock('../../src/api/client', () => ({
	apiFetch: vi.fn(),
}));

const { apiFetch } = await import('../../src/api/client');
const { listActs } = await import('../../src/api/acts');

describe('acts API', () => {
	beforeEach(() => {
		vi.mocked(apiFetch).mockReset();
	});

	it('listActs calls GET /api/acts with event param', async () => {
		const response = { acts: [{ id: '1', country: 'Sweden' }] };
		vi.mocked(apiFetch).mockResolvedValue(response);

		const result = await listActs('grandfinal');

		expect(apiFetch).toHaveBeenCalledWith('/api/acts', {
			params: { event: 'grandfinal' },
		});
		expect(result).toEqual(response);
	});
});
