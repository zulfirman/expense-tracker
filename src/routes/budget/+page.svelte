<script>
    import api from '$lib/api';
    import Swal from 'sweetalert2';
    import {onMount} from 'svelte';
    import '$lib/styles/shared.css';

    let categories = [];
  let budgets = {};
  let currentMonth = new Date().toISOString().slice(0, 7); // YYYY-MM format
  let loading = false;
  let categoriesLoading = true;
  let expensesByCategory = {};
  let budgetCards = [];
  let showBudgetForm = false;
  let editingCategory = null;
  let budgetAmount = '';

  onMount(async () => {
    // Ensure categories are loaded before budgets and expenses
    await loadCategories();
    await Promise.all([loadBudgets(), loadExpenses()]);
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories');
      // Only show active categories in budget page
      categories = response.data
        .filter(cat => cat.isActive !== false)
        .map(cat => ({
          id: cat.id,
          label: cat.name,
          slug: cat.slug
        }));
    } catch (error) {
      console.error('Error loading categories:', error);
      categories = [];
    } finally {
      categoriesLoading = false;
    }
  }

  async function loadBudgets() {
    try {
      const response = await api.get(`/budgets?month=${currentMonth}`);
      const budgetsData = response.data || [];
      budgets = {};
      budgetsData.forEach(budget => {
        budgets[budget.categoryId] = budget.amount;
      });
    } catch (error) {
      console.error('Error loading budgets:', error);
      budgets = {};
    }
  }

  async function loadExpenses() {
    try {
      const response = await api.get(`/expenses/month/${currentMonth}`);
      expensesByCategory = {};
      if (response.data && response.data.categories) {
        response.data.categories.forEach(item => {
          // Backend now returns categoryId directly, so we can map by ID
          const categoryId = item.categoryId;
          if (categoryId) {
            expensesByCategory[categoryId] = (expensesByCategory[categoryId] || 0) + item.total;
          } else {
            console.warn('Missing categoryId for expense category aggregate:', item);
          }
        });
      }
    } catch (error) {
      console.error('Error loading expenses:', error);
      expensesByCategory = {};
    }
  }

  function formatCurrency(amount) {
    if (!amount && amount !== 0) return '';
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }

  function getBudgetProgress(budget, spent) {
    if (!budget) return 0;
    return Math.min((spent / budget) * 100, 100);
  }

  function getBudgetStatus(progress) {
    if (progress >= 100) return 'exceeded';
    if (progress >= 75) return 'warning';
    if (progress >= 50) return 'caution';
    return 'good';
  }

  function getStatusColor(status) {
    switch (status) {
      case 'exceeded': return 'var(--danger)';
      case 'warning': return '#f59e0b';
      case 'caution': return '#fbbf24';
      default: return 'var(--success)';
    }
  }

  function openBudgetForm(categoryId = null) {
    editingCategory = categoryId;
    budgetAmount = categoryId && budgets[categoryId] ? budgets[categoryId].toString() : '';
    showBudgetForm = true;
  }

  function closeBudgetForm() {
    showBudgetForm = false;
    editingCategory = null;
    budgetAmount = '';
  }

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

  async function handleMonthChange() {
    await Promise.all([loadBudgets(), loadExpenses()]);
  }

  function handleAmountInput(e) {
    const value = e.target.value;
    budgetAmount = value.replace(/\D/g, '');
  }

  // Derive budget cards reactively so progress updates when data changes
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

  // Check for budget alerts
  function checkBudgetAlerts() {
    categories.forEach(category => {
      const status = getBudgetStatus(category.id);
      const progress = getBudgetProgress(category.id);
      const budget = budgets[category.id] || 0;
      const spent = expensesByCategory[category.id] || 0;

      if (budget > 0) {
        if (progress >= 100) {
          Swal.fire({
            icon: 'error',
            title: 'Budget Exceeded!',
            html: `<strong>${category.label}</strong><br>You've exceeded your budget of ${formatCurrency(budget)}. Current spending: ${formatCurrency(spent)}`,
            timer: 5000,
            zIndex: 9999
          });
        } else if (progress >= 75 && progress < 100) {
          Swal.fire({
            icon: 'warning',
            title: 'Budget Warning',
            html: `<strong>${category.label}</strong><br>You've used ${Math.round(progress)}% of your budget. Remaining: ${formatCurrency(budget - spent)}`,
            timer: 4000,
            zIndex: 9999
          });
        } else if (progress >= 50 && progress < 75) {
          Swal.fire({
            icon: 'info',
            title: 'Budget Alert',
            html: `<strong>${category.label}</strong><br>You've used ${Math.round(progress)}% of your budget.`,
            timer: 3000,
            zIndex: 9999
          });
        }
      }
    });
  }

  // Auto-check alerts when expenses change
  $: if (Object.keys(expensesByCategory).length > 0) {
    // Only show alerts if user has budgets set
    if (Object.keys(budgets).length > 0) {
      // Debounce to avoid too many alerts
      setTimeout(() => {
        // checkBudgetAlerts(); // Commented out to avoid annoying alerts, can be enabled if needed
      }, 1000);
    }
  }
</script>

<div class="budget-page">
  <h1>Monthly Budget Planning</h1>

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

  {#if categoriesLoading}
    <div class="loading">Loading categories...</div>
  {:else if categories.length === 0}
    <div class="no-categories">No categories available</div>
  {:else}
    <div class="budget-list">
      {#each budgetCards as card}
        <div class="budget-card">
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
  <br><br>
</div>

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

