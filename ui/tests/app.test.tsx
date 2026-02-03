/// <reference types="vitest" />

import { render, screen } from '@testing-library/react';
import { RouterProvider, createMemoryRouter } from 'react-router-dom';
import App from '../src/App';
import { AuthContext } from '../src/hooks/useAuth';
import { routes } from '../src/routes';

function renderWithRouter(initialEntries: string[], authValue = { user: null, loading: false }) {
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
		// Rendering App ensures the RouterProvider using the browser history mounts without crashing.
		expect(() =>
			render(
				<AuthContext.Provider value={{ user: null, loading: false }}>
					<App />
				</AuthContext.Provider>,
			),
		).not.toThrow();
	});
});
