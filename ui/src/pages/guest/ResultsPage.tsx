import { useParams } from 'react-router-dom';

export function ResultsPage() {
	const { code } = useParams();

	return (
		<section className="space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-lg shadow-indigo-500/10 backdrop-blur">
			<header>
				<p className="text-sm uppercase tracking-wide text-indigo-200/80">
					Results
				</p>
				<h2 className="mt-2 text-3xl font-semibold text-white">
					Party {code?.toUpperCase()} Results
				</h2>
			</header>
			<p className="text-sm text-slate-200">
				The final Eurovision-style scoreboard and winner announcements will be
				shown here once the host closes the voting round.
			</p>
		</section>
	);
}

export default ResultsPage;
