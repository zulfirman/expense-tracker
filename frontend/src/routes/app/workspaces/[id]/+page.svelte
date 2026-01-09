<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { workspace } from '$lib/stores/workspace';
  import { get } from 'svelte/store';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';
  import { theme } from '$lib/stores/theme';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import WorkspaceModal from '$lib/components/WorkspaceModal.svelte';
  import CurrentWorkspaceBadge from '$lib/components/CurrentWorkspaceBadge.svelte';
  import Swal from "sweetalert2";

  let workspaceData = null;
  let loading = true;
  let showEditModal = false;

  $: pageCode = getPageCode($page.url.pathname);
  $: workspaceId = $page.params.id ? Number($page.params.id) : null;

  onMount(async () => {
    // Ensure theme is initialized
    theme.init();
    
    const ok = await requireAuthWithSleep();
    if (!ok) return;

    if (!workspaceId) {
      goto('/app/profile');
      return;
    }

    await loadWorkspace();
  });

  async function loadWorkspace() {
    loading = true;
    try {
      workspaceData = await workspace.getById(workspaceId);
    } catch (error) {
      // Workspace not found or error
      goto('/app/profile');
    } finally {
      loading = false;
    }
  }

  function openEditModal() {
    showEditModal = true;
  }

  function closeEditModal() {
    showEditModal = false;
  }

  function handleWorkspaceSaved() {
    closeEditModal();
    loadWorkspace();
    // Refresh workspace list
    workspace.init();
  }

  function handleManageCategories() {
    goto(`/app/workspaces/${workspaceId}/categories`);
  }

  function handleManageCurrency() {
    goto(`/app/workspaces/${workspaceId}/currency`);
  }

  let showDeleteModal = false;
  let deleteConfirmText = '';
  let deleting = false;

  function openDeleteModal() {
    showDeleteModal = true;
    deleteConfirmText = '';
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    deleteConfirmText = '';
  }

  async function handleDeleteWorkspace() {
    if (deleteConfirmText !== 'Delete') {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Invalid Confirmation',
          text: 'You must type "Delete" to confirm',
        });
      }, 50);
      return;
    }

    // Check if this is the current workspace
    const wsState = get(workspace);
    if (wsState.currentId === workspaceId) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Cannot Delete Active Workspace',
          text: 'Please switch to another workspace before deleting this one.',
        });
      }, 50);
      closeDeleteModal();
      return;
    }

    deleting = true;
    try {
      await api.delete(`/workspaces/${workspaceId}`);
      
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Workspace Deleted',
          text: 'The workspace has been permanently deleted',
          timer: 2000,
          showConfirmButton: false,
        });
      }, 50);
      
      // Refresh workspace list and redirect to profile
      await workspace.init();
      goto('/app/profile');
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Delete Failed',
          text: error.response?.data?.message || 'Failed to delete workspace',
        });
      }, 50);
    } finally {
      deleting = false;
      closeDeleteModal();
    }
  }
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title={workspaceData?.name || 'Workspace'}
    subtitle="Manage workspace settings and preferences."
    pageCode={pageCode}
  />
  <CurrentWorkspaceBadge />

  {#if loading}
    <div class="flex justify-center py-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if workspaceData}
    <!-- Workspace Info -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body">
        <div class="flex items-start justify-between gap-4 mb-4">
          <div class="flex-1">
            <h2 class="text-2xl font-bold mb-2">{workspaceData.name}</h2>
            {#if workspaceData.description}
              <p class="text-sm text-base-content/70 whitespace-pre-wrap">
                {workspaceData.description}
              </p>
            {:else}
              <p class="text-sm text-base-content/50 italic">No description</p>
            {/if}
            <div class="mt-3 text-xs text-base-content/60">
              Created: {new Date(workspaceData.createdAt).toLocaleDateString()}
            </div>
          </div>
          <button class="btn btn-primary btn-sm" on:click={openEditModal}>
            Edit
          </button>
        </div>
      </div>
    </div>

    <!-- Workspace Settings -->
    <div class="card bg-base-100 shadow-xl border border-base-300">
      <div class="card-body space-y-4">
        <h3 class="text-lg font-semibold">Workspace Settings</h3>
        <p class="text-xs text-base-content/70">
          Manage categories, currency, and other settings for this workspace.
        </p>

        <div class="grid gap-3 md:grid-cols-2">
          <button
            class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base"
            on:click={handleManageCategories}
          >
            <span class="text-left">
              <span class="block font-semibold">Manage Categories</span>
              <span class="block text-xs text-base-content/70">
                Manage and reorder your income & expense categories.
              </span>
            </span>
            <span class="text-lg">›</span>
          </button>

          <button
            class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base"
            on:click={handleManageCurrency}
          >
            <span class="text-left">
              <span class="block font-semibold">Currency & Quick Amount</span>
              <span class="block text-xs text-base-content/70">
                Configure display currency and quick amount shortcuts.
              </span>
            </span>
            <span class="text-lg">›</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Danger Zone -->
    <div class="card bg-base-100 shadow-xl border border-error/20">
      <div class="card-body space-y-4">
        <div>
          <h3 class="text-lg font-semibold text-error">Danger Zone</h3>
          <p class="text-xs text-base-content/70">
            Permanently delete this workspace and all its data. This action cannot be undone.
          </p>
        </div>
        <button
          class="btn btn-error btn-sm w-full sm:w-auto"
          on:click={openDeleteModal}
          disabled={loading || (workspaceData && get(workspace).currentId === workspaceId)}
          title={workspaceData && get(workspace).currentId === workspaceId ? 'Cannot delete the active workspace. Switch to another workspace first.' : ''}
        >
          Delete Workspace
        </button>
        {#if workspaceData && get(workspace).currentId === workspaceId}
          <p class="text-xs text-error mt-2">
            ⚠️ You cannot delete the active workspace. Switch to another workspace first.
          </p>
        {/if}
      </div>
    </div>
  {/if}
</PageContainer>

<WorkspaceModal
  bind:open={showEditModal}
  workspaceId={workspaceId}
  initialName={workspaceData?.name || ''}
  initialDescription={workspaceData?.description || ''}
  on:close={closeEditModal}
  on:saved={handleWorkspaceSaved}
/>

<!-- Delete Workspace Modal -->
{#if showDeleteModal}
  <div class="modal modal-open z-[3000]" on:click={closeDeleteModal}>
    <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-mono text-base-content/30 opacity-50">{pageCode}</span>
            <p class="text-xs uppercase tracking-wide text-base-content/60">Danger Zone</p>
          </div>
          <h3 class="text-xl font-bold text-error">Delete Workspace</h3>
        </div>
        <button class="btn btn-ghost btn-sm" on:click={closeDeleteModal}>✕</button>
      </div>

      <div class="alert alert-error mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <h3 class="font-bold">Warning!</h3>
          <div class="text-xs">This action cannot be undone. All data in this workspace will be permanently deleted.</div>
        </div>
      </div>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Type "Delete" to confirm</legend>
        <input
          id="delete-confirm-input"
          type="text"
          bind:value={deleteConfirmText}
          placeholder="Type 'Delete' here"
          class="input input-bordered w-full border-2"
          disabled={deleting}
          on:keydown={(e) => { if (e.key === 'Enter' && !deleting && deleteConfirmText === 'Delete') { handleDeleteWorkspace(); } }}
          autofocus
        />
      </fieldset>

      <div class="modal-action">
        <button class="btn btn-soft" on:click={closeDeleteModal} disabled={deleting}>
          Cancel
        </button>
        <button 
          class="btn btn-error" 
          on:click={handleDeleteWorkspace} 
          disabled={deleting || deleteConfirmText !== 'Delete'}
        >
          {#if deleting}
            <span class="loading loading-spinner loading-sm"></span>
            Deleting...
          {:else}
            Delete Workspace
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

