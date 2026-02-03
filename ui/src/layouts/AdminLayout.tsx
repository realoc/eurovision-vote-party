import { NavLink, Outlet } from 'react-router-dom';

const navLinks = [
	{ to: '/admin', label: 'Dashboard' },
	{ to: '/admin/profile', label: 'Profile' },
	{ to: '/admin/party/new', label: 'Create Party' },
];

export function AdminLayout() {
	return (
		<div className="min-h-screen bg-slate-100 text-slate-900">
			<header className="border-b border-slate-200 bg-white shadow-sm">
				<div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
					<div>
						<p className="text-xs font-semibold uppercase tracking-wider text-indigo-500">
							Eurovision Vote Party
						</p>
						<h1 className="text-lg font-semibold text-slate-900">
							Admin Panel
						</h1>
					</div>
					<nav className="flex items-center gap-4 text-sm font-medium text-slate-600">
						{navLinks.map((link) => (
							<NavLink
								key={link.to}
								to={link.to}
								className={({ isActive }) =>
									[
										'rounded-full px-4 py-2 transition',
										isActive
											? 'bg-indigo-500 text-white shadow-sm'
											: 'hover:bg-indigo-50 hover:text-indigo-600',
									].join(' ')
								}
								end={link.to === '/admin'}
							>
								{link.label}
							</NavLink>
						))}
					</nav>
				</div>
			</header>
			<main className="mx-auto max-w-6xl px-6 py-10">
				<Outlet />
			</main>
		</div>
	);
}

export default AdminLayout;
