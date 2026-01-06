<script>
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { onMount } from 'svelte';
  import InputExpenses from '$components/InputExpenses.svelte';
  import DatePicker from '$lib/components/DatePicker.svelte';

  // Simplified for husband/wife budget tracking
  let categories = [];
  let selectedCategories = [];
  let expenseDate = new Date().toISOString().split('T')[0];
  let notes = '';
  let amount = '';
  let loading = false;
  let categoriesLoading = true;

  // Simplified quick amounts for common expenses
  const quickAmounts = [10000, 25000, 50000, 100000, 200000];

  onMount(async () => {
    await loadCategories();
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories');
      // Only show active categories in expense input
      categories = response.data
        .filter(cat => cat.isActive !== false)
        .map(cat => ({
          label: cat.name
        }));
    } catch (error) {
      // Simple default categories for family budget
      categories = [
        { id: 'daily', label: 'Daily Needs' },
        { id: 'monthly', label: 'Monthly Bills' },
        { id: 'others', label: 'Others' }
      ];
    } finally {
      categoriesLoading = false;
    }
  }

  function setQuickAmount(quickAmount) {
    amount = quickAmount.toString();
  }

  function formatCurrency(value) {
    if (!value && value !== 0) return '';
    let numericValue;
    if (typeof value === 'string') {
      numericValue = value.replace(/\D/g, '');
      if (!numericValue) return '';
      numericValue = parseFloat(numericValue);
    } else {
      numericValue = value;
    }
    if (isNaN(numericValue)) return '';
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(numericValue);
  }

  function handleAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    amount = numericValue;
  }

  function toggleCategory(categoryId) {
    if (selectedCategories.includes(categoryId)) {
      selectedCategories = selectedCategories.filter(id => id !== categoryId);
    } else {
      selectedCategories = [...selectedCategories, categoryId];
    }
    
    if (selectedCategories.length > 0) {
      setTimeout(() => {
        const dateInput = document.getElementById('date');
        if (dateInput) {
          dateInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
      }, 100);
    }
  }

  function handleDateChange() {
    setTimeout(() => {
      const amountInput = document.getElementById('amount');
      if (amountInput) {
        amountInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    }, 100);
  }

  async function handleSubmit() {
    if (selectedCategories.length === 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Select Category',
        text: 'Please select at least one category',
        zIndex: 9999
      });
      return;
    }

    if (!amount || parseFloat(amount) <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Enter Amount',
        text: 'Please enter a valid amount',
        zIndex: 9999
      });
      return;
    }

    loading = true;
    
    try {
      const payload = {
        categories: selectedCategories,
        date: expenseDate,
        notes: notes,
        amount: parseFloat(amount)
      };
      await api.post('/expenses', payload);
      
      const categoryNames = categories
        .filter(cat => selectedCategories.includes(cat.id))
        .map(cat => cat.label)
        .join(', ');
      
      Swal.fire({
        icon: 'success',
        title: 'Success!',
        html: `
          <div style="text-align: left;">
            <p><strong>Category:</strong> ${categoryNames}</p>
            <p><strong>Date:</strong> ${new Date(expenseDate).toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</p>
            <p><strong>Amount:</strong> ${formatCurrency(amount)}</p>
            ${notes ? `<p><strong>Notes:</strong> ${notes}</p>` : ''}
          </div>
        `,
        confirmButtonText: 'OK',
        zIndex: 9999
      });

      // Reset form
      selectedCategories = [];
      notes = '';
      amount = '';
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to save expense',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleClear() {
    selectedCategories = [];
    expenseDate = new Date().toISOString().split('T')[0];
    notes = '';
    amount = '';
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  }
</script>

