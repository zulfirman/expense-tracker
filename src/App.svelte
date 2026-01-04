<script>
  import InputExpenses from './components/InputExpenses.svelte';
  import ExpensesHistory from './components/ExpensesHistory.svelte';

  let activeTab = 'input';

  function switchTab(tab) {
    activeTab = tab;
  }
</script>

<main>
  {#if activeTab === 'input'}
    <InputExpenses />
  {:else}
    <ExpensesHistory />
  {/if}
</main>

<nav class="bottom-nav">
  <button 
    class="nav-item" 
    class:active={activeTab === 'input'}
    on:click={() => switchTab('input')}
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <line x1="12" y1="5" x2="12" y2="19"></line>
      <line x1="5" y1="12" x2="19" y2="12"></line>
    </svg>
    <span>Input Expenses</span>
  </button>
  <button 
    class="nav-item" 
    class:active={activeTab === 'history'}
    on:click={() => switchTab('history')}
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
      <line x1="16" y1="2" x2="16" y2="6"></line>
      <line x1="8" y1="2" x2="8" y2="6"></line>
      <line x1="3" y1="10" x2="21" y2="10"></line>
    </svg>
    <span>History</span>
  </button>
</nav>

<style>
  main {
    flex: 1;
    padding: 1rem;
    max-width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    min-height: 0;
    height: 0; /* Force flex child to respect parent height */
    -webkit-overflow-scrolling: touch;
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
    padding: 0.5rem 0;
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
  }

  .nav-item svg {
    width: 24px;
    height: 24px;
  }

  .nav-item.active {
    color: var(--primary-color);
  }

  .nav-item:active {
    opacity: 0.7;
  }
</style>

