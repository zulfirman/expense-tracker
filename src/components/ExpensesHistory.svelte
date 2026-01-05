<script>
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { onMount } from 'svelte';
  import CalendarView from './CalendarView.svelte';
  import MonthDetailModal from './MonthDetailModal.svelte';
  import DateExpensesModal from './DateExpensesModal.svelte';

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
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(numericValue);
  }

  onMount(async () => {
    await Promise.all([loadMonths(), loadBalance()]);
  });

  async function loadBalance() {
    try {
      const response = await axios.get('/api/income/balance');
      balance = response.data;
    } catch (error) {
      console.error('Error loading balance:', error);
    } finally {
      balanceLoading = false;
    }
  }

  async function loadMonths() {
    loading = true;
    try {
      const response = await axios.get('/api/expenses/months');
      months = response.data;
    } catch (error) {
      console.error('Error loading months:', error);
    } finally {
      loading = false;
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
  <h1>Expenses History</h1>

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
      // Refresh data but preserve scroll position
      const currentScroll = calendarWrapper ? calendarWrapper.scrollTop : 0;
      await Promise.all([loadMonths(), loadBalance()]);
      setTimeout(() => {
        if (calendarWrapper) {
          calendarWrapper.scrollTop = currentScroll;
        }
      }, 50);
    }}
  />
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
    font-size: 2.5rem;
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
</style>

