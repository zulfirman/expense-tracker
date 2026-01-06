<script>
  import '../app.css';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import { currency } from '$lib/stores/currency';
  import { onMount } from 'svelte';
  import Swal from 'sweetalert2';

  let showProfileMenu = false;
  let showAccountMenu = false;
  let showAddAccountModal = false;
  let addAccountEmail = '';
  let addAccountPassword = '';
  let addAccountLoading = false;
  let authInitialized = false;

  onMount(async () => {
    theme.init();
    accounts.init();
    auth.init();
    await currency.init();
    authInitialized = true;
    
    // Check authentication and redirect if needed (only after initialization)
    const pathname = $page.url.pathname;
    const publicRoutes = ['/login', '/signup'];
    
    if (!$auth.isAuthenticated && !publicRoutes.includes(pathname) && pathname !== '/') {
      goto('/login');
    }
    
    if ($auth.isAuthenticated && publicRoutes.includes(pathname)) {
      goto('/expenses');
    }
    
    // Close profile and account menus when clicking outside
    function handleClickOutside(event) {
      const profileContainer = event.target.closest('.profile-menu-container');
      const accountContainer = event.target.closest('.account-menu-container');
      if (!profileContainer && showProfileMenu) {
        showProfileMenu = false;
      }
      if (!accountContainer && showAccountMenu) {
        showAccountMenu = false;
      }
    }
    
    document.addEventListener('click', handleClickOutside);
    
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  });

  $: isAuthenticated = $auth.isAuthenticated;
  
  // Watch for auth changes and redirect (only after initialization)
  $: if (authInitialized && !isAuthenticated && $page && $page.url) {
    const pathname = $page.url.pathname;
    const publicRoutes = ['/login', '/signup'];
    if (!publicRoutes.includes(pathname) && pathname !== '/') {
      goto('/login');
    }
  }
  
  // Redirect authenticated users away from login/signup (only after initialization)
  // Only redirect if explicitly on login/signup pages, never redirect from other pages
  $: if (authInitialized && isAuthenticated && $page && $page.url) {
    const pathname = $page.url.pathname;
    // Only redirect from login/signup pages, don't interfere with other pages like /budget
    if (pathname === '/login' || pathname === '/signup') {
      goto('/expenses');
    }
  }
  $: currentUser = $auth.user;
  $: accountList = $accounts.accounts.sort((a, b) => {
    // Active user first
    if (a.id === $accounts.currentAccountId) return -1;
    if (b.id === $accounts.currentAccountId) return 1;
    // Then sort alphabetically by name
    return a.name.localeCompare(b.name);
  });
  $: currentAccountId = $accounts.currentAccountId;

  function toggleProfileMenu(event) {
    event.stopPropagation();
    showProfileMenu = !showProfileMenu;
  }

  function handleProfileClick() {
    showProfileMenu = false;
    goto('/profile');
  }

  function handlePreferencesClick() {
    showProfileMenu = false;
    goto('/preferences');
  }

  function handleChangePasswordClick() {
    showProfileMenu = false;
    goto('/change-password');
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
    if (accountId === currentAccountId) {
      // Clicking current account: just close menu, no reload
      showAccountMenu = false;
      return;
    }
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
      reverseButtons: true,
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

  function openAddAccountModal() {
    showAddAccountModal = true;
    showProfileMenu = false;
    addAccountEmail = '';
    addAccountPassword = '';
  }

  function closeAddAccountModal() {
    showAddAccountModal = false;
    addAccountEmail = '';
    addAccountPassword = '';
  }

  async function handleAddAccount() {
    if (!addAccountEmail || !addAccountPassword) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Fields',
        text: 'Please enter both email and password',
        zIndex: 9999
      });
      return;
    }

    addAccountLoading = true;
    try {
      const response = await fetch('/api/apps/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: addAccountEmail,
          password: addAccountPassword
        })
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Login failed');
      }

      const data = await response.json();
      
      // Add account to accounts store (this is explicitly adding a new account)
      // Don't clear existing accounts when using "Add Account" modal
      auth.login(data.user, data.token, data.refreshToken, false);
      
      Swal.fire({
        icon: 'success',
        title: 'Account Added!',
        text: `Logged in as ${data.user.name}`,
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });

      closeAddAccountModal();
      // Reload to refresh data
      window.location.reload();
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Login Failed',
        text: error.message || 'Invalid email or password',
        zIndex: 9999
      });
    } finally {
      addAccountLoading = false;
    }
  }

  function handleAddAccountKeyDown(e) {
    if (e.key === 'Enter') {
      handleAddAccount();
    }
  }
