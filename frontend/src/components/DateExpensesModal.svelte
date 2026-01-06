<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { formatCurrency } from '$lib/utils/currency';
  import InputExpenses from './InputExpenses.svelte';
  import InputIncome from './InputIncome.svelte';

  const dispatch = createEventDispatcher();

  export let date = null;

  let expenses = [];
  let incomes = [];
  let loading = true;
  let editingExpense = null;
  let editingIncome = null;
  let showEditExpenseForm = false;
  let showEditIncomeForm = false;
  let showAddExpenseForm = false;
  let showAddIncomeForm = false;

  onMount(async () => {
    await Promise.all([loadExpenses(), loadIncomes()]);
  });

  async function loadIncomes() {
    if (!date) return;
    
    try {
      const response = await api.get(`/income/date/${date.date}`);
      incomes = response.data;
    } catch (error) {
      // Incomes failed to load
    }
  }

  async function loadExpenses() {
    if (!date) return;
    
    loading = true;
    try {
      const response = await api.get(`/expenses/date/${date.date}`);
      expenses = response.data;
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'Failed to load expenses',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  async function loadAll() {
    await Promise.all([loadExpenses(), loadIncomes()]);
  }


  function formatDate(dateString) {
    const d = new Date(dateString);
    return d.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  }

  function close() {
    dispatch('close');
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      close();
    }
  }

  function startAddExpense() {
    showAddExpenseForm = true;
    showAddIncomeForm = false;
    showEditExpenseForm = false;
    showEditIncomeForm = false;
  }

  function startAddIncome() {
    showAddIncomeForm = true;
    showAddExpenseForm = false;
    showEditExpenseForm = false;
    showEditIncomeForm = false;
  }

  function cancelAdd() {
    showAddExpenseForm = false;
    showAddIncomeForm = false;
  }

  function startEditExpense(expense) {
    editingExpense = expense;
    editingIncome = null;
    showEditExpenseForm = true;
    showEditIncomeForm = false;
    showAddExpenseForm = false;
    showAddIncomeForm = false;
  }

  function startEditIncome(income) {
    editingIncome = income;
    editingExpense = null;
    showEditIncomeForm = true;
    showEditExpenseForm = false;
    showAddExpenseForm = false;
    showAddIncomeForm = false;
  }

  function cancelEdit() {
    editingExpense = null;
    editingIncome = null;
    showEditExpenseForm = false;
    showEditIncomeForm = false;
  }

  async function handleDeleteExpense(expenseId) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'Delete Expense?',
      text: 'Are you sure you want to delete this expense?',
      showCancelButton: true,
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      confirmButtonColor: '#ef4444',
      zIndex: 9999
    });

    if (result.isConfirmed) {
      try {
        await api.delete(`/expenses/${expenseId}`);
        
        Swal.fire({
          icon: 'success',
          title: 'Deleted!',
          text: 'Expense deleted successfully',
          zIndex: 9999
        });

        await loadAll();
        dispatch('refresh');
      } catch (error) {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to delete expense',
          zIndex: 9999
        });
      }
    }
  }

  async function handleDeleteIncome(incomeId) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'Delete Income?',
      text: 'Are you sure you want to delete this income?',
      showCancelButton: true,
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      confirmButtonColor: '#ef4444',
      zIndex: 9999
    });

    if (result.isConfirmed) {
      try {
        await api.delete(`/income/${incomeId}`);
        
        Swal.fire({
          icon: 'success',
          title: 'Deleted!',
          text: 'Income deleted successfully',
          zIndex: 9999
        });

        await loadAll();
        dispatch('refresh');
      } catch (error) {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to delete income',
          zIndex: 9999
        });
      }
    }
  }

  async function handleExpenseSuccess() {
    await loadAll();
    dispatch('refresh');
    showAddExpenseForm = false;
    showEditExpenseForm = false;
    editingExpense = null;
  }

  async function handleIncomeSuccess() {
    await loadAll();
    dispatch('refresh');
    showAddIncomeForm = false;
    showEditIncomeForm = false;
    editingIncome = null;
  }
</script>

