# U10: Results Page

## Status
- [ ] Not started

## Dependencies
- U8 (Party Overview Page)

## Tasks
- [ ] Display acts sorted by total votes (best to worst)
- [ ] Show point totals per act
- [ ] Indicate voting is closed

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Results                         │
│     My Eurovision Party             │
│     Grand Final                     │
│                                     │
│     10 voters participated          │
│                                     │
│  🥇 1. Sweden                       │
│     Artist - Song                   │
│     120 points                      │
│                                     │
│  🥈 2. Italy                        │
│     Artist - Song                   │
│     95 points                       │
│                                     │
│  🥉 3. Ukraine                      │
│     Artist - Song                   │
│     88 points                       │
│                                     │
│  4. France                          │
│     Artist - Song                   │
│     72 points                       │
│                                     │
│  ... (remaining acts)               │
│                                     │
│  [Back to Party]                    │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/guest/ResultsPage.tsx
export function ResultsPage() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const [results, setResults] = useState<VoteResults | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadResults()
  }, [code])

  const loadResults = async () => {
    try {
      const party = await partyApi.getByCode(code!)

      if (party.status !== 'closed') {
        setError('Voting is still open. Results will be available after voting ends.')
        setLoading(false)
        return
      }

      const data = await votesApi.getResults(party.id)
      setResults(data)
    } catch (err) {
      setError('Failed to load results')
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage>{error}</ErrorMessage>

  return (
    <div>
      <h1>Results</h1>
      <h2>{results?.partyName}</h2>
      <p>{results?.totalVoters} voters participated</p>

      <ResultsList results={results?.results || []} />

      <Button onClick={() => navigate(`/party/${code}`)}>
        Back to Party
      </Button>
    </div>
  )
}
```

### ResultsList Component
```typescript
// src/components/ResultsList.tsx
interface ResultsListProps {
  results: VoteResult[]
}

export function ResultsList({ results }: ResultsListProps) {
  const getMedal = (rank: number) => {
    switch (rank) {
      case 1: return '🥇'
      case 2: return '🥈'
      case 3: return '🥉'
      default: return null
    }
  }

  return (
    <div className="space-y-4">
      {results.map(result => (
        <Card key={result.actId} className={result.rank <= 3 ? 'border-2 border-yellow-400' : ''}>
          <div className="flex items-center gap-4">
            <span className="text-2xl">{getMedal(result.rank)}</span>
            <span className="font-bold">{result.rank}.</span>
            <div className="flex-1">
              <h3 className="font-bold">{result.country}</h3>
              <p className="text-sm">{result.artist} - {result.song}</p>
            </div>
            <div className="text-right">
              <span className="text-xl font-bold">{result.totalPoints}</span>
              <span className="text-sm"> points</span>
            </div>
          </div>
        </Card>
      ))}
    </div>
  )
}
```

### Logic
- Results only shown if party status is "closed"
- Show error message if voting still open
- Display results sorted by points (server returns sorted)
- Highlight top 3 with medals
- Show 0 points for acts that received no votes

## TDD Approach
1. Write tests for results display
2. Write tests for voting-still-open state
3. Write tests for ResultsList component
4. Implement page and components
5. Verify with `pnpm test`

## Verification
- Shows results sorted by points
- Top 3 highlighted with medals
- Shows total voter count
- Handles voting-still-open gracefully
- Back button navigates to party overview
