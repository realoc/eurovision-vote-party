# U14: Party Creation Page

## Status
- [ ] Not started

## Dependencies
- U13 (Admin Dashboard)

## Tasks
- [ ] Party name input
- [ ] Event type selector (Semifinal 1, Semifinal 2, Grand Final)
- [ ] Create button
- [ ] Display party code with copy-to-clipboard button
- [ ] Navigate to party overview button

## Details

### Page Design (Form)
```
┌─────────────────────────────────────┐
│  Create New Party                   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Party Name                   │   │
│  │ [Eurovision Watch Party  ]   │   │
│  └─────────────────────────────┘   │
│                                     │
│  Event Type                         │
│  ○ Semifinal 1                      │
│  ○ Semifinal 2                      │
│  ● Grand Final                      │
│                                     │
│  ┌─────────────────────────────┐   │
│  │       Create Party           │   │
│  └─────────────────────────────┘   │
│                                     │
│  [Cancel]                           │
└─────────────────────────────────────┘
```

### Page Design (Success)
```
┌─────────────────────────────────────┐
│  Party Created! 🎉                  │
│                                     │
│  Share this code with your guests:  │
│                                     │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │       ABC123                │   │
│  │                             │   │
│  │                    [📋 Copy]│   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │    Go to Party Overview      │   │
│  └─────────────────────────────┘   │
│                                     │
│  [Create Another Party]             │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/CreatePartyPage.tsx
export function CreatePartyPage() {
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [eventType, setEventType] = useState<EventType>('grandfinal')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [createdParty, setCreatedParty] = useState<Party | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()

    if (!name.trim()) {
      setError('Party name is required')
      return
    }

    setError(null)
    setLoading(true)

    try {
      const party = await partyApi.create({ name: name.trim(), eventType })
      setCreatedParty(party)
    } catch (err) {
      setError('Failed to create party. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleReset = () => {
    setCreatedParty(null)
    setName('')
    setEventType('grandfinal')
  }

  if (createdParty) {
    return (
      <div>
        <h1>Party Created! 🎉</h1>
        <p>Share this code with your guests:</p>

        <Card className="text-center">
          <span className="text-4xl font-mono font-bold">{createdParty.code}</span>
          <CopyButton text={createdParty.code} />
        </Card>

        <Button onClick={() => navigate(`/admin/party/${createdParty.id}`)}>
          Go to Party Overview
        </Button>

        <Button variant="secondary" onClick={handleReset}>
          Create Another Party
        </Button>
      </div>
    )
  }

  return (
    <div>
      <h1>Create New Party</h1>

      <form onSubmit={handleSubmit}>
        <Input
          label="Party Name"
          value={name}
          onChange={setName}
          placeholder="Eurovision Watch Party"
          maxLength={100}
        />

        <fieldset>
          <legend>Event Type</legend>
          {[
            { value: 'semifinal1', label: 'Semifinal 1' },
            { value: 'semifinal2', label: 'Semifinal 2' },
            { value: 'grandfinal', label: 'Grand Final' },
          ].map(option => (
            <label key={option.value} className="flex items-center gap-2">
              <input
                type="radio"
                name="eventType"
                value={option.value}
                checked={eventType === option.value}
                onChange={() => setEventType(option.value as EventType)}
              />
              {option.label}
            </label>
          ))}
        </fieldset>

        {error && <ErrorMessage>{error}</ErrorMessage>}

        <Button type="submit" disabled={loading || !name.trim()}>
          {loading ? 'Creating...' : 'Create Party'}
        </Button>

        <Button variant="secondary" onClick={() => navigate('/admin')}>
          Cancel
        </Button>
      </form>
    </div>
  )
}
```

### CopyButton Component
```typescript
// src/components/ui/CopyButton.tsx
interface CopyButtonProps {
  text: string
}

export function CopyButton({ text }: CopyButtonProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <button onClick={handleCopy} className="...">
      {copied ? '✓ Copied!' : '📋 Copy'}
    </button>
  )
}
```

## TDD Approach
1. Write tests for form validation
2. Write tests for submission flow
3. Write tests for success state display
4. Write tests for CopyButton component
5. Implement page and components
6. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Form validates party name
- Event type selector works
- Creates party successfully
- Displays party code after creation
- Copy button copies code to clipboard
- Navigation buttons work correctly
