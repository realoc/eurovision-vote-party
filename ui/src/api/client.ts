import { auth } from '../config/firebase';

export class ApiError extends Error {
	status: number;
	statusText: string;

	constructor(status: number, statusText: string, message: string) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.statusText = statusText;
	}
}

type ApiFetchOptions = {
	method?: string;
	body?: unknown;
	authenticated?: boolean;
	params?: Record<string, string>;
};

export async function apiFetch<T>(
	path: string,
	options: ApiFetchOptions = {},
): Promise<T> {
	const { method = 'GET', body, authenticated = false, params } = options;

	const baseUrl = import.meta.env.VITE_API_URL ?? '';
	let url = `${baseUrl}${path}`;

	if (params) {
		const searchParams = new URLSearchParams(params);
		url += `?${searchParams.toString()}`;
	}

	const headers = new Headers();

	if (body !== undefined) {
		headers.set('Content-Type', 'application/json');
	}

	if (authenticated) {
		const currentUser = auth.currentUser;
		if (!currentUser) {
			throw new ApiError(401, 'Unauthorized', 'Not authenticated');
		}
		const token = await currentUser.getIdToken();
		headers.set('Authorization', `Bearer ${token}`);
	}

	const response = await fetch(url, {
		method,
		headers,
		body: body !== undefined ? JSON.stringify(body) : undefined,
	});

	if (!response.ok) {
		let message: string;
		try {
			const errorBody = await response.json();
			message = errorBody.error ?? response.statusText;
		} catch {
			message = await response.text();
		}
		throw new ApiError(response.status, response.statusText, message);
	}

	if (response.status === 204) {
		return null as T;
	}

	return response.json() as Promise<T>;
}
