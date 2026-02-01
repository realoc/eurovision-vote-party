# U3: Firebase Auth Setup (UI)

## Status
- [ ] Not started

## Dependencies
- U2 (Routing & Layout)

## Tasks
- [ ] Add Firebase SDK
- [ ] Create auth context/provider
- [ ] Implement login/logout functionality
- [ ] Create useAuth hook
- [ ] Persist auth state

## Details

### Installation
```bash
pnpm add firebase
```

### Firebase Configuration
```typescript
// src/config/firebase.ts
import { initializeApp } from 'firebase/app'
import { getAuth } from 'firebase/auth'

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
}

export const app = initializeApp(firebaseConfig)
export const auth = getAuth(app)
```

### Auth Context
```typescript
// src/context/AuthContext.tsx
interface AuthContextType {
  user: User | null
  loading: boolean
  signInWithEmail: (email: string, password: string) => Promise<void>
  signInWithGoogle: () => Promise<void>
  signOut: () => Promise<void>
  getIdToken: () => Promise<string | null>
}

export const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user)
      setLoading(false)
    })
    return unsubscribe
  }, [])

  // ... implement sign in methods
}
```

### useAuth Hook
```typescript
// src/hooks/useAuth.ts
export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
```

### Sign In Methods
- `signInWithEmailAndPassword` - Email/password login
- `signInWithPopup` + `GoogleAuthProvider` - Google login
- `signOut` - Logout

### Environment Variables
Create `.env.local`:
```
VITE_FIREBASE_API_KEY=xxx
VITE_FIREBASE_AUTH_DOMAIN=xxx
VITE_FIREBASE_PROJECT_ID=xxx
```

## TDD Approach
1. Write tests for AuthProvider with mocked Firebase
2. Write tests for useAuth hook
3. Implement auth context and provider
4. Verify with `pnpm test`

## Verification
- Auth state persists across page reloads
- Login with email/password works
- Login with Google works
- Logout works
- Loading state handled correctly
