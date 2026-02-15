# U18: Integration Tests & Test Infrastructure

## Status
- [ ] Not started

## Dependencies
- U6 (Guest Entry Page)
- U7 (Guest Waiting Page)
- U8 (Party Overview Page)
- U9 (Voting Page)
- U10 (Results Page)
- U11 (Admin Login Page)
- U12 (Admin Profile Setup)
- U13 (Admin Dashboard)
- U14 (Party Creation Page)
- U15 (Admin Party Overview Page)
- U16 (Join Requests Page)

## Note
Component unit tests are written as part of each task (TDD). This task is for integration tests and test infrastructure.

## Tasks
- [ ] Setup MSW (Mock Service Worker) for API mocking
- [ ] End-to-end integration tests for full user flows
- [ ] Test guest flow: entry → waiting → approved → vote → results
- [ ] Test admin flow: login → create party → approve guests → end voting
- [ ] Add test scripts to package.json (`pnpm test`, `pnpm test:integration`)

## Details

### MSW Setup
```bash
pnpm add -D msw
```

```typescript
// tests/mocks/handlers.ts
import { http, HttpResponse } from 'msw'

export const handlers = [
  // Health
  http.get('/api/health', () => {
    return HttpResponse.json({ status: 'ok' })
  }),

  // Parties
  http.post('/api/parties', async ({ request }) => {
    const body = await request.json()
    return HttpResponse.json({
      id: 'party-1',
      code: 'ABC123',
      name: body.name,
      eventType: body.eventType,
      status: 'active',
      createdAt: new Date().toISOString(),
    }, { status: 201 })
  }),

  // Guests
  http.post('/api/parties/:code/join', async ({ request }) => {
    const body = await request.json()
    return HttpResponse.json({
      id: 'guest-1',
      username: body.username,
      status: 'pending',
      createdAt: new Date().toISOString(),
    }, { status: 201 })
  }),

  // ... more handlers
]
```

```typescript
// tests/mocks/server.ts
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
```

```typescript
// tests/setup.ts
import '@testing-library/jest-dom'
import { server } from './mocks/server'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
```

### Integration Test: Guest Flow
```typescript
// tests/integration/guest-flow.test.tsx
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { App } from '../../src/App'

describe('Guest Flow', () => {
  it('completes full guest journey: join → wait → vote → results', async () => {
    const user = userEvent.setup()
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    )

    // Enter party code and username
    await user.type(screen.getByLabelText(/party code/i), 'ABC123')
    await user.type(screen.getByLabelText(/your name/i), 'TestUser')
    await user.click(screen.getByRole('button', { name: /join/i }))

    // Should be on waiting page
    await waitFor(() => {
      expect(screen.getByText(/waiting for approval/i)).toBeInTheDocument()
    })

    // Simulate approval (update mock to return approved)
    server.use(
      http.get('/api/parties/:code/guest-status', () => {
        return HttpResponse.json({ status: 'approved' })
      })
    )

    // Should redirect to party overview
    await waitFor(() => {
      expect(screen.getByText(/vote now/i)).toBeInTheDocument()
    })

    // Navigate to voting
    await user.click(screen.getByRole('button', { name: /vote/i }))

    // Complete voting form
    // ... select acts for each point value

    // Submit vote
    await user.click(screen.getByRole('button', { name: /submit/i }))

    // Should be back on overview with votes shown
    await waitFor(() => {
      expect(screen.getByText(/your votes/i)).toBeInTheDocument()
    })
  })
})
```

### Integration Test: Admin Flow
```typescript
// tests/integration/admin-flow.test.tsx
describe('Admin Flow', () => {
  it('completes admin journey: login → create party → approve → end voting', async () => {
    const user = userEvent.setup()
    render(
      <MemoryRouter initialEntries={['/admin/login']}>
        <App />
      </MemoryRouter>
    )

    // Mock Firebase auth
    // ...

    // Login
    await user.type(screen.getByLabelText(/email/i), 'admin@test.com')
    await user.type(screen.getByLabelText(/password/i), 'password123')
    await user.click(screen.getByRole('button', { name: /sign in/i }))

    // Should be on dashboard
    await waitFor(() => {
      expect(screen.getByText(/your parties/i)).toBeInTheDocument()
    })

    // Create party
    await user.click(screen.getByRole('button', { name: /create/i }))
    await user.type(screen.getByLabelText(/party name/i), 'Test Party')
    await user.click(screen.getByRole('button', { name: /create party/i }))

    // Should show party code
    await waitFor(() => {
      expect(screen.getByText(/party created/i)).toBeInTheDocument()
    })

    // Navigate to party overview
    await user.click(screen.getByRole('button', { name: /go to party/i }))

    // Approve join request
    await user.click(screen.getByText(/join requests/i))
    await user.click(screen.getByRole('button', { name: /approve/i }))

    // End voting
    await user.click(screen.getByRole('button', { name: /end voting/i }))
    await user.click(screen.getByRole('button', { name: /confirm/i }))

    // Should show results option
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /view results/i })).toBeInTheDocument()
    })
  })
})
```

### Package.json Scripts
```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "test:integration": "vitest --config vitest.integration.config.ts"
  }
}
```

### Integration Vitest Config
```typescript
// vitest.integration.config.ts
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './tests/setup.ts',
    include: ['tests/integration/**/*.test.{ts,tsx}'],
    testTimeout: 10000,
  },
})
```

## TDD Approach
1. Setup MSW infrastructure
2. Create mock handlers for all endpoints
3. Write guest flow integration test
4. Write admin flow integration test
5. Add edge case tests
6. Verify with `pnpm test:integration`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- MSW intercepts all API calls
- Guest flow test passes
- Admin flow test passes
- Tests run in CI
- Coverage report includes integration tests
