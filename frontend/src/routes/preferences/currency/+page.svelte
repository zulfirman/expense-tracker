<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { currency } from '$lib/stores/currency';
  import { quickAmounts } from '$lib/stores/quickAmounts';
  import Swal from 'sweetalert2';
  import '$lib/styles/shared.css';

  let selectedCurrency = 'IDR';
  let amounts = [];
  let newAmount = '';
  let loading = false;
  let lastCurrency = null;

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

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
    if (loading) return;

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

      Swal.fire({
        icon: 'success',
        title: 'Preferences saved',
        text: 'Currency and quick amounts have been updated.',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Save failed',
        text: error?.response?.data?.message || 'Failed to save preferences',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function resetToDefaults() {
    const defaults = quickAmounts.getDefaultsByCurrency(selectedCurrency || 'IDR');
    amounts = defaults;
  }
</script>

<div class="page">
  <div class="header">
    <div>
      <p class="eyebrow">Preferences</p>
      <h1>Currency & Quick Amount</h1>
      <p class="muted">
        Configure your preferred currency and quick amount shortcuts.
      </p>
    </div>
    <button class="btn ghost" on:click={() => goto('/preferences')}>Back</button>
  </div>

  <!-- Currency Selection -->
  <div class="card">
    <div class="card-header">
      <div>
        <h3>Currency</h3>
        <p class="muted">Choose your preferred currency for displaying amounts</p>
      </div>
    </div>

    <div class="currency-options">
      <label class="currency-option">
        <input
          type="radio"
          bind:group={selectedCurrency}
          value="IDR"
          disabled={loading}
        />
        <span>Indonesian Rupiah (Rp.)</span>
      </label>
      <label class="currency-option">
        <input
          type="radio"
          bind:group={selectedCurrency}
          value="USD"
          disabled={loading}
        />
        <span>US Dollar ($)</span>
      </label>
      <label class="currency-option">
        <input
          type="radio"
          bind:group={selectedCurrency}
          value="EUR"
          disabled={loading}
        />
        <span>Euro (€)</span>
      </label>
      <label class="currency-option">
        <input
          type="radio"
          bind:group={selectedCurrency}
          value="JPY"
          disabled={loading}
        />
        <span>Japanese Yen (¥)</span>
      </label>
    </div>
  </div>

  <!-- Quick Amounts -->
  <div class="card">
    <div class="card-header">
      <div>
        <h3>Quick Amounts</h3>
        <p class="muted">
          Customize the quick amount buttons used when adding expenses
          (currency: {selectedCurrency}).
        </p>
      </div>
      <div class="header-actions">
        <button class="btn ghost" on:click={resetToDefaults} disabled={loading}>
          Use currency defaults
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
      <button class="btn secondary" on:click={addAmount} disabled={loading}>Add</button>
    </div>

    {#if amounts.length === 0}
      <div class="empty">No quick amounts yet. Add one above or use defaults.</div>
    {:else}
      <div class="amount-pills">
        {#each amounts as amt}
          <div class="pill">
            <span>{amt.toLocaleString()}</span>
            <button class="remove" on:click={() => removeAmount(amt)} disabled={loading}>×</button>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="footer-actions">
    <button class="btn primary" on:click={handleSave} disabled={loading}>
      {#if loading}
        <span class="spinner"></span> Saving...
      {:else}
        Save
      {/if}
    </button>
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
    margin-bottom: 1rem;
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

  .currency-options {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .currency-option {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .currency-option:hover {
    background: var(--background);
    border-color: var(--primary-color);
  }

  .currency-option input[type="radio"] {
    cursor: pointer;
  }

  .currency-option input[type="radio"]:checked + span {
    font-weight: 600;
    color: var(--primary-color);
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

  .footer-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 0.5rem;
  }
</style>


