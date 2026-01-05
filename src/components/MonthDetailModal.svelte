<script>
  import { createEventDispatcher, onMount, onDestroy, afterUpdate } from 'svelte';
  import api from '$lib/api';
  import { Chart, registerables } from 'chart.js';

  Chart.register(...registerables);

  const dispatch = createEventDispatcher();

  export let month = null;

  let chartCanvas;
  let chartInstance = null;
  let categoryChartCanvas;
  let categoryChartInstance = null;
  let loading = true;
  let expensesByCategory = [];
  let dailyExpenses = [];

  onMount(async () => {
    await loadMonthDetails();
  });

  onDestroy(() => {
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }
    if (categoryChartInstance) {
      categoryChartInstance.destroy();
      categoryChartInstance = null;
    }
  });

  afterUpdate(() => {
    // Recreate charts after data updates
    if (chartCanvas && dailyExpenses.length > 0 && !chartInstance) {
      createDailyChart();
    }
    if (categoryChartCanvas && expensesByCategory.length > 0 && !categoryChartInstance) {
      createCategoryChart();
    }
  });

  async function loadMonthDetails() {
    if (!month) return;
    
    loading = true;
    try {
      const response = await api.get(`/expenses/month/${month.month}`);
      expensesByCategory = response.data.categories || [];
      dailyExpenses = response.data.daily || [];
      
      // Destroy existing charts
      if (chartInstance) {
        chartInstance.destroy();
        chartInstance = null;
      }
      if (categoryChartInstance) {
        categoryChartInstance.destroy();
        categoryChartInstance = null;
      }

      // Create charts after a small delay to ensure canvas is ready
      setTimeout(() => {
        if (dailyExpenses.length > 0 && chartCanvas) {
          createDailyChart();
        }
        if (expensesByCategory.length > 0 && categoryChartCanvas) {
          createCategoryChart();
        }
      }, 100);
    } catch (error) {
      console.error('Error loading month details:', error);
    } finally {
      loading = false;
    }
  }

  function createDailyChart() {
    if (!chartCanvas || dailyExpenses.length === 0) return;

    const ctx = chartCanvas.getContext('2d');
    
    // Sort daily expenses by date
    const sortedDaily = [...dailyExpenses].sort((a, b) => {
      return new Date(a.date) - new Date(b.date);
    });

    const labels = sortedDaily.map(item => {
      const date = new Date(item.date);
      return date.getDate().toString();
    });
    
    const incomeData = sortedDaily.map(item => item.income || 0);
    const expenseData = sortedDaily.map(item => item.expense || 0);
    const netData = sortedDaily.map(item => item.netTotal || 0);

    chartInstance = new Chart(ctx, {
      type: 'line',
      data: {
        labels: labels,
        datasets: [
          {
            label: 'Income',
            data: incomeData,
            borderColor: 'rgb(16, 185, 129)',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            borderWidth: 2,
            fill: false,
            tension: 0.4,
            pointRadius: 4,
            pointHoverRadius: 6,
            pointBackgroundColor: 'rgb(16, 185, 129)',
            pointBorderColor: '#fff',
            pointBorderWidth: 2
          },
          {
            label: 'Expenses',
            data: expenseData,
            borderColor: 'rgb(239, 68, 68)',
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderWidth: 2,
            fill: false,
            tension: 0.4,
            pointRadius: 4,
            pointHoverRadius: 6,
            pointBackgroundColor: 'rgb(239, 68, 68)',
            pointBorderColor: '#fff',
            pointBorderWidth: 2
          },
          {
            label: 'Net (Income - Expense)',
            data: netData,
            borderColor: 'rgb(59, 130, 246)',
            backgroundColor: 'rgba(59, 130, 246, 0.2)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointRadius: 5,
            pointHoverRadius: 7,
            pointBackgroundColor: 'rgb(59, 130, 246)',
            pointBorderColor: '#fff',
            pointBorderWidth: 2,
            borderDash: [5, 5]
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: 'index',
          intersect: false
        },
        plugins: {
          legend: {
            display: true,
            position: 'top',
            labels: {
              usePointStyle: true,
              padding: 15,
              font: {
                size: 12
              }
            }
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                let label = context.dataset.label || '';
                if (label) {
                  label += ': ';
                }
                label += formatCurrency(context.parsed.y);
                return label;
              },
              title: function(context) {
                const index = context[0].dataIndex;
                const date = new Date(sortedDaily[index].date);
                return `Day ${date.getDate()} - ${date.toLocaleDateString('en-US', { month: 'long' })}`;
              }
            }
          }
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Day of Month',
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            grid: {
              display: true,
              color: 'rgba(0, 0, 0, 0.05)'
            }
          },
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Amount (IDR)',
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                if (value >= 1000000) {
                  return (value / 1000000).toFixed(1) + 'M';
                } else if (value >= 1000) {
                  return (value / 1000).toFixed(0) + 'K';
                }
                return value.toString();
              }
            },
            grid: {
              display: true,
              color: 'rgba(0, 0, 0, 0.05)'
            }
          }
        }
      }
    });
  }

  function createCategoryChart() {
    if (!categoryChartCanvas || expensesByCategory.length === 0) return;

    const ctx = categoryChartCanvas.getContext('2d');
    
    // Sort categories by total (descending)
    const sortedCategories = [...expensesByCategory].sort((a, b) => b.total - a.total);

    const colors = [
      'rgba(79, 70, 229, 0.8)',
      'rgba(124, 58, 237, 0.8)',
      'rgba(168, 85, 247, 0.8)',
      'rgba(192, 132, 252, 0.8)',
      'rgba(217, 119, 6, 0.8)',
      'rgba(234, 179, 8, 0.8)'
    ];

    categoryChartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: sortedCategories.map(item => item.category),
        datasets: [{
          label: 'Expenses by Category',
          data: sortedCategories.map(item => item.total),
          backgroundColor: sortedCategories.map((_, index) => colors[index % colors.length]),
          borderColor: sortedCategories.map((_, index) => colors[index % colors.length].replace('0.8', '1')),
          borderWidth: 2,
          borderRadius: 4
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                return formatCurrency(context.parsed.y);
              }
            }
          }
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Category',
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            grid: {
              display: false
            }
          },
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Amount (IDR)',
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                if (value >= 1000000) {
                  return (value / 1000000).toFixed(1) + 'M';
                } else if (value >= 1000) {
                  return (value / 1000).toFixed(0) + 'K';
                }
                return value.toString();
              }
            },
            grid: {
              display: true,
              color: 'rgba(0, 0, 0, 0.05)'
            }
          }
        }
      }
    });
  }

  function formatCurrency(amount) {
    return 'Rp. ' + new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }

  function formatMonthYear(monthYear) {
    const [year, monthNum] = monthYear.split('-');
    const date = new Date(year, monthNum - 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  }

  function close() {
    dispatch('close');
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      close();
    }
  }
