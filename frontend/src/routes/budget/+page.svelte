<script>
  // ============================================================================
  // IMPORTS
  // ============================================================================
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { onMount } from 'svelte';
  import '$lib/styles/shared.css';
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

<!-- ============================================================================
     TEMPLATE
     ============================================================================ -->
<div class="budget-page">
  <h1>Monthly Budget Planning</h1>

  <!-- Month Selector -->
  <div class="month-selector">
    <label for="month">Select Month:</label>
    <input
      id="month"
      type="month"
      bind:value={currentMonth}
      on:change={handleMonthChange}
      class="month-input"
    />
  </div>

  <!-- Loading State -->
  {#if categoriesLoading}
    <div class="loading">Loading categories...</div>
  <!-- Empty State -->
  {:else if categories.length === 0}
    <div class="no-categories">No categories available</div>
  <!-- Budget List -->
  {:else}
    <div class="budget-list">
      {#each budgetCards as card}
        <div class="budget-card">
          <!-- Card Header -->
          <div class="budget-header">
            <div class="category-info">
              <h3>{card.category.label}</h3>
              {#if card.budget > 0}
                <span class="budget-amount">Budget: {formatCurrency(card.budget)}</span>
              {:else}
                <span class="no-budget">No budget set</span>
              {/if}
            </div>
            <div class="budget-actions">
              {#if card.budget > 0}
                <button
                  class="btn-text btn-edit"
                  on:click={() => openBudgetForm(card.category.id)}
                >
                  Edit
                </button>
                <button
                  class="btn-text btn-delete"
                  on:click={() => deleteBudget(card.category.id)}
                >
                  Delete
                </button>
              {:else}
                <button
                  class="btn-add"
                  on:click={() => openBudgetForm(card.category.id)}
                >
                  + Set Budget
                </button>
              {/if}
            </div>
          </div>

          <!-- Progress Bar (only shown if budget is set) -->
          {#if card.budget > 0}
            <div class="budget-progress">
              <div class="progress-info">
                <span class="spent">Spent: {formatCurrency(card.spent)}</span>
                <span class="remaining" class:negative={card.remaining < 0}>
                  {card.remaining >= 0 ? 'Remaining: ' : 'Over by: '}{formatCurrency(Math.abs(card.remaining))}
                </span>
              </div>
              <div class="progress-bar-container">
                <div
                  class="progress-bar"
                  class:exceeded={card.progress >= 100}
                  class:warning={card.progress >= 75 && card.progress < 100}
                  class:caution={card.progress >= 50 && card.progress < 75}
                  style="width: {Math.min(card.progress, 100)}%; background-color: {getStatusColor(card.status)}"
                ></div>
              </div>
              <div class="progress-percentage">
                {Math.round(card.progress)}% used
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
  <div class="space-xl"></div>
</div>

<!-- Budget Form Modal -->
{#if showBudgetForm}
  <div class="modal-backdrop" on:click={closeBudgetForm}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h2>{editingCategory ? categories.find(c => c.id === editingCategory)?.label : 'Set Budget'}</h2>
        <button class="close-btn" on:click={closeBudgetForm}>×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label for="budget-amount">Monthly Budget Amount (Rp.)</label>
          <input
            id="budget-amount"
            type="text"
            bind:value={budgetAmount}
            on:input={handleAmountInput}
            on:keydown={(e) => { if (e.key === 'Enter') saveBudget(); }}
            placeholder="0"
            class="form-input"
          />
          {#if budgetAmount}
            <div class="amount-preview">{formatCurrency(budgetAmount)}</div>
          {/if}
        </div>
        <div class="button-group">
          <button class="btn btn-secondary" on:click={closeBudgetForm} disabled={loading}>Cancel</button>
          <button class="btn btn-primary" on:click={saveBudget} disabled={loading}>
            {#if loading}
              <span class="spinner"></span> Saving...
            {:else}
              Save Budget
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- ============================================================================
     STYLES
     ============================================================================ -->
<style>
  .budget-page {
    max-width: 800px;
    margin: 0 auto;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
  }

  .month-selector {
    margin-bottom: 2rem;
    padding: 1rem;
    background: var(--surface);
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }

  .month-selector label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .month-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
  }

  .loading, .no-categories {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .budget-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .budget-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    padding: 1.5rem;
    transition: box-shadow 0.2s;
  }

  .budget-card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .budget-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1rem;
  }

  .category-info h3 {
    font-size: 1.125rem;
    margin-bottom: 0.25rem;
    color: var(--text-primary);
  }

  .budget-amount {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .no-budget {
    font-size: 0.875rem;
    color: var(--text-secondary);
    font-style: italic;
  }

  .budget-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .btn-add {
    padding: 0.5rem 1rem;
    background: var(--primary-color);
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-add:hover {
    background: #4338ca;
  }

  .budget-progress {
    margin-top: 1rem;
  }

  .progress-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
  }

  .spent {
    color: var(--text-secondary);
  }

  .remaining {
    color: var(--success);
    font-weight: 500;
  }

  .remaining.negative {
    color: var(--danger);
  }

  .progress-bar-container {
    width: 100%;
    height: 8px;
    background: var(--border);
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 0.5rem;
  }

  .progress-bar {
    height: 100%;
    transition: width 0.3s ease;
    border-radius: 4px;
  }

  .progress-bar.exceeded {
    background-color: var(--danger) !important;
  }

  .progress-bar.warning {
    background-color: #f59e0b !important;
  }

  .progress-bar.caution {
    background-color: #fbbf24 !important;
  }

  .progress-percentage {
    font-size: 0.75rem;
    color: var(--text-secondary);
    text-align: right;
  }

  .amount-preview {
    margin-top: 0.5rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--primary-color);
  }
</style>
