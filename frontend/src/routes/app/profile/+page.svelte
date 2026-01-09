<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import ActionButton from '$lib/components/ActionButton.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { getInitials } from '$lib/utils/initials';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';

  let name = '';
  let email = '';
  let loading = false;
  let saving = false;
  let showEditModal = false;
  let editName = '';
  let editEmail = '';
  let showAddAccountModal = false;
  let addAccountEmail = '';
  let addAccountPassword = '';
  let addAccountLoading = false;
  
  $: pageCode = getPageCode($page.url.pathname);
  $: currentUser = $auth.user;
  $: accountList = $accounts.accounts.sort((a, b) => {
    if (a.id === $accounts.currentAccountId) return -1;
    if (b.id === $accounts.currentAccountId) return 1;
    return a.name.localeCompare(b.name);
  });
  $: currentAccountId = $accounts.currentAccountId;
  $: initials = getInitials(name || currentUser?.name || '');

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;

    await loadProfile();
  });

  async function loadProfile() {
    loading = true;
    try {
      const response = await api.get('/auth/profile');
      name = response.data.name;
      email = response.data.email;
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to load profile',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function openEditModal() {
    editName = name;
    editEmail = email;
    showEditModal = true;
  }

  function closeEditModal() {
    showEditModal = false;
    editName = '';
    editEmail = '';
  }

  async function handleSave() {
    if (saving) return; // Prevent double submission
    
    if (!editName || !editEmail) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Fields',
          text: 'Please fill in all fields',
          zIndex: 9999
        });
      }, 50);
      return;
    }

    saving = true;
    try {
      const response = await api.put('/auth/profile', {
        name: editName,
        email: editEmail
      });

      // Update auth store
      auth.login(response.data, $auth.token);
      name = editName;
      email = editEmail;
      closeEditModal();
      
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Profile Updated',
          text: 'Your profile has been updated successfully',
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Update Failed',
          text: error.response?.data?.message || 'Failed to update profile',
          zIndex: 9999
        });
      }, 50);
    } finally {
      saving = false;
    }
  }

  function handleEditModalKeyDown(e) {
    if (e.key === 'Enter' && !saving) {
      e.preventDefault();
      setTimeout(() => {
        handleSave();
      }, 50);
    }
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
      zIndex: 9999
    }).then((result) => {
      if (result.isConfirmed) {
        accounts.removeAccount(accountId);
        if (accountId === currentAccountId) {
          auth.logout();
          goto('/app/login');
        } else {
          loadProfile();
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

  function handleLogout() {
    Swal.fire({
      icon: 'question',
      title: 'Logout?',
      text: 'Are you sure you want to logout?',
      showCancelButton: true,
      confirmButtonText: 'Logout',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      zIndex: 9999
    }).then((result) => {
      if (result.isConfirmed) {
        auth.logout();
        goto('/app/login');
      }
    });
  }

  function handleAddAccountKeyDown(e) {
    if (e.key === 'Enter' && !addAccountLoading) {
      e.preventDefault();
      setTimeout(() => {
        handleAddAccount();
      }, 50);
    }
  }
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="Profile"
    subtitle="Manage your account and preferences."
    pageCode={pageCode}
  />

  {#if loading}
    <div class="flex justify-center py-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <!-- Profile Card -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body">
        <div class="flex flex-col items-center gap-4 mb-4">
          <div class="avatar placeholder">
            <div class="bg-primary text-primary-content rounded-full w-24 h-24 text-3xl font-bold flex items-center justify-center">
              <span>{initials}</span>
            </div>
          </div>
          <div class="text-center">
            <h2 class="text-2xl font-bold">{name}</h2>
            <p class="text-sm text-base-content/70 mt-1">{email}</p>
          </div>
          <button class="btn btn-primary" on:click={openEditModal}>
            Edit Profile
          </button>
        </div>

        <!-- Account Info (Readonly) -->
        <div class="divider"></div>
        <div class="space-y-3">
          <div>
            <p class="text-xs uppercase tracking-wide text-base-content/60 mb-1">Name</p>
            <p class="text-base font-medium">{name}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-wide text-base-content/60 mb-1">Email</p>
            <p class="text-base font-medium">{email}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Preferences Section -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body">
        <h3 class="text-lg font-semibold mb-2">Preferences</h3>
        <p class="text-sm text-base-content/70 mb-4">
          Customize your app settings including categories, currency, quick amounts, and password. Manage how the app looks and behaves for your account.
        </p>
        <div class="space-y-2">
          <a href="/app/preferences" class="btn btn-soft justify-between w-full normal-case">
            <span>Open Preferences</span>
            <span>›</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Account Management -->
    {#if accountList.length > 1}
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h3 class="text-lg font-semibold mb-4">Switch Account</h3>
          <div class="space-y-2">
            {#each accountList as account}
              <div class="flex items-center justify-between p-3 rounded-lg border border-base-300">
                <div class="flex flex-col">
                  <span class="font-medium text-sm">{account.name}</span>
                  <span class="text-xs text-base-content/70">{account.email}</span>
                </div>
                <div class="flex items-center gap-2">
                  {#if account.id === currentAccountId}
                    <span class="badge badge-primary badge-sm">Current</span>
                  {:else}
                    <button class="btn btn-xs btn-primary" on:click={() => handleSwitchAccount(account.id)}>
                      Switch
                    </button>
                  {/if}
                  <button class="btn btn-xs btn-soft text-error" on:click={(e) => handleRemoveAccount(account.id, e)}>
                    ×
                  </button>
                </div>
              </div>
            {/each}
            <button class="btn btn-primary btn-sm w-full mt-2" on:click={openAddAccountModal}>
              + Add Account
            </button>
          </div>
        </div>
      </div>
    {:else}
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h3 class="text-lg font-semibold mb-4">Account Management</h3>
          <button class="btn btn-primary w-full" on:click={openAddAccountModal}>
            + Add Account
          </button>
        </div>
      </div>
    {/if}

    <!-- Logout -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body">
        <button class="btn btn-error w-full" on:click={handleLogout}>
          Logout
        </button>
      </div>
    </div>
  {/if}
</PageContainer>

<!-- Edit Profile Modal -->
{#if showEditModal}
  <div class="modal modal-open z-[3000]" on:click={closeEditModal}>
    <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-mono text-base-content/30">{pageCode}</span>
            <p class="text-xs uppercase tracking-wide text-base-content/60">Profile</p>
          </div>
          <h3 class="text-xl font-bold">Edit Profile</h3>
        </div>
        <button class="btn btn-ghost btn-sm" on:click={closeEditModal}>✕</button>
      </div>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Name</legend>
        <input
          id="edit-name"
          type="text"
          bind:value={editName}
          placeholder="Your name"
          class="input input-bordered w-full border-2"
          disabled={saving}
          on:keydown={handleEditModalKeyDown}
        />
      </fieldset>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Email</legend>
        <input
          id="edit-email"
          type="email"
          bind:value={editEmail}
          placeholder="your@email.com"
          class="input input-bordered w-full border-2"
          disabled={saving}
          on:keydown={handleEditModalKeyDown}
        />
      </fieldset>

      <div class="modal-action">
        <button class="btn btn-soft" on:click={closeEditModal} disabled={saving}>Cancel</button>
        <button class="btn btn-primary" on:click={handleSave} disabled={saving}>
          {#if saving}
            <span class="loading loading-spinner loading-sm"></span>
            Saving...
          {:else}
            Save Changes
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Add Account Modal -->
{#if showAddAccountModal}
  <div class="modal modal-open z-[3000]" on:click={closeAddAccountModal}>
    <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-mono text-base-content/30">{pageCode}</span>
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
        <button class="btn btn-soft" on:click={closeAddAccountModal} disabled={addAccountLoading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleAddAccount} disabled={addAccountLoading}>
          {#if addAccountLoading}
            <span class="loading loading-spinner loading-sm"></span>
            Adding...
          {:else}
            Add Account
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}



