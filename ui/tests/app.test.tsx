/// <reference types="vitest" />

import { render, screen } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import { vi } from 'vitest';
import { AuthContext } from '../src/hooks/useAuth';
import { routes } from '../src/routes';

vi.mock('../src/config/firebase', () => ({
	auth: { currentUser: null },
}));

vi.mock('firebase/auth', () => ({
	onAuthStateChanged: vi.fn((_auth, callback) => {
		callback(null);
		return vi.fn();
	}),
	signInWithEmailAndPassword: vi.fn(),
	signInWithPopup: vi.fn(),
	signOut: vi.fn(),
	GoogleAuthProvider: vi.fn(),
}));

const { default: App } = await import('../src/App');

function renderWithRouter(
	initialEntries: string[],
	authValue = { user: null, loading: false },
) {
	const router = createMemoryRouter(routes, { initialEntries });

	return render(
		<AuthContext.Provider value={authValue}>
			<RouterProvider router={router} />
		</AuthContext.Provider>,
	);
}

describe('App', () => {
	it('renders the entry page content on the root route', () => {
		renderWithRouter(['/']);

		expect(
			screen.getByRole('heading', {
				name: /frontend scaffolding ready to go/i,
			}),
		).toBeInTheDocument();
	});

	it('redirects to the login page for protected routes when unauthenticated', async () => {
		renderWithRouter(['/admin']);

		expect(
			await screen.findByRole('heading', { name: /admin login/i }),
		).toBeInTheDocument();
	});

	it('uses the browser router in the App component', () => {
		expect(() => render(<App />)).not.toThrow();
	});
});
