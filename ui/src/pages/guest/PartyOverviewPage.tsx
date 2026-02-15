import { useParams } from 'react-router-dom';

export function PartyOverviewPage() {
	const { code } = useParams();

	return (
		<section className="space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-lg shadow-indigo-500/10 backdrop-blur">
			<header>
				<p className="text-sm uppercase tracking-wide text-indigo-200/80">
					Party Code
				</p>
				<h2 className="mt-2 text-3xl font-semibold text-white">
					Eurovision Party {code?.toUpperCase()}
				</h2>
			</header>
			<p className="text-sm text-slate-200">
				This page will show the current acts, running order, and live scores
				once the host starts the party.
			</p>
		</section>
	);
}

export default PartyOverviewPage;
