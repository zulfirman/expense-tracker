<script>
  import { createEventDispatcher } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { page } from '$app/stores';

  export let open = false;
  export let categoryType = 'expense'; // 'expense' or 'income'

  const dispatch = createEventDispatcher();

  let name = '';
  let saving = false;

  $: pageCode = getPageCode($page.url.pathname);

  $: if (open) {
    name = '';
  }

  function close() {
    open = false;
    name = '';
    dispatch('close');
  }

  async function handleSave() {
    if (saving) return;

    if (!name || !name.trim()) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Name',
          text: 'Please enter a category name',

        });
      }, 50);
      return;
    }

    saving = true;
    try {
      await api.post('/categories', {
        name: name.trim(),
        type: categoryType
      });
      
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Category Created',
          text: 'Category has been created successfully',
          timer: 1500,
          showConfirmButton: false,

        });
      }, 50);
      
      open = false;
      name = '';
      dispatch('saved');
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Creation Failed',
          text: error.response?.data?.message || 'Failed to create category',

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
            <p class="text-xs uppercase tracking-wide text-base-content/60">Category</p>
          </div>
          <h3 class="text-xl font-bold">Create {categoryType === 'expense' ? 'Expense' : 'Income'} Category</h3>
        </div>
        <button class="btn btn-ghost btn-sm" on:click={close}>✕</button>
      </div>

      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Name *</legend>
        <input
          id="category-name"
          type="text"
          bind:value={name}
          placeholder="e.g., Food, Transport, Salary"
          class="input input-bordered w-full border-2"
          disabled={saving}
          on:keydown={handleKeyDown}
          autofocus
        />
      </fieldset>

      <div class="modal-action">
        <button class="btn btn-soft" on:click={close} disabled={saving}>
          Cancel
        </button>
        <button class="btn btn-primary" on:click={handleSave} disabled={saving}>
          {#if saving}
            <span class="loading loading-spinner loading-sm"></span>
            Creating...
          {:else}
            Create Category
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

