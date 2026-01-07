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

  function getToneClasses(dateType) {
    if (dateType === 'income') {
      return 'bg-success/20 border-success/60 text-success-content';
    }
    if (dateType === 'expense') {
      return 'bg-error/20 border-error/60 text-error-content';
    }
    if (dateType === 'both') {
      return 'bg-info/20 border-info/60 text-info-content';
    }
    return 'bg-base-200/60 border-base-300';
  }
</script>

<div class="calendar-container" bind:this={scrollContainer} on:scroll={handleScroll}>
  {#each months as month}
    {@const [year, monthNum] = month.month.split('-').map(Number)}
    {@const allDays = getAllDaysInMonth(year, monthNum - 1)}

    <div class="card bg-base-100 shadow-md border border-base-300 mb-4">
      <div class="card-body p-4 space-y-3">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
          <div>
            <p class="text-xs uppercase tracking-wide text-base-content/60">Month</p>
            <h2 class="text-xl font-bold">{formatMonthYear(month.month)}</h2>
            <p class="text-sm text-base-content/60">Tap to open details</p>
          </div>
          <button
            class="btn btn-primary gap-2 shadow"
            on:click={() => handleMonthClick(month)}
          >
            {formatCurrency(month.netTotal ?? month.total)}
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </button>
        </div>

        <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 gap-2">
          {#each allDays.slice().reverse() as day}
            {@const dateKey = formatDateKey(day)}
            {@const dateData = month.dates?.find(d => d.date === dateKey)}
            {@const total = dateData?.total || 0}
            {@const hasIncome = dateData?.hasIncome || false}
            {@const hasExpense = dateData?.hasExpense || false}
            {@const hasData = hasIncome || hasExpense}
            {@const isToday = dateKey === new Date().toISOString().split('T')[0]}
            {@const dateType = hasIncome && hasExpense ? 'both' : (hasIncome ? 'income' : (hasExpense ? 'expense' : 'none'))}

            <div
              role="button"
              class={`card card-compact w-full cursor-pointer border shadow-sm transition text-sm ${getToneClasses(dateType)}`}
              class:outline={isToday}
              class:opacity-60={!hasData}
              on:click={() => handleDateClick(dateKey, true)}
            >
              <div class="card-body items-center text-center p-2">
                <div class="text-base font-semibold">{day.getDate()}</div>
                <div class="text-[11px] text-base-content/70">
                  {#if hasData}
                    {formatCurrency(Math.abs(total))}
                  {:else}
                    -
                  {/if}
                </div>
                {#if isToday}
                  <span class="badge px-3 py-2 badge-outline badge-xs mt-1">Today</span>
                {/if}
                <div class="flex gap-1 mt-1">
                  {#if dateType === 'income'}
                    <span class="badge px-3 py-2 badge-success badge-sm">Income</span>
                  {:else if dateType === 'expense'}
                    <span class="badge px-3 py-2 badge-error badge-sm">Expense</span>
                  {:else if dateType === 'both'}
                    <span class="badge px-3 py-2 badge-info badge-sm">Both</span>
                  {/if}
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/each}
</div>

<style>
  .calendar-container {
    height: 100%;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
    scroll-behavior: smooth;
  }
</style>