<script>
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { onMount } from 'svelte';
  import CalendarView from './CalendarView.svelte';
  import MonthDetailModal from './MonthDetailModal.svelte';
  import DateExpensesModal from './DateExpensesModal.svelte';
  import DatePicker from '$lib/components/DatePicker.svelte';
  import '$lib/styles/shared.css';

  let months = [];
  let selectedMonth = null;
  let selectedDate = null;
  let showMonthDetail = false;
  let showDateModal = false;
  let loading = false;
  let balance = { amount: 0 };
  let balanceLoading = true;
  let savedScrollPosition = 0;
  let calendarWrapper;
  
  // Search functionality
  let searchQuery = '';
  let searchCategory = '';
  let dateFrom = '';
  let dateTo = '';
  let showSearch = false;
  let allExpenses = [];
  let filteredExpenses = [];
  let categories = [];
  let searchLoading = false;

  function formatCurrency(value) {
    if (!value && value !== 0) return '';
    let numericValue;
    if (typeof value === 'string') {
      numericValue = value.replace(/\D/g, '');
      if (!numericValue) return '';
      numericValue = parseFloat(numericValue);
    } else {
      numericValue = value;
    }
    if (isNaN(numericValue)) return '';
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(numericValue);
  }

  onMount(async () => {
    await Promise.all([loadMonths(), loadBalance(), loadCategories()]);
  });

  async function loadCategories() {
    try {
      const response = await api.get('/categories');
      // Only show active categories in history filter
      categories = response.data
        .filter(cat => cat.isActive !== false)
        .map(cat => ({
          id: cat.id,
          label: cat.name,
          slug: cat.slug
        }));
    } catch (error) {
      categories = [];
    }
  }

  async function searchExpenses() {
    if (!searchQuery && !searchCategory && !dateFrom && !dateTo) {
      filteredExpenses = [];
      return;
    }

    searchLoading = true;
    try {
      const params = new URLSearchParams();
      if (searchQuery) params.append('q', searchQuery);
      if (searchCategory) params.append('categoryId', searchCategory);
      if (dateFrom) params.append('dateFrom', dateFrom);
      if (dateTo) params.append('dateTo', dateTo);

      const response = await api.get(`/expenses/search?${params.toString()}`);
      filteredExpenses = response.data || [];
    } catch (error) {
      filteredExpenses = [];
    } finally {
      searchLoading = false;
    }
  }

  function clearSearch() {
    searchQuery = '';
    searchCategory = '';
    dateFrom = '';
    dateTo = '';
    filteredExpenses = [];
    showSearch = false;
  }

  function toggleSearch() {
    showSearch = !showSearch;
    if (!showSearch) {
      clearSearch();
    }
  }

  $: if (searchQuery || searchCategory || dateFrom || dateTo) {
    searchExpenses();
  }

  async function loadBalance() {
    try {
      const response = await api.get('/income/balance');
      balance = response.data;
    } catch (error) {
      // Balance failed to load
    } finally {
      balanceLoading = false;
    }
  }

  async function loadMonths(silent = false) {
    if (!silent) {
      loading = true;
    }
    try {
      const response = await api.get('/expenses/months');
      months = response.data;
    } catch (error) {
      // Months failed to load
    } finally {
      if (!silent) {
        loading = false;
      }
    }
  }

  function getTotalExpenses() {
    if (!months || months.length === 0) {
      return 0;
    }
    const total = months.reduce((sum, month) => {
      // Ensure month.total is parsed as a number
      const monthTotal = typeof month.total === 'number' ? month.total : parseFloat(month.total) || 0;
      return sum + monthTotal;
    }, 0);
    return total;
  }

  function getRemainingBalance() {
    if (balanceLoading || !balance) {
      return 0;
    }
    const totalExpenses = getTotalExpenses();
    // Ensure balance.amount is parsed as a number
    const balanceAmount = typeof balance.amount === 'number' ? balance.amount : parseFloat(balance.amount) || 0;
    const remaining = balanceAmount - totalExpenses;
    return remaining;
  }

  function handleMonthClick(month) {
    // Save scroll position before opening modal
    if (calendarWrapper) {
      savedScrollPosition = calendarWrapper.scrollTop;
    }
    selectedMonth = month;
    showMonthDetail = true;
  }

  function handleDateClick(date) {
    // Save scroll position before opening modal
    if (calendarWrapper) {
      savedScrollPosition = calendarWrapper.scrollTop;
    }
    selectedDate = date;
    showDateModal = true;
  }

  function closeMonthDetail() {
    showMonthDetail = false;
    selectedMonth = null;
    // Restore scroll position after closing modal
    setTimeout(() => {
      if (calendarWrapper) {
        calendarWrapper.scrollTop = savedScrollPosition;
      }
    }, 10);
  }

  async function closeDateModal() {
    showDateModal = false;
    selectedDate = null;
    // Reload data in background to update balance, but preserve scroll position
    const currentScroll = calendarWrapper ? calendarWrapper.scrollTop : 0;
    // await Promise.all([loadMonths(), loadBalance()]);
    // Restore scroll position after data reload
    setTimeout(() => {
      if (calendarWrapper) {
        calendarWrapper.scrollTop = currentScroll;
      }
    }, 10);
  }

  function handleLoadMore(direction) {
    // Don't reload, just let user continue scrolling
    // The data is already loaded from the initial fetch (3 months back + current)
  }
</script>

