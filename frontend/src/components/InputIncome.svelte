<script>
  import { createEventDispatcher } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import DatePicker from '$lib/components/DatePicker.svelte';
  import { onMount } from 'svelte';
  import { formatCurrency } from '$lib/utils/currency';

  const dispatch = createEventDispatcher();

  // Props
  export let incomeId = null; // If provided, component is in edit mode
  export let initialData = null; // Initial data for edit mode: { date, notes, amount }
  export let fixedDate = null; // If provided, date field is fixed and not editable
  export let showTitle = true; // Whether to show the "Input Income" title
  export let showCancel = false; // Whether to show cancel button
  export let submitLabel = 'Submit'; // Label for submit button
  export let onSuccess = null; // Optional callback on success
  export let onCancel = null; // Optional callback on cancel

  let expenseDate = fixedDate || new Date().toISOString().split('T')[0];
  let notes = '';
  let amount = '';
  let loading = false;
  let categories = [];
  let selectedCategoryIds = [];
  let categoriesLoading = false;

  const quickAmounts = [100000, 250000, 500000, 750000, 1000000];

  onMount(async () => {
    await loadCategories();
    // If in edit mode, populate form with initial data
    if (incomeId && initialData) {
      expenseDate = initialData.date || expenseDate;
      notes = initialData.notes || '';
      amount = initialData.amount ? initialData.amount.toString() : '';
      // Ensure categoryIds are numbers for consistent comparison and filter out invalid ones
      const categoryIds = (initialData.categoryIds || [])
        .map(id => Number(id))
        .filter(id => !isNaN(id) && id > 0);
      selectedCategoryIds = [...categoryIds];
    }
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories?type=income');
      // Only show active income categories
      categories = response.data
        .filter(cat => cat.isActive !== false)
        .map(cat => ({
          id: cat.id,
          label: cat.name
        }));
    } catch (error) {
      categories = [];
    } finally {
      categoriesLoading = false;
    }
  }

  function setQuickAmount(quickAmount) {
    amount = quickAmount.toString();
  }

  function toggleCategory(categoryId, event) {
    event?.preventDefault();
    event?.stopPropagation();
    
    // Ensure categoryId is valid
    if (!categoryId && categoryId !== 0) {
      console.error('Invalid categoryId:', categoryId);
      return;
    }
    
    // Normalize to number for consistent comparison
    const id = Number(categoryId);
    
    // Check if ID is valid number
    if (isNaN(id)) {
      console.error('Category ID is not a valid number:', categoryId);
      return;
    }
    
    // Check if already selected
    const isSelected = selectedCategoryIds.some(cid => Number(cid) === id);
    
    if (isSelected) {
      // Remove if already selected
      selectedCategoryIds = selectedCategoryIds.filter(cid => Number(cid) !== id);
    } else {
      // Add if not selected
      selectedCategoryIds = [...selectedCategoryIds, id];
    }
  }


  function handleAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    amount = numericValue;
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
    if (selectedCategoryIds.length === 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please select at least one category',
        zIndex: 9999
      });
      return;
    }

    if (!amount || parseFloat(amount) <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please enter a valid amount',
        zIndex: 9999
      });
      return;
    }

    loading = true;
    
    try {
      const payload = {
        categoryIds: selectedCategoryIds,
        date: expenseDate,
        notes: notes,
        amount: parseFloat(amount)
      };

      if (incomeId) {
        // Edit mode
        await api.put(`/income/${incomeId}`, payload);
        
        Swal.fire({
          icon: 'success',
          title: 'Success!',
          text: 'Income updated successfully',
          zIndex: 9999,
          timer: 1500,
          showConfirmButton: false
        });
      } else {
        // Add mode
        const response = await api.post('/income', payload);
        
        const categoryNames = categories
          .filter(cat => selectedCategoryIds.includes(cat.id))
          .map(cat => cat.label)
          .join(', ');
        
        Swal.fire({
          icon: 'success',
          title: 'Success!',
          html: `
            <div style="text-align: left;">
              <p><strong>Type:</strong> Income</p>
              <p><strong>Categories:</strong> ${categoryNames}</p>
              <p><strong>Date:</strong> ${new Date(expenseDate).toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</p>
              <p><strong>Amount:</strong> ${formatCurrency(amount)}</p>
              <p><strong>Notes:</strong> ${notes || '-'}</p>
            </div>
          `,
          confirmButtonText: 'OK',
          zIndex: 9999
        });

        // Call success callback if provided
        if (onSuccess) {
          onSuccess();
        }
        
        // Dispatch success event
        dispatch('success', { incomeId: response?.data?.id });
      }

      // For edit mode, call success callback and dispatch event
      if (incomeId) {
        if (onSuccess) {
          onSuccess();
        }
        dispatch('success', { incomeId });
      }
      
      // Reset form only if not in edit mode
      if (!incomeId) {
        notes = '';
        amount = '';
        selectedCategoryIds = [];
      }
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || (incomeId ? 'Failed to update income' : 'Failed to save income'),
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleClear() {
    expenseDate = fixedDate || new Date().toISOString().split('T')[0];
    notes = '';
    amount = '';
  }

  function handleCancel() {
    if (onCancel) {
      onCancel();
    }
    dispatch('cancel');
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  }
</script>

<form class="input-income" on:submit|preventDefault={handleSubmit}>
  {#if showTitle}
    <h1>{incomeId ? 'Edit Income' : 'Input Income'}</h1>
  {/if}

  <div class="form-group">
    <label>Category {#if selectedCategoryIds.length > 0}<span class="selected-count">({selectedCategoryIds.length} selected)</span>{/if}</label>
    {#if categoriesLoading}
      <div class="loading-categories">Loading categories...</div>
    {:else}
      <div class="category-pills">
        {#each categories as category}
          <button
            type="button"
            class="category-pill"
            class:selected={selectedCategoryIds.some(id => {
              const selectedId = Number(id);
              const catId = Number(category.id);
              return !isNaN(selectedId) && !isNaN(catId) && selectedId === catId;
            })}
            on:click={(e) => toggleCategory(category.id, e)}
          >
            {category.label}
          </button>
        {/each}
      </div>
    {/if}
  </div>

  {#if !fixedDate}
    <div class="form-group">
      <label for="date">Date</label>
      <DatePicker
        id="date"
        bind:value={expenseDate}
        placeholder="Select date"
        on:dateChange={handleDateChange}
      />
    </div>
  {/if}

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
    
    <!-- Quick Amount Buttons -->
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
    <label for="notes">Notes</label>
    <textarea
      id="notes"
      bind:value={notes}
      placeholder="Add notes (optional)"
      class="form-textarea"
      rows="3"
      on:keydown={handleKeyDown}
    ></textarea>
  </div>

  <div class="button-group">
    {#if showCancel}
      <button type="button" class="btn btn-secondary" on:click={handleCancel} disabled={loading}>Cancel</button>
    {:else if !incomeId}
      <button type="button" class="btn btn-secondary" on:click={handleClear} disabled={loading}>Clear</button>
    {/if}
    <button type="submit" class="btn btn-primary" disabled={loading}>
      {#if loading}
        <span class="spinner"></span> {incomeId ? 'Updating...' : 'Submitting...'}
      {:else}
        {submitLabel}
      {/if}
    </button>
  </div>
  <div class="space-xl"></div>
</form>

<style>
  .input-income {
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

  /* Quick Amount Buttons */
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

  /* Category Pills */
  .selected-count {
    color: var(--primary-color);
    font-weight: 600;
  }

  .loading-categories {
    padding: 1rem;
    text-align: center;
    color: var(--text-secondary);
  }

  .category-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }

  .category-pill {
    padding: 0.5rem 1rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: 1.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .category-pill:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    transform: translateY(-1px);
  }

  .category-pill.selected {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  /* Mobile Optimizations */
  @media (max-width: 768px) {
    .input-income {
      padding: 0.5rem;
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

    .category-pill {
      padding: 0.75rem 1rem;
      font-size: 0.875rem;
    }
  }
</style>

