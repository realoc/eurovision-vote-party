/// <reference types="vitest" />

import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { useAuth } from '../../src/hooks/useAuth';

// --- Firebase mocks ---

const mockUnsubscribe = vi.fn();
let authStateCallback: ((user: unknown) => void) | null = null;

vi.mock('../../src/config/firebase', () => ({
	auth: { currentUser: null },
}));

vi.mock('firebase/auth', () => ({
	onAuthStateChanged: vi.fn((_auth, callback) => {
		authStateCallback = callback;
		return mockUnsubscribe;
	}),
	signInWithEmailAndPassword: vi.fn(),
	signInWithPopup: vi.fn(),
	signOut: vi.fn(),
	GoogleAuthProvider: vi.fn(),
}));

// Import after mocks
const { default: AuthProvider } = await import(
	'../../src/context/AuthProvider'
);
const firebaseAuth = await import('firebase/auth');
const firebaseConfig = await import('../../src/config/firebase');

function TestConsumer() {
	const { user, loading } = useAuth();
	return (
		<div>
			<span data-testid="loading">{String(loading)}</span>
			<span data-testid="user">{user ? user.email : 'none'}</span>
		</div>
	);
}

function SignInConsumer() {
	const { signInWithEmail, signInWithGoogle, signOut: signOutFn } = useAuth();
	return (
		<div>
			<button
				type="button"
				onClick={() => signInWithEmail?.('a@b.com', 'pass')}
			>
				email
			</button>
			<button type="button" onClick={() => signInWithGoogle?.()}>
				google
			</button>
			<button type="button" onClick={() => signOutFn?.()}>
				logout
			</button>
		</div>
	);
}

function TokenConsumer() {
	const { getIdToken } = useAuth();
	return (
		<button
			type="button"
			onClick={async () => {
				const token = await getIdToken?.();
				document.title = token ?? 'no-token';
			}}
		>
			token
		</button>
	);
}

describe('AuthProvider', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		authStateCallback = null;
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser =
			null;
	});

	it('starts with loading true', () => {
		render(
			<AuthProvider>
				<TestConsumer />
			</AuthProvider>,
		);

		expect(screen.getByTestId('loading').textContent).toBe('true');
	});

	it('sets mapped AuthUser when Firebase reports a signed-in user', () => {
		render(
			<AuthProvider>
				<TestConsumer />
			</AuthProvider>,
		);

		act(() => {
			authStateCallback?.({
				uid: 'fb-1',
				email: 'test@example.com',
				displayName: 'Test User',
			});
		});

		expect(screen.getByTestId('loading').textContent).toBe('false');
		expect(screen.getByTestId('user').textContent).toBe('test@example.com');
	});

	it('sets user to null when Firebase reports no user', () => {
		render(
			<AuthProvider>
				<TestConsumer />
			</AuthProvider>,
		);

		act(() => {
			authStateCallback?.(null);
		});

		expect(screen.getByTestId('loading').textContent).toBe('false');
		expect(screen.getByTestId('user').textContent).toBe('none');
	});

	it('signInWithEmail delegates to signInWithEmailAndPassword', async () => {
		const user = userEvent.setup();

		render(
			<AuthProvider>
				<SignInConsumer />
			</AuthProvider>,
		);

		await user.click(screen.getByText('email'));

		expect(firebaseAuth.signInWithEmailAndPassword).toHaveBeenCalledWith(
			firebaseConfig.auth,
			'a@b.com',
			'pass',
		);
	});

	it('signInWithGoogle delegates to signInWithPopup with GoogleAuthProvider', async () => {
		const user = userEvent.setup();

		render(
			<AuthProvider>
				<SignInConsumer />
			</AuthProvider>,
		);

		await user.click(screen.getByText('google'));

		expect(firebaseAuth.signInWithPopup).toHaveBeenCalledWith(
			firebaseConfig.auth,
			expect.any(Object),
		);
	});

	it('signOut delegates to Firebase signOut', async () => {
		const user = userEvent.setup();

		render(
			<AuthProvider>
				<SignInConsumer />
			</AuthProvider>,
		);

		await user.click(screen.getByText('logout'));

		expect(firebaseAuth.signOut).toHaveBeenCalledWith(firebaseConfig.auth);
	});

	it('getIdToken returns token when user is signed in', async () => {
		const user = userEvent.setup();
		const mockGetIdToken = vi.fn().mockResolvedValue('mock-token');
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser = {
			getIdToken: mockGetIdToken,
		};

		render(
			<AuthProvider>
				<TokenConsumer />
			</AuthProvider>,
		);

		await user.click(screen.getByText('token'));

		await waitFor(() => {
			expect(document.title).toBe('mock-token');
		});
	});

	it('getIdToken returns null when no user is signed in', async () => {
		const user = userEvent.setup();
		(firebaseConfig as { auth: { currentUser: unknown } }).auth.currentUser =
			null;

		render(
			<AuthProvider>
				<TokenConsumer />
			</AuthProvider>,
		);

		await user.click(screen.getByText('token'));

		await waitFor(() => {
			expect(document.title).toBe('no-token');
		});
	});

	it('unsubscribes from onAuthStateChanged on unmount', () => {
		const { unmount } = render(
			<AuthProvider>
				<TestConsumer />
			</AuthProvider>,
		);

		unmount();

		expect(mockUnsubscribe).toHaveBeenCalled();
	});
});
