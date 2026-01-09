<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import Swal from 'sweetalert2';
  import * as helper from '$lib/utils/helper';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';
  import {capitalizeFirst} from "../../../lib/utils/helper.js";

  const availableThemes = ['cupcake', 'night','valentine','dracula'];
  let selectedTheme = 'cupcake';
  let showAddAccountModal = false;
  let addAccountEmail = '';
  let addAccountPassword = '';
  let addAccountLoading = false;
  
  $: pageCode = getPageCode($page.url.pathname);
  $: accountList = $accounts.accounts.sort((a, b) => {
    if (a.id === $accounts.currentAccountId) return -1;
    if (b.id === $accounts.currentAccountId) return 1;
    return a.name.localeCompare(b.name);
  });
  $: currentAccountId = $accounts.currentAccountId;

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;
    selectedTheme = $theme || 'cupcake';
  });

  $: if ($theme) {
    selectedTheme = $theme;
  }

  function handleThemeChange() {
    theme.setTheme(selectedTheme);
  }
  
  function handleSwitchAccount(accountId) {
    if (accountId === currentAccountId) return;
    auth.switchAccount(accountId);
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

    }).then((result) => {
      if (result.isConfirmed) {
        accounts.removeAccount(accountId);
        if (accountId === currentAccountId) {
          auth.logout();
          goto('/app/login');
        }
      }
    });
  }

  function openAddAccountModal() {
    showAddAccountModal = true;
    addAccountEmail = '';
    addAccountPassword = '';
  }

  function closeAddAccountModal() {
    showAddAccountModal = false;
    addAccountEmail = '';
    addAccountPassword = '';
  }

  async function handleAddAccount() {
    if (addAccountLoading) return;

    if (!addAccountEmail || !addAccountPassword) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Fields',
          text: 'Please enter both email and password',

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
      auth.login(data.user, data.token, data.refreshToken, false);

      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Account Added!',
          text: `Logged in as ${data.user.name}`,
          timer: 1500,
          showConfirmButton: false,

        });
      }, 50);

      closeAddAccountModal();
      window.location.reload();
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Login Failed',
          text: error.message || 'Invalid email or password',

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

  function handleLogout() {
    Swal.fire({
      icon: 'question',
      title: 'Logout?',
      text: 'Are you sure you want to logout?',
      showCancelButton: true,
      confirmButtonText: 'Logout',
      cancelButtonText: 'Cancel',
      reverseButtons: true,

    }).then((result) => {
      if (result.isConfirmed) {
        auth.logout();
        goto('/app/login');
      }
    });
  }
</script>

<div class="max-w-4xl mx-auto space-y-4">
  <PageHeader
    title="Preferences"
    subtitle="Manage how the app looks and behaves for your account."
    pageCode={pageCode}
  />
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body space-y-4">
      <!-- Navigation links -->
      <div class="grid gap-3 md:grid-cols-1">
        <a href="/app/preferences/change-password"
           class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base">
          <span class="text-left">
            <span class="block font-semibold">Change Password</span>
            <span class="block text-xs text-base-content/70">
              Update your account password securely.
            </span>
          </span>
          <span class="text-lg">›</span>
        </a>
      </div>

      <div class="alert alert-info">
        <span class="text-sm">
          <strong>Note:</strong> Category and currency management are now available in your workspace settings.
          Visit your workspace detail page to manage them.
        </span>
      </div>

      <!-- Appearance -->
      <div class="divider my-2"></div>

      <div class="space-y-2 max-w-md">
        <h3 class="font-semibold text-base">Appearance</h3>
        <p class="text-xs text-base-content/70">
          Choose a DaisyUI theme for the app.
        </p>
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Theme</legend>
          <select
            class="select select-bordered w-full border-2"
            bind:value={selectedTheme}
            on:change={handleThemeChange}
          >
            {#each availableThemes as t}
              <option value={t}>{helper.capitalizeFirst(t)}</option>
            {/each}
          </select>
        </fieldset>
      </div>
    </div>
  </div>
  
  <!-- Account & Security -->
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body space-y-4">
      <h3 class="text-lg font-semibold">Account & Security</h3>
      <p class="text-xs text-base-content/70">
        Manage linked accounts, switch between them, and logout securely.
      </p>

      {#if accountList.length > 0}
        <div class="space-y-2">
          {#each accountList as account}
            <div class="flex items-center justify-between p-3 rounded-lg border border-base-300">
              <div class="flex flex-col">
                <span class="font-medium text-sm">{account.name}</span>
                <span class="text-xs text-base-content/70 truncate">{account.email}</span>
              </div>
              <div class="flex items-center gap-2">
                {#if account.id === currentAccountId}
                  <span class="badge badge-primary badge-sm">Current</span>
                {:else}
                  <button
                    class="btn btn-xs btn-primary"
                    on:click={() => handleSwitchAccount(account.id)}
                  >
                    Switch
                  </button>
                {/if}
                <button
                  class="btn btn-xs btn-soft text-error"
                  on:click={(e) => handleRemoveAccount(account.id, e)}
                >
                  ×
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      <div class="flex flex-col sm:flex-row gap-2 mt-2">
        <button class="btn btn-soft flex-1" on:click={openAddAccountModal}>
          + Add Account
        </button>
        <button class="btn btn-error flex-1" on:click={handleLogout}>
          Logout
        </button>
      </div>
    </div>
  </div>

  {#if showAddAccountModal}
    <div class="modal modal-open z-[3000]" on:click={closeAddAccountModal}>
      <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
        <div class="flex items-start justify-between gap-3 mb-4">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xs font-mono text-base-content/30 opacity-50">{pageCode}</span>
              <p class="text-xs uppercase tracking-wide text-base-content/60">Accounts</p>
            </div>
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

        <div class="modal-action">
          <button class="btn btn-soft" on:click={closeAddAccountModal} disabled={addAccountLoading}>
            Cancel
          </button>
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
</div>