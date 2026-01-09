<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { currency } from '$lib/stores/currency';
  import { quickAmounts } from '$lib/stores/quickAmounts';
  import Swal from 'sweetalert2';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';

  let selectedCurrency = 'IDR';
  let amounts = [];
  let newAmount = '';
  let loading = false;
  let lastCurrency = null;
  
  $: pageCode = getPageCode($page.url.pathname);

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;

    // Initialize currency & quick amounts
    await currency.init();
    await quickAmounts.init();

    selectedCurrency = $currency || 'IDR';
    amounts = $quickAmounts || [];
    lastCurrency = selectedCurrency;
  });

  // When currency changes (by user interaction), reset quick amounts to that currency’s defaults
  $: if (selectedCurrency && selectedCurrency !== lastCurrency) {
    const defaults = quickAmounts.getDefaultsByCurrency(selectedCurrency);
    amounts = defaults;
    lastCurrency = selectedCurrency;
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

  async function handleSave() {
    if (loading) return; // Prevent double submission

    const result = await Swal.fire({
      title: 'Save preferences?',
      text: 'This will update your currency and quick amounts.',
      icon: 'question',
      showCancelButton: true,
      confirmButtonText: 'Yes, save',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      zIndex: 9999
    });

    if (!result.isConfirmed) return;

    loading = true;
    try {
      await currency.setCurrency(selectedCurrency);
      await quickAmounts.setAmounts(amounts);

      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Preferences saved',
          text: 'Currency and quick amounts have been updated.',
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Save failed',
          text: error?.response?.data?.message || 'Failed to save preferences',
          zIndex: 9999
        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  function resetToDefaults() {
    const defaults = quickAmounts.getDefaultsByCurrency(selectedCurrency || 'IDR');
    amounts = defaults;
  }
</script>

<div class="max-w-3xl mx-auto space-y-4">
  <PageHeader
    title="Currency & Quick Amount"
    subtitle="Configure your preferred currency and quick amount shortcuts."
    pageCode={pageCode}
    actions={true}
  >
    <svelte:fragment slot="actions">
      <button class="btn btn-soft btn-sm" on:click={() => goto('/app/preferences')}>Back</button>
    </svelte:fragment>
  </PageHeader>

  <!-- Currency Selection -->
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body gap-4">
      <fieldset class="fieldset">
        <legend class="fieldset-legend">Currency</legend>
        <div class="grid gap-2">
          <label class="flex items-center gap-3 p-2 rounded-lg border border-base-300 cursor-pointer hover:bg-base-200 transition">
            <input
              type="radio"
              bind:group={selectedCurrency}
              value="IDR"
              disabled={loading}
              class="radio radio-primary radio-xs"
            />
            <span class="text-sm">Indonesian Rupiah (Rp.)</span>
          </label>
          <label class="flex items-center gap-3 p-2 rounded-lg border border-base-300 cursor-pointer hover:bg-base-200 transition">
            <input
              type="radio"
              bind:group={selectedCurrency}
              value="USD"
              disabled={loading}
              class="radio radio-primary radio-xs"
            />
            <span class="text-sm">US Dollar ($)</span>
          </label>
          <label class="flex items-center gap-3 p-2 rounded-lg border border-base-300 cursor-pointer hover:bg-base-200 transition">
            <input
              type="radio"
              bind:group={selectedCurrency}
              value="EUR"
              disabled={loading}
              class="radio radio-primary radio-xs"
            />
            <span class="text-sm">Euro (€)</span>
          </label>
          <label class="flex items-center gap-3 p-2 rounded-lg border border-base-300 cursor-pointer hover:bg-base-200 transition">
            <input
              type="radio"
              bind:group={selectedCurrency}
              value="JPY"
              disabled={loading}
              class="radio radio-primary radio-xs"
            />
            <span class="text-sm">Japanese Yen (¥)</span>
          </label>
        </div>
      </fieldset>
    </div>
  </div>

  <!-- Quick Amounts -->
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body gap-4">
      <div class="flex items-start justify-between gap-3">
        <div>
          <h3 class="font-semibold text-base">Quick Amounts</h3>
          <p class="text-xs text-base-content/70">
            Customize the quick amount buttons used when adding expenses
            (currency: {selectedCurrency}).
          </p>
        </div>
        <button class="btn btn-soft btn-xs" on:click={resetToDefaults} disabled={loading}>
          Use currency defaults
        </button>
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

      <div class="flex justify-end pt-2">
        <button class="btn btn-primary" on:click={handleSave} disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm mr-1"></span>
            Saving...
          {:else}
            Save
          {/if}
        </button>
      </div>
    </div>
  </div>
</div>

