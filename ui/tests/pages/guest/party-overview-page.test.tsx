/// <reference types="vitest" />

vi.mock('../../../src/api/parties', () => ({
	getPartyByCode: vi.fn(),
}));

vi.mock('../../../src/api/acts', () => ({
	listActs: vi.fn(),
}));

vi.mock('../../../src/api/guests', () => ({
	listApprovedGuests: vi.fn(),
}));

vi.mock('../../../src/api/votes', () => ({
	getGuestVotes: vi.fn(),
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
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import { listActs } from '../../../src/api/acts';
import { listApprovedGuests } from '../../../src/api/guests';
import { getPartyByCode } from '../../../src/api/parties';
import { getGuestVotes } from '../../../src/api/votes';
import PartyOverviewPage from '../../../src/pages/guest/PartyOverviewPage';
import type {
	Act,
	ActsResponse,
	Guest,
	PublicParty,
	Vote,
} from '../../../src/types/api';

const activeParty: PublicParty = {
	id: 'p1',
	name: 'Eurovision Finals',
	code: 'ABCDEF',
	eventType: 'grandfinal',
	status: 'active',
};

const closedParty: PublicParty = {
	...activeParty,
	status: 'closed',
};

const guests: Guest[] = [
	{
		id: 'g1',
		partyId: 'p1',
		username: 'Alice',
		status: 'approved',
		createdAt: '',
	},
	{
		id: 'g2',
		partyId: 'p1',
		username: 'Bob',
		status: 'approved',
		createdAt: '',
	},
];

const acts: Act[] = [
	{
		id: 'act-1',
		country: 'Sweden',
		artist: 'Loreen',
		song: 'Tattoo',
		runningOrder: 1,
		eventType: 'grandfinal',
	},
	{
		id: 'act-2',
		country: 'Finland',
		artist: 'Käärijä',
		song: 'Cha Cha Cha',
		runningOrder: 2,
		eventType: 'grandfinal',
	},
];

const actsResponse: ActsResponse = { acts };

const myVotes: Vote = {
	id: 'v1',
	guestId: 'guest-123',
	partyId: 'p1',
	votes: { '12': 'act-1', '10': 'act-2' },
	createdAt: '',
};

function renderPage() {
	const result = render(
		<MemoryRouter>
			<PartyOverviewPage />
		</MemoryRouter>,
	);
	return { unmount: result.unmount };
}

describe('PartyOverviewPage', () => {
	let getItemSpy: ReturnType<typeof vi.spyOn>;

	beforeEach(() => {
		vi.useFakeTimers({ shouldAdvanceTime: true });
		vi.clearAllMocks();
		getItemSpy = vi
			.spyOn(Storage.prototype, 'getItem')
			.mockImplementation((key: string) =>
				key === 'guest_ABCDEF' ? 'guest-123' : null,
			);
		vi.mocked(getPartyByCode).mockResolvedValue(activeParty);
		vi.mocked(listApprovedGuests).mockResolvedValue(guests);
		vi.mocked(listActs).mockResolvedValue(actsResponse);
		vi.mocked(getGuestVotes).mockRejectedValue(new Error('No votes'));
	});

	afterEach(() => {
		vi.useRealTimers();
		getItemSpy.mockRestore();
	});

	it('redirects to / if no guestId in localStorage', async () => {
		getItemSpy.mockReturnValue(null);

		renderPage();

		await waitFor(() => {
			expect(mockNavigate).toHaveBeenCalledWith('/');
		});
	});

	it('shows loading spinner initially', () => {
		renderPage();
		expect(screen.getByLabelText('Loading')).toBeInTheDocument();
	});

	it('shows party name and event type badge', async () => {
		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('heading', { name: /Eurovision Finals/i }),
			).toBeInTheDocument();
		});

		expect(screen.getByText('grandfinal')).toBeInTheDocument();
	});

	it('shows approved guest list with usernames', async () => {
		renderPage();

		await waitFor(() => {
			expect(screen.getByText('Alice')).toBeInTheDocument();
			expect(screen.getByText('Bob')).toBeInTheDocument();
		});
	});

	it('shows "Vote Now" button when voting active and no votes yet', async () => {
		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('button', { name: /vote now/i }),
			).toBeInTheDocument();
		});
	});

	it('shows "Edit Votes" button when voting active and votes exist', async () => {
		vi.mocked(getGuestVotes).mockResolvedValue(myVotes);

		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('button', { name: /edit votes/i }),
			).toBeInTheDocument();
		});
	});

	it('shows own votes section when votes submitted', async () => {
		vi.mocked(getGuestVotes).mockResolvedValue(myVotes);

		renderPage();

		await waitFor(() => {
			// Heading appears
			expect(screen.getByText(/^Your Votes$/)).toBeInTheDocument();
			// Vote entries appear (anchored at start to avoid matching acts section)
			expect(screen.getByText(/^12 pts/)).toBeInTheDocument();
			expect(screen.getByText(/^10 pts/)).toBeInTheDocument();
		});
		// Country names are shown in vote entries
		expect(screen.getByText(/^12 pts/)).toHaveTextContent('Sweden');
		expect(screen.getByText(/^10 pts/)).toHaveTextContent('Finland');
	});

	it('shows acts in running order', async () => {
		renderPage();

		await waitFor(() => {
			expect(screen.getByText(/Loreen/)).toBeInTheDocument();
			expect(screen.getByText(/Tattoo/)).toBeInTheDocument();
			expect(screen.getByText(/Käärijä/)).toBeInTheDocument();
			expect(screen.getByText(/Cha Cha Cha/)).toBeInTheDocument();
		});
	});

	it('shows "Voting Closed" badge when party status is closed', async () => {
		vi.mocked(getPartyByCode).mockResolvedValue(closedParty);

		renderPage();

		await waitFor(() => {
			expect(screen.getByText(/voting closed/i)).toBeInTheDocument();
		});
	});

	it('shows "View Results" button when voting closed', async () => {
		vi.mocked(getPartyByCode).mockResolvedValue(closedParty);

		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('button', { name: /view results/i }),
			).toBeInTheDocument();
		});
	});

	it('navigates to vote page when "Vote Now" clicked', async () => {
		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('button', { name: /vote now/i }),
			).toBeInTheDocument();
		});

		screen.getByRole('button', { name: /vote now/i }).click();

		expect(mockNavigate).toHaveBeenCalledWith('/party/ABCDEF/vote');
	});

	it('navigates to results page when "View Results" clicked', async () => {
		vi.mocked(getPartyByCode).mockResolvedValue(closedParty);

		renderPage();

		await waitFor(() => {
			expect(
				screen.getByRole('button', { name: /view results/i }),
			).toBeInTheDocument();
		});

		screen.getByRole('button', { name: /view results/i }).click();

		expect(mockNavigate).toHaveBeenCalledWith('/party/ABCDEF/results');
	});

	it('polls every 10 seconds for updates', async () => {
		renderPage();

		await waitFor(() => {
			expect(getPartyByCode).toHaveBeenCalledTimes(1);
		});

		await act(async () => {
			vi.advanceTimersByTime(10000);
		});

		await waitFor(() => {
			expect(getPartyByCode).toHaveBeenCalledTimes(2);
		});

		await act(async () => {
			vi.advanceTimersByTime(10000);
		});

		await waitFor(() => {
			expect(getPartyByCode).toHaveBeenCalledTimes(3);
		});
	});

	it('stops polling on unmount', async () => {
		const { unmount } = renderPage();

		await waitFor(() => {
			expect(getPartyByCode).toHaveBeenCalledTimes(1);
		});

		unmount();

		vi.advanceTimersByTime(10000);

		expect(getPartyByCode).toHaveBeenCalledTimes(1);
	});
});
