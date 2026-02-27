# U8: Party Overview Page (Guest View)

## Status
- [x] Done

## Dependencies
- U7 (Guest Waiting Page)

## Tasks
- [x] List of approved guests
- [x] List of acts with running order
- [x] Display guest's own votes (if submitted)
- [x] "Vote" / "Edit Votes" button
- [x] Handle voting closed state

## Details

### Page Design (Voting Open)
```
┌─────────────────────────────────────┐
│     My Eurovision Party             │
│     Grand Final                     │
│                                     │
│  Guests (5)                         │
│  ├ John ✓ (voted)                   │
│  ├ Jane                             │
│  ├ Bob ✓ (voted)                    │
│  └ Alice                            │
│                                     │
│  ┌─────────────────────────────┐   │
│  │        Vote Now              │   │
│  └─────────────────────────────┘   │
│                                     │
│  Your Votes:                        │
│  12 pts - Sweden                    │
│  10 pts - Italy                     │
│  ...                                │
│                                     │
│  Acts (26)                          │
│  1. Sweden - Artist - Song          │
│  2. Italy - Artist - Song           │
│  ...                                │
└─────────────────────────────────────┘
```

### Page Design (Voting Closed)
```
┌─────────────────────────────────────┐
│     My Eurovision Party             │
│     Grand Final - VOTING CLOSED     │
│                                     │
│  ┌─────────────────────────────┐   │
│  │       View Results           │   │
│  └─────────────────────────────┘   │
│                                     │
│  ... (same as above)                │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/guest/PartyOverviewPage.tsx
export function PartyOverviewPage() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const [party, setParty] = useState<Party | null>(null)
  const [guests, setGuests] = useState<Guest[]>([])
  const [acts, setActs] = useState<Act[]>([])
  const [myVotes, setMyVotes] = useState<Vote | null>(null)

  const guestId = localStorage.getItem(`guest_${code}`)

  useEffect(() => {
    fetchData()
    // Poll for updates every 10 seconds
    const interval = setInterval(fetchData, 10000)
    return () => clearInterval(interval)
  }, [code])

  const fetchData = async () => {
    const [partyData, guestsData, actsData] = await Promise.all([
      partyApi.getByCode(code!),
      guestApi.list(code!),
      actsApi.list(partyData?.eventType),
    ])
    setParty(partyData)
    setGuests(guestsData.guests.filter(g => g.status === 'approved'))
    setActs(actsData.acts)

    // Fetch own votes if submitted
    try {
      const votes = await votesApi.get(partyData.id, guestId!)
      setMyVotes(votes)
    } catch {
      // No votes yet
    }
  }

  const handleVote = () => {
    navigate(`/party/${code}/vote`)
  }

  const handleViewResults = () => {
    navigate(`/party/${code}/results`)
  }

  const isVotingClosed = party?.status === 'closed'

  return (
    <div>
      <h1>{party?.name}</h1>
      <Badge>{party?.eventType}</Badge>
      {isVotingClosed && <Badge variant="warning">Voting Closed</Badge>}

      <GuestList guests={guests} votes={/* map of guestId to hasVoted */} />

      {isVotingClosed ? (
        <Button onClick={handleViewResults}>View Results</Button>
      ) : (
        <Button onClick={handleVote}>
          {myVotes ? 'Edit Votes' : 'Vote Now'}
        </Button>
      )}

      {myVotes && <MyVotesDisplay votes={myVotes} acts={acts} />}

      <ActList acts={acts} />
    </div>
  )
}
```

### Sub-components
- `GuestList` - Display guests with vote status
- `ActList` - Display acts in running order
- `MyVotesDisplay` - Show guest's submitted votes

### Data Fetching
- Fetch party, guests, acts on mount
- Poll for updates (voting status, guest list)
- Fetch own votes if guestId exists

## TDD Approach
1. Write tests for data fetching
2. Write tests for vote button states
3. Write tests for voting closed state
4. Implement page and sub-components
5. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Shows party name and event type
- Lists approved guests
- Shows vote status per guest
- Vote button shows correct text
- Handles voting closed state
- Shows own votes if submitted
