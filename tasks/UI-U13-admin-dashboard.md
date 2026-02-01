# U13: Admin Dashboard

## Status
- [ ] Not started

## Dependencies
- U12 (Admin Profile Setup)

## Tasks
- [ ] List of admin's parties
- [ ] Create new party button
- [ ] Delete party option
- [ ] Navigate to party overview

## Details

### Page Design
```
┌─────────────────────────────────────┐
│  Eurovision Vote Party      [Logout]│
│  Welcome, John!                     │
│                                     │
│  Your Parties                       │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ 🎉 Eurovision 2025           │   │
│  │ Grand Final · 5 guests       │   │
│  │ Code: ABC123                 │   │
│  │ Status: Active               │   │
│  │                    [🗑️]      │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ 🎉 Semi-Final Watch          │   │
│  │ Semifinal 1 · 3 guests       │   │
│  │ Code: XYZ789                 │   │
│  │ Status: Closed               │   │
│  │                    [🗑️]      │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │      + Create New Party      │   │
│  └─────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/DashboardPage.tsx
export function DashboardPage() {
  const { user, signOut } = useAuth()
  const navigate = useNavigate()
  const [profile, setProfile] = useState<User | null>(null)
  const [parties, setParties] = useState<Party[]>([])
  const [loading, setLoading] = useState(true)
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [profileData, partiesData] = await Promise.all([
        userApi.getProfile(),
        partyApi.list(),
      ])
      setProfile(profileData)
      setParties(partiesData.parties)
    } catch (err) {
      // Handle error
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (partyId: string) => {
    try {
      await partyApi.delete(partyId)
      setParties(prev => prev.filter(p => p.id !== partyId))
      setDeleteConfirm(null)
    } catch (err) {
      // Handle error
    }
  }

  const handleLogout = async () => {
    await signOut()
    navigate('/')
  }

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <header className="flex justify-between items-center">
        <h1>Eurovision Vote Party</h1>
        <Button variant="secondary" onClick={handleLogout}>Logout</Button>
      </header>

      <p>Welcome, {profile?.username}!</p>

      <h2>Your Parties</h2>

      {parties.length === 0 ? (
        <p>You haven't created any parties yet.</p>
      ) : (
        <div className="space-y-4">
          {parties.map(party => (
            <PartyCard
              key={party.id}
              party={party}
              onClick={() => navigate(`/admin/party/${party.id}`)}
              onDelete={() => setDeleteConfirm(party.id)}
            />
          ))}
        </div>
      )}

      <Button onClick={() => navigate('/admin/party/new')}>
        + Create New Party
      </Button>

      {deleteConfirm && (
        <ConfirmDialog
          title="Delete Party?"
          message="This will permanently delete the party and all associated data."
          onConfirm={() => handleDelete(deleteConfirm)}
          onCancel={() => setDeleteConfirm(null)}
        />
      )}
    </div>
  )
}
```

### PartyCard Component
```typescript
// src/components/PartyCard.tsx
interface PartyCardProps {
  party: Party
  onClick: () => void
  onDelete: () => void
}

export function PartyCard({ party, onClick, onDelete }: PartyCardProps) {
  const eventTypeLabels: Record<EventType, string> = {
    semifinal1: 'Semifinal 1',
    semifinal2: 'Semifinal 2',
    grandfinal: 'Grand Final',
  }

  return (
    <Card onClick={onClick} className="cursor-pointer hover:bg-gray-50">
      <div className="flex justify-between">
        <div>
          <h3 className="font-bold">{party.name}</h3>
          <p className="text-sm">{eventTypeLabels[party.eventType]}</p>
          <p className="text-sm font-mono">Code: {party.code}</p>
          <Badge variant={party.status === 'active' ? 'success' : 'secondary'}>
            {party.status}
          </Badge>
        </div>
        <button
          onClick={(e) => { e.stopPropagation(); onDelete(); }}
          className="text-red-500 hover:text-red-700"
        >
          🗑️
        </button>
      </div>
    </Card>
  )
}
```

## TDD Approach
1. Write tests for data loading
2. Write tests for party list display
3. Write tests for delete confirmation flow
4. Write tests for navigation
5. Implement page and components
6. Verify with `pnpm test`

## Verification
- Shows welcome message with username
- Lists all admin's parties
- Party cards are clickable
- Delete requires confirmation
- Create new party button navigates correctly
- Logout works