<div class="modal-backdrop" on:click={handleBackdropClick}>
  <div class="modal-content">
    <div class="modal-header">
      <h2>{date ? formatDate(date.date) : ''}</h2>
      <button class="close-btn" on:click={close}>×</button>
    </div>

    <div class="modal-body">
      {#if loading}
        <div class="loading">Loading...</div>
      {:else if showAddExpenseForm}
        <InputExpenses
          fixedDate={date?.date}
          showTitle={false}
          showCancel={true}
          submitLabel="Add"
          onSuccess={handleExpenseSuccess}
          onCancel={cancelAdd}
        />
      {:else if showAddIncomeForm}
        <InputIncome
          fixedDate={date?.date}
          showTitle={false}
          showBalance={false}
          showCancel={true}
          submitLabel="Add"
          onSuccess={handleIncomeSuccess}
          onCancel={cancelAdd}
        />
      {:else if showEditExpenseForm && editingExpense}
        <InputExpenses
          expenseId={editingExpense.id}
          initialData={{
            categoryIds: editingExpense.categoryIds || editingExpense.categories?.map(c => typeof c === 'object' ? c.id : c) || [],
            date: date?.date || editingExpense.date,
            notes: editingExpense.notes || '',
            amount: editingExpense.amount
          }}
          fixedDate={date?.date}
          showTitle={false}
          showCancel={true}
          submitLabel="Update"
          onSuccess={handleExpenseSuccess}
          onCancel={cancelEdit}
        />
      {:else if showEditIncomeForm && editingIncome}
        <InputIncome
          incomeId={editingIncome.id}
          initialData={{
            date: date?.date || editingIncome.date,
            notes: editingIncome.notes || '',
            amount: editingIncome.amount
          }}
          fixedDate={date?.date}
          showTitle={false}
          showBalance={false}
          showCancel={true}
          submitLabel="Update"
          onSuccess={handleIncomeSuccess}
          onCancel={cancelEdit}
        />
      {:else}
        {@const totalIncome = incomes.reduce((sum, inc) => sum + inc.amount, 0)}
        {@const totalExpenses = expenses.reduce((sum, exp) => sum + exp.amount, 0)}
        {@const netTotal = totalIncome - totalExpenses}
        
        <div class="date-summary">
          <div class="summary-item income">
            <span class="summary-label">Total Income</span>
            <span class="summary-amount">{formatCurrency(totalIncome)}</span>
          </div>
          <div class="summary-item expense">
            <span class="summary-label">Total Expenses</span>
            <span class="summary-amount">{formatCurrency(totalExpenses)}</span>
          </div>
          <div class="summary-item net" class:positive={netTotal > 0} class:negative={netTotal < 0}>
            <span class="summary-label">Net Total</span>
            <span class="summary-amount">{formatCurrency(Math.abs(netTotal))}</span>
          </div>
        </div>

        <div class="transactions-section">
          <div class="section-header">
            <h3>Income</h3>
            <button class="btn btn-primary add-btn" on:click={startAddIncome}>+ Add</button>
          </div>
          {#if incomes.length === 0}
            <div class="no-items">No income for this date</div>
          {:else}
            {#each incomes as income}
              <div class="transaction-item income-item">
                <div class="transaction-info">
                  <div class="transaction-amount income-amount">{formatCurrency(income.amount)}</div>
                  {#if income.notes}
                    <div class="transaction-notes">{income.notes}</div>
                  {/if}
                </div>
                <div class="transaction-actions">
                  <button class="action-btn edit-btn" on:click={() => startEditIncome(income)}>
                    Edit
                  </button>
                  <button class="action-btn delete-btn" on:click={() => handleDeleteIncome(income.id)}>
                    Delete
                  </button>
                </div>
              </div>
            {/each}
          {/if}
        </div>

        <div class="transactions-section">
          <div class="section-header">
            <h3>Expenses</h3>
            <button class="btn btn-primary add-btn" on:click={startAddExpense}>+ Add</button>
          </div>
          {#if expenses.length === 0}
            <div class="no-items">No expenses for this date</div>
          {:else}
            {#each expenses as expense}
              <div class="transaction-item expense-item">
                <div class="transaction-info">
                  <div class="transaction-categories">
                    {#each expense.categories as cat}
                      <span class="category-badge">{cat}</span>
                    {/each}
                  </div>
                  <div class="transaction-amount expense-amount">{formatCurrency(expense.amount)}</div>
                  {#if expense.notes}
                    <div class="transaction-notes">{expense.notes}</div>
                  {/if}
                </div>
                <div class="transaction-actions">
                  <button class="action-btn edit-btn" on:click={() => startEditExpense(expense)}>
                    Edit
                  </button>
                  <button class="action-btn delete-btn" on:click={() => handleDeleteExpense(expense.id)}>
                    Delete
                  </button>
                </div>
              </div>
            {/each}
          {/if}
        </div>
      {/if}
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
    border-radius: 1.25rem;
    width: 100%;
    max-width: 650px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    overflow: hidden;
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
    text-transform: capitalize;
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
    overflow-y: auto;
    flex: 1;
    scroll-behavior: smooth;
  }

  .loading,
  .no-data {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .no-data button {
    margin-top: 1rem;
  }

  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .list-header h3 {
    font-size: 1.125rem;
    color: var(--text-primary);
    margin: 0;
  }


  .expenses-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .date-summary {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 2rem;
    padding: 1.25rem;
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    border-radius: 1rem;
    border: 1px solid var(--border);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }

  .summary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.875rem 1rem;
    border-radius: 0.625rem;
    transition: all 0.2s;
  }

  .summary-item.income {
    background: rgba(16, 185, 129, 0.08);
    border-left: 3px solid #10b981;
  }

  .summary-item.expense {
    background: rgba(239, 68, 68, 0.08);
    border-left: 3px solid #ef4444;
  }

  .summary-item.net {
    background: var(--surface);
    border: 2px solid var(--border);
    font-weight: 600;
    margin-top: 0.5rem;
    padding: 1rem;
  }

  .summary-item.net.positive {
    border-color: #10b981;
    background: rgba(16, 185, 129, 0.1);
  }

  .summary-item.net.negative {
    border-color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  .summary-label {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .summary-amount {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--text-primary);
  }

  .summary-item.net .summary-amount {
    font-size: 1.25rem;
  }

  .transactions-section {
    margin-bottom: 2rem;
  }

  .transactions-section:last-child {
    margin-bottom: 0;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.25rem;
    padding-bottom: 0.75rem;
    border-bottom: 2px solid var(--border);
  }

  .section-header h3 {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .add-btn {
    padding: 0.375rem 0.75rem;
    font-size: 0.75rem;
    min-width: auto;
    font-weight: 500;
      max-width: 50%;
  }

  .no-items {
    padding: 2rem 1.5rem;
    text-align: center;
    color: var(--text-secondary);
    background: var(--background);
    border-radius: 0.75rem;
    border: 2px dashed var(--border);
    font-size: 0.875rem;
  }

  .transaction-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 1.25rem;
    background: var(--surface);
    border-radius: 0.75rem;
    border: 1px solid var(--border);
    margin-bottom: 1rem;
    transition: all 0.2s;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  }

  .transaction-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-color: var(--primary-color);
  }

  .transaction-item.income-item {
    border-left: 4px solid #10b981;
    background: linear-gradient(to right, rgba(16, 185, 129, 0.05), var(--surface));
  }

  .transaction-item.expense-item {
    border-left: 4px solid #ef4444;
    background: linear-gradient(to right, rgba(239, 68, 68, 0.05), var(--surface));
  }

  .transaction-info {
    flex: 1;
  }

  .transaction-categories {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .transaction-amount {
    font-size: 1.25rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
  }

  .transaction-amount.income-amount {
    color: #10b981;
  }

  .transaction-amount.expense-amount {
    color: #ef4444;
  }

  .transaction-notes {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 0.25rem;
  }

  .transaction-actions {
    display: flex;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .action-btn:hover {
    transform: scale(1.05);
  }

  .expense-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 1rem;
    background: var(--background);
    border-radius: 0.5rem;
    border: 1px solid var(--border);
  }

  .expense-info {
    flex: 1;
  }

  .expense-categories {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .category-badge {
    padding: 0.25rem 0.75rem;
    background: var(--primary-color);
    color: white;
    border-radius: 1rem;
    font-size: 0.75rem;
    text-transform: capitalize;
  }

  .expense-amount {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--primary-color);
    margin-bottom: 0.25rem;
  }

  .expense-notes {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 0.25rem;
  }

  .expense-actions {
    display: flex;
    gap: 0.5rem;
  }

  .edit-btn {
    background: var(--primary-color);
    color: white;
  }

  .edit-btn:hover {
    background: #4338ca;
  }

  .delete-btn {
    background: var(--danger);
    color: white;
  }

  .delete-btn:hover {
    background: #dc2626;
  }

  .add-form,
  .edit-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .add-form h3,
  .edit-form h3 {
    font-size: 1.125rem;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-weight: 500;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 0.5rem;
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

  .form-input,
  .form-textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
  }

  .form-input:focus,
  .form-textarea:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .amount-preview {
    font-size: 1rem;
    font-weight: 600;
    color: var(--primary-color);
  }

  .button-group {
    display: flex;
    gap: 1rem;
    margin-top: 0.5rem;
  }

  .btn {
    flex: 1;
    padding: 0.75rem;
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

  .btn-secondary {
    background-color: var(--surface);
    color: var(--text-primary);
    border: 1px solid var(--border);
  }

  .btn-secondary:hover {
    background-color: var(--background);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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

  /* Suggestions */
  .suggestions-label {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
  }

  .suggestions-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .suggestion-btn {
    padding: 0.5rem 1rem;
    background: rgba(79, 70, 229, 0.1);
    border: 1px solid var(--primary-color);
    border-radius: 1.5rem;
    color: var(--primary-color);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .suggestion-btn:hover {
    background: var(--primary-color);
    color: white;
    transform: scale(1.05);
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
</style>
