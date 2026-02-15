import {
	signOut as firebaseSignOut,
	GoogleAuthProvider,
	onAuthStateChanged,
	signInWithEmailAndPassword,
	signInWithPopup,
} from 'firebase/auth';
import {
	type ReactNode,
	useCallback,
	useEffect,
	useMemo,
	useState,
} from 'react';
import { auth } from '../config/firebase';
import { AuthContext, type AuthUser } from '../hooks/useAuth';

export default function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<AuthUser | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		const unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
			if (firebaseUser) {
				setUser({
					id: firebaseUser.uid,
					email: firebaseUser.email,
					displayName: firebaseUser.displayName,
				});
			} else {
				setUser(null);
			}
			setLoading(false);
		});
		return unsubscribe;
	}, []);

	const signInWithEmail = useCallback(
		async (email: string, password: string) => {
			await signInWithEmailAndPassword(auth, email, password);
		},
		[],
	);

	const signInWithGoogle = useCallback(async () => {
		const provider = new GoogleAuthProvider();
		await signInWithPopup(auth, provider);
	}, []);

	const signOut = useCallback(async () => {
		await firebaseSignOut(auth);
	}, []);

	const getIdToken = useCallback(async (): Promise<string | null> => {
		return (await auth.currentUser?.getIdToken()) ?? null;
	}, []);

	const value = useMemo(
		() => ({
			user,
			loading,
			signInWithEmail,
			signInWithGoogle,
			signOut,
			getIdToken,
		}),
		[user, loading, signInWithEmail, signInWithGoogle, signOut, getIdToken],
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
