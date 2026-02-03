import { createBrowserRouter, type RouteObject } from 'react-router-dom';
import AdminLayout from '../layouts/AdminLayout';
import GuestLayout from '../layouts/GuestLayout';
import AdminPartyOverviewPage from '../pages/admin/AdminPartyOverviewPage';
import CreatePartyPage from '../pages/admin/CreatePartyPage';
import DashboardPage from '../pages/admin/DashboardPage';
import JoinRequestsPage from '../pages/admin/JoinRequestsPage';
import LoginPage from '../pages/admin/LoginPage';
import ProfileSetupPage from '../pages/admin/ProfileSetupPage';
import EntryPage from '../pages/guest/EntryPage';
import PartyOverviewPage from '../pages/guest/PartyOverviewPage';
import ResultsPage from '../pages/guest/ResultsPage';
import WaitingPage from '../pages/guest/WaitingPage';
import VotingPage from '../pages/guest/VotingPage';
import ProtectedRoute from './ProtectedRoute';

export const routes: RouteObject[] = [
	{
		path: '/',
		element: <GuestLayout />,
		children: [
			{
				index: true,
				element: <EntryPage />,
			},
			{
				path: 'waiting/:code',
				element: <WaitingPage />,
			},
			{
				path: 'party/:code',
				element: <PartyOverviewPage />,
			},
			{
				path: 'party/:code/vote',
				element: <VotingPage />,
			},
			{
				path: 'party/:code/results',
				element: <ResultsPage />,
			},
			{
				path: 'admin/login',
				element: <LoginPage />,
			},
		],
	},
	{
		path: '/admin',
		element: (
			<ProtectedRoute>
				<AdminLayout />
			</ProtectedRoute>
		),
		children: [
			{
				index: true,
				element: <DashboardPage />,
			},
			{
				path: 'profile',
				element: <ProfileSetupPage />,
			},
			{
				path: 'party/new',
				element: <CreatePartyPage />,
			},
			{
				path: 'party/:id',
				element: <AdminPartyOverviewPage />,
			},
			{
				path: 'party/:id/requests',
				element: <JoinRequestsPage />,
			},
		],
	},
];

export const router = createBrowserRouter(routes);
