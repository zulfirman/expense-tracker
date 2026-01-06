<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import api from '$lib/api';
  import '$lib/styles/shared.css';

  const dispatch = createEventDispatcher();

  export let month = null; // YYYY-MM format

  let loading = true;

  function formatCurrency(amount) {
    if (!amount && amount !== 0) return '';
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }

  onMount(async () => {
    loading = false;
  });

  function formatMonthYear(monthStr) {
    const [year, monthNum] = monthStr.split('-').map(Number);
    const date = new Date(year, monthNum - 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  }

  function closeModal() {
    dispatch('close');
  }


</script>

<div class="modal-backdrop" on:click={closeModal}>
  <div class="modal-content analytics-modal" on:click|stopPropagation>
    <div class="modal-header">
      <h2>Advanced Analytics & Insights</h2>
      <button class="close-btn" on:click={closeModal}>×</button>
    </div>
    <div class="modal-body">
      {#if month}
        <div class="month-info">
          <h3>{formatMonthYear(month)}</h3>
        </div>
      {/if}

      {#if loading}
        <div class="loading">Loading analytics...</div>
      {:else}
        <div class="no-data">
          <p>No analytics data available.</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .analytics-modal {
    max-width: 900px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
  }

  .month-info {
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--border);
  }

  .month-info h3 {
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 0;
  }

  .loading {
    text-align: center;
    padding: 3rem;
    color: var(--text-secondary);
  }

  .no-data {
    text-align: center;
    padding: 3rem;
    color: var(--text-secondary);
  }

  .no-data .hint {
    font-size: 0.875rem;
    margin-top: 0.5rem;
    opacity: 0.7;
  }

</style>

