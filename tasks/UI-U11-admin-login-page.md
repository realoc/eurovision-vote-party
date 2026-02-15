# U11: Admin Login Page

## Status
- [ ] Not started

## Dependencies
- U3 (Firebase Auth Setup)

## Tasks
- [ ] Firebase Auth login UI with email/password form
- [ ] Google sign-in button
- [ ] Error handling for login failures
- [ ] Redirect to admin dashboard on success

## Details

### Page Design
```
┌─────────────────────────────────────┐
│     Admin Login                     │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Email                        │   │
│  │ [admin@example.com       ]   │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Password                     │   │
│  │ [••••••••                ]   │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │         Sign In              │   │
│  └─────────────────────────────┘   │
│                                     │
│         - or -                      │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ 🔵 Sign in with Google       │   │
│  └─────────────────────────────┘   │
│                                     │
│  [← Back to Home]                   │
│                                     │
│  Don't have an account?             │
│  [Create Account]                   │
│                                     │
└─────────────────────────────────────┘
```

### Component Structure
```typescript
// src/pages/admin/LoginPage.tsx
export function LoginPage() {
  const { signInWithEmail, signInWithGoogle, user } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [isSignUp, setIsSignUp] = useState(false)

  // Redirect if already logged in
  useEffect(() => {
    if (user) {
      checkProfileAndRedirect()
    }
  }, [user])

  const checkProfileAndRedirect = async () => {
    try {
      await userApi.getProfile()
      navigate('/admin')
    } catch {
      // No profile, redirect to setup
      navigate('/admin/profile')
    }
  }

  const handleEmailSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      if (isSignUp) {
        await createUserWithEmailAndPassword(auth, email, password)
      } else {
        await signInWithEmail(email, password)
      }
      // Auth state change will trigger redirect
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleGoogleSignIn = async () => {
    setError(null)
    try {
      await signInWithGoogle()
      // Auth state change will trigger redirect
    } catch (err) {
      setError(getErrorMessage(err))
    }
  }

  return (
    <div>
      <h1>{isSignUp ? 'Create Account' : 'Admin Login'}</h1>

      <form onSubmit={handleEmailSubmit}>
        <Input
          type="email"
          label="Email"
          value={email}
          onChange={setEmail}
          required
        />
        <Input
          type="password"
          label="Password"
          value={password}
          onChange={setPassword}
          required
          minLength={6}
        />
        {error && <ErrorMessage>{error}</ErrorMessage>}
        <Button type="submit" disabled={loading}>
          {loading ? 'Please wait...' : (isSignUp ? 'Create Account' : 'Sign In')}
        </Button>
      </form>

      <Divider>or</Divider>

      <Button variant="google" onClick={handleGoogleSignIn}>
        Sign in with Google
      </Button>

      <Link to="/">← Back to Home</Link>

      <p>
        {isSignUp ? 'Already have an account?' : "Don't have an account?"}
        <button onClick={() => setIsSignUp(!isSignUp)}>
          {isSignUp ? 'Sign In' : 'Create Account'}
        </button>
      </p>
    </div>
  )
}
```

### Error Handling
```typescript
function getErrorMessage(error: unknown): string {
  if (error instanceof FirebaseError) {
    switch (error.code) {
      case 'auth/invalid-email':
        return 'Invalid email address'
      case 'auth/user-disabled':
        return 'This account has been disabled'
      case 'auth/user-not-found':
        return 'No account found with this email'
      case 'auth/wrong-password':
        return 'Incorrect password'
      case 'auth/email-already-in-use':
        return 'An account already exists with this email'
      case 'auth/weak-password':
        return 'Password should be at least 6 characters'
      default:
        return 'An error occurred. Please try again.'
    }
  }
  return 'An error occurred. Please try again.'
}
```

### Flow
1. If user already logged in → check profile → redirect
2. Email/password sign in → check profile → redirect
3. Google sign in → check profile → redirect
4. If no profile exists → redirect to profile setup

## TDD Approach
1. Write tests for email sign in flow
2. Write tests for Google sign in flow
3. Write tests for error handling
4. Write tests for redirect logic
5. Implement page component
6. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- Email/password login works
- Google login works
- Error messages display correctly
- Redirects to profile setup if no profile
- Redirects to dashboard if profile exists
- Sign up mode works