</script>

{#if isAuthenticated || $page.url.pathname === '/login' || $page.url.pathname === '/signup' || $page.url.pathname === '/'}
<div class="app-container">
  <header class="top-header">
    <div class="header-actions">
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
              <button class="profile-menu-item" on:click={handleProfileClick}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                Profile
              </button>
              <button class="profile-menu-item" on:click={handlePreferencesClick}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M12 1v6m0 6v6m9-9h-6m-6 0H3"></path>
                </svg>
                Preferences
              </button>
              <button class="profile-menu-item" on:click={openAddAccountModal}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
                Add Account
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

{#if showAddAccountModal}
  <div class="modal-backdrop" on:click={closeAddAccountModal}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Add Account</h2>
        <button class="close-btn" on:click={closeAddAccountModal}>×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label for="add-account-email">Email</label>
          <input
            id="add-account-email"
            type="email"
            bind:value={addAccountEmail}
            placeholder="your@email.com"
            class="form-input"
            on:keydown={handleAddAccountKeyDown}
            disabled={addAccountLoading}
          />
        </div>
        <div class="form-group">
          <label for="add-account-password">Password</label>
          <input
            id="add-account-password"
            type="password"
            bind:value={addAccountPassword}
            placeholder="Enter password"
            class="form-input"
            on:keydown={handleAddAccountKeyDown}
            disabled={addAccountLoading}
          />
        </div>
        <div class="button-group">
          <button class="btn btn-secondary" on:click={closeAddAccountModal} disabled={addAccountLoading}>Cancel</button>
          <button class="btn btn-primary" on:click={handleAddAccount} disabled={addAccountLoading}>
            {#if addAccountLoading}
              <span class="spinner"></span> Logging in...
            {:else}
              Add Account
            {/if}
          </button>
        </div>
      </div>
    </div>
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


  .profile-menu-container {
    position: relative;
  }

  .profile-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 0.75rem;
    border: 1.5px solid var(--primary-color);
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
    font-weight: 700;
    color: var(--text-primary);
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

  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3000;
    padding: 1rem;
  }

  .modal-content {
    background: var(--surface);
    border-radius: 1rem;
    width: 100%;
    max-width: 400px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--border);
  }

  .modal-header h2 {
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 0;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 2rem;
    color: var(--text-secondary);
    cursor: pointer;
    line-height: 1;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    color: var(--text-primary);
  }

  .modal-body {
    padding: 1.5rem;
  }

  .modal-body .form-group {
    margin-bottom: 1.5rem;
  }

  .modal-body .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .modal-body .form-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
    background: var(--background);
    color: var(--text-primary);
  }

  .modal-body .form-input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .modal-body .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .modal-body .button-group {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .modal-body .btn {
    flex: 1;
    padding: 0.875rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .modal-body .btn-primary {
    background-color: var(--primary-color);
    color: white;
  }

  .modal-body .btn-primary:hover:not(:disabled) {
    background-color: #4338ca;
  }

  .modal-body .btn-secondary {
    background-color: var(--surface);
    color: var(--text-primary);
    border: 1px solid var(--border);
  }

  .modal-body .btn-secondary:hover:not(:disabled) {
    background-color: var(--background);
  }

  .modal-body .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .modal-body .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 0.6s linear infinite;
    margin-right: 0.5rem;
    vertical-align: middle;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
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
