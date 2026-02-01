# U1: UI Project Setup

## Status
- [ ] Not started

## Dependencies
None

## Tasks
- [ ] Initialize Vite project with React + TypeScript template using pnpm
- [ ] Configure Biome for linting/formatting
- [ ] Configure Vitest for testing
- [ ] Setup project structure: `components/`, `pages/`, `api/`, `hooks/`, `types/`
- [ ] Add Tailwind CSS for styling
- [ ] Add `.nvmrc` or `.node-version` for Node.js version

## Details

### Project Initialization
```bash
pnpm create vite ui --template react-ts
cd ui
pnpm install
```

### Biome Configuration
```bash
pnpm add -D @biomejs/biome
pnpm biome init
```

**biome.json**:
```json
{
  "$schema": "https://biomejs.dev/schemas/1.9.0/schema.json",
  "organizeImports": { "enabled": true },
  "linter": {
    "enabled": true,
    "rules": { "recommended": true }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "tab"
  }
}
```

### Vitest Configuration
```bash
pnpm add -D vitest @testing-library/react @testing-library/jest-dom jsdom
```

**vitest.config.ts**:
```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './tests/setup.ts',
  },
})
```

### Tailwind CSS Setup
```bash
pnpm add -D tailwindcss postcss autoprefixer
pnpm tailwindcss init -p
```

### Project Structure
```
ui/
├── src/
│   ├── components/
│   │   └── ui/
│   ├── pages/
│   │   ├── guest/
│   │   └── admin/
│   ├── api/
│   ├── hooks/
│   ├── context/
│   ├── types/
│   └── routes/
├── tests/
│   └── setup.ts
└── docs/
```

### Node Version
**.node-version**:
```
22
```

## TDD Approach
1. Write a simple test to verify setup works
2. Verify Biome linting works
3. Verify Vitest runs
4. Verify Tailwind compiles

## Verification
- `pnpm dev` starts dev server
- `pnpm test` runs tests
- `pnpm biome check .` passes
- Tailwind styles work
