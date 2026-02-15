# X1: CORS & API Configuration

## Status
- [ ] Not started

## Dependencies
- Server-S1 (Server Project Setup)
- UI-U1 (UI Project Setup)

## Tasks
- [ ] Configure CORS on server for UI origin
- [ ] Environment-based API URL configuration
- [ ] Handle development vs production URLs

## Details

### Server CORS Configuration
```go
// middleware/cors.go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            // Check if origin is allowed
            for _, allowed := range allowedOrigins {
                if origin == allowed || allowed == "*" {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    break
                }
            }

            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")

            // Handle preflight
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusNoContent)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Server Environment Configuration
```go
// config/config.go
type Config struct {
    Port           string
    AllowedOrigins []string
    FirestoreHost  string // For emulator
    Environment    string // "development" or "production"
}

func LoadConfig() *Config {
    env := os.Getenv("ENVIRONMENT")
    if env == "" {
        env = "development"
    }

    allowedOrigins := []string{"http://localhost:5173"}
    if env == "production" {
        allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
    }

    return &Config{
        Port:           getEnv("PORT", "8080"),
        AllowedOrigins: allowedOrigins,
        FirestoreHost:  os.Getenv("FIRESTORE_EMULATOR_HOST"),
        Environment:    env,
    }
}
```

### UI Environment Configuration
```typescript
// src/config/env.ts
export const config = {
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  firebaseConfig: {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  },
}
```

### Environment Files

**UI .env.development**
```
VITE_API_URL=http://localhost:8080/api
VITE_FIREBASE_API_KEY=demo-key
VITE_FIREBASE_AUTH_DOMAIN=localhost
VITE_FIREBASE_PROJECT_ID=demo-project
```

**UI .env.production**
```
VITE_API_URL=https://api.eurovision-party.example.com/api
VITE_FIREBASE_API_KEY=<production-key>
VITE_FIREBASE_AUTH_DOMAIN=<project>.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=<project>
```

**Server .env.development**
```
PORT=8080
ENVIRONMENT=development
FIRESTORE_EMULATOR_HOST=localhost:8081
```

### Development Setup
1. Server runs on `localhost:8080`
2. UI runs on `localhost:5173` (Vite default)
3. CORS allows `localhost:5173`
4. Firestore emulator on `localhost:8081`

### Production Setup
1. Server deployed with `ALLOWED_ORIGINS` set
2. UI built with production env vars
3. Real Firestore connection
4. Firebase Auth with real config

## TDD Approach
1. Write tests for CORS middleware
2. Write tests for config loading
3. Implement CORS middleware
4. Implement config module
5. Verify with `go test ./...`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- CORS headers set correctly in responses
- Preflight requests handled
- API calls from UI succeed
- Environment switching works
