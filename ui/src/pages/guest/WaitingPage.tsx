import { useParams } from 'react-router-dom';

export function WaitingPage() {
	const { code } = useParams();

	return (
		<section className="space-y-4 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-lg shadow-indigo-500/10 backdrop-blur">
			<h2 className="text-3xl font-semibold text-white">Waiting Room</h2>
			<p className="text-sm text-indigo-200/80">
				Share the party code with your friends and wait for the host to start the
				show.
			</p>
			<div className="rounded-2xl border border-white/10 bg-white/10 px-6 py-4 text-2xl font-semibold tracking-widest text-white">
				{code?.toUpperCase() ?? '----'}
			</div>
		</section>
	);
}

export default WaitingPage;
