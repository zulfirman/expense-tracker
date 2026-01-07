# DaisyUI Setup Guide

This project uses **DaisyUI 5** with **Tailwind CSS 3.4.19** for styling.

## Current Setup

### Installed Packages
- `tailwindcss@^3.4.19`
- `daisyui@^5.5.14`
- `@tailwindcss/vite@latest` (for Tailwind CSS v4 compatibility, but we're using v3)

### Configuration Files

#### `tailwind.config.cjs`
```javascript
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: ['pastel', 'synthwave'],
    darkTheme: 'synthwave',
    base: true,
    styled: true,
    utils: true,
    prefix: '',
    logs: true,
    themeRoot: ':root',
  },
}
```

#### `postcss.config.js`
```javascript
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

#### `src/app.css`
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

## Themes

The app uses two DaisyUI themes:
- **pastel** - Light theme (default)
- **synthwave** - Dark theme

Theme switching is handled via the `theme` store in `src/lib/stores/theme.js`, which sets the `data-theme` attribute on the `<html>` element.

## Using DaisyUI Components

### Basic Rules

1. **Use DaisyUI class names** - Prefer DaisyUI component classes over custom CSS
2. **Customize with Tailwind utilities** - Use Tailwind utility classes for additional styling
3. **Use `!` for overrides** - Only when necessary to override DaisyUI styles
4. **Semantic colors** - Use DaisyUI color names (`primary`, `secondary`, `base-100`, etc.) instead of Tailwind color names for theme-aware colors

### Common Components

#### Buttons
```html
<button class="btn btn-primary">Primary Button</button>
<button class="btn btn-secondary btn-outline">Secondary Outline</button>
<button class="btn btn-ghost btn-sm">Ghost Small</button>
```

#### Cards
```html
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Card Title</h2>
    <p>Card content</p>
  </div>
</div>
```

#### Inputs
```html
<input type="text" class="input input-bordered input-primary w-full" placeholder="Type here" />
```

#### Forms
```html
<div class="form-control">
  <label class="label">
    <span class="label-text">Label</span>
  </label>
  <input type="text" class="input input-bordered" />
</div>
```

#### Modals
```html
<dialog id="my_modal" class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Modal Title</h3>
    <p class="py-4">Modal content</p>
    <div class="modal-action">
      <form method="dialog">
        <button class="btn">Close</button>
      </form>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
```

#### Loading States
```html
<span class="loading loading-spinner loading-sm"></span>
```

#### Dividers
```html
<div class="divider">OR</div>
```

## Color System

### DaisyUI Color Names
- `primary` - Primary brand color
- `primary-content` - Text color on primary background
- `secondary` - Secondary brand color
- `secondary-content` - Text color on secondary background
- `accent` - Accent color
- `accent-content` - Text color on accent background
- `base-100` - Base surface color (page background)
- `base-200` - Darker shade for elevations
- `base-300` - Even darker shade
- `base-content` - Text color on base background
- `info`, `success`, `warning`, `error` - Status colors
- `*-content` - Text colors for each status

### Usage Examples
```html
<!-- Background colors -->
<div class="bg-primary text-primary-content">Primary background</div>
<div class="bg-base-100 text-base-content">Base background</div>

<!-- Button colors -->
<button class="btn btn-primary">Primary</button>
<button class="btn btn-success">Success</button>
<button class="btn btn-error">Error</button>

<!-- Input colors -->
<input class="input input-primary" />
<input class="input input-success" />
```

## Migration Status

### ✅ Completed
- Login page - Fully migrated to DaisyUI
- Signup page - Fully migrated to DaisyUI
- Theme system - Configured with pastel/synthwave themes
- Tailwind CSS setup - Configured with PostCSS

### 🔄 In Progress
- Other pages - To be migrated one by one

### 📋 To Do
- Expenses page
- Income page
- History page
- Budget page
- Categories page
- Preferences page
- Profile page
- All modals and components

## Best Practices

1. **Always use DaisyUI components** when available instead of custom CSS
2. **Use semantic color names** (`primary`, `base-100`, etc.) for theme compatibility
3. **Avoid custom CSS** - Use Tailwind utilities or DaisyUI modifiers instead
4. **Responsive design** - Use Tailwind responsive prefixes (`sm:`, `md:`, `lg:`, etc.)
5. **Accessibility** - DaisyUI components include accessibility features by default
6. **Consistency** - Use the same component patterns across the app

## Resources

- [DaisyUI Documentation](https://daisyui.com)
- [DaisyUI Components](https://daisyui.com/components/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [DaisyUI Themes](https://daisyui.com/docs/themes/)


