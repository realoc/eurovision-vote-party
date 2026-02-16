/// <reference types="vitest" />

import { type Mock, vi } from 'vitest';

// Mock firebase config BEFORE importing client
vi.mock('../../src/config/firebase', () => ({
	auth: { currentUser: null },
}));

const firebaseConfig = await import('../../src/config/firebase');

// Import after mocks
const { ApiError, apiFetch } = await import('../../src/api/client');

describe('ApiError', () => {
	it('is an instance of Error', () => {
		const error = new ApiError(404, 'Not Found', 'resource not found');
		expect(error).toBeInstanceOf(Error);
		expect(error.status).toBe(404);
		expect(error.statusText).toBe('Not Found');
		expect(error.message).toBe('resource not found');
	});

	it('has the correct name', () => {
		const error = new ApiError(500, 'Internal Server Error', 'oops');
		expect(error.name).toBe('ApiError');
	});
});

describe('apiFetch', () => {
	beforeEach(() => {
		vi.stubEnv('VITE_API_URL', 'http://localhost:8080');
		vi.restoreAllMocks();
		globalThis.fetch = vi.fn();
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser =
			null;
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it('prepends VITE_API_URL to the path', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({ data: 'test' }),
		});

		await apiFetch('/api/parties');

		expect(globalThis.fetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/parties',
			expect.objectContaining({ method: 'GET' }),
		);
	});

	it('defaults to GET method', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({}),
		});

		await apiFetch('/api/test');

		expect(globalThis.fetch).toHaveBeenCalledWith(
			expect.any(String),
			expect.objectContaining({ method: 'GET' }),
		);
	});

	it('sends JSON Content-Type and stringified body when body is present', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 201,
			json: () => Promise.resolve({ id: '1' }),
		});

		await apiFetch('/api/parties', {
			method: 'POST',
			body: { name: 'Party', eventType: 'grandfinal' },
		});

		const [, options] = (globalThis.fetch as Mock).mock.calls[0];
		expect(options.headers.get('Content-Type')).toBe('application/json');
		expect(options.body).toBe(
			JSON.stringify({ name: 'Party', eventType: 'grandfinal' }),
		);
	});

	it('does not send Content-Type when no body', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve([]),
		});

		await apiFetch('/api/parties');

		const [, options] = (globalThis.fetch as Mock).mock.calls[0];
		expect(options.headers.has('Content-Type')).toBe(false);
	});

	it('returns parsed JSON for ok responses', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({ id: '1', name: 'Party' }),
		});

		const result = await apiFetch('/api/parties/1');
		expect(result).toEqual({ id: '1', name: 'Party' });
	});

	it('returns null for 204 No Content', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 204,
		});

		const result = await apiFetch('/api/parties/1', { method: 'DELETE' });
		expect(result).toBeNull();
	});

	it('throws ApiError with JSON error body for non-ok responses', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: false,
			status: 400,
			statusText: 'Bad Request',
			json: () => Promise.resolve({ error: 'invalid name' }),
			text: () => Promise.resolve(''),
		});

		await expect(apiFetch('/api/parties')).rejects.toThrow(ApiError);
		try {
			await apiFetch('/api/parties');
		} catch (e) {
			expect((e as ApiError).status).toBe(400);
			expect((e as ApiError).message).toBe('invalid name');
		}
	});

	it('throws ApiError with plain text body for non-ok responses', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: false,
			status: 404,
			statusText: 'Not Found',
			json: () => Promise.reject(new Error('not json')),
			text: () => Promise.resolve('party not found'),
		});

		try {
			await apiFetch('/api/parties/abc');
		} catch (e) {
			expect(e).toBeInstanceOf(ApiError);
			expect((e as ApiError).status).toBe(404);
			expect((e as ApiError).message).toBe('party not found');
		}
	});

	it('adds Authorization header when authenticated and user is signed in', async () => {
		const mockGetIdToken = vi.fn().mockResolvedValue('firebase-token-123');
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser = {
			getIdToken: mockGetIdToken,
		};

		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve([]),
		});

		await apiFetch('/api/parties', { authenticated: true });

		const [, options] = (globalThis.fetch as Mock).mock.calls[0];
		expect(options.headers.get('Authorization')).toBe(
			'Bearer firebase-token-123',
		);
	});

	it('omits Authorization header when authenticated is false', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({}),
		});

		await apiFetch('/api/test', { authenticated: false });

		const [, options] = (globalThis.fetch as Mock).mock.calls[0];
		expect(options.headers.has('Authorization')).toBe(false);
	});

	it('omits Authorization header when no user is signed in and authenticated is not set', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({}),
		});

		await apiFetch('/api/test');

		const [, options] = (globalThis.fetch as Mock).mock.calls[0];
		expect(options.headers.has('Authorization')).toBe(false);
	});

	it('throws ApiError(401) when authenticated is true but no user is signed in', async () => {
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser =
			null;

		await expect(
			apiFetch('/api/parties', { authenticated: true }),
		).rejects.toThrow(ApiError);

		try {
			await apiFetch('/api/parties', { authenticated: true });
		} catch (e) {
			expect((e as ApiError).status).toBe(401);
		}
	});

	it('appends query params to the URL', async () => {
		(globalThis.fetch as Mock).mockResolvedValue({
			ok: true,
			status: 200,
			json: () => Promise.resolve({}),
		});

		await apiFetch('/api/parties/code/guest-status', {
			params: { guestId: 'g-123' },
		});

		expect(globalThis.fetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/parties/code/guest-status?guestId=g-123',
			expect.any(Object),
		);
	});
});
