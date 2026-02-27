import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { listActs } from '../../api/acts';
import { listApprovedGuests } from '../../api/guests';
import { getPartyByCode } from '../../api/parties';
import { getGuestVotes } from '../../api/votes';
import { LoadingSpinner } from '../../components/ui/LoadingSpinner';
import type { Act, Guest, PublicParty, Vote } from '../../types/api';
import { EUROVISION_POINTS } from '../../types/api';

export function PartyOverviewPage() {
	const { code } = useParams();
	const navigate = useNavigate();
	const [party, setParty] = useState<PublicParty | null>(null);
	const [guests, setGuests] = useState<Guest[]>([]);
	const [acts, setActs] = useState<Act[]>([]);
	const [myVotes, setMyVotes] = useState<Vote | null>(null);
	const [loading, setLoading] = useState(true);
	const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

	const guestId = code ? localStorage.getItem(`guest_${code}`) : null;

	const fetchData = useCallback(async () => {
		if (!code || !guestId) return;
		try {
			const partyData = await getPartyByCode(code);
			const [guestsData, actsData] = await Promise.all([
				listApprovedGuests(partyData.id),
				listActs(partyData.eventType),
			]);
			setParty(partyData);
			setGuests(guestsData.filter((g) => g.status === 'approved'));
			setActs(actsData.acts);
			try {
				const v = await getGuestVotes(partyData.id, guestId);
				setMyVotes(v);
			} catch {
				// no votes yet — leave as null
			}
		} finally {
			setLoading(false);
		}
	}, [code, guestId]);

	useEffect(() => {
		if (!code || !guestId) {
			navigate('/');
			return undefined;
		}

		fetchData();
		intervalRef.current = setInterval(fetchData, 10000);

		return () => {
			if (intervalRef.current) {
				clearInterval(intervalRef.current);
				intervalRef.current = null;
			}
		};
	}, [code, guestId, navigate, fetchData]);

	if (loading) {
		return (
			<section className="mx-auto max-w-2xl space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-lg shadow-indigo-500/10 backdrop-blur">
				<LoadingSpinner />
			</section>
		);
	}

	if (!party) return null;

	const isVotingClosed = party.status === 'closed';
	const sortedActs = [...acts].sort((a, b) => a.runningOrder - b.runningOrder);
	const actById = new Map(acts.map((a) => [a.id, a]));

	const sortedVoteEntries = myVotes
		? EUROVISION_POINTS.map((pts) => ({
				pts,
				actId: myVotes.votes[String(pts)],
			})).filter((e) => e.actId != null)
		: [];

	return (
		<section className="mx-auto max-w-2xl space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-lg shadow-indigo-500/10 backdrop-blur">
			<header className="space-y-2">
				<h2 className="text-3xl font-semibold text-white">{party.name}</h2>
				<div className="flex flex-wrap gap-2">
					<span className="rounded-full bg-indigo-500/20 px-3 py-1 text-xs font-medium uppercase tracking-wide text-indigo-300">
						{party.eventType}
					</span>
					{isVotingClosed && (
						<span className="rounded-full bg-amber-500/20 px-3 py-1 text-xs font-medium uppercase tracking-wide text-amber-300">
							Voting Closed
						</span>
					)}
				</div>
			</header>

			<div className="space-y-2">
				<h3 className="text-sm font-semibold uppercase tracking-wide text-indigo-200/60">
					Guests ({guests.length})
				</h3>
				<ul className="space-y-1">
					{guests.map((g) => (
						<li key={g.id} className="text-sm text-white/80">
							{g.username}
						</li>
					))}
				</ul>
			</div>

			{isVotingClosed ? (
				<button
					type="button"
					onClick={() => navigate(`/party/${code}/results`)}
					className="w-full rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white transition hover:bg-indigo-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900"
				>
					View Results
				</button>
			) : (
				<button
					type="button"
					onClick={() => navigate(`/party/${code}/vote`)}
					className="w-full rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white transition hover:bg-indigo-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900"
				>
					{myVotes ? 'Edit Votes' : 'Vote Now'}
				</button>
			)}

			{myVotes && sortedVoteEntries.length > 0 && (
				<div className="space-y-2">
					<h3 className="text-sm font-semibold uppercase tracking-wide text-indigo-200/60">
						Your Votes
					</h3>
					<ul className="space-y-1">
						{sortedVoteEntries.map(({ pts, actId }) => {
							const act = actById.get(actId);
							return (
								<li key={pts} className="text-sm text-white/80">
									{pts} pts — {act ? act.country : actId}
								</li>
							);
						})}
					</ul>
				</div>
			)}

			<div className="space-y-2">
				<h3 className="text-sm font-semibold uppercase tracking-wide text-indigo-200/60">
					Acts ({sortedActs.length})
				</h3>
				<ol className="space-y-1">
					{sortedActs.map((act) => (
						<li key={act.id} className="text-sm text-white/80">
							{act.runningOrder}. {act.country} — {act.artist} — {act.song}
						</li>
					))}
				</ol>
			</div>
		</section>
	);
}

export default PartyOverviewPage;
