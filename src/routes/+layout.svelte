<script>
  import '../app.css';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
</script>

<div class="app-container">
  <main>
    <slot />
  </main>

  <nav class="bottom-nav">
    <a 
      href="/expenses"
      data-sveltekit-preload-data="hover"
      class="nav-item" 
      class:active={$page.url.pathname === '/expenses' || $page.url.pathname === '/'}
      on:click|preventDefault={() => goto('/expenses')}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
      <span>Expenses</span>
    </a>
    <a 
      href="/income"
      data-sveltekit-preload-data="hover"
      class="nav-item" 
      class:active={$page.url.pathname === '/income'}
      on:click|preventDefault={() => goto('/income')}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="19" x2="12" y2="5"></line>
        <polyline points="5 12 12 5 19 12"></polyline>
      </svg>
      <span>Income</span>
    </a>
    <a 
      href="/history"
      data-sveltekit-preload-data="hover"
      class="nav-item" 
      class:active={$page.url.pathname === '/history'}
      on:click|preventDefault={() => goto('/history')}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
        <line x1="16" y1="2" x2="16" y2="6"></line>
        <line x1="8" y1="2" x2="8" y2="6"></line>
        <line x1="3" y1="10" x2="21" y2="10"></line>
      </svg>
      <span>History</span>
    </a>
  </nav>
</div>

<style>
  .app-container {
    height: 100vh;
    max-height: 100vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
  }

  main {
    flex: 1;
    padding: 1rem;
    max-width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    min-height: 0;
    height: 0;
    -webkit-overflow-scrolling: touch;
    padding-bottom: 80px;
  }

  .bottom-nav {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: var(--surface);
    border-top: 1px solid var(--border);
    display: flex;
    justify-content: space-around;
    padding: 0.5rem 0 calc(0.5rem + env(safe-area-inset-bottom));
    z-index: 1000;
    box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  }

  .nav-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    padding: 0.5rem 1rem;
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    transition: color 0.2s;
    font-size: 0.75rem;
    flex: 1;
    min-width: 0;
    text-decoration: none;
  }

  .nav-item svg {
    width: 24px;
    height: 24px;
    flex-shrink: 0;
  }

  .nav-item.active {
    color: var(--primary-color);
  }

  .nav-item:active {
    opacity: 0.7;
  }

  .nav-item span {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
  }

  @media (max-width: 768px) {
    .nav-item {
      font-size: 0.7rem;
      padding: 0.5rem 0.5rem;
    }

    .nav-item svg {
      width: 22px;
      height: 22px;
    }
  }
</style>
