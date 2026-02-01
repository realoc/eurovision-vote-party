# U15: Party Overview Page (Admin View)

## Status
- [ ] Not started

## Dependencies
- U14 (Party Creation Page)

## Tasks
- [ ] List of guests with remove option
- [ ] List of acts in running order
- [ ] Join requests icon with count badge
- [ ] Vote/Edit Votes button (same as guest)
- [ ] End Voting button
- [ ] Party code display with copy button

## Details

### Page Design
```
┌─────────────────────────────────────┐
│  [← Back]     My Eurovision Party   │
│                                     │
│  Party Code: ABC123  [📋]           │
│  Grand Final · Active               │
│                                     │
│  ┌─────────────────────────────┐   │
│  │  Join Requests  🔴 3         │   │
│  └─────────────────────────────┘   │
│                                     │
│  Guests (5)                         │
│  ├ John ✓ (voted)         [Remove] │
│  ├ Jane                   [Remove] │
│  ├ Bob ✓ (voted)          [Remove] │
│  └ Alice                  [Remove] │
│                                     │
│  ┌────────────┐  ┌──────────────┐  │
│  │ Vote/Edit  │  │  End Voting  │  │
│  └────────────┘  └──────────────┘  │
│                                     │
│  Acts (26)                          │
│  1. Sweden - Artist - Song          │
│  2. Italy - Artist - Song           │
│  ...                                │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/PartyOverviewPage.tsx
export function AdminPartyOverviewPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [party, setParty] = useState<Party | null>(null)
  const [guests, setGuests] = useState<Guest[]>([])
  const [pendingCount, setPendingCount] = useState(0)
  const [acts, setActs] = useState<Act[]>([])
  const [myVotes, setMyVotes] = useState<Vote | null>(null)
  const [loading, setLoading] = useState(true)
  const [endVotingConfirm, setEndVotingConfirm] = useState(false)

  useEffect(() => {
    loadData()
    // Poll for updates
    const interval = setInterval(loadData, 10000)
    return () => clearInterval(interval)
  }, [id])

  const loadData = async () => {
    try {
      const partyData = await partyApi.getById(id!)
      setParty(partyData)

      const [guestsData, actsData, joinRequests] = await Promise.all([
        guestApi.list(id!),
        actsApi.list(partyData.eventType),
        guestApi.getJoinRequests(id!),
      ])

      setGuests(guestsData.guests.filter(g => g.status === 'approved'))
      setActs(actsData.acts)
      setPendingCount(joinRequests.requests.length)

      // Get admin's own votes
      try {
        const profile = await userApi.getProfile()
        const votes = await votesApi.get(id!, profile.id)
        setMyVotes(votes)
      } catch {
        // No votes yet
      }
    } catch (err) {
      // Handle error
    } finally {
      setLoading(false)
    }
  }

  const handleRemoveGuest = async (guestId: string) => {
    try {
      await guestApi.remove(id!, guestId)
      setGuests(prev => prev.filter(g => g.id !== guestId))
    } catch (err) {
      // Handle error
    }
  }

  const handleEndVoting = async () => {
    try {
      await votesApi.endVoting(id!)
      setParty(prev => prev ? { ...prev, status: 'closed' } : null)
      setEndVotingConfirm(false)
    } catch (err) {
      // Handle error
    }
  }

  const handleVote = () => {
    navigate(`/party/${party?.code}/vote`)
  }

  if (loading) return <LoadingSpinner />

  const isVotingClosed = party?.status === 'closed'

  return (
    <div>
      <header className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => navigate('/admin')}>
          ← Back
        </Button>
        <h1>{party?.name}</h1>
      </header>

      <div className="flex items-center gap-2">
        <span>Party Code: <span className="font-mono">{party?.code}</span></span>
        <CopyButton text={party?.code || ''} />
      </div>

      <Badge>{party?.eventType}</Badge>
      <Badge variant={isVotingClosed ? 'secondary' : 'success'}>
        {party?.status}
      </Badge>

      <Button onClick={() => navigate(`/admin/party/${id}/requests`)}>
        Join Requests
        {pendingCount > 0 && <Badge variant="danger">{pendingCount}</Badge>}
      </Button>

      <h2>Guests ({guests.length})</h2>
      <GuestList
        guests={guests}
        showRemove
        onRemove={handleRemoveGuest}
      />

      <div className="flex gap-4">
        {!isVotingClosed && (
          <>
            <Button onClick={handleVote}>
              {myVotes ? 'Edit Votes' : 'Vote'}
            </Button>
            <Button variant="danger" onClick={() => setEndVotingConfirm(true)}>
              End Voting
            </Button>
          </>
        )}
        {isVotingClosed && (
          <Button onClick={() => navigate(`/party/${party?.code}/results`)}>
            View Results
          </Button>
        )}
      </div>

      <h2>Acts ({acts.length})</h2>
      <ActList acts={acts} />

      {endVotingConfirm && (
        <ConfirmDialog
          title="End Voting?"
          message="This will close voting and calculate final results. This cannot be undone."
          onConfirm={handleEndVoting}
          onCancel={() => setEndVotingConfirm(false)}
        />
      )}
    </div>
  )
}
```

## TDD Approach
1. Write tests for data loading
2. Write tests for guest removal
3. Write tests for end voting flow
4. Write tests for join requests badge
5. Implement page component
6. Verify with `pnpm test`

## Verification
- Shows party details with code
- Copy button works
- Join requests badge shows count
- Guest list with remove buttons
- Admin can vote (uses party code)
- End voting requires confirmation
- View results button after voting closed
