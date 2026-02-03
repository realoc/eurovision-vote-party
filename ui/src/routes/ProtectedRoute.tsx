import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { useAuth } from '../hooks/useAuth';

type ProtectedRouteProps = {
	children: ReactNode;
};

export function ProtectedRoute({ children }: ProtectedRouteProps) {
	const { user, loading } = useAuth();

	if (loading) {
		return (
			<div className="flex min-h-[40vh] items-center justify-center">
				<LoadingSpinner />
			</div>
		);
	}

	if (!user) {
		return <Navigate to="/admin/login" replace />;
	}

	return <>{children}</>;
}

export default ProtectedRoute;
