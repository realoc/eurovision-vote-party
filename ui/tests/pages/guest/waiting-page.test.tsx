/// <reference types="vitest" />

vi.mock('../../../src/api/guests', () => ({
	getGuestStatus: vi.fn(),
}));

vi.mock('../../../src/api/parties', () => ({
	getPartyByCode: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
	const actual = await vi.importActual('react-router-dom');
	return {
		...actual,
		useNavigate: () => mockNavigate,
		useParams: () => ({ code: 'ABCDEF' }),
	};
});

import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import { getGuestStatus } from '../../../src/api/guests';
import { getPartyByCode } from '../../../src/api/parties';
import WaitingPage from '../../../src/pages/guest/WaitingPage';
import type { Guest, PublicParty } from '../../../src/types/api';

const pendingGuest: Guest = {
	id: 'guest-123',
	partyId: 'p1',
	username: 'Alice',
	status: 'pending',
	createdAt: '',
};

const approvedGuest: Guest = {
	...pendingGuest,
	status: 'approved',
};

const rejectedGuest: Guest = {
	...pendingGuest,
	status: 'rejected',
};

const party: PublicParty = {
	id: 'p1',
	name: 'Eurovision Finals',
	code: 'ABCDEF',
	eventType: 'grandfinal',
	status: 'active',
};

function renderWaitingPage() {
	const result = render(
		<MemoryRouter>
			<WaitingPage />
		</MemoryRouter>,
	);
	return { unmount: result.unmount };
}

describe('WaitingPage', () => {
	let getItemSpy: ReturnType<typeof vi.spyOn>;
	let removeItemSpy: ReturnType<typeof vi.spyOn>;

	beforeEach(() => {
		vi.useFakeTimers({ shouldAdvanceTime: true });
		vi.clearAllMocks();
		getItemSpy = vi
			.spyOn(Storage.prototype, 'getItem')
			.mockImplementation((key: string) =>
				key === 'guest_ABCDEF' ? 'guest-123' : null,
			);
		removeItemSpy = vi.spyOn(Storage.prototype, 'removeItem');
		vi.mocked(getGuestStatus).mockResolvedValue(pendingGuest);
		vi.mocked(getPartyByCode).mockResolvedValue(party);
	});

	afterEach(() => {
		vi.useRealTimers();
		getItemSpy.mockRestore();
		removeItemSpy.mockRestore();
	});

	it('redirects to / if no guest ID in localStorage', async () => {
		getItemSpy.mockReturnValue(null);

		renderWaitingPage();

		await waitFor(() => {
			expect(mockNavigate).toHaveBeenCalledWith('/');
		});
	});

	it('renders waiting message and spinner', async () => {
		renderWaitingPage();

		expect(
			screen.getByRole('heading', { name: /waiting for approval/i }),
		).toBeInTheDocument();

		expect(screen.getByLabelText('Loading')).toBeInTheDocument();
	});

	it('fetches party name on mount and displays it', async () => {
		renderWaitingPage();

		await waitFor(() => {
			expect(getPartyByCode).toHaveBeenCalledWith('ABCDEF');
		});

		await waitFor(() => {
			expect(screen.getByText(/you've requested to join/i)).toBeInTheDocument();
			expect(screen.getByText(/Eurovision Finals/)).toBeInTheDocument();
		});
	});

	it('polls getGuestStatus every 3 seconds', async () => {
		renderWaitingPage();

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(1);
		});

		await act(async () => {
			vi.advanceTimersByTime(3000);
		});

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(2);
		});

		await act(async () => {
			vi.advanceTimersByTime(3000);
		});

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(3);
		});
	});

	it('navigates to /party/ABCDEF on approved status', async () => {
		vi.mocked(getGuestStatus).mockResolvedValue(approvedGuest);

		renderWaitingPage();

		await waitFor(() => {
			expect(mockNavigate).toHaveBeenCalledWith('/party/ABCDEF');
		});
	});

	it('shows rejection message when rejected', async () => {
		vi.mocked(getGuestStatus).mockResolvedValue(rejectedGuest);

		renderWaitingPage();

		await waitFor(() => {
			expect(
				screen.getByRole('heading', { name: /request rejected/i }),
			).toBeInTheDocument();
		});
	});

	it('clears localStorage and redirects home after rejection', async () => {
		vi.mocked(getGuestStatus).mockResolvedValue(rejectedGuest);

		renderWaitingPage();

		await waitFor(() => {
			expect(
				screen.getByRole('heading', { name: /request rejected/i }),
			).toBeInTheDocument();
		});

		expect(removeItemSpy).toHaveBeenCalledWith('guest_ABCDEF');

		await vi.advanceTimersByTimeAsync(3000);

		await waitFor(() => {
			expect(mockNavigate).toHaveBeenCalledWith('/');
		});
	});

	it('cancel button clears localStorage and navigates home', async () => {
		const user = userEvent.setup({
			advanceTimers: vi.advanceTimersByTime,
		});

		renderWaitingPage();

		const cancelButton = screen.getByRole('button', { name: /cancel/i });
		await user.click(cancelButton);

		expect(removeItemSpy).toHaveBeenCalledWith('guest_ABCDEF');
		expect(mockNavigate).toHaveBeenCalledWith('/');
	});

	it('stops polling on unmount', async () => {
		const { unmount } = renderWaitingPage();

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(1);
		});

		unmount();

		vi.advanceTimersByTime(3000);

		expect(getGuestStatus).toHaveBeenCalledTimes(1);
	});

	it('continues polling on network errors', async () => {
		vi.mocked(getGuestStatus)
			.mockRejectedValueOnce(new Error('Network error'))
			.mockResolvedValue(pendingGuest);

		renderWaitingPage();

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(1);
		});

		await act(async () => {
			vi.advanceTimersByTime(3000);
		});

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(2);
		});

		await act(async () => {
			vi.advanceTimersByTime(3000);
		});

		await waitFor(() => {
			expect(getGuestStatus).toHaveBeenCalledTimes(3);
		});
	});
});
