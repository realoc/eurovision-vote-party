# U16: Join Requests Page (Admin)

## Status
- [ ] Not started

## Dependencies
- U15 (Party Overview Page)

## Tasks
- [ ] List of pending join requests
- [ ] Username display for each request
- [ ] Approve button per request
- [ ] Reject button per request
- [ ] Back to party overview button

## Details

### Page Design
```
┌─────────────────────────────────────┐
│  [← Back]     Join Requests         │
│                                     │
│  3 pending requests                 │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Sarah                        │   │
│  │ Requested 2 minutes ago      │   │
│  │ [✓ Approve]  [✗ Reject]     │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Mike                         │   │
│  │ Requested 5 minutes ago      │   │
│  │ [✓ Approve]  [✗ Reject]     │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Lisa                         │   │
│  │ Requested 10 minutes ago     │   │
│  │ [✓ Approve]  [✗ Reject]     │   │
│  └─────────────────────────────┘   │
│                                     │
│  --- No more requests ---           │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/JoinRequestsPage.tsx
export function JoinRequestsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [requests, setRequests] = useState<Guest[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadRequests()
    // Poll for new requests
    const interval = setInterval(loadRequests, 5000)
    return () => clearInterval(interval)
  }, [id])

  const loadRequests = async () => {
    try {
      const data = await guestApi.getJoinRequests(id!)
      setRequests(data.requests)
    } catch (err) {
      // Handle error
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (guestId: string) => {
    try {
      await guestApi.approve(id!, guestId)
      setRequests(prev => prev.filter(r => r.id !== guestId))
    } catch (err) {
      // Handle error
    }
  }

  const handleReject = async (guestId: string) => {
    try {
      await guestApi.reject(id!, guestId)
      setRequests(prev => prev.filter(r => r.id !== guestId))
    } catch (err) {
      // Handle error
    }
  }

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <header className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => navigate(`/admin/party/${id}`)}>
          ← Back
        </Button>
        <h1>Join Requests</h1>
      </header>

      <p>{requests.length} pending request{requests.length !== 1 ? 's' : ''}</p>

      {requests.length === 0 ? (
        <p className="text-center text-gray-500">No pending requests</p>
      ) : (
        <div className="space-y-4">
          {requests.map(request => (
            <JoinRequestCard
              key={request.id}
              request={request}
              onApprove={() => handleApprove(request.id)}
              onReject={() => handleReject(request.id)}
            />
          ))}
        </div>
      )}
    </div>
  )
}
```

### JoinRequestCard Component
```typescript
// src/components/JoinRequestCard.tsx
interface JoinRequestCardProps {
  request: Guest
  onApprove: () => void
  onReject: () => void
}

export function JoinRequestCard({ request, onApprove, onReject }: JoinRequestCardProps) {
  const timeAgo = formatTimeAgo(request.createdAt)

  return (
    <Card>
      <div className="flex justify-between items-center">
        <div>
          <h3 className="font-bold">{request.username}</h3>
          <p className="text-sm text-gray-500">Requested {timeAgo}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="success" size="sm" onClick={onApprove}>
            ✓ Approve
          </Button>
          <Button variant="danger" size="sm" onClick={onReject}>
            ✗ Reject
          </Button>
        </div>
      </div>
    </Card>
  )
}
```

### Time Formatting
```typescript
// src/utils/time.ts
export function formatTimeAgo(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`
  return `${Math.floor(seconds / 86400)} days ago`
}
```

## TDD Approach
1. Write tests for request list display
2. Write tests for approve action
3. Write tests for reject action
4. Write tests for empty state
5. Write tests for time formatting utility
6. Implement page and components
7. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Shows list of pending requests
- Shows "no requests" when empty
- Approve button works and removes from list
- Reject button works and removes from list
- Time ago displays correctly
- Polls for new requests
- Back button navigates correctly
