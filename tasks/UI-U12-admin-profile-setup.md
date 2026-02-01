# U12: Admin Profile Setup

## Status
- [ ] Not started

## Dependencies
- U11 (Admin Login Page)

## Tasks
- [ ] Username input for first-time setup
- [ ] Check if username exists on login
- [ ] Redirect to dashboard after setup

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Complete Your Profile           │
│                                     │
│  Welcome! Before you can create     │
│  parties, we need a display name.   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Display Name                 │   │
│  │ [John                    ]   │   │
│  └─────────────────────────────┘   │
│  3-30 characters                    │
│                                     │
│  ┌─────────────────────────────┐   │
│  │      Complete Setup          │   │
│  └─────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/ProfileSetupPage.tsx
export function ProfileSetupPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  // Check if profile already exists
  useEffect(() => {
    checkExistingProfile()
  }, [])

  const checkExistingProfile = async () => {
    try {
      await userApi.getProfile()
      // Profile exists, redirect to dashboard
      navigate('/admin')
    } catch {
      // No profile, stay on setup page
    }
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()

    if (username.length < 3 || username.length > 30) {
      setError('Display name must be 3-30 characters')
      return
    }

    if (!/^[a-zA-Z0-9_]+$/.test(username)) {
      setError('Display name can only contain letters, numbers, and underscores')
      return
    }

    setError(null)
    setLoading(true)

    try {
      await userApi.updateProfile(username)
      navigate('/admin')
    } catch (err) {
      setError('Failed to save profile. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const isValid = username.length >= 3 && username.length <= 30 && /^[a-zA-Z0-9_]+$/.test(username)

  return (
    <div>
      <h1>Complete Your Profile</h1>
      <p>Welcome! Before you can create parties, we need a display name.</p>

      <form onSubmit={handleSubmit}>
        <Input
          label="Display Name"
          value={username}
          onChange={setUsername}
          placeholder="Enter your name"
          maxLength={30}
        />
        <p className="text-sm text-gray-500">3-30 characters, letters, numbers, and underscores only</p>

        {error && <ErrorMessage>{error}</ErrorMessage>}

        <Button type="submit" disabled={!isValid || loading}>
          {loading ? 'Saving...' : 'Complete Setup'}
        </Button>
      </form>
    </div>
  )
}
```

### Validation Rules
- Required field
- 3-30 characters
- Alphanumeric and underscores only: `/^[a-zA-Z0-9_]+$/`

### Flow
1. Check if profile already exists → redirect to dashboard
2. User enters username
3. Validate input
4. Submit to API
5. Redirect to dashboard on success

## TDD Approach
1. Write tests for validation logic
2. Write tests for existing profile check
3. Write tests for form submission
4. Write tests for redirect on success
5. Implement page component
6. Verify with `pnpm test`

## Verification
- Redirects if profile already exists
- Validates username format
- Shows validation errors
- Successfully creates profile
- Redirects to dashboard after setup