</script>

<div class="modal-backdrop" on:click={handleBackdropClick}>
  <div class="modal-content">
    <div class="modal-header">
      <div>
        <h2>Advanced Analytics & Insights</h2>
        <p class="modal-subtitle">{formatMonthYear(month.month)}</p>
      </div>
      <button class="close-btn" on:click={close}>×</button>
    </div>

    <div class="modal-body">
      {#if loading}
        <div class="loading">Loading...</div>
      {:else if dailyExpenses.length > 0 || expensesByCategory.length > 0}
        {#if dailyExpenses.length > 0}
          <div class="chart-section">
            <h3>Daily Income vs Expenses Trend</h3>
            <div class="chart-container">
              <canvas bind:this={chartCanvas}></canvas>
            </div>
          </div>
        {/if}

        {#if expensesByCategory.length > 0}
          <div class="chart-section">
            <h3>Expenses by Category</h3>
            <div class="chart-container">
              <canvas bind:this={categoryChartCanvas}></canvas>
            </div>
          </div>

          <div class="category-list">
            <h4>Category Breakdown</h4>
            {#each expensesByCategory.sort((a, b) => b.total - a.total) as item}
              <div class="category-item">
                <span class="category-name">{item.category}</span>
                <span class="category-amount">{formatCurrency(item.total)}</span>
              </div>
            {/each}
          </div>
        {:else if dailyExpenses.length === 0}
          <div class="no-data">No expenses data available</div>
        {/if}
      {:else}
        <div class="no-data">No data available</div>
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
    border-radius: 1rem;
    width: 100%;
    max-width: 600px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--border);
  }

  .modal-header h2 {
    font-size: 1.5rem;
    color: var(--text-primary);
    margin: 0;
    font-weight: 700;
  }

  .modal-subtitle {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin: 0.25rem 0 0 0;
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
  }

  .loading,
  .no-data {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .chart-section {
    margin-bottom: 2.5rem;
  }

  .chart-section h3 {
    font-size: 1.25rem;
    color: var(--text-primary);
    margin-bottom: 1rem;
    font-weight: 600;
  }

  .chart-container {
    position: relative;
    height: 350px;
    width: 100%;
    background: var(--background);
    border-radius: 0.75rem;
    padding: 1rem;
    border: 1px solid var(--border);
  }

  .category-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 1.5rem;
  }

  .category-list h4 {
    font-size: 1rem;
    color: var(--text-primary);
    margin-bottom: 0.75rem;
    font-weight: 600;
  }

  .category-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: var(--background);
    border-radius: 0.5rem;
  }

  .category-name {
    font-weight: 500;
    color: var(--text-primary);
    text-transform: capitalize;
  }

  .category-amount {
    font-weight: 600;
    color: var(--primary-color);
  }
</style>

