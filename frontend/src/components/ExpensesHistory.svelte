<script>
  // ============================================================================
  // IMPORTS
  // ============================================================================
  import api from '$lib/api';
  import { onMount } from 'svelte';
  import CalendarView from './CalendarView.svelte';
  import MonthDetailModal from './MonthDetailModal.svelte';
  import DateExpensesModal from './DateExpensesModal.svelte';
  import DatePicker from '$lib/components/DatePicker.svelte';
  import { formatCurrency } from '$lib/utils/currency';

  // ============================================================================
  // STATE VARIABLES - Calendar & Months
  // ============================================================================
  let months = [];                      // Array of month data from API
  let selectedMonth = null;              // Currently selected month object
  let selectedDate = null;               // Currently selected date string
  let showMonthDetail = false;           // Month detail modal visibility
  let showDateModal = false;             // Date expenses modal visibility
  let showAnalytics = false;             // Analytics modal visibility
  let analyticsMonth = null;             // Month string (YYYY-MM) for analytics
  let loading = false;                   // Initial loading state
  let calendarWrapper = null;            // Reference to calendar container element
  let isLoadingMore = false;             // Loading more months state
  let hasMoreMonths = true;             // Whether more months can be loaded
  let oldestLoadedMonth = null;         // Oldest month loaded (for pagination)
  let savedScrollPosition = 0;           // Saved scroll position for modals

  // ============================================================================
  // STATE VARIABLES - Balance
  // ============================================================================
  let balance = { amount: 0 };           // User's current balance
  let balanceLoading = true;             // Balance loading state

  // ============================================================================
  // STATE VARIABLES - Search
  // ============================================================================
  let searchQuery = '';                  // Text search query
  let searchCategory = '';               // Category filter ID
  let dateFrom = '';                     // Start date for search
  let dateTo = '';                       // End date for search
  let showSearch = false;                // Search modal visibility
  let filteredExpenses = [];             // Filtered search results
  let categories = [];                   // Available categories for filter
  let searchLoading = false;             // Search loading state

  // ============================================================================
  // DATA LOADING FUNCTIONS
  // ============================================================================
  
  /**
   * Load categories from API
   * Only loads active categories for the search filter
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
    }
  }

  /**
   * Load user's current balance from API
   */
  async function loadBalance() {
    try {
      const response = await api.get('/income/balance');
      balance = response.data;
    } catch (error) {
      // Balance failed to load - silently fail
    } finally {
      balanceLoading = false;
    }
  }

  /**
   * Load months data from API
   * Supports pagination via 'before' parameter
   * @param {boolean} silent - If true, don't show loading indicator
   * @param {string|null} beforeMonth - Month string (YYYY-MM) to load months before
   */
  async function loadMonths(silent = false, beforeMonth = null) {
    // Set loading states
    if (!silent && !beforeMonth) {
      loading = true;
    }
    if (beforeMonth) {
      isLoadingMore = true;
    }
    
    try {
      // Build API URL with optional pagination parameter
      let url = '/expenses/months';
      if (beforeMonth) {
        url += `?before=${beforeMonth}`;
      }
      
      const response = await api.get(url);
      const newMonths = response.data || [];
      
      // Append or replace months based on pagination
      if (beforeMonth) {
        // Append older months for infinite scroll
        months = [...months, ...newMonths];
      } else {
        // Initial load - replace all months
        months = newMonths;
      }
      
      // Update oldest loaded month for next pagination
      if (newMonths.length > 0) {
        oldestLoadedMonth = newMonths[newMonths.length - 1].month;
      }
      
      // Determine if more months can be loaded
      // Initial load: 4 months (current + 3 back)
      // Pagination: 6 months per page
      const expectedCount = beforeMonth ? 6 : 4;
      hasMoreMonths = newMonths.length >= expectedCount;
    } catch (error) {
      hasMoreMonths = false;
    } finally {
      // Reset loading states
      if (!silent && !beforeMonth) {
        loading = false;
      }
      if (beforeMonth) {
        isLoadingMore = false;
      }
    }
  }

  /**
   * Load more months for infinite scroll
   * Called when user scrolls near bottom
   */
  async function loadMoreMonths() {
    if (isLoadingMore || !hasMoreMonths || !oldestLoadedMonth) return;
    await loadMonths(true, oldestLoadedMonth);
  }

  // ============================================================================
  // CALCULATION FUNCTIONS
  // ============================================================================
  
  /**
   * Calculate total expenses across all loaded months
   * @returns {number} Total expense amount
   */
  function getTotalExpenses() {
    // Kept for potential display, but overall balance is computed on the backend now
    if (!months || months.length === 0) {
      return 0;
    }
    return months.reduce((sum, month) => {
      const monthTotal = typeof month.total === 'number'
        ? month.total
        : parseFloat(month.total) || 0;
      return sum + monthTotal;
    }, 0);
  }

  /**
   * Calculate remaining balance (balance - total expenses)
   * @returns {number} Remaining balance amount
   */
  function getRemainingBalance() {
    if (balanceLoading || !balance) {
      return 0;
    }
    const balanceAmount = typeof balance.amount === 'number'
      ? balance.amount
      : parseFloat(balance.amount) || 0;
    return balanceAmount;
  }

  // ============================================================================
  // SEARCH FUNCTIONS
  // ============================================================================
  
  /**
   * Search expenses based on query, category, and date range
   * Called automatically when search parameters change (reactive)
   */
  async function searchExpenses() {
    // Clear results if no search criteria
    if (!searchQuery && !searchCategory && !dateFrom && !dateTo) {
      filteredExpenses = [];
      return;
    }

    searchLoading = true;
    try {
      // Build query parameters
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

  /**
   * Clear all search filters and results
   */
  function clearSearch() {
    searchQuery = '';
    searchCategory = '';
    dateFrom = '';
    dateTo = '';
    filteredExpenses = [];
    showSearch = false;
  }

  /**
   * Toggle search modal visibility
   */
  function toggleSearch() {
    showSearch = !showSearch;
    if (!showSearch) {
      clearSearch();
    }
  }

  // ============================================================================
  // EVENT HANDLERS - Month & Date Clicks
  // ============================================================================
  
  /**
   * Handle month header click
   * Opens month detail modal
   * @param {object} month - Month object from calendar
   */
  function handleMonthClick(month) {
    // Save scroll position before opening modal
    if (calendarWrapper) {
      savedScrollPosition = calendarWrapper.scrollTop;
    }
    selectedMonth = month;
    showMonthDetail = true;
  }

  /**
   * Handle date click in calendar
   * Opens date expenses modal
   * @param {string} date - Date string (YYYY-MM-DD)
   */
  function handleDateClick(date) {
    // Save scroll position before opening modal
    if (calendarWrapper) {
      savedScrollPosition = calendarWrapper.scrollTop;
    }
    selectedDate = date;
    showDateModal = true;
  }

  /**
   * Handle load more months request (infinite scroll)
   */
  function handleLoadMore() {
    loadMoreMonths();
  }

  // ============================================================================
  // EVENT HANDLERS - Modal Close
  // ============================================================================
  
  /**
   * Close month detail modal and restore scroll position
   */
  function closeMonthDetail() {
    showMonthDetail = false;
    selectedMonth = null;
    
    // Restore scroll position after modal closes
    setTimeout(() => {
      if (calendarWrapper) {
        calendarWrapper.scrollTop = savedScrollPosition;
      }
    }, 10);
  }

  /**
   * Close date expenses modal
   * Preserves scroll position
   */
  async function closeDateModal() {
    showDateModal = false;
    selectedDate = null;
    
    // Preserve scroll position
    const currentScroll = calendarWrapper ? calendarWrapper.scrollTop : 0;
    setTimeout(() => {
      if (calendarWrapper) {
        calendarWrapper.scrollTop = currentScroll;
      }
    }, 10);
  }

  /**
   * Close analytics modal
   */
  function closeAnalytics() {
    showAnalytics = false;
    analyticsMonth = null;
  }

  // ============================================================================
  // LIFECYCLE HOOKS
  // ============================================================================
  
  /**
   * Component mount: Load initial data
   */
  onMount(async () => {
    await Promise.all([loadMonths(), loadBalance(), loadCategories()]);
  });

  // ============================================================================
  // REACTIVE STATEMENTS
  // ============================================================================
  
  /**
   * Automatically search when search parameters change
   * This is a Svelte reactive statement (runs when dependencies change)
   */
  $: if (searchQuery || searchCategory || dateFrom || dateTo) {
    searchExpenses();
  }
</script>

<!-- ============================================================================
     TEMPLATE
     ============================================================================ -->
<div class="expenses-history">
  <!-- Header Section -->
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

  <!-- Balance Display (compact) -->
  {#if !balanceLoading && !loading}
    <div class="balance-section compact">
      <div class="balance-card compact">
        <div class="balance-label">Balance</div>
        <div class="balance-amount">{formatCurrency(balance.amount)}</div>
      </div>
    </div>
  {/if}

  <!-- Calendar View -->
  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="calendar-wrapper" bind:this={calendarWrapper}>
      <CalendarView
        {months}
        {hasMoreMonths}
        on:monthClick={(e) => handleMonthClick(e.detail)}
        on:dateClick={(e) => handleDateClick(e.detail)}
        on:loadMore={(e) => handleLoadMore()}
      />
      {#if isLoadingMore}
        <div class="loading-more">
          <div class="spinner"></div>
          <span>Loading more months...</span>
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Month Detail Modal -->
{#if showMonthDetail && selectedMonth}
  <MonthDetailModal
    month={selectedMonth}
    on:close={closeMonthDetail}
  />
{/if}

<!-- Date Expenses Modal -->
{#if showDateModal && selectedDate}
  <DateExpensesModal
    date={selectedDate}
    on:close={closeDateModal}
    on:refresh={async () => {
      // Refresh data but preserve scroll position
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

<!-- Search Modal -->
{#if showSearch}
  <div class="modal-backdrop" on:click={() => showSearch = false}>
    <div class="modal-content search-modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Search Expenses</h2>
        <button class="close-btn" on:click={() => showSearch = false}>×</button>
      </div>
      <div class="modal-body">
        <div class="search-form">
          <!-- Search Query Input -->
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
          
          <!-- Category Filter -->
          <div class="form-group">
            <label for="search-category">Category</label>
            <select id="search-category" bind:value={searchCategory} class="form-input">
              <option value="">All Categories</option>
              {#each categories as category}
                <option value={category.id}>{category.label}</option>
              {/each}
            </select>
          </div>
          
          <!-- Date Range -->
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
          
          <!-- Action Buttons -->
          <div class="button-group">
            <button type="button" class="btn btn-secondary" on:click={clearSearch}>Clear</button>
            <button type="button" class="btn btn-primary" on:click={searchExpenses}>Search</button>
          </div>
        </div>
        
        <!-- Search Results -->
        {#if searchLoading}
          <div class="search-loading">Searching...</div>
        {:else if filteredExpenses.length > 0}
          <div class="search-results">
            <h3>Search Results ({filteredExpenses.length})</h3>
            <div class="expenses-list">
              {#each filteredExpenses as expense}
                <div class="expense-item">
                  <div class="expense-info">
                    <div class="expense-date">
                      {new Date(expense.date).toLocaleDateString('en-US', { 
                        month: 'short', 
                        day: 'numeric', 
                        year: 'numeric' 
                      })}
                    </div>
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

<!-- ============================================================================
     COMPONENT-SPECIFIC STYLES
     ============================================================================ -->
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

  .balance-section.compact {
    margin-bottom: 0.75rem;
    flex-shrink: 0;
  }

  .balance-card.compact {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    padding: 0.6rem 0.9rem;
    box-shadow: none;
  }

  .balance-label {
    font-size: 0.85rem;
    color: var(--text-secondary);
    font-weight: 600;
  }

  .balance-amount {
    font-size: 1rem;
    font-weight: 700;
    color: var(--text-primary);
    line-height: 1.2;
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

  .loading-more {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    gap: 0.75rem;
    color: var(--text-secondary);
    font-size: 0.875rem;
  }

  .loading-more .spinner {
    width: 24px;
    height: 24px;
    border: 3px solid var(--border);
    border-top-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
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
