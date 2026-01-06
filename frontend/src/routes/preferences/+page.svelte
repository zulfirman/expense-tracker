<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import '$lib/styles/shared.css';

  let currentTheme = 'light';

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }
    currentTheme = $theme;
  });

  $: if ($theme) {
    currentTheme = $theme;
  }

  function handleThemeToggle() {
    theme.toggle();
  }
</script>

<div class="preferences-page">
  <div class="header">
    <h1>Preferences</h1>
  </div>

  <div class="nav-grid">
    <a class="nav-card" href="/preferences/categories">
      <div>
        <h3>Categories</h3>
        <p>Manage and reorder your expense categories.</p>
      </div>
      <span class="nav-arrow">›</span>
    </a>
    <a class="nav-card" href="/preferences/currency">
      <div>
        <h3>Currency & Quick Amount</h3>
        <p>Configure currency and quick amount shortcuts.</p>
      </div>
      <span class="nav-arrow">›</span>
    </a>
    <a class="nav-card" href="/preferences/change-password">
      <div>
        <h3>Change Password</h3>
        <p>Update your account password securely.</p>
      </div>
      <span class="nav-arrow">›</span>
    </a>
  </div>

  <!-- Theme Toggle -->
  <div class="preference-section">
    <h2>Appearance</h2>
    <p class="preference-description">Choose between light and dark mode</p>
    <div class="theme-toggle-section">
      <span class="theme-label">Theme: {currentTheme === 'dark' ? 'Dark' : 'Light'}</span>
      <button
        class="theme-toggle-btn"
        on:click={handleThemeToggle}
        title="Toggle Theme"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="sun-icon">
          <circle cx="12" cy="12" r="5"></circle>
          <line x1="12" y1="1" x2="12" y2="3"></line>
          <line x1="12" y1="21" x2="12" y2="23"></line>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
          <line x1="1" y1="12" x2="3" y2="12"></line>
          <line x1="21" y1="12" x2="23" y2="12"></line>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
        </svg>
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="moon-icon">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
        </svg>
      </button>
    </div>
  </div>
</div>

<style>
  .preferences-page {
    max-width: 600px;
    margin: 0 auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
  }

  .preference-section {
    background: var(--surface);
    border-radius: 0.75rem;
    padding: 1.5rem;
    border: 1px solid var(--border);
    margin-bottom: 1.5rem;
  }

  .nav-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .nav-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    background: var(--surface);
    color: var(--text-primary);
    text-decoration: none;
    transition: all 0.2s;
  }

  .nav-card:hover {
    border-color: var(--primary-color);
    background: var(--background);
  }

  .nav-card.active {
    border-color: var(--primary-color);
  }

  .nav-card h3 {
    margin: 0;
    font-size: 1rem;
  }

  .nav-card p {
    margin: 0.25rem 0 0;
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .nav-arrow {
    font-size: 1.25rem;
    color: var(--text-secondary);
  }

  .preference-section:last-child {
    margin-bottom: 0;
  }

  .preference-section h2 {
    font-size: 1.125rem;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }

  .preference-description {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }

  .currency-options {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .currency-option {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .currency-option:hover {
    background: var(--background);
    border-color: var(--primary-color);
  }

  .currency-option input[type="radio"] {
    cursor: pointer;
  }

  .currency-option input[type="radio"]:checked + span {
    font-weight: 600;
    color: var(--primary-color);
  }

  .currency-option:has(input[type="radio"]:checked) {
    border-color: var(--primary-color);
    background: rgba(79, 70, 229, 0.05);
  }


  .theme-toggle-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
  }

  .theme-label {
    font-weight: 500;
    color: var(--text-primary);
  }

  .theme-toggle-btn {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    border: 1px solid var(--border);
    background: var(--surface);
    color: var(--text-primary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .theme-toggle-btn:hover {
    background: var(--background);
    transform: scale(1.05);
  }

  .theme-toggle-btn svg {
    width: 20px;
    height: 20px;
  }

  .moon-icon {
    display: none;
  }

  :global(.dark) .sun-icon {
    display: none;
  }

  :global(.dark) .moon-icon {
    display: block;
  }
</style>

