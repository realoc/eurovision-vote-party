import { useParams } from 'react-router-dom';

export function VotingPage() {
	const { code } = useParams();

	return (
		<section className="space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-lg shadow-indigo-500/10 backdrop-blur">
			<header>
				<p className="text-sm uppercase tracking-wide text-indigo-200/80">
					Voting
				</p>
				<h2 className="mt-2 text-3xl font-semibold text-white">
					Cast your votes for party {code?.toUpperCase()}
				</h2>
			</header>
			<p className="text-sm text-slate-200">
				The interactive voting interface will appear here. Guests can rank their
				favourite acts while keeping track of the scoreboard in real time.
			</p>
		</section>
	);
}

export default VotingPage;
