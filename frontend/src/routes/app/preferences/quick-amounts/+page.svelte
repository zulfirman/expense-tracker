<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { quickAmounts } from '$lib/stores/quickAmounts';
  import { currency } from '$lib/stores/currency';
  import { auth } from '$lib/stores/auth';
  import Swal from 'sweetalert2';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';

  let amounts = [];
  let newAmount = '';
  let loading = false;
  
  $: pageCode = getPageCode($page.url.pathname);

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;
    await quickAmounts.init();
    amounts = $quickAmounts || [];
  });

  $: if ($quickAmounts) {
    amounts = $quickAmounts;
  }

  function addAmount() {
    const value = Number(newAmount);
    if (isNaN(value) || value <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Invalid amount',
        text: 'Enter a positive number',
        zIndex: 9999
      });
      return;
    }
    amounts = Array.from(new Set([...amounts, value])).sort((a, b) => a - b);
    newAmount = '';
  }

  function removeAmount(value) {
    amounts = amounts.filter((v) => v !== value);
  }

  async function save() {
    if (loading) return; // Prevent double submission
    
    loading = true;
    try {
      await quickAmounts.setAmounts(amounts);
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Quick amounts saved',
          timer: 1200,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
    } catch (err) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Save failed',
          text: 'Could not save quick amounts',
          zIndex: 9999
        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  function resetToDefaults() {
    const defaults = quickAmounts.getDefaultsByCurrency($currency || 'IDR');
    amounts = defaults;
  }
</script>

<div class="max-w-3xl mx-auto space-y-4">
  <PageHeader
    title="Quick Amount"
    subtitle="Customize your quick amount shortcuts (currency: {$currency || 'IDR'})."
    pageCode={pageCode}
    actions={true}
  >
    <svelte:fragment slot="actions">
      <button class="btn btn-soft btn-sm" on:click={() => goto('/app/preferences')}>Back</button>
    </svelte:fragment>
  </PageHeader>

  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body gap-4">
      <div class="flex items-start justify-between gap-3">
        <div>
          <h3 class="font-semibold text-base">Manage Quick Amounts</h3>
          <p class="text-xs text-base-content/70">
            Add, remove, and reorder the quick amount buttons shown in expense input.
          </p>
        </div>
        <div class="flex gap-2">
          <button class="btn btn-soft btn-xs" on:click={resetToDefaults} disabled={loading}>
            Use currency defaults
          </button>
          <button class="btn btn-primary btn-xs" on:click={save} disabled={loading}>
            {#if loading}
              <span class="loading loading-spinner loading-xs mr-1"></span>
              Saving...
            {:else}
              Save
            {/if}
          </button>
        </div>
      </div>

      <fieldset class="fieldset">
        <legend class="fieldset-legend">Add quick amount</legend>
        <div class="flex gap-2">
          <input
            type="number"
            min="0"
            step="1"
            placeholder="Add amount"
            bind:value={newAmount}
            class="input input-bordered w-full border-2"
            on:keydown={(e) => { if (e.key === 'Enter') addAmount(); }}
          />
          <button class="btn btn-primary" on:click={addAmount} disabled={loading}>Add</button>
        </div>
      </fieldset>

      {#if amounts.length === 0}
        <div class="text-sm text-base-content/60 py-2">
          No quick amounts yet. Add one above or use defaults.
        </div>
      {:else}
        <div class="flex flex-wrap gap-2">
          {#each amounts as amt}
            <div class="badge px-3 py-2 badge-outline gap-1 px-3 py-3 border-1">
              <span class="text-sm">{amt.toLocaleString()}</span>
              <button
                class="btn btn-xs btn-ghost text-error min-h-0 h-5 w-5 p-0"
                on:click={() => removeAmount(amt)}
                disabled={loading}
              >
                ×
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</div>

