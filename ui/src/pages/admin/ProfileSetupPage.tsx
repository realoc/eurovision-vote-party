export function ProfileSetupPage() {
	return (
		<section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
			<header>
				<h2 className="text-2xl font-semibold text-slate-900">Profile Setup</h2>
				<p className="text-sm text-slate-500">
					Configure your host profile and connect your Eurovision parties.
				</p>
			</header>
			<p className="text-sm text-slate-600">
				Profile configuration will be added after Firebase authentication is in
				place.
			</p>
		</section>
	);
}

export default ProfileSetupPage;
