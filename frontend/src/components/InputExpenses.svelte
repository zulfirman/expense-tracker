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
      return true
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

<div class="card-body">
  {#if showTitle}
    <h2 class="card-title text-2xl mb-4">{expenseId ? 'Edit Expense' : 'Input Expenses'}</h2>
  {/if}

  <form on:submit|preventDefault={handleSubmit}>
    <!-- Templates Section -->
    {#if templates.length > 0}
      <div class="form-control mb-4">
        <div class="collapse collapse-arrow bg-base-200">
          <input type="checkbox" bind:checked={showTemplates} />
          <div class="collapse-title font-semibold">
            Templates ({templates.length})
          </div>
          <div class="collapse-content">
            <div class="flex flex-col gap-2">
              {#each templates as template}
                <button
                  type="button"
                  class="btn btn-soft justify-between"
                  on:click={() => applyTemplate(template)}
                >
                  <span>{template.name}</span>
                  <span class="font-semibold text-primary">{formatCurrency(template.amount)}</span>
                </button>
              {/each}
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Category Selection -->
    <div class="form-control mb-4">
      <label class="label">
            <span class="label-text font-semibold">
              Category
              {#if selectedCategoryIds.length > 0}
                ({selectedCategoryIds.length} Selected)
              {/if}
            </span>
      </label>
      {#if categoriesLoading}
        <div class="flex justify-center py-4">
          <span class="loading loading-spinner loading-md"></span>
        </div>
      {:else}
        <div class="flex flex-wrap gap-2 mt-2">
          {#each categories as category}
            {@const isSelected = selectedCategoryIds.some(id => {
                const selectedId = Number(id);
                const catId = Number(category.id);
                return !isNaN(selectedId) && !isNaN(catId) && selectedId === catId;
            })}
            <button
              type="button"
              class="btn btn-sm h-auto py-2 px-4 rounded-full transition-all"
              class:btn-primary={isSelected}
              class:btn-soft={!isSelected}
              on:click={(e) => toggleCategory(category.id, e)}
            >
              {category.label}
            </button>
          {/each}
        </div>
      {/if}
    </div>


    <!-- Date Selection -->
    {#if !fixedDate}
      <div class="form-control mb-4">
        <DatePicker
          id="date"
          bind:value={expenseDate}
          placeholder="Select date"
          label="Date"
          on:dateChange={handleDateChange}
        />
      </div>
    {/if}

    <!-- Amount Input -->
    <div class="form-control mb-4">
      <fieldset class="fieldset">
        <legend class="fieldset-legend">Amount (Rp.)</legend>
        <input
          id="amount"
          type="text"
          bind:value={amount}
          on:input={handleAmountInput}
          placeholder="0"
          class="input input-bordered w-full text-lg border-2"
          on:keydown={handleKeyDown}
          inputmode="numeric"
        />
      </fieldset>
      {#if amount}
        <div class="label">
              <span class="label-text-alt text-primary font-semibold text-xl mt-2">
                {formatCurrency(amount)}
              </span>
        </div>
      {/if}

      <!-- Quick Amount Buttons -->
      {#if quickAmountsList.length > 0}
        <div class="mt-4 pt-4 border-t border-base-300">
          <label class="label">
            <span class="label-text text-sm text-base-content/70">Quick Amount:</span>
          </label>
          <div class="grid grid-cols-3 gap-2 mt-2">
            {#each quickAmountsList as quickAmount}
              <button
                type="button"
                class="btn btn-sm"
                class:btn-primary={amount === quickAmount.toString()}
                class:btn-soft={amount !== quickAmount.toString()}
                on:click={() => setQuickAmount(quickAmount)}
              >
                {formatCurrency(quickAmount.toString())}
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Notes -->
    <div class="form-control mb-6">
      <label class="label" for="notes">
        <span class="label-text font-semibold">Notes</span>
      </label>
      <textarea
        id="notes"
        bind:value={notes}
        placeholder="Add notes (optional)"
        class="textarea textarea-bordered w-full border-2"
        rows="3"
        on:keydown={handleKeyDown}
      ></textarea>
    </div>

    <!-- Action Buttons -->
    <div class="flex gap-2">
      {#if showCancel}
        <button type="button" class="btn btn-secondary flex-1" on:click={handleCancel} disabled={loading}>
          Cancel
        </button>
      {:else if !expenseId}
        <button type="button" class="btn btn-soft flex-1" on:click={handleClear} disabled={loading}>
          Clear
        </button>
      {/if}
      <button type="submit" class="btn btn-primary flex-1" disabled={loading}>
        {#if loading}
          <span class="loading loading-spinner loading-sm"></span>
          {expenseId ? 'Updating...' : 'Submitting...'}
        {:else}
          {submitLabel}
        {/if}
      </button>
    </div>
  </form>
</div>

