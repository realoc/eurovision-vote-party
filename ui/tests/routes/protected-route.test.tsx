/// <reference types="vitest" />

import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import ProtectedRoute from '../../src/routes/ProtectedRoute';
import { AuthContext } from '../../src/hooks/useAuth';

describe('ProtectedRoute', () => {
	it('shows a loading spinner while authentication state is loading', () => {
		render(
			<AuthContext.Provider value={{ user: null, loading: true }}>
				<MemoryRouter initialEntries={['/admin']}>
					<ProtectedRoute>
						<div>Private Area</div>
					</ProtectedRoute>
				</MemoryRouter>
			</AuthContext.Provider>,
		);

		expect(screen.getByLabelText(/loading/i)).toBeInTheDocument();
	});

	it('redirects to the login route when there is no authenticated user', () => {
		render(
			<AuthContext.Provider value={{ user: null, loading: false }}>
				<MemoryRouter initialEntries={['/admin']}>
					<Routes>
						<Route
							path="/admin"
							element={
								<ProtectedRoute>
									<div>Private Area</div>
								</ProtectedRoute>
							}
						/>
						<Route path="/admin/login" element={<div>Login Page</div>} />
					</Routes>
				</MemoryRouter>
			</AuthContext.Provider>,
		);

		expect(screen.getByText(/login page/i)).toBeInTheDocument();
	});

	it('renders the protected content when an authenticated user exists', () => {
		render(
			<AuthContext.Provider
				value={{
					user: { id: 'admin-1', email: 'admin@example.com' },
					loading: false,
				}}
			>
				<MemoryRouter initialEntries={['/admin']}>
					<Routes>
						<Route
							path="/admin"
							element={
								<ProtectedRoute>
									<div>Private Area</div>
								</ProtectedRoute>
							}
						/>
						<Route path="/admin/login" element={<div>Login Page</div>} />
					</Routes>
				</MemoryRouter>
			</AuthContext.Provider>,
		);

		expect(screen.getByText(/private area/i)).toBeInTheDocument();
	});
});
