export function EntryPage() {
	return (
		<section className="w-full rounded-3xl border border-white/10 bg-white/5 p-10 text-center shadow-2xl shadow-indigo-500/10 backdrop-blur">
			<p className="text-sm font-semibold uppercase tracking-widest text-indigo-300">
				Eurovision Vote Party
			</p>
			<h2 className="mt-4 text-4xl font-bold tracking-tight text-white">
				Frontend scaffolding ready to go
			</h2>
			<p className="mt-4 text-base leading-7 text-slate-200">
				Tailwind CSS, Biome, and Vitest are wired up so you can focus on
				building the guest and admin experiences.
			</p>
			<div className="mt-8 flex flex-wrap items-center justify-center gap-4">
				<a
					className="inline-flex items-center rounded-full bg-indigo-500 px-5 py-2 text-sm font-semibold text-white transition hover:bg-indigo-400 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900"
					href="https://tailwindcss.com/docs"
					target="_blank"
					rel="noreferrer"
				>
					Tailwind Docs
				</a>
				<a
					className="inline-flex items-center rounded-full border border-white/20 px-5 py-2 text-sm font-semibold text-slate-200 transition hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900"
					href="https://vitest.dev/guide/"
					target="_blank"
					rel="noreferrer"
				>
					Vitest Guide
				</a>
			</div>
		</section>
	);
}

export default EntryPage;
