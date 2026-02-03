import type { ReactNode } from 'react';
import { Outlet } from 'react-router-dom';

type GuestLayoutProps = {
	children?: ReactNode;
};

export function GuestLayout({ children }: GuestLayoutProps) {
	return (
		<div className="min-h-screen bg-gradient-to-b from-indigo-950 via-purple-950 to-blue-950 text-slate-100">
			<header className="border-b border-white/10 bg-white/5 py-6 backdrop-blur">
				<div className="mx-auto flex max-w-5xl items-center justify-between px-6">
					<p className="text-sm font-semibold uppercase tracking-widest text-indigo-200">
						Eurovision Vote Party
					</p>
					<span className="text-xs text-indigo-200/70">
						Join the party and cast your votes
					</span>
				</div>
			</header>
			<main className="mx-auto flex max-w-5xl flex-1 items-start justify-center px-6 py-10">
				<div className="w-full">
					{children ?? <Outlet />}
				</div>
			</main>
		</div>
	);
}

export default GuestLayout;
