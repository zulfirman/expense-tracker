# DaisyUI Implementation Guide

This document provides guidelines for migrating the Expenses Tracker app to use DaisyUI components.

## Current Status

✅ **Completed:**
- Tailwind CSS 3.4.19 + DaisyUI 5.5.14 installed and configured
- Theme system: Pastel (light) and Synthwave (dark)
- Login page migrated to DaisyUI
- Signup page migrated to DaisyUI
- Navigation hidden on auth pages

🔄 **In Progress:**
- Migrating remaining pages one by one

## DaisyUI Component Reference

Based on [DaisyUI Documentation](https://daisyui.com/llms.txt), here are the key components and their usage:

### Buttons (`btn`)
```html
<!-- Basic button -->
<button class="btn btn-primary">Primary</button>

<!-- Sizes -->
<button class="btn btn-sm">Small</button>
<button class="btn btn-lg">Large</button>

<!-- Styles -->
<button class="btn btn-outline">Outline</button>
<button class="btn btn-ghost">Ghost</button>
<button class="btn btn-link">Link</button>

<!-- With loading -->
<button class="btn btn-primary" disabled>
  <span class="loading loading-spinner loading-sm"></span>
  Loading...
</button>
```

### Cards (`card`)
```html
<div class="card bg-base-100 shadow-xl">
  <figure>
    <img src="image.jpg" alt="Image" />
  </figure>
  <div class="card-body">
    <h2 class="card-title">Title</h2>
    <p>Content</p>
    <div class="card-actions justify-end">
      <button class="btn btn-primary">Action</button>
    </div>
  </div>
</div>
```

### Forms (`form-control`, `input`, `label`)
```html
<div class="form-control w-full">
  <label class="label">
    <span class="label-text">Label</span>
  </label>
  <input type="text" class="input input-bordered input-primary w-full" placeholder="Type here" />
  <label class="label">
    <span class="label-text-alt">Helper text</span>
  </label>
</div>
```

### Modals (`modal`)
```html
<dialog id="my_modal" class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Title</h3>
    <p class="py-4">Content</p>
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

### Tabs (`tabs`, `tab`)
```html
<div role="tablist" class="tabs tabs-boxed">
  <button role="tab" class="tab tab-active">Tab 1</button>
  <button role="tab" class="tab">Tab 2</button>
</div>
```

### Tables (`table`)
```html
<div class="overflow-x-auto">
  <table class="table">
    <thead>
      <tr>
        <th>Header</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Data</td>
      </tr>
    </tbody>
  </table>
</div>
```

### Badges (`badge`)
```html
<span class="badge badge-primary">Primary</span>
<span class="badge badge-outline">Outline</span>
<span class="badge badge-ghost">Ghost</span>
```

### Alerts (`alert`)
```html
<div role="alert" class="alert alert-info">
  <span>Info message</span>
</div>
```

### Loading (`loading`)
```html
<span class="loading loading-spinner"></span>
<span class="loading loading-dots"></span>
<span class="loading loading-ring"></span>
```

### Dividers (`divider`)
```html
<div class="divider">OR</div>
<div class="divider divider-vertical"></div>
```

## Color System

### Semantic Colors
Always use DaisyUI semantic color names for theme compatibility:

- `primary` / `primary-content`
- `secondary` / `secondary-content`
- `accent` / `accent-content`
- `base-100`, `base-200`, `base-300` / `base-content`
- `info` / `info-content`
- `success` / `success-content`
- `warning` / `warning-content`
- `error` / `error-content`

### Usage
```html
<!-- Backgrounds -->
<div class="bg-primary text-primary-content">Primary</div>
<div class="bg-base-100 text-base-content">Base</div>

<!-- Buttons -->
<button class="btn btn-primary">Primary</button>
<button class="btn btn-success">Success</button>

<!-- Inputs -->
<input class="input input-primary" />
<input class="input input-error" />
```

## Migration Checklist

When migrating a page:

1. ✅ Replace custom CSS classes with DaisyUI components
2. ✅ Use semantic color names (`primary`, `base-100`, etc.)
3. ✅ Replace custom buttons with `btn` classes
4. ✅ Replace custom cards with `card` component
5. ✅ Replace custom inputs with `input` classes
6. ✅ Replace custom modals with `modal` component
7. ✅ Use `form-control` for form layouts
8. ✅ Replace loading spinners with `loading` component
9. ✅ Use `divider` for separators
10. ✅ Remove inline styles and custom CSS where possible

## Common Patterns

### Page Layout
```html
<div class="min-h-screen bg-base-200 p-4">
  <div class="container mx-auto">
    <!-- Content -->
  </div>
</div>
```

### Form Layout
```html
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Form Title</h2>
    <form class="space-y-4">
      <div class="form-control">
        <!-- Form fields -->
      </div>
      <div class="form-control mt-6">
        <button class="btn btn-primary w-full">Submit</button>
      </div>
    </form>
  </div>
</div>
```

### List/Table
```html
<div class="overflow-x-auto">
  <table class="table table-zebra">
    <thead>
      <tr>
        <th>Column 1</th>
        <th>Column 2</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Data 1</td>
        <td>Data 2</td>
      </tr>
    </tbody>
  </table>
</div>
```

## Best Practices

1. **Always prefer DaisyUI components** over custom CSS
2. **Use semantic colors** for theme compatibility
3. **Use Tailwind utilities** for spacing, sizing, and layout
4. **Avoid `!important`** - Use `!` suffix on Tailwind classes if needed
5. **Responsive design** - Use Tailwind responsive prefixes
6. **Accessibility** - DaisyUI components include ARIA attributes
7. **Consistency** - Use the same patterns across all pages

## Resources

- [DaisyUI Components](https://daisyui.com/components/)
- [DaisyUI Themes](https://daisyui.com/docs/themes/)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [DaisyUI LLM Documentation](https://daisyui.com/llms.txt)










