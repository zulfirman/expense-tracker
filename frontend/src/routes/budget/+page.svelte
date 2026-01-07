<script>
    // ============================================================================
    // IMPORTS
    // ============================================================================
    import api from '$lib/api';
    import Swal from 'sweetalert2';
    import { onMount } from 'svelte';
    import { formatCurrency } from '$lib/utils/currency';

    // ============================================================================
    // STATE VARIABLES
    // ============================================================================
    let categories = [];              // List of available categories
    let budgets = {};                 // Map of categoryId -> budget amount
    let currentMonth = new Date().toISOString().slice(0, 7); // YYYY-MM format
    let loading = false;              // Loading state for async operations
    let categoriesLoading = true;    // Loading state for categories
    let expensesByCategory = {};     // Map of categoryId -> total expenses
    let budgetCards = [];            // Computed array of budget cards with progress
    let showBudgetForm = false;      // Modal visibility for budget form
    let editingCategory = null;      // Category ID being edited (null = new)
    let budgetAmount = '';           // Budget amount input value

    // ============================================================================
    // LIFECYCLE HOOKS
    // ============================================================================
    onMount(async () => {
        // Load categories first, then budgets and expenses in parallel
        await loadCategories();
        await Promise.all([loadBudgets(), loadExpenses()]);

        // Check if we should offer to import from latest budget
        await checkAndOfferImport();
    });

    // ============================================================================
    // UTILITY FUNCTIONS
    // ============================================================================

    /**
     * Format month string (YYYY-MM) to readable format (e.g., "January 2024")
     */
    function formatMonthYear(monthStr) {
        const [year, monthNum] = monthStr.split('-').map(Number);
        const date = new Date(year, monthNum - 1);
        return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    }


    /**
     * Calculate budget progress percentage (0-100)
     */
    function getBudgetProgress(budget, spent) {
        if (!budget) return 0;
        return Math.min((spent / budget) * 100, 100);
    }

    /**
     * Get budget status based on progress percentage
     * Returns: 'exceeded', 'warning', 'caution', or 'good'
     */
    function getBudgetStatus(progress) {
        if (progress >= 100) return 'exceeded';
        if (progress >= 75) return 'warning';
        if (progress >= 50) return 'caution';
        return 'good';
    }

    /**
     * Get color for budget status
     */
    function getStatusColor(status) {
        switch (status) {
            case 'exceeded': return 'var(--danger)';
            case 'warning': return '#f59e0b';
            case 'caution': return '#fbbf24';
            default: return 'var(--success)';
        }
    }

    // ============================================================================
    // DATA LOADING FUNCTIONS
    // ============================================================================

    /**
     * Load categories from API
     * Only shows active categories
     */
    async function loadCategories() {
        try {
            const response = await api.get('/categories?type=expense');
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

    /**
     * Load budgets for current month
     * Stores as map: categoryId -> amount
     */
    async function loadBudgets() {
        try {
            const response = await api.get(`/budgets?month=${currentMonth}`);
            const budgetsData = response.data || [];
            budgets = {};
            budgetsData.forEach(budget => {
                budgets[budget.categoryId] = budget.amount;
            });
        } catch (error) {
            budgets = {};
        }
    }

    /**
     * Load expenses for current month grouped by category
     */
    async function loadExpenses() {
        try {
            const response = await api.get(`/expenses/month/${currentMonth}`);
            expensesByCategory = {};
            if (response.data && response.data.categories) {
                response.data.categories.forEach(item => {
                    const categoryId = item.categoryId;
                    if (categoryId) {
                        expensesByCategory[categoryId] = (expensesByCategory[categoryId] || 0) + item.total;
                    }
                });
            }
        } catch (error) {
            expensesByCategory = {};
        }
    }

    /**
     * Check if current month has no budgets and offer to import from latest budget month
     * Uses /api/budgets/latest endpoint to get the most recent month with budgets
     */
    async function checkAndOfferImport() {
        const hasBudgets = Object.keys(budgets).length > 0;

        if (!hasBudgets) {
            try {
                // Get the latest month that has budgets (most recently inputted)
                const latestResponse = await api.get('/budgets/latest');
                const latestMonth = latestResponse.data?.month;

                if (latestMonth && latestMonth !== currentMonth) {
                    // Load budgets from latest month to get count
                    const latestBudgetsResponse = await api.get(`/budgets?month=${latestMonth}`);
                    const latestBudgets = latestBudgetsResponse.data || [];

                    if (latestBudgets.length > 0) {
                        const currentMonthName = formatMonthYear(currentMonth);
                        const latestMonthName = formatMonthYear(latestMonth);

                        const result = await Swal.fire({
                            icon: 'question',
                            title: 'Import Latest Budget?',
                            html: `No budgets found for <strong>${currentMonthName}</strong>.<br>Would you like to import budgets from <strong>${latestMonthName}</strong> (your latest budget)?`,
                            showCancelButton: true,
                            confirmButtonText: 'Yes, Import',
                            cancelButtonText: 'No, Thanks',
                            reverseButtons: true,
                            zIndex: 9999
                        });

                        if (result.isConfirmed) {
                            loading = true;
                            try {
                                await api.post('/budgets/copy', {
                                    fromMonth: latestMonth,
                                    toMonth: currentMonth
                                });

                                await loadBudgets();

                                Swal.fire({
                                    icon: 'success',
                                    title: 'Budgets Imported!',
                                    text: `Successfully imported ${latestBudgets.length} budget(s) from ${latestMonthName}`,
                                    timer: 2000,
                                    showConfirmButton: false,
                                    zIndex: 9999
                                });
                            } catch (error) {
                                Swal.fire({
                                    icon: 'error',
                                    title: 'Import Failed',
                                    text: error.response?.data?.message || 'Failed to import budgets',
                                    zIndex: 9999
                                });
                            } finally {
                                loading = false;
                            }
                        }
                    }
                }
            } catch (error) {
                // Silently fail if we can't check latest budget month
            }
        }
    }

    // ============================================================================
    // BUDGET FORM FUNCTIONS
    // ============================================================================

    /**
     * Open budget form modal
     * @param {number|null} categoryId - Category ID to edit, or null for new budget
     */
    function openBudgetForm(categoryId = null) {
        editingCategory = categoryId;
        budgetAmount = categoryId && budgets[categoryId] ? budgets[categoryId].toString() : '';
        showBudgetForm = true;
    }

    /**
     * Close budget form modal and reset state
     */
    function closeBudgetForm() {
        showBudgetForm = false;
        editingCategory = null;
        budgetAmount = '';
    }

    /**
     * Save budget to API
     * Validates amount and creates/updates budget
     */
    async function saveBudget() {
        if (!editingCategory) return;

        const amount = parseFloat(budgetAmount.replace(/\D/g, ''));
        if (isNaN(amount) || amount <= 0) {
            Swal.fire({
                icon: 'warning',
                title: 'Invalid Amount',
                text: 'Please enter a valid budget amount',
                zIndex: 9999
            });
            return;
        }

        loading = true;
        try {
            await api.post('/budgets', {
                categoryId: editingCategory,
                amount: amount,
                month: currentMonth
            });

            await loadBudgets();
            closeBudgetForm();

            Swal.fire({
                icon: 'success',
                title: 'Budget Saved',
                text: 'Budget has been saved successfully',
                timer: 1500,
                showConfirmButton: false,
                zIndex: 9999
            });
        } catch (error) {
            Swal.fire({
                icon: 'error',
                title: 'Error',
                text: error.response?.data?.message || 'Failed to save budget',
                zIndex: 9999
            });
        } finally {
            loading = false;
        }
    }

    /**
     * Delete budget for a category
     * Shows confirmation dialog first
     */
    async function deleteBudget(categoryId) {
        const result = await Swal.fire({
            icon: 'warning',
            title: 'Delete Budget?',
            text: 'Are you sure you want to delete this budget?',
            showCancelButton: true,
            confirmButtonText: 'Delete',
            cancelButtonText: 'Cancel',
            reverseButtons: true,
            zIndex: 9999
        });

        if (result.isConfirmed) {
            loading = true;
            try {
                await api.delete(`/budgets/${categoryId}?month=${currentMonth}`);
                await loadBudgets();

                Swal.fire({
                    icon: 'success',
                    title: 'Budget Deleted',
                    timer: 1500,
                    showConfirmButton: false,
                    zIndex: 9999
                });
            } catch (error) {
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.response?.data?.message || 'Failed to delete budget',
                    zIndex: 9999
                });
            } finally {
                loading = false;
            }
        }
    }

    /**
     * Handle month selector change
     * Reloads budgets and expenses for new month
     */
    async function handleMonthChange() {
        await Promise.all([loadBudgets(), loadExpenses()]);
        // Small delay to ensure budgets are loaded before checking
        setTimeout(() => {
            checkAndOfferImport();
        }, 100);
    }

    /**
     * Handle budget amount input
     * Only allows numeric characters
     */
    function handleAmountInput(e) {
        const value = e.target.value;
        budgetAmount = value.replace(/\D/g, '');
    }

    // ============================================================================
    // REACTIVE COMPUTATIONS
    // ============================================================================

    /**
     * Compute budget cards with progress data
     * This runs automatically when categories, budgets, or expenses change
     */
    $: budgetCards = categories.map((category) => {
        const budget = budgets[category.id] || 0;
        const spent = expensesByCategory[category.id] || 0;
        const progress = getBudgetProgress(budget, spent);
        const status = getBudgetStatus(progress);
        const remaining = budget - spent;

        return {
            category,
            budget,
            spent,
            progress,
            status,
            remaining
        };
    });
