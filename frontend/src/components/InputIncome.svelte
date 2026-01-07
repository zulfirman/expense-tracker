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
    categoriesLoading = true;
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
    selectedCategoryIds = [];
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

<div class="max-w-2xl mx-auto">
  <div class="card-body">
    {#if showTitle}
      <h2 class="card-title text-2xl mb-4">{incomeId ? 'Edit Income' : 'Input Income'}</h2>
    {/if}

    <form on:submit|preventDefault={handleSubmit}>
      <!-- Category Selection -->
      <div class="form-control mb-4">
        <label class="label">
            <span class="label-text font-semibold">
              Category
              {#if selectedCategoryIds.length > 0}
                <span class="badge badge-primary badge-sm ml-2">{selectedCategoryIds.length} Selected</span>
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
                class:btn-outline={!isSelected}
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
          <label class="label" for="date">
            <span class="label-text font-semibold">Date</span>
          </label>
          <DatePicker
            id="date"
            bind:value={expenseDate}
            placeholder="Select date"
            on:dateChange={handleDateChange}
          />
        </div>
      {/if}

      <!-- Amount Input -->
      <div class="form-control mb-4">
        <label class="label" for="amount">
          <span class="label-text font-semibold">Amount (Rp.)</span>
        </label>
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
        {#if amount}
          <div class="label">
              <span class="label-text-alt text-primary font-semibold text-xl mt-2">
                {formatCurrency(amount)}
              </span>
          </div>
        {/if}

        <!-- Quick Amount Buttons -->
        <div class="mt-4 pt-4 border-t border-base-300">
          <label class="label">
            <span class="label-text text-sm text-base-content/70">Quick Amount:</span>
          </label>
          <div class="grid grid-cols-3 gap-2 mt-2">
            {#each quickAmounts as quickAmount}
              <button
                type="button"
                class="btn btn-sm"
                class:btn-primary={amount === quickAmount.toString()}
                class:btn-outline={amount !== quickAmount.toString()}
                on:click={() => setQuickAmount(quickAmount)}
              >
                {formatCurrency(quickAmount.toString())}
              </button>
            {/each}
          </div>
        </div>
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
        {:else if !incomeId}
          <button type="button" class="btn btn-soft flex-1" on:click={handleClear} disabled={loading}>
            Clear
          </button>
        {/if}
        <button type="submit" class="btn btn-primary flex-1" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            {incomeId ? 'Updating...' : 'Submitting...'}
          {:else}
            {submitLabel}
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>
