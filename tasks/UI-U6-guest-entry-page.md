# U6: Guest Entry Page

## Status
- [x] Done

## Dependencies
- U4 (API Client)
- U5 (Type Definitions)

## Tasks
- [x] Party code input field
- [x] Username input field
- [x] Submit button
- [x] Error handling (invalid code, etc.)
- [x] Redirect to waiting page on success

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Eurovision Vote Party           │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Party Code                   │   │
│  │ [ABC123            ]         │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Your Name                    │   │
│  │ [                  ]         │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │        Join Party            │   │
│  └─────────────────────────────┘   │
│                                     │
│         - or -                      │
│                                     │
│  [Admin Login]                      │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/guest/EntryPage.tsx
export function EntryPage() {
  const [code, setCode] = useState('')
  const [username, setUsername] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      const guest = await guestApi.join(code.toUpperCase(), username)
      // Store guest ID for polling
      localStorage.setItem(`guest_${code}`, guest.id)
      navigate(`/waiting/${code}`)
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.status === 404) {
          setError('Party not found. Check your code.')
        } else {
          setError('Something went wrong. Please try again.')
        }
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <Input
        label="Party Code"
        value={code}
        onChange={setCode}
        placeholder="ABC123"
        maxLength={6}
      />
      <Input
        label="Your Name"
        value={username}
        onChange={setUsername}
        placeholder="Enter your name"
      />
      {error && <ErrorMessage>{error}</ErrorMessage>}
      <Button type="submit" disabled={loading || !code || !username}>
        {loading ? 'Joining...' : 'Join Party'}
      </Button>
      <Link to="/admin/login">Admin Login</Link>
    </form>
  )
}
```

### Validation
- Party code: 6 characters, alphanumeric, converted to uppercase
- Username: 3-30 characters, required

### Error States
- Invalid party code format
- Party not found (404)
- Network error
- Server error

### Local Storage
Store guest ID in localStorage for polling:
- Key: `guest_{code}`
- Value: guest ID

## TDD Approach
1. Write tests for form validation
2. Write tests for submission flow (success and error cases)
3. Write tests for navigation
4. Implement page component
5. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Form validates inputs
- Shows error for invalid party code
- Navigates to waiting page on success
- Guest ID stored in localStorage
- Admin login link works
