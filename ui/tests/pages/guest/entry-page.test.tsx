/// <reference types="vitest" />

vi.mock('../../../src/api/guests', () => ({
	joinParty: vi.fn(),
}));

vi.mock('../../../src/config/firebase', () => ({
	auth: { currentUser: null },
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
	const actual = await vi.importActual('react-router-dom');
	return { ...actual, useNavigate: () => mockNavigate };
});

import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import { ApiError } from '../../../src/api/client';
import { joinParty } from '../../../src/api/guests';
import EntryPage from '../../../src/pages/guest/EntryPage';
import type { Guest } from '../../../src/types/api';

function renderEntryPage() {
	render(
		<MemoryRouter>
			<EntryPage />
		</MemoryRouter>,
	);
	return { user: userEvent.setup() };
}

describe('EntryPage', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders a "Join a Party" heading', () => {
		renderEntryPage();

		expect(
			screen.getByRole('heading', { name: /join a party/i }),
		).toBeInTheDocument();
	});

	it('renders code and username inputs and a submit button', () => {
		renderEntryPage();

		expect(screen.getByLabelText(/party code/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/your name/i)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /join/i })).toBeInTheDocument();
	});

	it('renders an admin login link to /admin/login', () => {
		renderEntryPage();

		const link = screen.getByRole('link', { name: /admin/i });
		expect(link).toBeInTheDocument();
		expect(link).toHaveAttribute('href', '/admin/login');
	});

	describe('code input behaviour', () => {
		it('auto-uppercases input', async () => {
			const { user } = renderEntryPage();
			const codeInput = screen.getByLabelText(/party code/i);

			await user.type(codeInput, 'abcdef');

			expect(codeInput).toHaveValue('ABCDEF');
		});

		it('strips non-alphanumeric characters', async () => {
			const { user } = renderEntryPage();
			const codeInput = screen.getByLabelText(/party code/i);

			await user.type(codeInput, 'AB!@C3');

			expect(codeInput).toHaveValue('ABC3');
		});

		it('limits input to 6 characters', async () => {
			const { user } = renderEntryPage();
			const codeInput = screen.getByLabelText(/party code/i);

			await user.type(codeInput, 'ABCDEFGH');

			expect(codeInput).toHaveValue('ABCDEF');
		});
	});

	describe('button state', () => {
		it('is disabled when both fields are empty', () => {
			renderEntryPage();

			expect(screen.getByRole('button', { name: /join/i })).toBeDisabled();
		});

		it('is disabled when code is shorter than 6 characters', async () => {
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDE');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');

			expect(screen.getByRole('button', { name: /join/i })).toBeDisabled();
		});

		it('is disabled when username is shorter than 3 characters', async () => {
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Al');

			expect(screen.getByRole('button', { name: /join/i })).toBeDisabled();
		});

		it('is enabled when code is 6 chars and username >= 3 chars', async () => {
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Ali');

			expect(screen.getByRole('button', { name: /join/i })).toBeEnabled();
		});
	});

	describe('successful submission', () => {
		const guestResponse: Guest = {
			id: 'guest-123',
			partyId: 'p1',
			username: 'Alice',
			status: 'pending',
			createdAt: '',
		};

		it('calls joinParty with uppercased code and { username }', async () => {
			vi.mocked(joinParty).mockResolvedValue(guestResponse);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'abcdef');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			expect(joinParty).toHaveBeenCalledWith('ABCDEF', {
				username: 'Alice',
			});
		});

		it('stores the guest ID in localStorage as guest_{CODE}', async () => {
			vi.mocked(joinParty).mockResolvedValue(guestResponse);
			const spy = vi.spyOn(Storage.prototype, 'setItem');
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			await waitFor(() => {
				expect(spy).toHaveBeenCalledWith('guest_ABCDEF', 'guest-123');
			});

			spy.mockRestore();
		});

		it('navigates to /waiting/{CODE}', async () => {
			vi.mocked(joinParty).mockResolvedValue(guestResponse);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			await waitFor(() => {
				expect(mockNavigate).toHaveBeenCalledWith('/waiting/ABCDEF');
			});
		});

		it('shows "Joining..." on the button while loading', async () => {
			let resolveJoin!: (value: Guest) => void;
			vi.mocked(joinParty).mockReturnValue(
				new Promise((r) => {
					resolveJoin = r;
				}),
			);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			expect(
				screen.getByRole('button', { name: /joining/i }),
			).toBeInTheDocument();

			resolveJoin(guestResponse);

			await waitFor(() => {
				expect(mockNavigate).toHaveBeenCalled();
			});
		});
	});

	describe('error handling', () => {
		it('shows "Party not found." message for 404 errors', async () => {
			vi.mocked(joinParty).mockRejectedValue(
				new ApiError(404, 'Not Found', 'not found'),
			);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			const alert = await screen.findByRole('alert');
			expect(alert).toHaveTextContent(/party not found/i);
		});

		it('shows a generic error message for other errors', async () => {
			vi.mocked(joinParty).mockRejectedValue(
				new ApiError(500, 'Internal Server Error', 'server error'),
			);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			const alert = await screen.findByRole('alert');
			expect(alert).toHaveTextContent(/something went wrong/i);
		});

		it('clears the error when user types in code field', async () => {
			vi.mocked(joinParty).mockRejectedValue(
				new ApiError(404, 'Not Found', 'not found'),
			);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			await screen.findByRole('alert');

			await user.clear(screen.getByLabelText(/party code/i));
			await user.type(screen.getByLabelText(/party code/i), 'X');

			expect(screen.queryByRole('alert')).not.toBeInTheDocument();
		});

		it('clears the error when user types in username field', async () => {
			vi.mocked(joinParty).mockRejectedValue(
				new ApiError(404, 'Not Found', 'not found'),
			);
			const { user } = renderEntryPage();

			await user.type(screen.getByLabelText(/party code/i), 'ABCDEF');
			await user.type(screen.getByLabelText(/your name/i), 'Alice');
			await user.click(screen.getByRole('button', { name: /join/i }));

			await screen.findByRole('alert');

			await user.clear(screen.getByLabelText(/your name/i));
			await user.type(screen.getByLabelText(/your name/i), 'B');

			expect(screen.queryByRole('alert')).not.toBeInTheDocument();
		});
	});
});
