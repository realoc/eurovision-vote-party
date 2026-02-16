import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getGuestStatus } from '../../api/guests';
import { getPartyByCode } from '../../api/parties';
import { LoadingSpinner } from '../../components/ui/LoadingSpinner';

export function WaitingPage() {
	const { code } = useParams();
	const navigate = useNavigate();
	const [partyName, setPartyName] = useState('');
	const [rejected, setRejected] = useState(false);
	const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
	const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

	const guestId = code ? localStorage.getItem(`guest_${code}`) : null;

	const cleanup = useCallback(() => {
		if (intervalRef.current) {
			clearInterval(intervalRef.current);
			intervalRef.current = null;
		}
		if (timeoutRef.current) {
			clearTimeout(timeoutRef.current);
			timeoutRef.current = null;
		}
	}, []);

	function handleCancel() {
		cleanup();
		if (code) localStorage.removeItem(`guest_${code}`);
		navigate('/');
	}

	useEffect(() => {
		if (!code || !guestId) {
			navigate('/');
			return undefined;
		}

		getPartyByCode(code)
			.then((p) => setPartyName(p.name))
			.catch(() => {});

		const currentCode = code;
		const currentGuestId = guestId;

		async function poll() {
			try {
				const guest = await getGuestStatus(currentCode, currentGuestId);
				if (guest.status === 'approved') {
					cleanup();
					navigate(`/party/${code}`);
				} else if (guest.status === 'rejected') {
					cleanup();
					setRejected(true);
					localStorage.removeItem(`guest_${code}`);
					timeoutRef.current = setTimeout(() => navigate('/'), 3000);
				}
			} catch {
				// continue polling on errors
			}
		}

		poll();
		intervalRef.current = setInterval(poll, 3000);

		return cleanup;
	}, [code, guestId, navigate, cleanup]);

	if (rejected) {
		return (
			<section className="mx-auto max-w-md space-y-4 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-lg shadow-indigo-500/10 backdrop-blur">
				<h2 className="text-3xl font-semibold text-white">Request Rejected</h2>
				<p className="text-sm text-indigo-200/80">
					The party admin has declined your request.
				</p>
				<p className="text-xs text-white/40">Redirecting...</p>
			</section>
		);
	}

	return (
		<section className="mx-auto max-w-md space-y-4 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-lg shadow-indigo-500/10 backdrop-blur">
			<LoadingSpinner />
			<h2 className="text-3xl font-semibold text-white">
				Waiting for approval...
			</h2>
			{partyName && (
				<p className="text-sm text-indigo-200/80">
					You've requested to join '{partyName}'
				</p>
			)}
			<p className="text-sm text-indigo-200/80">
				The party admin will review your request shortly.
			</p>
			<button
				type="button"
				onClick={handleCancel}
				className="w-full rounded-xl border border-white/10 bg-white/5 px-6 py-3 text-sm font-semibold text-white transition hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-900"
			>
				Cancel
			</button>
		</section>
	);
}

export default WaitingPage;
