import { createContext, useContext } from 'react';

export type AuthUser = {
	id: string;
	email?: string | null;
	displayName?: string | null;
};

export type AuthContextValue = {
	user: AuthUser | null;
	loading: boolean;
	signInWithEmail?: (email: string, password: string) => Promise<void>;
	signInWithGoogle?: () => Promise<void>;
	signOut?: () => Promise<void>;
	getIdToken?: () => Promise<string | null>;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

export function useAuth(): AuthContextValue {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return context;
}
