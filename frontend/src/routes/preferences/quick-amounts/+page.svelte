<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { quickAmounts } from '$lib/stores/quickAmounts';
  import { currency } from '$lib/stores/currency';
  import { auth } from '$lib/stores/auth';
  import Swal from 'sweetalert2';

  let amounts = [];
  let newAmount = '';
  let loading = false;

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }
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
    loading = true;
    try {
      await quickAmounts.setAmounts(amounts);
      Swal.fire({
        icon: 'success',
        title: 'Quick amounts saved',
        timer: 1200,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (err) {
      Swal.fire({
        icon: 'error',
        title: 'Save failed',
        text: 'Could not save quick amounts',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function resetToDefaults() {
    const defaults = quickAmounts.getDefaultsByCurrency($currency || 'IDR');
    amounts = defaults;
  }
</script>

<div class="page">
  <div class="header">
    <div>
      <p class="eyebrow">Preferences</p>
      <h1>Quick Amount</h1>
      <p class="muted">Customize your quick amount shortcuts (currency: {$currency || 'IDR'})</p>
    </div>
    <button class="btn ghost" on:click={() => goto('/preferences')}>Back</button>
  </div>

  <div class="card">
    <div class="card-header">
      <div>
        <h3>Manage Quick Amounts</h3>
        <p class="muted">Add, remove, and reorder the quick amount buttons shown in expense input</p>
      </div>
      <div class="header-actions">
        <button class="btn ghost" on:click={resetToDefaults} disabled={loading}>Use currency defaults</button>
        <button class="btn primary" on:click={save} disabled={loading}>
          {#if loading}
            <span class="spinner"></span> Saving...
          {:else}
            Save
          {/if}
        </button>
      </div>
    </div>

    <div class="form-row">
      <input
        type="number"
        min="0"
        step="1"
        placeholder="Add amount"
        bind:value={newAmount}
        class="form-input"
        on:keydown={(e) => { if (e.key === 'Enter') addAmount(); }}
      />
      <button class="btn secondary" on:click={addAmount}>Add</button>
    </div>

    {#if amounts.length === 0}
      <div class="empty">No quick amounts yet. Add one above or use defaults.</div>
    {:else}
      <div class="amount-pills">
        {#each amounts as amt}
          <div class="pill">
            <span>{amt.toLocaleString()}</span>
            <button class="remove" on:click={() => removeAmount(amt)}>×</button>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .page {
    max-width: 800px;
    margin: 0 auto;
    padding: 1.5rem;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }
  .eyebrow {
    font-size: 0.85rem;
    color: var(--text-secondary);
    margin: 0;
  }
  h1 {
    margin: 0.2rem 0;
    color: var(--text-primary);
    font-size: 1.5rem;
  }
  .muted {
    color: var(--text-secondary);
    margin: 0.2rem 0 0;
  }
  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    padding: 1.25rem;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  .header-actions {
    display: flex;
    gap: 0.5rem;
  }
  .form-row {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  .form-input {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    background: var(--background);
    color: var(--text-primary);
  }
  .amount-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .pill {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--border);
    border-radius: 1rem;
    background: var(--background);
  }
  .pill .remove {
    border: none;
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 0.9rem;
  }
  .empty {
    color: var(--text-secondary);
    padding: 0.75rem 0;
  }
</style>


