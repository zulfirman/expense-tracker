<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { getInitials } from '$lib/utils/initials';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';
  import { workspace } from '$lib/stores/workspace';
  import CurrentWorkspaceBadge from '$lib/components/CurrentWorkspaceBadge.svelte';
  import WorkspaceModal from '$lib/components/WorkspaceModal.svelte';

  let name = '';
  let email = '';
  let loading = false;
  let saving = false;
  let showEditModal = false;
  let editName = '';
  let editEmail = '';
  let workspaceSearch = '';
  let showWorkspaceModal = false;
  let editingWorkspaceId = null;
  let editingWorkspaceName = '';
  let editingWorkspaceDescription = '';

  $: pageCode = getPageCode($page.url.pathname);
  $: currentUser = $auth.user;
  $: initials = getInitials(name || currentUser?.name || '');
  $: wsState = $workspace;
  $: workspaces = wsState.list || [];
  $: currentWorkspaceId = wsState.currentId;
  $: filteredWorkspaces =
    workspaceSearch.trim().length === 0
      ? workspaces
      : workspaces.filter((ws) =>
          ws.name.toLowerCase().includes(workspaceSearch.trim().toLowerCase())
        );

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;

    await loadProfile();
    await workspace.init();
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
    if (saving) return;

    if (!editName || !editEmail) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Fields',
          text: 'Please fill in all fields',

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

        });
      }, 50);
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Update Failed',
          text: error.response?.data?.message || 'Failed to update profile',

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

  function openCreateWorkspaceModal() {
    editingWorkspaceId = null;
    editingWorkspaceName = '';
    editingWorkspaceDescription = '';
    showWorkspaceModal = true;
  }

  async function handleSwitchWorkspace(ws) {
    if (ws.id === currentWorkspaceId) {
      goto('/app/expenses');
      return;
    }

    const result = await Swal.fire({
      icon: 'question',
      title: 'Switch Workspace?',
      text: `Switch to "${ws.name}"? All data will change to this workspace.`,
      showCancelButton: true,
      confirmButtonText: 'Yes, Switch',
      cancelButtonText: 'Cancel',
      reverseButtons: true,

    });

    if (result.isConfirmed) {
      workspace.setCurrent(ws.id);
      goto('/app/expenses');
      setTimeout(() => {
          Swal.fire({
              icon: 'success',
              title: 'Switched to '+ws.name,
              timer: 1500,
              showConfirmButton: false,
          });
      }, 50);
    }
  }

  function openEditWorkspaceModal(ws) {
    editingWorkspaceId = ws.id;
    editingWorkspaceName = ws.name;
    editingWorkspaceDescription = ws.description || '';
    showWorkspaceModal = true;
  }

  function closeWorkspaceModal() {
    showWorkspaceModal = false;
    editingWorkspaceId = null;
    editingWorkspaceName = '';
    editingWorkspaceDescription = '';
  }

  function handleWorkspaceSaved() {
    closeWorkspaceModal();
    // Reload workspaces
    workspace.init();
  }
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="Profile"
    subtitle="Your account overview and workspaces."
    pageCode={pageCode}
  />
  <CurrentWorkspaceBadge />

  {#if loading}
    <div class="flex justify-center py-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <!-- Profile header -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body">
        <div class="flex flex-col sm:flex-row items-center sm:items-start gap-4">
          <div class="avatar placeholder">
            <div class="bg-primary text-primary-content rounded-full w-24 h-24 flex items-center justify-center">
              <span class="text-3xl font-bold">{initials}</span>
            </div>
          </div>
          <div class="flex-1 min-w-0 text-center sm:text-left">
            <h2 class="text-2xl font-bold truncate">{name}</h2>
            <p class="text-sm text-base-content/70 truncate">{email}</p>
            <div class="mt-2 flex flex-wrap gap-2 justify-center sm:justify-start text-xs text-base-content/60">
              <span>Workspaces: {workspaces.length}</span>
              {#if currentWorkspaceId}
                <span>
                  Active: {workspaces.find((w) => w.id === currentWorkspaceId)?.name}
                </span>
              {/if}
            </div>
          </div>
          <div class="flex flex-col gap-2 w-full sm:w-auto">
            <button class="btn btn-primary btn-sm w-full" on:click={openEditModal}>
              Edit Profile
            </button>
            <button
              class="btn btn-soft btn-sm w-full"
              on:click={() => goto('/app/preferences')}
            >
              Open Preferences
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Workspaces section -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body space-y-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold">Workspaces</h3>
            <p class="text-xs text-base-content/70">
              Each workspace keeps its own expenses, income, categories, and budgets.
            </p>
          </div>
          <button class="btn btn-primary btn-sm whitespace-nowrap" on:click={openCreateWorkspaceModal}>
            + New
          </button>
        </div>

        <div class="form-control">
          <input
            type="search"
            class="input input-bordered w-full border-2"
            placeholder="Search workspace..."
            bind:value={workspaceSearch}
          />
        </div>

        {#if workspaces.length === 0}
          <div class="py-8 text-center text-sm text-base-content/70">
            No workspaces yet. Create your first one to start tracking.
          </div>
        {:else}
          <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
            {#each filteredWorkspaces as ws}
              <div
                class={`group relative aspect-square rounded-2xl border border-base-300 bg-base-100 overflow-hidden flex flex-col items-center justify-center text-center px-2 transition-all duration-150 hover:border-primary hover:bg-base-200 ${
                  ws.id === currentWorkspaceId ? 'bg-primary/10 border-primary' : ''
                }`}
              >
                <button
                  type="button"
                  class="absolute inset-0 z-10"
                  on:click={() => {
                    goto(`/app/workspaces/${ws.id}`);
                  }}
                ></button>
                <div class="space-y-1 w-full">
                  <div class="font-semibold text-xs sm:text-sm truncate">
                    {ws.name}
                  </div>
                  <div class="text-[0.6rem] text-base-content/60">
                    {#if ws.createdAt}
                      {new Date(ws.createdAt).toLocaleDateString()}
                    {:else}
                      &nbsp;
                    {/if}
                  </div>
                  {#if ws.id === currentWorkspaceId}
                    <div class="mt-1">
                      <span class="badge badge-primary badge-xs px-3 py-2 rounded-full">
                        Active
                      </span>
                    </div>
                  {/if}
                </div>
                {#if ws.id !== currentWorkspaceId}
                  <button
                    class="absolute top-1 right-1 btn btn-xs btn-primary z-20"
                    on:click|stopPropagation={() => handleSwitchWorkspace(ws)}
                    title="Switch to this workspace"
                  ><!--opacity-0 group-hover:opacity-100-->
                    Switch
                  </button>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
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
            <span class="text-xs font-mono text-base-content/30 opacity-50">{pageCode}</span>
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
        <button class="btn btn-soft" on:click={closeEditModal} disabled={saving}>
          Cancel
        </button>
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

<WorkspaceModal
  bind:open={showWorkspaceModal}
  workspaceId={editingWorkspaceId}
  initialName={editingWorkspaceName}
  initialDescription={editingWorkspaceDescription}
  on:close={closeWorkspaceModal}
  on:saved={handleWorkspaceSaved}
/>

<!-- Edit Profile Modal -->
{#if showEditModal}
  <div class="modal modal-open z-[3000]" on:click={closeEditModal}>
    <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-mono text-base-content/30 opacity-50">{pageCode}</span>
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
