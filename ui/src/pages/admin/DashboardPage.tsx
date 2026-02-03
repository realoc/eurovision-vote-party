export function DashboardPage() {
	return (
		<section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
			<header>
				<h2 className="text-2xl font-semibold text-slate-900">
					Admin Dashboard
				</h2>
				<p className="text-sm text-slate-500">
					Manage parties, review join requests, and monitor the live voting
					status.
				</p>
			</header>
			<p className="text-sm text-slate-600">
				Dashboard widgets and insights will be added once data models are in
				place.
			</p>
		</section>
	);
}

export default DashboardPage;
