# U2: Routing & Layout

## Status
- [x] Completed

## Dependencies
- U1 (Project Setup)

## Tasks
- [x] Install React Router
- [x] Create route structure for all pages
- [x] Create base layout component with navigation
- [x] Setup protected routes for admin pages

## Details

### Installation
```bash
pnpm add react-router-dom
```

### Route Structure
```typescript
// src/routes/index.tsx
const routes = [
  // Guest routes
  { path: '/', element: <EntryPage /> },
  { path: '/waiting/:code', element: <WaitingPage /> },
  { path: '/party/:code', element: <PartyOverviewPage /> },
  { path: '/party/:code/vote', element: <VotingPage /> },
  { path: '/party/:code/results', element: <ResultsPage /> },

  // Admin routes (protected)
  { path: '/admin/login', element: <LoginPage /> },
  { path: '/admin/profile', element: <ProfileSetupPage /> },
  { path: '/admin', element: <DashboardPage /> },
  { path: '/admin/party/new', element: <CreatePartyPage /> },
  { path: '/admin/party/:id', element: <AdminPartyOverviewPage /> },
  { path: '/admin/party/:id/requests', element: <JoinRequestsPage /> },
]
```

### Protected Route Component
```typescript
// src/routes/ProtectedRoute.tsx
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()

  if (loading) return <LoadingSpinner />
  if (!user) return <Navigate to="/admin/login" />

  return <>{children}</>
}
```

### Layout Components

**GuestLayout** - Simple layout for guest pages
```typescript
function GuestLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gradient-to-b from-purple-900 to-blue-900">
      <header>Eurovision Vote Party</header>
      <main>{children}</main>
    </div>
  )
}
```

**AdminLayout** - Layout with navigation for admin pages
```typescript
function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-100">
      <nav>Admin Navigation</nav>
      <main>{children}</main>
    </div>
  )
}
```

## TDD Approach
1. Write tests for ProtectedRoute component
2. Write tests for route rendering
3. Implement routing setup
4. Verify with `pnpm test`

## Verification
- All routes render correct pages
- Protected routes redirect to login
- Layouts display correctly
- Navigation works
