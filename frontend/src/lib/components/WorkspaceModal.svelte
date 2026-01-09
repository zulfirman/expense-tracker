<script>
  import { createEventDispatcher } from 'svelte';
  import Swal from 'sweetalert2';
  import { workspace } from '$lib/stores/workspace';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { page } from '$app/stores';

  export let open = false;
  export let workspaceId = null; // null = create, number = edit
  export let initialName = '';
  export let initialDescription = '';
  export let required = false; // If true, user MUST create workspace (can't cancel)

  const dispatch = createEventDispatcher();

  let name = '';
  let description = '';
  let saving = false;

  $: pageCode = getPageCode($page.url.pathname);

  $: if (open) {
    name = initialName || '';
    description = initialDescription || '';
  }

  function close() {
    if (required) return; // Can't close if required
    open = false;
    name = '';
    description = '';
    dispatch('close');
  }

  async function handleSave() {
    if (saving) return;

    if (!name || !name.trim()) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Name',
          text: 'Please enter a workspace name',

        });
      }, 50);
      return;
    }

    saving = true;
    try {
      if (workspaceId) {
        // Update existing workspace
        await workspace.update(workspaceId, name.trim(), description.trim());
        setTimeout(() => {
          Swal.fire({
            icon: 'success',
            title: 'Workspace Updated',
            text: 'Workspace has been updated successfully',
            timer: 1500,
            showConfirmButton: false,

          });
        }, 50);
      } else {
        // Create new workspace
        const ws = await workspace.create(name.trim(), description.trim());
        setTimeout(() => {
          Swal.fire({
            icon: 'success',
            title: 'Workspace Created',
            text: 'Your workspace has been created successfully',
            timer: 1500,
            showConfirmButton: false,

          });
        }, 50);
      }
      open = false;
      name = '';
      description = '';
      dispatch('saved');
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: workspaceId ? 'Update Failed' : 'Creation Failed',
          text: error.response?.data?.message || 'Something went wrong',

        });
      }, 50);
    } finally {
      saving = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !saving) {
      e.preventDefault();
      handleSave();
    }
  }
</script>

{#if open}
  <div class="modal modal-open z-[3000]" on:click={close}>
    <div class="modal-box w-11/12 max-w-md" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-mono text-base-content/30 opacity-50">{pageCode}</span>
            <p class="text-xs uppercase tracking-wide text-base-content/60">Workspace</p>
          </div>
          <h3 class="text-xl font-bold">
            {workspaceId ? 'Edit Workspace' : 'Create Workspace'}
          </h3>
        </div>
        {#if !required}
          <button class="btn btn-ghost btn-sm" on:click={close}>✕</button>
        {/if}
      </div>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Name *</legend>
        <input
          id="workspace-name"
          type="text"
          bind:value={name}
          placeholder="e.g., Personal, Business, Family"
          class="input input-bordered w-full border-2"
          disabled={saving}
          on:keydown={handleKeyDown}
          autofocus
        />
      </fieldset>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Description</legend>
        <textarea
          id="workspace-description"
          bind:value={description}
          placeholder="Optional description for this workspace..."
          class="textarea textarea-bordered w-full border-2 min-h-[100px]"
          disabled={saving}
          on:keydown={handleKeyDown}
        />
      </fieldset>

      <div class="modal-action">
        {#if !required}
          <button class="btn btn-soft" on:click={close} disabled={saving}>
            Cancel
          </button>
        {/if}
        <button class="btn btn-primary" on:click={handleSave} disabled={saving}>
          {#if saving}
            <span class="loading loading-spinner loading-sm"></span>
            {workspaceId ? 'Saving...' : 'Creating...'}
          {:else}
            {workspaceId ? 'Save Changes' : 'Create Workspace'}
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