</script>

<div class="max-w-3xl mx-auto space-y-6">
  <div class="flex items-center justify-between gap-2">
    <div>
      <h1 class="text-2xl font-bold">Monthly Budget Planning</h1>
      <p class="text-sm text-base-content/70 mt-1">
        Track how much you plan to spend in each category this month.
      </p>
    </div>
  </div>

  <!-- Month Selector -->
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body flex flex-col md:flex-row md:items-end gap-4">
      <div class="form-control w-full md:w-auto">
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Select month</legend>
          <input
            id="month"
            type="month"
            bind:value={currentMonth}
            on:change={handleMonthChange}
            class="input input-bordered w-full md:w-56 border-2"
          />
        </fieldset>
      </div>
      <div class="text-sm text-base-content/70 md:ml-auto">
        <span>Viewing budgets for </span>
        <span class="font-semibold">{formatMonthYear(currentMonth)}</span>
      </div>
    </div>
  </div>

  <!-- Loading / Empty / List -->
  {#if categoriesLoading}
    <div class="flex justify-center py-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if categories.length === 0}
    <div class="alert alert-info">
      <span>No expense categories available yet.</span>
    </div>
  {:else}
    <div class="space-y-4">
      {#each budgetCards as card}
        <div class="card bg-base-100 shadow-xl border-1">
          <div class="card-body gap-3">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h3 class="font-semibold text-base">{card.category.label}</h3>
                {#if card.budget > 0}
                  <p class="text-xs text-base-content/70">
                    Budget: <span class="font-medium">{formatCurrency(card.budget)}</span>
                  </p>
                {:else}
                  <p class="text-xs italic text-base-content/60">
                    No budget set for this category.
                  </p>
                {/if}
              </div>
              <div class="flex gap-2">
                {#if card.budget > 0}
                  <button
                    type="button"
                    class="btn btn-sm btn-success"
                    on:click={() => openBudgetForm(card.category.id)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="btn btn-sm btn-error"
                    on:click={() => deleteBudget(card.category.id)}
                  >
                    Delete
                  </button>
                {:else}
                  <button
                    type="button"
                    class="btn btn-sm btn-primary"
                    on:click={() => openBudgetForm(card.category.id)}
                  >
                    + Set budget
                  </button>
                {/if}
              </div>
            </div>

            {#if card.budget > 0}
              <div class="space-y-2 mt-1">
                <div class="flex items-center justify-between text-xs">
                  <span class="text-base-content/70">
                    Spent: {formatCurrency(card.spent)}
                  </span>
                  <span class="remaining" class:negative={card.remaining < 0}>
                  {card.remaining >= 0 ? 'Remaining: ' : 'Over by: '}{formatCurrency(Math.abs(card.remaining))}
                </span>
                </div>
                <progress
                  class="progress w-full {card.status === 'exceeded'
                    ? 'progress-error'
                    : card.status === 'warning'
                    ? 'progress-warning'
                    : card.status === 'caution'
                    ? 'progress-info'
                    : 'progress-success'}"
                  value={Math.min(card.progress, 100)}
                  max="100"
                />
                <div class="text-right text-[11px] text-base-content/60">
                  {Math.round(card.progress)}% of budget used
                </div>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- Budget Form Modal -->
{#if showBudgetForm}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-semibold text-lg mb-2">
        {editingCategory ? categories.find(c => c.id === editingCategory)?.label : 'Set Budget'}
      </h3>
      <p class="text-sm text-base-content/70 mb-4">
        Enter monthly budget amount for {formatMonthYear(currentMonth)}.
      </p>
      <div class="form-control mb-4">
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Monthly Budget Amount (Rp.)</legend>
          <input
            id="budget-amount"
            type="text"
            bind:value={budgetAmount}
            on:input={handleAmountInput}
            on:keydown={(e) => { if (e.key === 'Enter') saveBudget(); }}
            placeholder="0"
            class="input input-bordered w-full border-2"
          />
        </fieldset>
        {#if budgetAmount}
          <span class="mt-2 text-sm font-semibold text-primary">
            {formatCurrency(budgetAmount)}
          </span>
        {/if}
      </div>
      <div class="modal-action">
        <button class="btn btn-soft" on:click={closeBudgetForm} disabled={loading}>
          Cancel
        </button>
        <button class="btn btn-primary" on:click={saveBudget} disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm mr-1"></span>
            Saving...
          {:else}
            Save budget
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
