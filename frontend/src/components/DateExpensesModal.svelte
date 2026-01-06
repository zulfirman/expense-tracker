<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';

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

  let categories = [];
  let templates = [];
  let categoriesLoading = true;
  let templatesLoading = true;
  let showTemplates = false;

  let selectedCategories = [];
  let expenseDate = '';
  let notes = '';
  let amount = '';
  let submitting = false;

  const quickAmounts = [10000, 25000, 50000, 75000, 100000];

  onMount(async () => {
    if (date) {
      expenseDate = date.date;
    }
    await Promise.all([loadExpenses(), loadIncomes(), loadCategories(), loadTemplates()]);
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories');
      // Only show active categories in expense input
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
      // Templates failed to load
    } finally {
      templatesLoading = false;
    }
  }

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

  function formatCurrency(amount) {
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
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
    selectedCategories = [];
    notes = '';
    amount = '';
    if (date) {
      expenseDate = date.date;
    }
  }

  function startAddIncome() {
    showAddIncomeForm = true;
    showAddExpenseForm = false;
    showEditExpenseForm = false;
    showEditIncomeForm = false;
    notes = '';
    amount = '';
    if (date) {
      expenseDate = date.date;
    }
  }

  function cancelAdd() {
    showAddExpenseForm = false;
    showAddIncomeForm = false;
    selectedCategories = [];
    notes = '';
    amount = '';
  }

  function startEditExpense(expense) {
    editingExpense = expense;
    editingIncome = null;
    showEditExpenseForm = true;
    showEditIncomeForm = false;
    showAddExpenseForm = false;
    showAddIncomeForm = false;
    editExpenseForm = {
      categories: [...(expense.categoryIds || expense.categories || [])],
      notes: expense.notes || '',
      amount: expense.amount.toString()
    };
  }

  function startEditIncome(income) {
    editingIncome = income;
    editingExpense = null;
    showEditIncomeForm = true;
    showEditExpenseForm = false;
    showAddExpenseForm = false;
    showAddIncomeForm = false;
    editIncomeForm = {
      notes: income.notes || '',
      amount: income.amount.toString()
    };
  }

  function cancelEdit() {
    editingExpense = null;
    editingIncome = null;
    showEditExpenseForm = false;
    showEditIncomeForm = false;
    editExpenseForm = {
      categories: [],
      notes: '',
      amount: ''
    };
    editIncomeForm = {
      notes: '',
      amount: ''
    };
  }

  function toggleCategory(categoryId) {
    if (selectedCategories.includes(categoryId)) {
      selectedCategories = selectedCategories.filter(id => id !== categoryId);
    } else {
      selectedCategories = [...selectedCategories, categoryId];
    }
  }


  function setQuickAmount(quickAmount) {
    amount = quickAmount.toString();
  }

  function applyTemplate(template) {
    selectedCategories = [...(template.categoryIds || template.categories || [])];
    amount = template.amount.toString();
    notes = template.notes || '';
    showTemplates = false;
  }

  function handleAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    amount = numericValue;
  }

  async function handleAdd() {
    if (selectedCategories.length === 0) {
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

    submitting = true;
    try {
      await api.post('/expenses', {
        categoryIds: selectedCategories,
        date: expenseDate,
        notes: notes,
        amount: parseFloat(amount)
      });

      Swal.fire({
        icon: 'success',
        title: 'Success!',
        text: 'Expense added successfully',
        zIndex: 9999
      });

      cancelAdd();
      await loadAll();
      dispatch('refresh');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to add expense',
        zIndex: 9999
      });
    } finally {
      submitting = false;
    }
  }

  let editExpenseForm = {
    categories: [],
    notes: '',
    amount: ''
  };

  let editIncomeForm = {
    notes: '',
    amount: ''
  };

  function toggleEditCategory(categoryId) {
    if (editExpenseForm.categories.includes(categoryId)) {
      editExpenseForm.categories = editExpenseForm.categories.filter(id => id !== categoryId);
    } else {
      editExpenseForm.categories = [...editExpenseForm.categories, categoryId];
    }
  }

  async function handleUpdateExpense() {
    if (editExpenseForm.categories.length === 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please select at least one category',
        zIndex: 9999
      });
      return;
    }

    if (!editExpenseForm.amount || parseFloat(editExpenseForm.amount) <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please enter a valid amount',
        zIndex: 9999
      });
      return;
    }

    try {
      await api.put(`/expenses/${editingExpense.id}`, {
        categoryIds: editExpenseForm.categories,
        notes: editExpenseForm.notes,
        amount: parseFloat(editExpenseForm.amount)
      });

      Swal.fire({
        icon: 'success',
        title: 'Success!',
        text: 'Expense updated successfully',
        zIndex: 9999
      });

      cancelEdit();
      await loadAll();
      dispatch('refresh');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to update expense',
        zIndex: 9999
      });
    }
  }

  async function handleUpdateIncome() {
    if (!editIncomeForm.amount || parseFloat(editIncomeForm.amount) <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please enter a valid amount',
        zIndex: 9999
      });
      return;
    }

    try {
      await api.put(`/income/${editingIncome.id}`, {
        date: date.date,
        notes: editIncomeForm.notes,
        amount: parseFloat(editIncomeForm.amount)
      });

      Swal.fire({
        icon: 'success',
        title: 'Success!',
        text: 'Income updated successfully',
        zIndex: 9999
      });

      cancelEdit();
      await loadAll();
      dispatch('refresh');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to update income',
        zIndex: 9999
      });
    }
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

  async function handleAddIncome() {
    if (!amount || parseFloat(amount) <= 0) {
      Swal.fire({
        icon: 'warning',
        title: 'Warning',
        text: 'Please enter a valid amount',
        zIndex: 9999
      });
      return;
    }

    submitting = true;
    try {
      await api.post('/income', {
        date: expenseDate,
        notes: notes,
        amount: parseFloat(amount)
      });

      Swal.fire({
        icon: 'success',
        title: 'Success!',
        text: 'Income added successfully',
        zIndex: 9999
      });

      cancelAdd();
      await loadAll();
      dispatch('refresh');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to add income',
        zIndex: 9999
      });
    } finally {
      submitting = false;
    }
  }

  function handleEditExpenseAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    editExpenseForm.amount = numericValue;
  }

  function handleEditIncomeAmountInput(e) {
    const value = e.target.value;
    const numericValue = value.replace(/\D/g, '');
    editIncomeForm.amount = numericValue;
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (showAddExpenseForm) {
        handleAdd();
      } else if (showAddIncomeForm) {
        handleAddIncome();
      } else if (showEditExpenseForm) {
        handleUpdateExpense();
      } else if (showEditIncomeForm) {
        handleUpdateIncome();
      }
    }
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
        <form class="add-form" on:submit|preventDefault={handleAdd}>
          <h3>Add Expense</h3>

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
            <label>Category {#if selectedCategories.length > 0}<span class="selected-count">({selectedCategories.length} selected)</span>{/if}</label>
            {#if categoriesLoading}
              <div class="loading-categories">Loading categories...</div>
            {:else}
              <div class="checkbox-group">
                {#each categories as category}
                  <label class="checkbox-label" class:selected={selectedCategories.includes(category.id)}>
                    <input
                      type="checkbox"
                      checked={selectedCategories.includes(category.id)}
                      on:change={() => toggleCategory(category.id)}
                    />
                    <span>{category.label}</span>
                  </label>
                {/each}
              </div>
            {/if}
          </div>

          <div class="form-group">
            <label for="add-amount">Amount (Rp.)</label>
            <input
              id="add-amount"
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
            <label for="add-notes">Notes</label>
            <textarea
              id="add-notes"
              bind:value={notes}
              placeholder="Add notes (optional)"
              class="form-textarea"
              rows="3"
              on:keydown={handleKeyDown}
            ></textarea>
          </div>

          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={cancelAdd}>Cancel</button>
            <button type="submit" class="btn btn-primary" disabled={submitting}>
              {#if submitting}
                <span class="spinner"></span> Adding...
              {:else}
                Add
              {/if}
            </button>
          </div>
        </form>
      {:else if showAddIncomeForm}
        <form class="add-form" on:submit|preventDefault={handleAddIncome}>
          <h3>Add Income</h3>

          <div class="form-group">
            <label for="add-income-amount">Amount (Rp.)</label>
            <input
              id="add-income-amount"
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
            <label for="add-income-notes">Notes</label>
            <textarea
              id="add-income-notes"
              bind:value={notes}
              placeholder="Add notes (optional)"
              class="form-textarea"
              rows="3"
              on:keydown={handleKeyDown}
            ></textarea>
          </div>

          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={cancelAdd}>Cancel</button>
            <button type="submit" class="btn btn-primary" disabled={submitting}>
              {#if submitting}
                <span class="spinner"></span> Adding...
              {:else}
                Add
              {/if}
            </button>
          </div>
        </form>
      {:else if showEditExpenseForm && editingExpense}
        <form class="edit-form" on:submit|preventDefault={handleUpdateExpense}>
          <h3>Edit Expense</h3>
          
          <div class="form-group">
            <label>Category {#if editExpenseForm.categories.length > 0}<span class="selected-count">({editExpenseForm.categories.length} selected)</span>{/if}</label>
            {#if categoriesLoading}
              <div class="loading-categories">Loading categories...</div>
            {:else}
              <div class="checkbox-group">
                {#each categories as category}
                  <label class="checkbox-label" class:selected={editExpenseForm.categories.includes(category.id)}>
                    <input
                      type="checkbox"
                      checked={editExpenseForm.categories.includes(category.id)}
                      on:change={() => toggleEditCategory(category.id)}
                    />
                    <span>{category.label}</span>
                  </label>
                {/each}
              </div>
            {/if}
          </div>

          <div class="form-group">
            <label for="edit-expense-amount">Amount (Rp.)</label>
            <input
              id="edit-expense-amount"
              type="text"
              bind:value={editExpenseForm.amount}
              on:input={handleEditExpenseAmountInput}
              on:keydown={handleKeyDown}
              class="form-input"
              inputmode="numeric"
            />
            {#if editExpenseForm.amount}
              <div class="amount-preview">{formatCurrency(editExpenseForm.amount)}</div>
            {/if}
          </div>

          <div class="form-group">
            <label for="edit-expense-notes">Notes</label>
            <textarea
              id="edit-expense-notes"
              bind:value={editExpenseForm.notes}
              on:keydown={handleKeyDown}
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>

          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={cancelEdit}>Cancel</button>
            <button type="submit" class="btn btn-primary">Update</button>
          </div>
        </form>
      {:else if showEditIncomeForm && editingIncome}
        <form class="edit-form" on:submit|preventDefault={handleUpdateIncome}>
          <h3>Edit Income</h3>

          <div class="form-group">
            <label for="edit-income-amount">Amount (Rp.)</label>
            <input
              id="edit-income-amount"
              type="text"
              bind:value={editIncomeForm.amount}
              on:input={handleEditIncomeAmountInput}
              on:keydown={handleKeyDown}
              class="form-input"
              inputmode="numeric"
            />
            {#if editIncomeForm.amount}
              <div class="amount-preview">{formatCurrency(editIncomeForm.amount)}</div>
            {/if}
          </div>

          <div class="form-group">
            <label for="edit-income-notes">Notes</label>
            <textarea
              id="edit-income-notes"
              bind:value={editIncomeForm.notes}
              on:keydown={handleKeyDown}
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>

          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={cancelEdit}>Cancel</button>
            <button type="submit" class="btn btn-primary">Update</button>
          </div>
        </form>
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

  .checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    padding: 0.5rem;
    border: 2px solid var(--border);
    border-radius: 0.5rem;
    transition: all 0.2s;
  }

  .checkbox-label:hover {
    background-color: var(--background);
    border-color: var(--primary-color);
  }

  .checkbox-label.selected {
    background-color: rgba(79, 70, 229, 0.1);
    border-color: var(--primary-color);
  }

  .checkbox-label input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
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
