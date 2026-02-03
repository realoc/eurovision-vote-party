import { createContext, useContext } from 'react';

export type AuthUser = {
	id: string;
	email?: string | null;
	displayName?: string | null;
};

export type AuthContextValue = {
	user: AuthUser | null;
	loading: boolean;
};

const defaultAuthContext: AuthContextValue = {
	user: null,
	loading: false,
};

export const AuthContext = createContext<AuthContextValue>(defaultAuthContext);

export function useAuth(): AuthContextValue {
	return useContext(AuthContext);
}
