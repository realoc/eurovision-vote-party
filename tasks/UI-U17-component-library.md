# U17: Component Library

## Status
- [ ] Not started

## Dependencies
- U1 (Project Setup)

## Tasks
- [ ] Button component (variants: primary, danger, disabled)
- [ ] Input component
- [ ] Card component
- [ ] List component
- [ ] Badge component (for notifications)
- [ ] Dropdown/Select component
- [ ] Copy-to-clipboard component

## Details

### Button Component
```typescript
// src/components/ui/Button.tsx
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost' | 'google'
  size?: 'sm' | 'md' | 'lg'
  children: React.ReactNode
}

export function Button({
  variant = 'primary',
  size = 'md',
  disabled,
  children,
  className,
  ...props
}: ButtonProps) {
  const baseStyles = 'font-medium rounded-lg transition-colors focus:outline-none focus:ring-2'

  const variantStyles = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
    secondary: 'bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-500',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
    success: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500',
    ghost: 'bg-transparent hover:bg-gray-100 focus:ring-gray-500',
    google: 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50',
  }

  const sizeStyles = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  }

  const disabledStyles = 'opacity-50 cursor-not-allowed'

  return (
    <button
      className={cn(
        baseStyles,
        variantStyles[variant],
        sizeStyles[size],
        disabled && disabledStyles,
        className
      )}
      disabled={disabled}
      {...props}
    >
      {children}
    </button>
  )
}
```

### Input Component
```typescript
// src/components/ui/Input.tsx
interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export function Input({ label, error, className, ...props }: InputProps) {
  return (
    <div className="flex flex-col gap-1">
      {label && (
        <label className="text-sm font-medium text-gray-700">{label}</label>
      )}
      <input
        className={cn(
          'px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500',
          error ? 'border-red-500' : 'border-gray-300',
          className
        )}
        {...props}
      />
      {error && <span className="text-sm text-red-500">{error}</span>}
    </div>
  )
}
```

### Card Component
```typescript
// src/components/ui/Card.tsx
interface CardProps {
  children: React.ReactNode
  className?: string
  onClick?: () => void
}

export function Card({ children, className, onClick }: CardProps) {
  return (
    <div
      className={cn(
        'bg-white rounded-lg shadow-md p-4',
        onClick && 'cursor-pointer hover:shadow-lg transition-shadow',
        className
      )}
      onClick={onClick}
    >
      {children}
    </div>
  )
}
```

### Badge Component
```typescript
// src/components/ui/Badge.tsx
interface BadgeProps {
  children: React.ReactNode
  variant?: 'default' | 'success' | 'danger' | 'warning' | 'secondary'
  className?: string
}

export function Badge({ children, variant = 'default', className }: BadgeProps) {
  const variantStyles = {
    default: 'bg-blue-100 text-blue-800',
    success: 'bg-green-100 text-green-800',
    danger: 'bg-red-100 text-red-800',
    warning: 'bg-yellow-100 text-yellow-800',
    secondary: 'bg-gray-100 text-gray-800',
  }

  return (
    <span
      className={cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
        variantStyles[variant],
        className
      )}
    >
      {children}
    </span>
  )
}
```

### Select Component
```typescript
// src/components/ui/Select.tsx
interface SelectProps {
  value?: string
  onChange: (value: string) => void
  placeholder?: string
  children: React.ReactNode
  className?: string
}

export function Select({ value, onChange, placeholder, children, className }: SelectProps) {
  return (
    <select
      value={value || ''}
      onChange={(e) => onChange(e.target.value)}
      className={cn(
        'px-3 py-2 border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500',
        !value && 'text-gray-400',
        className
      )}
    >
      {placeholder && (
        <option value="" disabled>
          {placeholder}
        </option>
      )}
      {children}
    </select>
  )
}
```

### CopyButton Component
```typescript
// src/components/ui/CopyButton.tsx
interface CopyButtonProps {
  text: string
  className?: string
}

export function CopyButton({ text, className }: CopyButtonProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <button
      onClick={handleCopy}
      className={cn(
        'inline-flex items-center gap-1 px-2 py-1 text-sm text-gray-600 hover:text-gray-800 transition-colors',
        className
      )}
    >
      {copied ? '✓ Copied!' : '📋 Copy'}
    </button>
  )
}
```

### Utility: cn (classnames)
```typescript
// src/lib/utils.ts
export function cn(...classes: (string | undefined | null | false)[]): string {
  return classes.filter(Boolean).join(' ')
}
```

## TDD Approach
1. Write tests for each component's rendering
2. Write tests for component interactions (click, change)
3. Write tests for variant styles
4. Write tests for disabled states
5. Implement components
6. Verify with `pnpm test`

## Verification
- `pnpm run lint` passes with zero errors and zero warnings
- All components render correctly
- Button variants display proper styles
- Input shows error state
- Badge variants display proper colors
- Select handles value changes
- CopyButton copies text and shows feedback
