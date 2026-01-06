<script>
  import { createEventDispatcher, afterUpdate } from 'svelte';
  import { formatCurrency } from '$lib/utils/currency';

  const dispatch = createEventDispatcher();

  export let months = [];
  export let hasMoreMonths = true;

  let scrollContainer;
  let isLoading = false;
  let savedScrollPosition = 0;
  let isInitialLoad = true;
  let previousMonthsLength = 0;

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
    
    // Check if user scrolled near bottom (within 300px)
    const scrollBottom = container.scrollHeight - container.scrollTop - container.clientHeight;
    if (scrollBottom < 300 && hasMoreMonths) {
      dispatch('loadMore', { direction: 'down' });
    }
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
        <div class="month-header-content">
          <div class="month-title-section">
            <h2>{formatMonthYear(month.month)}</h2>
            <span class="click-hint">Click for details</span>
          </div>
          <div class="total-amount-wrapper">
            <span class="total-amount">
              {formatCurrency(month.netTotal ?? month.total)}
            </span>
            <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </div>
        </div>
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
    border-radius: 1rem;
    margin-bottom: 1.5rem;
    padding: 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    border: 1px solid var(--border);
    transition: box-shadow 0.2s;
  }

  .month-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .month-header {
    width: 100%;
    padding: 1.25rem 1.5rem;
    background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
    border: none;
    border-radius: 0.875rem;
    margin-bottom: 1rem;
    box-shadow: 0 2px 8px rgba(79, 70, 229, 0.25);
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
  }

  .month-header::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: left 0.5s;
  }

  .month-header:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(79, 70, 229, 0.35);
  }

  .month-header:hover::before {
    left: 100%;
  }

  .month-header:active {
    transform: translateY(0);
    box-shadow: 0 2px 8px rgba(79, 70, 229, 0.25);
  }

  .month-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    position: relative;
    z-index: 1;
  }

  .month-title-section {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .month-header h2 {
    font-size: 1.5rem;
    color: white;
    margin: 0;
    text-transform: capitalize;
    font-weight: 700;
    letter-spacing: 0.5px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .click-hint {
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.85);
    font-weight: 500;
    letter-spacing: 0.3px;
    opacity: 0.9;
  }

  .total-amount-wrapper {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: rgba(255, 255, 255, 0.15);
    padding: 0.5rem 1rem;
    border-radius: 0.625rem;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    transition: all 0.3s;
  }

  .month-header:hover .total-amount-wrapper {
    background: rgba(255, 255, 255, 0.25);
    transform: scale(1.05);
  }

  .total-amount {
    font-size: 1.375rem;
    font-weight: 700;
    color: white;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .arrow-icon {
    color: white;
    opacity: 0.9;
    transition: transform 0.3s;
  }

  .month-header:hover .arrow-icon {
    transform: translateX(4px);
  }


  .dates-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(90px, 1fr));
    gap: 0.75rem;
  }

  .date-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.875rem 0.5rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: 0.75rem;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.875rem;
    min-height: 80px;
    justify-content: center;
    position: relative;
    overflow: hidden;
  }

  .date-item::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: transparent;
    transition: background 0.2s;
  }

  .date-item:hover {
    background: var(--surface);
    border-color: var(--primary-color);
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
  }

  .date-item:active {
    transform: translateY(-1px) scale(0.98);
  }

  .date-item.has-data {
    background: rgba(79, 70, 229, 0.08);
    border-color: var(--primary-color);
    border-width: 2px;
  }

  .date-item.has-data::before {
    background: var(--primary-color);
  }

  .date-item.is-today {
    border-color: var(--success);
    border-width: 3px;
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.15);
    background: rgba(16, 185, 129, 0.05);
  }

  .date-item.is-today::before {
    background: var(--success);
    height: 4px;
  }

  .date-item.no-data {
    opacity: 0.5;
    border-style: dashed;
  }

  .date-item.no-data:hover {
    opacity: 0.8;
    border-style: solid;
  }

  .date-number {
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 0.375rem;
    font-size: 1rem;
  }

  .date-item.is-today .date-number {
    color: var(--success);
    font-size: 1.125rem;
  }

  .date-amount {
    font-size: 0.7rem;
    color: var(--text-secondary);
    font-weight: 600;
    text-align: center;
    line-height: 1.2;
    word-break: break-word;
  }

  .date-amount.empty {
    color: var(--text-secondary);
    opacity: 0.4;
    font-size: 0.65rem;
  }

  .date-item.income-only {
    background: rgba(16, 185, 129, 0.12);
    border-color: #10b981;
    border-width: 2px;
  }

  .date-item.income-only::before {
    background: #10b981;
  }

  .date-item.expense-only {
    background: rgba(239, 68, 68, 0.12);
    border-color: #ef4444;
    border-width: 2px;
  }

  .date-item.expense-only::before {
    background: #ef4444;
  }

  .date-item.both {
    background: rgba(59, 130, 246, 0.12);
    border-color: #3b82f6;
    border-width: 2px;
  }

  .date-item.both::before {
    background: linear-gradient(90deg, #10b981 0%, #3b82f6 50%, #ef4444 100%);
  }

  .date-item.income-only .date-amount {
    color: #10b981;
    font-weight: 700;
    font-size: 0.75rem;
  }

  .date-item.expense-only .date-amount {
    color: #ef4444;
    font-weight: 700;
    font-size: 0.75rem;
  }

  .date-item.both .date-amount {
    color: #3b82f6;
    font-weight: 700;
    font-size: 0.75rem;
  }

  @media (max-width: 768px) {
    .dates-list {
      grid-template-columns: repeat(auto-fill, minmax(75px, 1fr));
      gap: 0.5rem;
    }

    .date-item {
      min-height: 70px;
      padding: 0.75rem 0.375rem;
    }

    .month-header {
      padding: 0.875rem 1rem;
    }

    .month-header h2 {
      font-size: 1.125rem;
    }

    .total-amount {
      font-size: 1rem;
      padding: 0.25rem 0.75rem;
    }


    .month-card {
      padding: 1.25rem;
      margin-bottom: 1.25rem;
    }
  }
</style>

