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

<div class="modal modal-open z-[2100]" on:click={handleBackdropClick}>
  <div class="modal-box w-11/12 max-w-4xl" on:click|stopPropagation>
    <div class="flex items-start justify-between gap-3 mb-4">
      <div>
        <p class="text-xs uppercase tracking-wide text-base-content/60">Date</p>
        <h2 class="text-2xl font-bold">{date ? formatDate(date.date) : ''}</h2>
      </div>
      <button class="btn btn-ghost btn-sm" on:click={close}>✕</button>
    </div>

    <div class="max-h-[70vh] overflow-y-auto pr-1">
      {#if loading}
        <div class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
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
            categoryIds: editingIncome.categoryIds || editingIncome.categories?.map(c => typeof c === 'object' ? c.id : c) || [],
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

        <div class="stats shadow-sm bg-base-200 border border-base-300 mb-6">
          <div class="stat">
            <div class="stat-title">Total Income</div>
            <div class="stat-value text-success text-xl">{formatCurrency(totalIncome)}</div>
          </div>
          <div class="stat">
            <div class="stat-title">Total Expenses</div>
            <div class="stat-value text-error text-xl">{formatCurrency(totalExpenses)}</div>
          </div>
          <div class="stat">
            <div class="stat-title">Net Total</div>
            <div class="stat-value text-xl" class:text-success={netTotal >= 0} class:text-error={netTotal < 0}>
              {formatCurrency(Math.abs(netTotal))}
            </div>
          </div>
        </div>

        <div class="card bg-base-100 border border-base-300 shadow-sm mb-6">
          <div class="card-body p-4 space-y-3">
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-lg">Income</h3>
              <button class="btn btn-primary btn-sm" on:click={startAddIncome}>+ Add</button>
            </div>
            {#if incomes.length === 0}
              <div class="alert alert-info"><span>No income for this date</span></div>
            {:else}
              <div class="space-y-2">
                {#each incomes as income}
                  <div class="flex items-start justify-between gap-3 p-3 rounded-lg border border-base-200 bg-base-200/50">
                    <div class="space-y-1">
                      <div class="text-lg font-semibold text-success">{formatCurrency(income.amount)}</div>
                      {#if income.notes}
                        <div class="text-sm text-base-content/70">{income.notes}</div>
                      {/if}
                    </div>
                    <div class="flex gap-2">
                      <button class="btn btn-soft btn-xs" on:click={() => startEditIncome(income)}>Edit</button>
                      <button class="btn btn-error btn-xs" on:click={() => handleDeleteIncome(income.id)}>Delete</button>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>

        <div class="card bg-base-100 border border-base-300 shadow-sm">
          <div class="card-body p-4 space-y-3">
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-lg">Expenses</h3>
              <button class="btn btn-primary btn-sm" on:click={startAddExpense}>+ Add</button>
            </div>
            {#if expenses.length === 0}
              <div class="alert alert-info"><span>No expenses for this date</span></div>
            {:else}
              <div class="space-y-2">
                {#each expenses as expense}
                  <div class="flex items-start justify-between gap-3 p-3 rounded-lg border border-base-200 bg-base-200/50">
                    <div class="space-y-1">
                      <div class="flex flex-wrap gap-1">
                        {#each expense.categories as cat}
                          <span class="badge badge-outline">
                            {#if typeof cat === 'string'}
                              {cat}
                            {:else}
                              {cat?.name}
                            {/if}
                          </span>
                        {/each}
                      </div>
                      <div class="text-lg font-semibold text-error">{formatCurrency(expense.amount)}</div>
                      {#if expense.notes}
                        <div class="text-sm text-base-content/70">{expense.notes}</div>
                      {/if}
                    </div>
                    <div class="flex gap-2">
                      <button class="btn btn-soft btn-xs" on:click={() => startEditExpense(expense)}>Edit</button>
                      <button class="btn btn-error btn-xs" on:click={() => handleDeleteExpense(expense.id)}>Delete</button>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .category-badge {
    text-transform: capitalize;
  }
</style>
