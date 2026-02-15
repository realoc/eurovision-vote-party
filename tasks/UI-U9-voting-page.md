# U9: Voting Page

## Status
- [ ] Not started

## Dependencies
- U8 (Party Overview Page)

## Tasks
- [ ] Display Eurovision points: 12, 10, 8, 7, 6, 5, 4, 3, 2, 1
- [ ] Dropdown/picker for each point with acts in running order
- [ ] Remove selected acts from other pickers
- [ ] Submit button (red/disabled → green/enabled when complete)
- [ ] Pre-fill with existing votes when editing
- [ ] Submit and redirect to overview

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Cast Your Votes                 │
│                                     │
│  12 points  [▼ Sweden - Loreen    ] │
│  10 points  [▼ Select country...  ] │
│   8 points  [▼ Select country...  ] │
│   7 points  [▼ Select country...  ] │
│   6 points  [▼ Select country...  ] │
│   5 points  [▼ Select country...  ] │
│   4 points  [▼ Select country...  ] │
│   3 points  [▼ Select country...  ] │
│   2 points  [▼ Select country...  ] │
│   1 point   [▼ Select country...  ] │
│                                     │
│  ┌─────────────────────────────┐   │
│  │     Submit Votes (8/10)      │   │ (red, disabled)
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │       Submit Votes           │   │ (green, enabled)
│  └─────────────────────────────┘   │
│                                     │
│  [Cancel]                           │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/guest/VotingPage.tsx
export function VotingPage() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const [acts, setActs] = useState<Act[]>([])
  const [votes, setVotes] = useState<VoteFormData>({})
  const [loading, setLoading] = useState(false)
  const [existingVote, setExistingVote] = useState<Vote | null>(null)

  const guestId = localStorage.getItem(`guest_${code}`)

  useEffect(() => {
    loadData()
  }, [code])

  const loadData = async () => {
    const party = await partyApi.getByCode(code!)
    const actsData = await actsApi.list(party.eventType)
    setActs(actsData.acts)

    // Load existing votes if editing
    try {
      const existing = await votesApi.get(party.id, guestId!)
      setExistingVote(existing)
      setVotes(existing.votes)
    } catch {
      // No existing votes
    }
  }

  const handleVoteChange = (points: EurovisionPoints, actId: string) => {
    setVotes(prev => ({ ...prev, [points]: actId }))
  }

  const getAvailableActs = (currentPoints: EurovisionPoints) => {
    const selectedActIds = Object.entries(votes)
      .filter(([p]) => Number(p) !== currentPoints)
      .map(([, actId]) => actId)

    return acts.filter(act => !selectedActIds.includes(act.id))
  }

  const isComplete = isVoteComplete(votes)
  const completedCount = Object.keys(votes).length

  const handleSubmit = async () => {
    if (!isComplete) return

    setLoading(true)
    try {
      const party = await partyApi.getByCode(code!)
      if (existingVote) {
        await votesApi.update(party.id, guestId!, votes)
      } else {
        await votesApi.submit(party.id, guestId!, votes)
      }
      navigate(`/party/${code}`)
    } catch (err) {
      // Handle error
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h1>Cast Your Votes</h1>

      {EUROVISION_POINTS.map(points => (
        <VoteSelect
          key={points}
          points={points}
          acts={getAvailableActs(points)}
          selectedActId={votes[points]}
          onChange={(actId) => handleVoteChange(points, actId)}
        />
      ))}

      <Button
        onClick={handleSubmit}
        disabled={!isComplete || loading}
        variant={isComplete ? 'success' : 'danger'}
      >
        {loading ? 'Submitting...' : `Submit Votes (${completedCount}/10)`}
      </Button>

      <Button variant="secondary" onClick={() => navigate(`/party/${code}`)}>
        Cancel
      </Button>
    </div>
  )
}
```

### VoteSelect Component
```typescript
// src/components/VoteSelect.tsx
interface VoteSelectProps {
  points: EurovisionPoints
  acts: Act[]
  selectedActId?: string
  onChange: (actId: string) => void
}

export function VoteSelect({ points, acts, selectedActId, onChange }: VoteSelectProps) {
  return (
    <div className="flex items-center gap-4">
      <span className="w-20 text-right font-bold">
        {points} {points === 1 ? 'point' : 'points'}
      </span>
      <Select
        value={selectedActId}
        onChange={onChange}
        placeholder="Select country..."
      >
        {acts.map(act => (
          <option key={act.id} value={act.id}>
            {act.country} - {act.artist}
          </option>
        ))}
      </Select>
    </div>
  )
}
```

### Logic
- Filter out already-selected acts from each dropdown
- Track completion status (all 10 points assigned)
- Show completion counter in button
- Button color: red/disabled until complete, then green/enabled

## TDD Approach
1. Write tests for act filtering logic
2. Write tests for completion detection
3. Write tests for pre-fill with existing votes
4. Write tests for submission flow
5. Implement page and components
6. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- All 10 point values shown
- Selected acts removed from other dropdowns
- Button disabled until all selections made
- Button changes from red to green when complete
- Pre-fills with existing votes when editing
- Successfully submits and navigates
