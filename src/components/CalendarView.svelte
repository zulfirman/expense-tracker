<script>
  import { createEventDispatcher, afterUpdate } from 'svelte';

  const dispatch = createEventDispatcher();

  export let months = [];

  let scrollContainer;
  let isLoading = false;
  let savedScrollPosition = 0;
  let isInitialLoad = true;
  let previousMonthsLength = 0;

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount);
  }

  function formatMonthYear(monthYear) {
    const [year, month] = monthYear.split('-');
    const date = new Date(year, month - 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  }

  function getAllDaysInMonth(year, month) {
    const date = new Date(year, month, 1);
    const days = [];
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    
    for (let day = 1; day <= daysInMonth; day++) {
      days.push(new Date(year, month, day));
    }
    
    return days;
  }

  function getDateTotal(dateStr, dateTotals) {
    const found = dateTotals.find(d => d.date === dateStr);
    return found ? found.total : 0;
  }

  function formatDateKey(date) {
    // Format date as YYYY-MM-DD without timezone conversion
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }

  function handleMonthClick(month) {
    dispatch('monthClick', month);
  }

  function handleDateClick(dateStr, hasData) {
    dispatch('dateClick', { date: dateStr, total: 0 });
  }

  function handleScroll(e) {
    const container = e.target;
    savedScrollPosition = container.scrollTop;
  }

  // Restore scroll position when months change (but not on initial load)
  afterUpdate(() => {
    if (scrollContainer && savedScrollPosition > 0 && !isInitialLoad) {
      // Only restore if months array actually changed
      if (months.length !== previousMonthsLength) {
        previousMonthsLength = months.length;
        // Use setTimeout to ensure DOM is fully updated
        setTimeout(() => {
          if (scrollContainer) {
            scrollContainer.scrollTop = savedScrollPosition;
          }
        }, 10);
      }
    }
    if (isInitialLoad && months.length > 0) {
      isInitialLoad = false;
      previousMonthsLength = months.length;
    }
  });
</script>

<div class="calendar-container" bind:this={scrollContainer} on:scroll={handleScroll}>
  {#each months as month}
    {@const [year, monthNum] = month.month.split('-').map(Number)}
    {@const allDays = getAllDaysInMonth(year, monthNum - 1)}
    
    <div class="month-card">
      <button class="month-header" on:click={() => handleMonthClick(month)}>
        <h2>{formatMonthYear(month.month)}</h2>
        <span class="total-amount">{formatCurrency(month.total)}</span>
      </button>
      
      <div class="dates-list">
        {#each allDays.slice().reverse() as day}
          {@const dateKey = formatDateKey(day)}
          {@const dateData = month.dates?.find(d => d.date === dateKey)}
          {@const total = dateData?.total || 0}
          {@const hasIncome = dateData?.hasIncome || false}
          {@const hasExpense = dateData?.hasExpense || false}
          {@const hasData = hasIncome || hasExpense}
          {@const isToday = dateKey === new Date().toISOString().split('T')[0]}
          {@const dateType = hasIncome && hasExpense ? 'both' : (hasIncome ? 'income' : (hasExpense ? 'expense' : 'none'))}
          
          <button 
            class="date-item"
            class:has-data={hasData}
            class:is-today={isToday}
            class:no-data={!hasData}
            class:income-only={dateType === 'income'}
            class:expense-only={dateType === 'expense'}
            class:both={dateType === 'both'}
            on:click={() => handleDateClick(dateKey, true)}
          >
            <span class="date-number">{day.getDate()}</span>
            {#if hasData}
              <span class="date-amount">{formatCurrency(Math.abs(total))}</span>
            {:else}
              <span class="date-amount empty">-</span>
            {/if}
          </button>
        {/each}
      </div>
    </div>
  {/each}
</div>

<style>
  .calendar-container {
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
    position: relative;
    scroll-behavior: smooth;
  }

  .month-card {
    background: var(--surface);
    border-radius: 0.75rem;
    margin-bottom: 1rem;
    padding: 1rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .month-header {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: none;
    border: none;
    cursor: pointer;
    border-radius: 0.5rem;
    transition: background-color 0.2s;
    margin-bottom: 0.75rem;
  }

  .month-header:hover {
    background-color: var(--background);
  }

  .month-header h2 {
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 0;
    text-transform: capitalize;
  }

  .total-amount {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--primary-color);
  }

  .dates-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 0.5rem;
  }

  .date-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.75rem 0.5rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.875rem;
    min-height: 70px;
    justify-content: center;
  }

  .date-item:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    transform: translateY(-2px);
  }

  .date-item:active {
    transform: scale(0.95);
  }

  .date-item.has-data {
    background: rgba(79, 70, 229, 0.05);
    border-color: var(--primary-color);
  }

  .date-item.is-today {
    border-color: var(--success);
    border-width: 2px;
    box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
  }

  .date-item.no-data {
    opacity: 0.6;
  }

  .date-item.no-data:hover {
    opacity: 1;
  }

  .date-number {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.25rem;
  }

  .date-amount {
    font-size: 0.75rem;
    color: var(--text-secondary);
    font-weight: 600;
  }

  .date-amount.empty {
    color: var(--text-secondary);
    opacity: 0.5;
  }

  .date-item.income-only {
    background: rgba(16, 185, 129, 0.1);
    border-color: #10b981;
  }

  .date-item.expense-only {
    background: rgba(239, 68, 68, 0.1);
    border-color: #ef4444;
  }

  .date-item.both {
    background: rgba(59, 130, 246, 0.1);
    border-color: #3b82f6;
  }

  .date-item.income-only .date-amount {
    color: #10b981;
    font-weight: 700;
  }

  .date-item.expense-only .date-amount {
    color: #ef4444;
    font-weight: 700;
  }

  .date-item.both .date-amount {
    color: #3b82f6;
    font-weight: 700;
  }
</style>

