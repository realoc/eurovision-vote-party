# U7: Guest Waiting Page

## Status
- [ ] Not started

## Dependencies
- U6 (Guest Entry Page)

## Tasks
- [ ] Display "Waiting for approval" message
- [ ] Poll server every 3-5 seconds for approval status
- [ ] Redirect to party overview on approval
- [ ] Handle rejection (show message, redirect to entry)

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Eurovision Vote Party           │
│                                     │
│         ⏳                          │
│                                     │
│   Waiting for approval...           │
│                                     │
│   You've requested to join          │
│   "John's Eurovision Party"         │
│                                     │
│   The party admin will review       │
│   your request shortly.             │
│                                     │
│   [Cancel]                          │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/guest/WaitingPage.tsx
export function WaitingPage() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const [status, setStatus] = useState<GuestStatus>('pending')
  const [partyName, setPartyName] = useState<string>('')

  const guestId = localStorage.getItem(`guest_${code}`)

  useEffect(() => {
    if (!guestId || !code) {
      navigate('/')
      return
    }

    // Initial fetch
    fetchStatus()

    // Poll every 3 seconds
    const interval = setInterval(fetchStatus, 3000)
    return () => clearInterval(interval)
  }, [code, guestId])

  const fetchStatus = async () => {
    try {
      const response = await guestApi.getStatus(code!, guestId!)
      setStatus(response.status)

      if (response.status === 'approved') {
        navigate(`/party/${code}`)
      } else if (response.status === 'rejected') {
        // Show rejection message briefly, then redirect
        setTimeout(() => {
          localStorage.removeItem(`guest_${code}`)
          navigate('/')
        }, 3000)
      }
    } catch (err) {
      console.error('Failed to fetch status:', err)
    }
  }

  const handleCancel = () => {
    localStorage.removeItem(`guest_${code}`)
    navigate('/')
  }

  if (status === 'rejected') {
    return (
      <div>
        <h2>Request Rejected</h2>
        <p>Sorry, your request to join was rejected.</p>
        <p>Redirecting to home page...</p>
      </div>
    )
  }

  return (
    <div>
      <Spinner />
      <h2>Waiting for approval...</h2>
      <p>You've requested to join "{partyName}"</p>
      <p>The party admin will review your request shortly.</p>
      <Button variant="secondary" onClick={handleCancel}>
        Cancel
      </Button>
    </div>
  )
}
```

### Polling Logic
- Start polling when component mounts
- Poll every 3 seconds
- Stop polling on unmount or status change
- Handle network errors gracefully (continue polling)

### Status Handling
- `pending`: Continue showing waiting message
- `approved`: Navigate to party overview
- `rejected`: Show rejection message, then redirect home

### Local Storage
- Read guest ID from `guest_{code}`
- Clear on cancel or rejection

## TDD Approach
1. Write tests for initial render
2. Write tests for polling behavior (mock timers)
3. Write tests for status transitions
4. Write tests for cancel functionality
5. Implement page component
6. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Shows waiting message on load
- Polls server every 3 seconds
- Redirects on approval
- Shows rejection message and redirects
- Cancel button works
