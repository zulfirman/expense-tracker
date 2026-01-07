<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { currency } from '$lib/stores/currency';
  import { theme } from '$lib/stores/theme';
  import api from '$lib/api';
  import Swal from 'sweetalert2';

  const dispatch = createEventDispatcher();

  let selectedCurrency = 'IDR';
  let currentTheme = 'light';
  let loading = false;

  onMount(async () => {
    // Initialize currency store
    await currency.init();
    selectedCurrency = $currency;
    currentTheme = $theme;
  });

  $: if ($currency) {
    selectedCurrency = $currency;
  }

  $: if ($theme) {
    currentTheme = $theme;
  }

  async function handleSaveCurrency() {
    if (selectedCurrency === $currency) return;
    
    loading = true;
    try {
      await currency.setCurrency(selectedCurrency);
      Swal.fire({
        icon: 'success',
        title: 'Currency Updated',
        text: 'Your currency preference has been saved',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Update Failed',
        text: error.response?.data?.message || 'Failed to update currency',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleThemeToggle() {
    theme.toggle();
  }

  function close() {
    dispatch('close');
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      close();
    }
  }
</script>

<div class="modal-backdrop" on:click={handleBackdropClick}>
  <div class="modal-content" on:click|stopPropagation>
    <div class="modal-header">
      <h2>Preferences</h2>
      <button class="close-btn" on:click={close}>×</button>
    </div>
    <div class="modal-body">
      <!-- Currency Selection -->
      <div class="preference-section">
        <h3>Currency</h3>
        <p class="preference-description">Choose your preferred currency for displaying amounts</p>
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Currency</legend>
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
        </fieldset>
        <button
          class="btn btn-primary"
          on:click={handleSaveCurrency}
          disabled={loading || selectedCurrency === $currency}
        >
          {loading ? 'Saving...' : 'Save Currency'}
        </button>
      </div>

      <!-- Theme Toggle -->
      <div class="preference-section">
        <h3>Appearance</h3>
        <p class="preference-description">Choose between light and dark mode</p>
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Theme</legend>
          <div class="theme-toggle-section">
            <span class="theme-label">Theme: {currentTheme === 'dark' ? 'Dark' : 'Light'}</span>
            <button
              class="theme-toggle-btn"
              on:click={handleThemeToggle}
              title="Toggle Theme"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="sun-icon">
                <circle cx="12" cy="12" r="5"></circle>
                <line x1="12" y1="1" x2="12" y2="3"></line>
                <line x1="12" y1="21" x2="12" y2="23"></line>
                <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
                <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
                <line x1="1" y1="12" x2="3" y2="12"></line>
                <line x1="21" y1="12" x2="23" y2="12"></line>
                <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
                <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
              </svg>
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="moon-icon">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
              </svg>
            </button>
          </div>
        </fieldset>
      </div>
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    padding: 1rem;
  }

  .modal-content {
    background: var(--surface);
    border-radius: 1rem;
    width: 100%;
    max-width: 500px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--border);
  }

  .modal-header h2 {
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 0;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 2rem;
    color: var(--text-secondary);
    cursor: pointer;
    line-height: 1;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    color: var(--text-primary);
  }

  .modal-body {
    padding: 1.5rem;
  }

  .preference-section {
    margin-bottom: 2rem;
  }

  .preference-section:last-child {
    margin-bottom: 0;
  }

  .preference-section h3 {
    font-size: 1.125rem;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }

  .preference-description {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1rem;
  }

  .currency-options {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1rem;
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

  .currency-option:has(input[type="radio"]:checked) {
    border-color: var(--primary-color);
    background: rgba(79, 70, 229, 0.05);
  }

  .btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background-color: #4338ca;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .theme-toggle-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
  }

  .theme-label {
    font-weight: 500;
    color: var(--text-primary);
  }

  .theme-toggle-btn {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    border: 1px solid var(--border);
    background: var(--surface);
    color: var(--text-primary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .theme-toggle-btn:hover {
    background: var(--background);
    transform: scale(1.05);
  }

  .theme-toggle-btn svg {
    width: 20px;
    height: 20px;
  }

  .moon-icon {
    display: none;
  }

  :global(.dark) .sun-icon {
    display: none;
  }

  :global(.dark) .moon-icon {
    display: block;
  }
</style>



