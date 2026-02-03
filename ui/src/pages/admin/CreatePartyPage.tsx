export function CreatePartyPage() {
	return (
		<section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
			<header>
				<h2 className="text-2xl font-semibold text-slate-900">
					Create a New Party
				</h2>
				<p className="text-sm text-slate-500">
					Set up a new Eurovision watch party and invite your guests.
				</p>
			</header>
			<p className="text-sm text-slate-600">
				The creation form will be connected to the backend once data endpoints
				are available.
			</p>
		</section>
	);
}

export default CreatePartyPage;