<div class="expenses-history">
  <div class="header-section">
    <h1>Expenses History</h1>
    <button class="search-toggle-btn" on:click={toggleSearch} class:active={showSearch}>
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.35-4.35"></path>
      </svg>
      Search
    </button>
  </div>


  <!-- Balance Display -->
  {#if !balanceLoading && !loading}
    {@const totalExpenses = getTotalExpenses()}
    {@const remainingBalance = getRemainingBalance()}
    <div class="balance-section">
      <div class="balance-card">
        <div class="balance-label">Current Balance</div>
        <div class="balance-amount">{formatCurrency(remainingBalance)}</div>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="calendar-wrapper" bind:this={calendarWrapper}>
      <CalendarView
        {months}
        on:monthClick={(e) => handleMonthClick(e.detail)}
        on:dateClick={(e) => handleDateClick(e.detail)}
        on:loadMore={(e) => handleLoadMore(e.detail)}
      />
    </div>
  {/if}
</div>

{#if showMonthDetail && selectedMonth}
  <MonthDetailModal
    month={selectedMonth}
    on:close={closeMonthDetail}
  />
{/if}

{#if showDateModal && selectedDate}
  <DateExpensesModal
    date={selectedDate}
    on:close={closeDateModal}
    on:refresh={async () => {
      // Refresh data but preserve scroll position and avoid resetting calendar scroll
      const currentScroll = calendarWrapper ? calendarWrapper.scrollTop : 0;
      await Promise.all([loadMonths(true), loadBalance()]);
      setTimeout(() => {
        if (calendarWrapper) {
          calendarWrapper.scrollTop = currentScroll;
        }
      }, 50);
    }}
  />
{/if}

{#if showSearch}
  <div class="modal-backdrop" on:click={() => showSearch = false}>
    <div class="modal-content search-modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Search Expenses</h2>
        <button class="close-btn" on:click={() => showSearch = false}>×</button>
      </div>
      <div class="modal-body">
        <div class="search-form">
          <div class="form-group">
            <label for="search-query">Search (notes)</label>
            <input
              id="search-query"
              type="text"
              bind:value={searchQuery}
              placeholder="Search expenses..."
              class="form-input"
              on:keydown={(e) => { if (e.key === 'Enter') searchExpenses(); }}
            />
          </div>
          <div class="form-group">
            <label for="search-category">Category</label>
            <select id="search-category" bind:value={searchCategory} class="form-input">
              <option value="">All Categories</option>
              {#each categories as category}
                <option value={category.id}>{category.label}</option>
              {/each}
            </select>
          </div>
          <div class="date-range">
            <div class="form-group">
              <label for="date-from">From Date</label>
              <DatePicker
                id="date-from"
                bind:value={dateFrom}
                placeholder="From date"
              />
            </div>
            <div class="form-group">
              <label for="date-to">To Date</label>
              <DatePicker
                id="date-to"
                bind:value={dateTo}
                placeholder="To date"
              />
            </div>
          </div>
          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={clearSearch}>Clear</button>
            <button type="button" class="btn btn-primary" on:click={searchExpenses}>Search</button>
          </div>
        </div>
        
        {#if searchLoading}
          <div class="search-loading">Searching...</div>
        {:else if filteredExpenses.length > 0}
          <div class="search-results">
            <h3>Search Results ({filteredExpenses.length})</h3>
            <div class="expenses-list">
              {#each filteredExpenses as expense}
                <div class="expense-item">
                  <div class="expense-info">
                    <div class="expense-date">{new Date(expense.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</div>
                    <div class="expense-category">{expense.category}</div>
                    {#if expense.notes}
                      <div class="expense-notes">{expense.notes}</div>
                    {/if}
                  </div>
                  <div class="expense-amount">{formatCurrency(expense.amount)}</div>
                </div>
              {/each}
            </div>
          </div>
        {:else if searchQuery || searchCategory || dateFrom || dateTo}
          <div class="no-results">No expenses found matching your search criteria.</div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .expenses-history {
    max-width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
    flex-shrink: 0;
  }

  .calendar-wrapper {
    flex: 1;
    overflow: hidden;
    min-height: 0;
    position: relative;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .balance-section {
    margin-bottom: 1.5rem;
    flex-shrink: 0;
  }

  .balance-card {
    background: linear-gradient(135deg, var(--primary-color) 0%, #6366f1 100%);
    border-radius: 1rem;
    padding: 2rem;
    box-shadow: 0 8px 16px rgba(79, 70, 229, 0.3);
    text-align: center;
  }

  .balance-label {
    font-size: 1rem;
    color: rgba(255, 255, 255, 0.9);
    font-weight: 500;
    margin-bottom: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .balance-amount {
    font-size: 1.75rem;
    font-weight: 700;
    color: white;
    line-height: 1.2;
  }

  .balance-details {
    margin-top: 1.5rem;
    padding-top: 1rem;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .balance-detail-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.875rem;
    color: rgba(255, 255, 255, 0.9);
  }

  .balance-detail-item .expense-text {
    color: rgba(255, 200, 200, 1);
    font-weight: 600;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    flex-shrink: 0;
  }

  .search-toggle-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    color: var(--text-primary);
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s;
  }

  .search-toggle-btn:hover {
    background: var(--background);
    border-color: var(--primary-color);
  }

  .search-toggle-btn.active {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .search-modal {
    max-width: 600px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .search-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .date-range {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .search-loading, .no-results {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .search-results {
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border);
  }

  .search-results h3 {
    font-size: 1rem;
    margin-bottom: 1rem;
    color: var(--text-primary);
  }

  .expenses-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
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

  .expense-date {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.25rem;
  }

  .expense-category {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 0.25rem;
  }

  .expense-notes {
    font-size: 0.875rem;
    color: var(--text-secondary);
    font-style: italic;
  }

  .expense-amount {
    font-weight: 600;
    color: var(--danger);
    font-size: 1.125rem;
  }

  @media (max-width: 768px) {
    .date-range {
      grid-template-columns: 1fr;
    }

    .expenses-list {
      max-height: 300px;
    }
  }
</style>

