import { vi } from 'vitest';
import type { AuthContextValue } from '../../src/hooks/useAuth';

export function createMockAuthContext(
	overrides?: Partial<AuthContextValue>,
): AuthContextValue {
	return {
		user: null,
		loading: false,
		signInWithEmail: vi.fn(),
		signInWithGoogle: vi.fn(),
		signOut: vi.fn(),
		getIdToken: vi.fn().mockResolvedValue(null),
		...overrides,
	};
}
