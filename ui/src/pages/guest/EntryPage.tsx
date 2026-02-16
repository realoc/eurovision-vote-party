import { type FormEvent, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ApiError } from '../../api/client';
import { joinParty } from '../../api/guests';

export function EntryPage() {
	const navigate = useNavigate();
	const [code, setCode] = useState('');
	const [username, setUsername] = useState('');
	const [error, setError] = useState('');
	const [loading, setLoading] = useState(false);

	const canSubmit = code.length === 6 && username.length >= 3 && !loading;

	function handleCodeChange(value: string) {
		setCode(
			value
				.replace(/[^a-zA-Z0-9]/g, '')
				.toUpperCase()
				.slice(0, 6),
		);
		setError('');
	}

	function handleUsernameChange(value: string) {
		setUsername(value.slice(0, 30));
		setError('');
	}

	async function handleSubmit(e: FormEvent) {
		e.preventDefault();
		setLoading(true);
		setError('');

		try {
			const guest = await joinParty(code, { username });
			localStorage.setItem(`guest_${code}`, guest.id);
			navigate(`/waiting/${code}`);
		} catch (err) {
			if (err instanceof ApiError && err.status === 404) {
				setError('Party not found. Check your code and try again.');
			} else {
				setError('Something went wrong. Please try again.');
			}
			setLoading(false);
		}
	}

	return (
		<section className="mx-auto max-w-md space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-lg shadow-indigo-500/10 backdrop-blur">
			<p className="text-sm font-semibold uppercase tracking-widest text-indigo-300">
				Eurovision Vote Party
			</p>
			<h2 className="text-3xl font-semibold text-white">Join a Party</h2>

			<form onSubmit={handleSubmit} className="space-y-4">
				<div className="text-left">
					<label
						htmlFor="party-code"
						className="mb-1 block text-sm font-medium text-indigo-100/80"
					>
						Party Code
					</label>
					<input
						id="party-code"
						type="text"
						value={code}
						onChange={(e) => handleCodeChange(e.target.value)}
						placeholder="e.g. ABC123"
						className="w-full rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-center text-lg font-mono tracking-widest text-white placeholder:text-white/30 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/30"
					/>
				</div>

				<div className="text-left">
					<label
						htmlFor="username"
						className="mb-1 block text-sm font-medium text-indigo-100/80"
					>
						Your Name
					</label>
					<input
						id="username"
						type="text"
						value={username}
						onChange={(e) => handleUsernameChange(e.target.value)}
						placeholder="Enter your name"
						className="w-full rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-white placeholder:text-white/30 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/30"
					/>
				</div>

				{error && (
					<div
						role="alert"
						className="rounded-xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-200"
					>
						{error}
					</div>
				)}

				<button
					type="submit"
					disabled={!canSubmit}
					className="w-full rounded-xl bg-indigo-500 px-6 py-3 text-sm font-semibold text-white transition hover:bg-indigo-400 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{loading ? 'Joining...' : 'Join'}
				</button>
			</form>

			<p className="text-xs text-white/40">
				Are you an admin?{' '}
				<Link
					to="/admin/login"
					className="text-indigo-300 underline hover:text-indigo-200"
				>
					Admin Login
				</Link>
			</p>
		</section>
	);
}

export default EntryPage;
