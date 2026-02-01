# X2: Error Handling & Loading States

## Status
- [ ] Not started

## Dependencies
- UI-U4 (API Client)

## Tasks
- [ ] Consistent error responses from server
- [ ] Error boundaries in UI
- [ ] Loading spinners/states for async operations

## Details

### Server Error Response Format
```go
// models/error.go
type APIError struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details any    `json:"details,omitempty"`
}

// Common error codes
const (
    ErrCodeNotFound      = "NOT_FOUND"
    ErrCodeBadRequest    = "BAD_REQUEST"
    ErrCodeUnauthorized  = "UNAUTHORIZED"
    ErrCodeForbidden     = "FORBIDDEN"
    ErrCodeInternal      = "INTERNAL_ERROR"
    ErrCodeConflict      = "CONFLICT"
    ErrCodeValidation    = "VALIDATION_ERROR"
)
```

```go
// handlers/error.go
func writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIError{
        Error: message,
        Code:  code,
    })
}

func handleNotFound(w http.ResponseWriter, message string) {
    writeError(w, http.StatusNotFound, ErrCodeNotFound, message)
}

func handleBadRequest(w http.ResponseWriter, message string) {
    writeError(w, http.StatusBadRequest, ErrCodeBadRequest, message)
}

func handleUnauthorized(w http.ResponseWriter) {
    writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Authentication required")
}

func handleForbidden(w http.ResponseWriter) {
    writeError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied")
}

func handleInternalError(w http.ResponseWriter, err error) {
    log.Printf("Internal error: %v", err)
    writeError(w, http.StatusInternalServerError, ErrCodeInternal, "An internal error occurred")
}
```

### UI Error Boundary
```typescript
// src/components/ErrorBoundary.tsx
interface Props {
  children: React.ReactNode
  fallback?: React.ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends React.Component<Props, State> {
  state: State = { hasError: false, error: null }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error boundary caught:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="p-8 text-center">
          <h2 className="text-xl font-bold text-red-600">Something went wrong</h2>
          <p className="text-gray-600 mt-2">Please refresh the page or try again later.</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded"
          >
            Refresh Page
          </button>
        </div>
      )
    }

    return this.props.children
  }
}
```

### UI Loading Spinner
```typescript
// src/components/ui/LoadingSpinner.tsx
interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function LoadingSpinner({ size = 'md', className }: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  }

  return (
    <div className={cn('flex justify-center items-center', className)}>
      <div
        className={cn(
          'animate-spin rounded-full border-2 border-gray-300 border-t-blue-600',
          sizeClasses[size]
        )}
      />
    </div>
  )
}
```

### UI Loading Page
```typescript
// src/components/LoadingPage.tsx
export function LoadingPage() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <LoadingSpinner size="lg" />
    </div>
  )
}
```

### UI Error Message Component
```typescript
// src/components/ui/ErrorMessage.tsx
interface ErrorMessageProps {
  children: React.ReactNode
  className?: string
}

export function ErrorMessage({ children, className }: ErrorMessageProps) {
  return (
    <div
      className={cn(
        'bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded',
        className
      )}
      role="alert"
    >
      {children}
    </div>
  )
}
```

### API Error Handling in Client
```typescript
// src/api/client.ts
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }

  static async fromResponse(response: Response): Promise<ApiError> {
    try {
      const body = await response.json()
      return new ApiError(response.status, body.code || 'UNKNOWN', body.error || 'Unknown error')
    } catch {
      return new ApiError(response.status, 'UNKNOWN', response.statusText)
    }
  }
}

// In fetch wrapper
if (!response.ok) {
  throw await ApiError.fromResponse(response)
}
```

### Custom Hooks for Loading State
```typescript
// src/hooks/useAsync.ts
interface AsyncState<T> {
  data: T | null
  loading: boolean
  error: Error | null
}

export function useAsync<T>(
  asyncFn: () => Promise<T>,
  deps: any[] = []
): AsyncState<T> & { refetch: () => void } {
  const [state, setState] = useState<AsyncState<T>>({
    data: null,
    loading: true,
    error: null,
  })

  const execute = useCallback(async () => {
    setState(s => ({ ...s, loading: true, error: null }))
    try {
      const data = await asyncFn()
      setState({ data, loading: false, error: null })
    } catch (error) {
      setState({ data: null, loading: false, error: error as Error })
    }
  }, deps)

  useEffect(() => {
    execute()
  }, [execute])

  return { ...state, refetch: execute }
}
```

## TDD Approach
1. Write tests for server error responses
2. Write tests for ErrorBoundary component
3. Write tests for LoadingSpinner component
4. Write tests for useAsync hook
5. Implement all components and utilities
6. Verify with tests

## Verification
- Server returns consistent error format
- Error boundary catches rendering errors
- Loading states display correctly
- API errors handled gracefully
- useAsync hook manages state correctly
