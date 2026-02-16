import type { ActsResponse, EventType } from '../types/api';
import { apiFetch } from './client';

export function listActs(event: EventType): Promise<ActsResponse> {
	return apiFetch('/api/acts', {
		params: { event },
	});
}
