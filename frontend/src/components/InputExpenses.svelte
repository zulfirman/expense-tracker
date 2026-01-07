<script>
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import DatePicker from '$lib/components/DatePicker.svelte';
  import { onMount } from 'svelte';
  import { formatCurrency } from '$lib/utils/currency';
  import { quickAmounts } from '$lib/stores/quickAmounts';
  import { currency } from '$lib/stores/currency';

  const dispatch = createEventDispatcher();

  // Props
  export let expenseId = null; // If provided, component is in edit mode
  export let initialData = null; // Initial data for edit mode: { categoryIds, date, notes, amount }
  export let fixedDate = null; // If provided, date field is fixed and not editable
  export let showTitle = true; // Whether to show the "Input Expenses" title
  export let showCancel = false; // Whether to show cancel button
  export let submitLabel = 'Submit'; // Label for submit button
  export let onSuccess = null; // Optional callback on success
  export let onCancel = null; // Optional callback on cancel

  let categories = [];
  let templates = [];
  let selectedCategoryIds = []; // Array of category IDs (numbers)
  let expenseDate = fixedDate || new Date().toISOString().split('T')[0];
  let notes = '';
  let amount = '';
  let loading = false;
  let categoriesLoading = true;
  let templatesLoading = true;
  let showTemplates = false;

  let quickAmountsList = [];
  $: quickAmountsList = $quickAmounts || [];
  $: if ((!quickAmountsList || quickAmountsList.length === 0) && $currency) {
    quickAmountsList = quickAmounts.getDefaultsByCurrency($currency);
  }

  onMount(async () => {
    await Promise.all([loadCategories(), loadTemplates(), quickAmounts.init()]);
    quickAmountsList = get(quickAmounts);
    
    // If in edit mode, populate form with initial data
    if (expenseId && initialData) {
      // Ensure categoryIds are numbers for consistent comparison and filter out invalid ones
      const categoryIds = (initialData.categoryIds || [])
        .map(id => Number(id))
        .filter(id => !isNaN(id) && id > 0);
      selectedCategoryIds = [...categoryIds];
      expenseDate = initialData.date || expenseDate;
      notes = initialData.notes || '';
      amount = initialData.amount ? initialData.amount.toString() : '';
    }
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories?type=expense');
      // Only show active expense categories in expense input
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

  async function loadTemplates() {
    try {
      const response = await api.get('/templates');
      templates = response.data;
    } catch (error) {
      // Templates failed to load, continue without them
    } finally {
      templatesLoading = false;
    }
  }


  function setQuickAmount(quickAmount) {
    amount = quickAmount.toString();
  }

  function applyTemplate(template) {
    // Normalize category IDs to numbers and filter out invalid ones
    const templateIds = (template.categoryIds || template.categories || [])
      .map(id => Number(id))
      .filter(id => !isNaN(id) && id > 0);
    selectedCategoryIds = [...templateIds];
    amount = template.amount.toString();
    notes = template.notes || '';
    showTemplates = false;
  }

  function handleAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    amount = numericValue;
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
    
    // If category selected, scroll to date field
    /*if (selectedCategoryIds.length > 0) {
      setTimeout(() => {
        const dateInput = document.getElementById('date');
        if (dateInput) {
          dateInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
      }, 100);
    }*/
  }

  function handleDateChange() {
    // After date selected, scroll to amount field (without focusing)
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
      // Filter out any invalid IDs and ensure they're numbers
      const validCategoryIds = selectedCategoryIds
        .map(id => Number(id))
        .filter(id => !isNaN(id) && id > 0);
      
      if (validCategoryIds.length === 0) {
        Swal.fire({
          icon: 'warning',
          title: 'Warning',
          text: 'Please select at least one valid category',
          zIndex: 9999
        });
        loading = false;
        return;
      }
      
      const payload = {
        categoryIds: validCategoryIds,
        date: expenseDate,
        notes: notes,
        amount: parseFloat(amount)
      };

      if (expenseId) {
        // Edit mode
        await api.put(`/expenses/${expenseId}`, payload);
        
        Swal.fire({
          icon: 'success',
          title: 'Success!',
          text: 'Expense updated successfully',
          zIndex: 9999,
          timer: 1500,
          showConfirmButton: false
        });
      } else {
        // Add mode
        const response = await api.post('/expenses', payload);
        
        const categoryNames = categories
          .filter(cat => selectedCategoryIds.includes(cat.id))
          .map(cat => cat.label)
          .join(', ');
        
        Swal.fire({
          icon: 'success',
          title: 'Success!',
          html: `
            <div style="text-align: left;">
              <p><strong>Type:</strong> Expense</p>
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
        dispatch('success', { expenseId: response?.data?.id });
      }

      // For edit mode, call success callback and dispatch event
      if (expenseId) {
        if (onSuccess) {
          onSuccess();
        }
        dispatch('success', { expenseId });
      }

      // Reset form only if not in edit mode
      if (!expenseId) {
        selectedCategoryIds = [];
        notes = '';
        amount = '';
      }
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || (expenseId ? 'Failed to update expense' : 'Failed to save expense'),
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleClear() {
    selectedCategoryIds = [];
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

<form class="input-expenses" on:submit|preventDefault={handleSubmit}>
  {#if showTitle}
    <h1>{expenseId ? 'Edit Expense' : 'Input Expenses'}</h1>
  {/if}

  <!-- Templates Section -->
  {#if templates.length > 0}
    <div class="form-group">
      <div class="templates-header">
        <label>Templates</label>
        <button 
          type="button" 
          class="toggle-btn"
          on:click={() => showTemplates = !showTemplates}
        >
          {showTemplates ? 'Hide' : 'Show'} ({templates.length})
        </button>
      </div>
      {#if showTemplates}
        <div class="templates-list">
          {#each templates as template}
            <button
              type="button"
              class="template-btn"
              on:click={() => applyTemplate(template)}
            >
              <span class="template-name">{template.name}</span>
              <span class="template-amount">{formatCurrency(template.amount)}</span>
            </button>
          {/each}
        </div>
      {/if}
    </div>
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
        {#each quickAmountsList as quickAmount}
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
    {:else if !expenseId}
      <button type="button" class="btn btn-secondary" on:click={handleClear} disabled={loading}>Clear</button>
    {/if}
    <button type="submit" class="btn btn-primary" disabled={loading}>
      {#if loading}
        <span class="spinner"></span> {expenseId ? 'Updating...' : 'Submitting...'}
      {:else}
        {submitLabel}
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

  /* Templates */
  .templates-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .toggle-btn {
    padding: 0.25rem 0.75rem;
    background: var(--background);
    border: 1px solid var(--border);
    border-radius: 0.375rem;
    font-size: 0.875rem;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s;
  }

  .toggle-btn:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    color: var(--primary-color);
  }

  .templates-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .template-btn {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .template-btn:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    transform: translateX(4px);
  }

  .template-name {
    font-weight: 500;
    color: var(--text-primary);
  }

  .template-amount {
    font-weight: 600;
    color: var(--primary-color);
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

  /* Mobile Optimizations */
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
      min-height: 40px; /* Better touch target */
    }

    .form-input,
    .form-textarea {
      font-size: 16px; /* Prevents zoom on iOS */
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

