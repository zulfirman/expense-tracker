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
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { page } from '$app/stores';
  import { getPageCode } from '$lib/utils/pageCodes';

  $: pageCode = getPageCode($page?.url?.pathname || '/history');

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
<div class="expenses-history space-y-4">
  <PageHeader
    title="Expenses History"
    subtitle="View and manage your expense history."
    pageCode={pageCode}
    actions={true}
  >
    <svelte:fragment slot="actions">
      <button
        class="btn btn-primary gap-2"
        class:btn-active={showSearch}
        on:click={toggleSearch}
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
        Search
      </button>
    </svelte:fragment>
  </PageHeader>

  {#if !balanceLoading && !loading}
    <div class="stats shadow-sm bg-base-100 border border-base-300">
      <div class="stat">
        <div class="stat-title">Balance</div>
        <div class="stat-value text-primary text-2xl">{formatCurrency(balance.amount)}</div>
      </div>
    </div>
  {/if}

  <div class="card bg-base-100 shadow-xl border border-base-300">
    <div class="card-body p-0">
      {#if loading}
        <div class="flex justify-center py-10">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else}
        <div class="calendar-wrapper px-4 pb-4" bind:this={calendarWrapper}>
          <CalendarView
            {months}
            {hasMoreMonths}
            on:monthClick={(e) => handleMonthClick(e.detail)}
            on:dateClick={(e) => handleDateClick(e.detail)}
            on:loadMore={(e) => handleLoadMore()}
          />
          {#if isLoadingMore}
            <div class="flex items-center gap-2 justify-center py-4 text-base-content/70">
              <span class="loading loading-spinner loading-sm"></span>
              <span>Loading more months...</span>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
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

{#if showSearch}
  <div class="modal modal-open z-[2100]" on:click={() => showSearch = false}>
    <div class="modal-box w-11/12 max-w-4xl" on:click|stopPropagation>
      <div class="flex items-start justify-between gap-3 mb-4">
        <div>
          <h3 class="font-bold text-xl">Search Expenses</h3>
          <p class="text-sm text-base-content/70">Filter by notes, category, or date range.</p>
        </div>
        <button class="btn btn-ghost btn-sm" on:click={() => showSearch = false}>✕</button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <fieldset class="fieldset w-full">
          <legend class="fieldset-legend">Search (notes)</legend>
          <input
            id="search-query"
            type="text"
            bind:value={searchQuery}
            placeholder="Search expenses..."
            class="input input-bordered w-full border-2"
            on:keydown={(e) => { if (e.key === 'Enter') searchExpenses(); }}
          />
        </fieldset>

        <fieldset class="fieldset w-full">
          <legend class="fieldset-legend">Category</legend>
          <select
            id="search-category"
            bind:value={searchCategory}
            class="select select-bordered w-full border-2"
          >
            <option value="">All Categories</option>
            {#each categories as category}
              <option value={category.id}>{category.label}</option>
            {/each}
          </select>
        </fieldset>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:col-span-2">
          <DatePicker
            id="date-from"
            bind:value={dateFrom}
            placeholder="From date"
            label="From date"
          />
          <DatePicker
            id="date-to"
            bind:value={dateTo}
            placeholder="To date"
            label="To date"
          />
        </div>
      </div>

      <div class="flex justify-end gap-2 mt-4">
        <button type="button" class="btn btn-soft" on:click={clearSearch}>Clear</button>
        <button type="button" class="btn btn-primary" on:click={searchExpenses}>Search</button>
      </div>

      {#if searchLoading}
        <div class="flex justify-center py-6 text-base-content/70">
          <span class="loading loading-spinner loading-md"></span>
        </div>
      {:else if filteredExpenses.length > 0}
        <div class="divider mt-6">Results ({filteredExpenses.length})</div>
        <div class="space-y-3 max-h-96 overflow-y-auto pr-1">
          {#each filteredExpenses as expense}
            <div class="card bg-base-100 border border-base-300 shadow-sm">
              <div class="card-body py-3 px-4 flex flex-col gap-1 md:flex-row md:items-center md:justify-between">
                <div class="space-y-1">
                  <p class="text-sm text-base-content/70">
                    {new Date(expense.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })} • {expense.category}
                  </p>
                  {#if expense.notes}
                    <p class="text-sm text-base-content/90">{expense.notes}</p>
                  {/if}
                </div>
                <div class="text-right font-semibold text-primary text-lg">{formatCurrency(expense.amount)}</div>
              </div>
            </div>
          {/each}
        </div>
      {:else if searchQuery || searchCategory || dateFrom || dateTo}
        <div class="alert alert-info mt-4">
          <span>No expenses found matching your search criteria.</span>
        </div>
      {/if}
    </div>
  </div>
{/if}

<!-- ============================================================================
     COMPONENT-SPECIFIC STYLES
     ============================================================================ -->
<style>
  .expenses-history {
    max-width: 100%;
  }

  .calendar-wrapper {
    height: 65vh;
    min-height: 420px;
  }

  @media (max-width: 768px) {
    .calendar-wrapper {
      height: 70vh;
    }
  }
</style>
