export function LoginPage() {
	return (
		<section className="mx-auto max-w-md space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-lg shadow-indigo-500/10 backdrop-blur">
			<h2 className="text-3xl font-semibold text-white">Admin Login</h2>
			<p className="text-sm text-indigo-100/80">
				Sign in to manage your Eurovision parties and keep track of guest votes.
			</p>
			<div className="rounded-2xl border border-white/10 bg-white/10 px-6 py-4 text-sm text-white/80">
				Authentication will be implemented in the next milestone.
			</div>
		</section>
	);
}

export default LoginPage;
