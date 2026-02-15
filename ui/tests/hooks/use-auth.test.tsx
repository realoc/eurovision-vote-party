/// <reference types="vitest" />

import { renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { AuthContext, useAuth } from '../../src/hooks/useAuth';
import { createMockAuthContext } from '../helpers/auth-test-utils';

describe('useAuth', () => {
	it('throws when used outside an AuthProvider', () => {
		expect(() => renderHook(() => useAuth())).toThrow(
			'useAuth must be used within an AuthProvider',
		);
	});

	it('returns the context value when inside a provider', () => {
		const mockValue = createMockAuthContext({
			user: { id: 'u1', email: 'a@b.com', displayName: 'Alice' },
		});

		const wrapper = ({ children }: { children: ReactNode }) => (
			<AuthContext.Provider value={mockValue}>{children}</AuthContext.Provider>
		);

		const { result } = renderHook(() => useAuth(), { wrapper });

		expect(result.current.user).toEqual({
			id: 'u1',
			email: 'a@b.com',
			displayName: 'Alice',
		});
		expect(result.current.loading).toBe(false);
		expect(result.current.signInWithEmail).toBeDefined();
		expect(result.current.signInWithGoogle).toBeDefined();
		expect(result.current.signOut).toBeDefined();
		expect(result.current.getIdToken).toBeDefined();
	});
});
