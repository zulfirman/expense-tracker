<script>
  import '../app.css';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import { onMount } from 'svelte';

  let showProfileMenu = false;
  let showAccountMenu = false;

  onMount(() => {
    theme.init();
    accounts.init();
    auth.init();
  });

  $: isAuthenticated = $auth.isAuthenticated;
  $: currentUser = $auth.user;
  $: accountList = $accounts.accounts;
  $: currentAccountId = $accounts.currentAccountId;

  function toggleProfileMenu() {
    showProfileMenu = !showProfileMenu;
  }

  function handleProfileClick() {
    showProfileMenu = false;
    goto('/profile');
  }

  function handleLogout() {
    auth.logout();
    showProfileMenu = false;
    goto('/login');
  }

  function toggleAccountMenu() {
    showAccountMenu = !showAccountMenu;
  }

  function handleSwitchAccount(accountId) {
    auth.switchAccount(accountId);
    showAccountMenu = false;
    // Reload page to refresh data
    window.location.reload();
  }

  function handleRemoveAccount(accountId, e) {
    e.stopPropagation();
    Swal.fire({
      icon: 'warning',
      title: 'Remove Account?',
      text: 'This will remove the account from your list. You can add it back by logging in again.',
      showCancelButton: true,
      confirmButtonText: 'Remove',
      cancelButtonText: 'Cancel',
      zIndex: 9999
    }).then((result) => {
      if (result.isConfirmed) {
        accounts.removeAccount(accountId);
        if (accountId === currentAccountId) {
          // If removed current account, redirect to login
          auth.logout();
          goto('/login');
        }
      }
    });
  }
</script>

{#if isAuthenticated || $page.url.pathname === '/login' || $page.url.pathname === '/signup'}
<div class="app-container">
  <header class="top-header">
    <div class="header-actions">
      <button class="theme-toggle" on:click={() => theme.toggle()} title="Toggle Dark Mode">
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
      {#if isAuthenticated}
        {#if accountList.length > 1}
          <div class="account-menu-container">
            <button class="account-toggle" on:click={toggleAccountMenu} title="Switch Account">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              </svg>
              <span class="account-count">{accountList.length} Accounts</span>
            </button>
            {#if showAccountMenu}
              <div class="account-menu">
                <div class="account-menu-header">Switch Account</div>
                {#each accountList as account}
                  <div 
                    class="account-item" 
                    class:active={account.id === currentAccountId}
                    on:click={() => handleSwitchAccount(account.id)}
                  >
                    <div class="account-info">
                      <div class="account-name">{account.name}</div>
                      <div class="account-email">{account.email}</div>
                    </div>
                    {#if account.id === currentAccountId}
                      <span class="current-badge">Current</span>
                    {/if}
                    <button 
                      class="remove-account-btn"
                      on:click={(e) => handleRemoveAccount(account.id, e)}
                      title="Remove Account"
                    >
                      ×
                    </button>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
        <div class="profile-menu-container">
          <button class="profile-toggle" on:click={toggleProfileMenu} title="Profile">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            {#if currentUser}
              <span class="user-name">{currentUser.name}</span>
            {/if}
          </button>
          {#if showProfileMenu}
            <div class="profile-menu">
              <a href="/categories" class="profile-menu-item" on:click={() => showProfileMenu = false}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="8" y1="6" x2="21" y2="6"></line>
                  <line x1="8" y1="12" x2="21" y2="12"></line>
                  <line x1="8" y1="18" x2="21" y2="18"></line>
                  <line x1="3" y1="6" x2="3.01" y2="6"></line>
                  <line x1="3" y1="12" x2="3.01" y2="12"></line>
                  <line x1="3" y1="18" x2="3.01" y2="18"></line>
                </svg>
                Categories
              </a>
              <button class="profile-menu-item" on:click={handleProfileClick}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                Profile
              </button>
              <button class="profile-menu-item" on:click={handleLogout}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                  <polyline points="16 17 21 12 16 7"></polyline>
                  <line x1="21" y1="12" x2="9" y2="12"></line>
                </svg>
                Logout
              </button>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </header>

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
    <a 
      href="/budget"
      data-sveltekit-preload-data="hover"
      class="nav-item" 
      class:active={$page.url.pathname === '/budget'}
      on:click|preventDefault={() => goto('/budget')}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="1" x2="12" y2="23"></line>
        <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
      </svg>
      <span>Budget</span>
    </a>
  </nav>
</div>
{/if}

<style>
  .app-container {
    height: 100vh;
    max-height: 100vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
  }

  .top-header {
    position: fixed;
    top: 0;
    right: 0;
    padding: 1rem;
    z-index: 1001;
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .theme-toggle {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: 1px solid var(--border);
    background: var(--surface);
    color: var(--text-primary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .theme-toggle:hover {
    background: var(--background);
    transform: scale(1.05);
  }

  .theme-toggle svg {
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

  .profile-menu-container {
    position: relative;
  }

  .profile-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    border: 1px solid var(--border);
    background: var(--surface);
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.875rem;
  }

  .profile-toggle:hover {
    background: var(--background);
  }

  .user-name {
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .profile-menu {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    min-width: 150px;
    overflow: hidden;
  }

  .profile-menu-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border: none;
    background: none;
    color: var(--text-primary);
    cursor: pointer;
    transition: background 0.2s;
    font-size: 0.875rem;
    text-align: left;
  }

  .profile-menu-item:hover {
    background: var(--background);
  }

  .profile-menu-item svg {
    flex-shrink: 0;
  }

  .account-menu-container {
    position: relative;
  }

  .account-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    border: 1px solid var(--border);
    background: var(--surface);
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.875rem;
  }

  .account-toggle:hover {
    background: var(--background);
  }

  .account-count {
    font-size: 0.75rem;
  }

  .account-menu {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    min-width: 250px;
    max-width: 300px;
    overflow: hidden;
  }

  .account-menu-header {
    padding: 0.75rem 1rem;
    font-weight: 600;
    color: var(--text-primary);
    border-bottom: 1px solid var(--border);
    font-size: 0.875rem;
  }

  .account-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background 0.2s;
    position: relative;
  }

  .account-item:hover {
    background: var(--background);
  }

  .account-item.active {
    background: var(--background);
  }

  .account-info {
    flex: 1;
    min-width: 0;
  }

  .account-name {
    font-weight: 500;
    color: var(--text-primary);
    font-size: 0.875rem;
    margin-bottom: 0.25rem;
  }

  .account-email {
    font-size: 0.75rem;
    color: var(--text-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .current-badge {
    font-size: 0.75rem;
    color: var(--primary-color);
    font-weight: 500;
    margin-right: 0.5rem;
  }

  .remove-account-btn {
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0.25rem 0.5rem;
    line-height: 1;
    opacity: 0.7;
    transition: all 0.2s;
  }

  .remove-account-btn:hover {
    opacity: 1;
    color: var(--danger);
  }

  .profile-menu-item {
    text-decoration: none;
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
    padding-top: 60px;
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