<form class="input-expenses" on:submit|preventDefault={handleSubmit}>
  <h1>Record Expense</h1>

  <div class="form-group">
    <label>Category {#if selectedCategories.length > 0}<span class="selected-count">({selectedCategories.length} selected)</span>{/if}</label>
    {#if categoriesLoading}
      <div class="loading-categories">Loading categories...</div>
    {:else}
      <div class="category-pills">
        {#each categories as category}
          <button
            type="button"
            class="category-pill"
            class:selected={selectedCategories.includes(category.id)}
            on:click={() => toggleCategory(category.id)}
          >
            {category.label}
          </button>
        {/each}
      </div>
    {/if}
  </div>

  <div class="form-group">
    <label for="date">Date</label>
    <DatePicker
      id="date"
      bind:value={expenseDate}
      placeholder="Select date"
      on:dateChange={handleDateChange}
    />
  </div>

  <div class="form-group">
    <label for="amount">Amount (Rp.)</label>
    <input
      id="amount"
      type="text"
      bind:value={amount}
      on:input={handleAmountInput}
      placeholder="0"
      class="form-input"
      on:keydown={handleKeyDown}
      inputmode="numeric"
    />
    {#if amount}
      <div class="amount-preview">{formatCurrency(amount)}</div>
    {/if}
    
    <div class="quick-amounts">
      <label class="quick-amounts-label">Quick Amount:</label>
      <div class="quick-amounts-buttons">
        {#each quickAmounts as quickAmount}
          <button
            type="button"
            class="quick-amount-btn"
            class:active={amount === quickAmount.toString()}
            on:click={() => setQuickAmount(quickAmount)}
          >
            {formatCurrency(quickAmount.toString())}
          </button>
        {/each}
      </div>
    </div>
  </div>

  <div class="form-group">
    <label for="notes">Notes (optional)</label>
    <textarea
      id="notes"
      bind:value={notes}
      placeholder="Add notes"
      class="form-textarea"
      rows="3"
      on:keydown={handleKeyDown}
    ></textarea>
  </div>

  <div class="button-group">
    <button type="button" class="btn btn-secondary" on:click={handleClear} disabled={loading}>Clear</button>
    <button type="submit" class="btn btn-primary" disabled={loading}>
      {#if loading}
        <span class="spinner"></span> Saving...
      {:else}
        Save
      {/if}
    </button>
  </div>
  <div class="space-xl"></div>
</form>

<style>
  .input-expenses {
    max-width: 600px;
    margin: 0 auto;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .category-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .category-pill {
    padding: 0.5rem 1rem;
    border: 2px solid var(--border);
    border-radius: 2rem;
    background: var(--surface);
    color: var(--text-primary);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .category-pill:hover {
    background-color: var(--background);
    border-color: var(--primary-color);
  }

  .category-pill.selected {
    background-color: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .selected-count {
    font-size: 0.875rem;
    font-weight: normal;
    color: var(--text-secondary);
  }

  .loading-categories {
    padding: 1rem;
    text-align: center;
    color: var(--text-secondary);
  }

  .form-input,
  .form-textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
    transition: border-color 0.2s;
  }

  .form-input:focus,
  .form-textarea:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .amount-preview {
    margin-top: 0.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--primary-color);
  }

  .button-group {
    display: flex;
    gap: 1rem;
    margin-top: 2rem;
  }

  .btn {
    flex: 1;
    padding: 0.875rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: white;
  }

  .btn-primary:hover {
    background-color: #4338ca;
  }

  .btn-primary:active {
    transform: scale(0.98);
  }

  .btn-secondary {
    background-color: var(--surface);
    color: var(--text-primary);
    border: 1px solid var(--border);
  }

  .btn-secondary:hover {
    background-color: var(--background);
  }

  .btn-secondary:active {
    transform: scale(0.98);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-primary:disabled {
    background-color: var(--text-secondary);
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 0.6s linear infinite;
    margin-right: 0.5rem;
    vertical-align: middle;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .quick-amounts {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border);
  }

  .quick-amounts-label {
    display: block;
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
  }

  .quick-amounts-buttons {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.5rem;
  }

  .quick-amount-btn {
    padding: 0.75rem 0.5rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s;
    text-align: center;
  }

  .quick-amount-btn:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    transform: translateY(-2px);
  }

  .quick-amount-btn.active {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  @media (max-width: 768px) {
    .input-expenses {
      padding: 0.5rem;
    }

    .category-pills {
      gap: 0.5rem;
    }

    .category-pill {
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
      min-height: 40px;
    }

    .form-input,
    .form-textarea {
      font-size: 16px;
      padding: 1rem;
      min-height: 48px;
    }

    .quick-amounts-buttons {
      grid-template-columns: repeat(2, 1fr);
    }

    .quick-amount-btn {
      padding: 1rem 0.5rem;
      font-size: 0.875rem;
    }

    .button-group {
      flex-direction: column;
    }

    .btn {
      min-height: 48px;
      font-size: 1rem;
    }
  }
</style>

