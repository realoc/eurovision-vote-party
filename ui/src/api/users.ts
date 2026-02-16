import type { UpdateProfileRequest, User } from '../types/api';
import { apiFetch } from './client';

export function getProfile(): Promise<User> {
	return apiFetch('/api/users/profile', {
		authenticated: true,
	});
}

export function updateProfile(req: UpdateProfileRequest): Promise<User> {
	return apiFetch('/api/users/profile', {
		method: 'PUT',
		body: req,
		authenticated: true,
	});
}
