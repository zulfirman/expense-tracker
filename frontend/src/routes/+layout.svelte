<script>
  import '../app.css';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { afterNavigate } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import { currency } from '$lib/stores/currency';
  import { onMount } from 'svelte';
  import Swal from 'sweetalert2';
  import { getInitials } from '$lib/utils/initials';

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
    const publicRoutes = ['/app/login', '/app/signup'];
    
    if (!$auth.isAuthenticated && !publicRoutes.includes(pathname) && pathname !== '/' && !pathname.startsWith('/app/login') && !pathname.startsWith('/app/signup')) {
      goto('/app/login');
    }
    
    if ($auth.isAuthenticated && (pathname === '/app/login' || pathname === '/app/signup')) {
      goto('/app/expenses');
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
    const publicRoutes = ['/app/login', '/app/signup'];
    if (!publicRoutes.includes(pathname) && pathname !== '/') {
      goto('/app/login');
    }
  }
  
  // Redirect authenticated users away from login/signup (only after initialization)
  // Only redirect if explicitly on login/signup pages, never redirect from other pages
  $: if (authInitialized && isAuthenticated && $page && $page.url) {
    const pathname = $page.url.pathname;
    // Only redirect from login/signup pages, don't interfere with other pages like /app/preferences
    if (pathname === '/app/login' || pathname === '/app/signup') {
      goto('/app/expenses');
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
      goto('/app/profile');
  }

  function handlePreferencesClick() {
    showProfileMenu = false;
      goto('/app/preferences');
  }

  function handleLogout() {
    auth.logout();
    showProfileMenu = false;
    goto('/app/login');
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
          goto('/app/login');
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
    if (addAccountLoading) return; // Prevent double submission
    
    if (!addAccountEmail || !addAccountPassword) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Fields',
          text: 'Please enter both email and password',
          zIndex: 9999
        });
      }, 50);
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
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Account Added!',
          text: `Logged in as ${data.user.name}`,
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
      // Reload to refresh data
      window.location.reload();
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Login Failed',
          text: error.message || 'Invalid email or password',
          zIndex: 9999
        });
      }, 50);
    } finally {
      addAccountLoading = false;
    }
  }

  function handleAddAccountKeyDown(e) {
    if (e.key === 'Enter' && !addAccountLoading) {
      e.preventDefault();
      setTimeout(() => {
        handleAddAccount();
      }, 50);
    }
  }

  import { getPageCode as getPageCodeUtil } from '$lib/utils/pageCodes';

  $: currentPath = $page?.url?.pathname || '/';
  $: pageCode = getPageCodeUtil(currentPath);
  
  // Reset scroll position on route change
  afterNavigate(() => {
    const mainContent = document.querySelector('.main-content');
    if (mainContent) {
      mainContent.scrollTop = 0;
    }
  });
</script>

{#if isAuthenticated || $page.url.pathname === '/app/login' || $page.url.pathname === '/app/signup' || $page.url.pathname === '/' || $page.url.pathname.startsWith('/app/')}
<div class="app-container">

  <main class="main-content">
    <slot />
  </main>

  {#if isAuthenticated && $page.url.pathname !== '/app/login' && $page.url.pathname !== '/app/signup' && $page.url.pathname !== '/app/login' && $page.url.pathname !== '/app/signup'}
  <div class="dock dock-sm fixed left-0 right-0 bottom-0 px-3 pb-[calc(0.5rem+env(safe-area-inset-bottom))]">
    <button
      class:dock-active={$page.url.pathname === '/app/expenses' || $page.url.pathname === '/expenses' || $page.url.pathname === '/'}
      on:click={() => goto('/app/expenses')}
    >
      <svg class="size-[1.4em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
      <span class="dock-label">Expenses</span>
    </button>

    <button
      class:dock-active={$page.url.pathname === '/app/income' || $page.url.pathname === '/income'}
      on:click={() => goto('/app/income')}
    >
      <svg class="size-[1.4em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="19" x2="12" y2="5"></line>
        <polyline points="5 12 12 5 19 12"></polyline>
      </svg>
      <span class="dock-label">Income</span>
    </button>

    <button
      class:dock-active={$page.url.pathname === '/app/history' || $page.url.pathname === '/history'}
      on:click={() => goto('/app/history')}
    >
      <svg class="size-[1.4em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
        <line x1="16" y1="2" x2="16" y2="6"></line>
        <line x1="8" y1="2" x2="8" y2="6"></line>
        <line x1="3" y1="10" x2="21" y2="10"></line>
      </svg>
      <span class="dock-label">History</span>
    </button>

    <button
      class:dock-active={$page.url.pathname === '/app/budget' || $page.url.pathname === '/budget'}
      on:click={() => goto('/app/budget')}
    >
      <svg class="size-[1.4em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="1" x2="12" y2="23"></line>
        <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
      </svg>
      <span class="dock-label">Budget</span>
    </button>

    <button
      class:dock-active={$page.url.pathname === '/app/profile' || $page.url.pathname === '/profile'}
      on:click={() => goto('/app/profile')}
      title={$auth.user?.name || 'Profile'}
    >
      {#if $auth.user?.name}
        <div class="avatar placeholder">
          <div class="bg-primary text-primary-content rounded-full w-8 h-8 text-xs font-bold flex items-center justify-center">
            <span>{getInitials($auth.user.name)}</span>
          </div>
        </div>
      {:else}
        <svg class="size-[1.4em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
      {/if}
      <!--<span class="dock-label">Profile</span>-->
    </button>
  </div>
  {/if}
</div>
{/if}


{#if showAddAccountModal}
  <div class="modal modal-open z-[3000]" on:click={closeAddAccountModal}>
    <div class="modal-box max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-2 mb-3">
        <div>
          <p class="text-xs uppercase tracking-wide text-base-content/60">Accounts</p>
          <h3 class="text-xl font-bold">Add Account</h3>
        </div>
        <button class="btn btn-ghost btn-sm" on:click={closeAddAccountModal}>✕</button>
      </div>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Email</legend>
        <input
          id="add-account-email"
          type="email"
          bind:value={addAccountEmail}
          placeholder="your@email.com"
          class="input input-bordered w-full border-2"
          on:keydown={handleAddAccountKeyDown}
          disabled={addAccountLoading}
        />
      </fieldset>
      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Password</legend>
        <input
          id="add-account-password"
          type="password"
          bind:value={addAccountPassword}
          placeholder="Enter password"
          class="input input-bordered w-full border-2"
          on:keydown={handleAddAccountKeyDown}
          disabled={addAccountLoading}
        />
      </fieldset>

      <div class="flex justify-end gap-2">
        <button class="btn btn-soft" on:click={closeAddAccountModal} disabled={addAccountLoading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleAddAccount} disabled={addAccountLoading}>
          {#if addAccountLoading}
            <span class="loading loading-spinner loading-sm"></span>
            Logging in...
          {:else}
            Add Account
          {/if}
        </button>
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
    z-index: 800;
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  main.main-content {
    flex: 1;
    padding: 1rem;
    max-width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    min-height: 0;
    height: 0;
    -webkit-overflow-scrolling: touch;
    padding-bottom: 80px;
    padding-top: 20px;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  main.main-content::-webkit-scrollbar {
    display: none;
  }

  /* Dock styling */
  :global(.dock button) {
    border-radius: 0.75rem;
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
