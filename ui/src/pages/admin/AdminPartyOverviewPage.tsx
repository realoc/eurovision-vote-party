import { useParams } from 'react-router-dom';

export function AdminPartyOverviewPage() {
	const { id } = useParams();

	return (
		<section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
			<header>
				<p className="text-xs uppercase tracking-wide text-indigo-500">
					Party Overview
				</p>
				<h2 className="text-2xl font-semibold text-slate-900">
					Party {id?.toUpperCase()}
				</h2>
			</header>
			<p className="text-sm text-slate-600">
				This page will display party details, guest list, and current status
				once connected to the backend.
			</p>
		</section>
	);
}

export default AdminPartyOverviewPage;
