<script>
  import { createEventDispatcher, onMount, onDestroy, afterUpdate } from 'svelte';
  import api from '$lib/api';
  import { Chart, registerables } from 'chart.js';
  import '$lib/styles/charts.css';
  import { formatCurrency } from '$lib/utils/currency';
  import { currency } from '$lib/stores/currency';

  Chart.register(...registerables);

  const dispatch = createEventDispatcher();

  export let month = null;

  let chartCanvas;
  let chartInstance = null;
  let categoryChartCanvas;
  let categoryChartInstance = null;
  let pieChartCanvas;
  let pieChartInstance = null;
  let cumulativeChartCanvas;
  let cumulativeChartInstance = null;
  let weeklyChartCanvas;
  let weeklyChartInstance = null;
  let budgetChartCanvas;
  let budgetChartInstance = null;
  let loading = true;
  let expensesByCategory = [];
  let dailyExpenses = [];
  let weeklyData = [];
  let cumulativeData = [];
  let budgets = [];

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
    if (pieChartInstance) {
      pieChartInstance.destroy();
      pieChartInstance = null;
    }
    if (cumulativeChartInstance) {
      cumulativeChartInstance.destroy();
      cumulativeChartInstance = null;
    }
    if (weeklyChartInstance) {
      weeklyChartInstance.destroy();
      weeklyChartInstance = null;
    }
    if (budgetChartInstance) {
      budgetChartInstance.destroy();
      budgetChartInstance = null;
    }
  });

  afterUpdate(() => {
    // Charts are created in loadMonthDetails after data is loaded
  });

  async function loadMonthDetails() {
    if (!month) return;
    
    loading = true;
    try {
      // Load expenses and budgets in parallel
      const [expensesResponse, budgetsResponse] = await Promise.all([
        api.get(`/expenses/month/${month.month}`),
        api.get(`/budgets?month=${month.month}`).catch(() => ({ data: [] }))
      ]);
      
      expensesByCategory = expensesResponse.data.categories || [];
      dailyExpenses = expensesResponse.data.daily || [];
      budgets = budgetsResponse.data || [];
      
      // Process data for additional charts
      processChartData();
      
      // Destroy existing charts
      if (chartInstance) {
        chartInstance.destroy();
        chartInstance = null;
      }
      if (categoryChartInstance) {
        categoryChartInstance.destroy();
        categoryChartInstance = null;
      }
      if (pieChartInstance) {
        pieChartInstance.destroy();
        pieChartInstance = null;
      }
      if (cumulativeChartInstance) {
        cumulativeChartInstance.destroy();
        cumulativeChartInstance = null;
      }
      if (weeklyChartInstance) {
        weeklyChartInstance.destroy();
        weeklyChartInstance = null;
      }
      if (budgetChartInstance) {
        budgetChartInstance.destroy();
        budgetChartInstance = null;
      }

      // Create charts after a small delay to ensure canvas is ready
      setTimeout(() => {
        if (dailyExpenses.length > 0 && chartCanvas) {
          createDailyChart();
        }
        if (expensesByCategory.length > 0) {
          if (categoryChartCanvas) {
            createCategoryChart();
          }
          if (pieChartCanvas) {
            createPieChart();
          }
        }
        if (cumulativeData.length > 0 && cumulativeChartCanvas) {
          createCumulativeChart();
        }
        if (weeklyData.length > 0 && weeklyChartCanvas) {
          createWeeklyChart();
        }
        if (budgets.length > 0 && budgetChartCanvas) {
          createBudgetChart();
        }
      }, 200);
    } catch (error) {
      // Month details failed to load
    } finally {
      loading = false;
    }
  }

  function processChartData() {
    // Process cumulative spending data
    if (dailyExpenses.length > 0) {
      const sortedDaily = [...dailyExpenses].sort((a, b) => {
        return new Date(a.date) - new Date(b.date);
      });
      
      let cumulative = 0;
      cumulativeData = sortedDaily.map(item => {
        cumulative += item.expense || 0;
        return {
          date: item.date,
          cumulative: cumulative
        };
      });
    }

    // Process weekly spending data
    if (dailyExpenses.length > 0) {
      const sortedDaily = [...dailyExpenses].sort((a, b) => {
        return new Date(a.date) - new Date(b.date);
      });
      
      const weeklyMap = new Map();
      sortedDaily.forEach(item => {
        const date = new Date(item.date);
        const weekStart = new Date(date);
        weekStart.setDate(date.getDate() - date.getDay()); // Start of week (Sunday)
        const weekKey = weekStart.toISOString().split('T')[0];
        
        if (!weeklyMap.has(weekKey)) {
          weeklyMap.set(weekKey, {
            week: `Week ${Math.ceil(date.getDate() / 7)}`,
            startDate: weekKey,
            expense: 0,
            income: 0
          });
        }
        
        const weekData = weeklyMap.get(weekKey);
        weekData.expense += item.expense || 0;
        weekData.income += item.income || 0;
      });
      
      weeklyData = Array.from(weeklyMap.values()).sort((a, b) => {
        return new Date(a.startDate) - new Date(b.startDate);
      });
    }
  }

  function createDailyChart() {
    if (!chartCanvas || dailyExpenses.length === 0) return;

    // Destroy existing chart if it exists
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }

    const ctx = chartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
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
                label += formatCurrency(context.parsed.y, currentCurrency);
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
              text: `Amount (${currentCurrency})`,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
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

    // Destroy existing chart if it exists
    if (categoryChartInstance) {
      categoryChartInstance.destroy();
      categoryChartInstance = null;
    }

    const ctx = categoryChartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
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
                return formatCurrency(context.parsed.y, currentCurrency);
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
              text: `Amount (${currentCurrency})`,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
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

  function createPieChart() {
    if (!pieChartCanvas || expensesByCategory.length === 0) return;

    // Destroy existing chart if it exists
    if (pieChartInstance) {
      pieChartInstance.destroy();
      pieChartInstance = null;
    }

    const ctx = pieChartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
    // Sort categories by total (descending)
    const sortedCategories = [...expensesByCategory].sort((a, b) => b.total - a.total);

    const colors = [
      'rgba(79, 70, 229, 0.8)',
      'rgba(124, 58, 237, 0.8)',
      'rgba(168, 85, 247, 0.8)',
      'rgba(192, 132, 252, 0.8)',
      'rgba(217, 119, 6, 0.8)',
      'rgba(234, 179, 8, 0.8)',
      'rgba(16, 185, 129, 0.8)',
      'rgba(59, 130, 246, 0.8)',
      'rgba(239, 68, 68, 0.8)',
      'rgba(245, 158, 11, 0.8)'
    ];

    const total = sortedCategories.reduce((sum, item) => sum + item.total, 0);

    pieChartInstance = new Chart(ctx, {
      type: 'pie',
      data: {
        labels: sortedCategories.map(item => item.category),
        datasets: [{
          data: sortedCategories.map(item => item.total),
          backgroundColor: sortedCategories.map((_, index) => colors[index % colors.length]),
          borderColor: sortedCategories.map((_, index) => colors[index % colors.length].replace('0.8', '1')),
          borderWidth: 2
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: true,
            position: 'right',
            labels: {
              usePointStyle: true,
              padding: 15,
              font: {
                size: 12
              },
              generateLabels: function(chart) {
                const data = chart.data;
                if (data.labels.length && data.datasets.length) {
                  return data.labels.map((label, i) => {
                    const value = data.datasets[0].data[i];
                    const percentage = ((value / total) * 100).toFixed(1);
                    return {
                      text: `${label} (${percentage}%)`,
                      fillStyle: data.datasets[0].backgroundColor[i],
                      strokeStyle: data.datasets[0].borderColor[i],
                      lineWidth: data.datasets[0].borderWidth,
                      hidden: false,
                      index: i
                    };
                  });
                }
                return [];
              }
            }
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.label || '';
                const value = context.parsed || 0;
                const percentage = ((value / total) * 100).toFixed(1);
                return `${label}: ${formatCurrency(value, currentCurrency)} (${percentage}%)`;
              }
            }
          }
        }
      }
    });
  }

  function createCumulativeChart() {
    if (!cumulativeChartCanvas || cumulativeData.length === 0) return;

    // Destroy existing chart if it exists
    if (cumulativeChartInstance) {
      cumulativeChartInstance.destroy();
      cumulativeChartInstance = null;
    }

    const ctx = cumulativeChartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
    const labels = cumulativeData.map(item => {
      const date = new Date(item.date);
      return date.getDate().toString();
    });
    
    const cumulativeValues = cumulativeData.map(item => item.cumulative);

    cumulativeChartInstance = new Chart(ctx, {
      type: 'line',
      data: {
        labels: labels,
        datasets: [{
          label: 'Cumulative Spending',
          data: cumulativeValues,
          borderColor: 'rgb(239, 68, 68)',
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          borderWidth: 3,
          fill: true,
          tension: 0.4,
          pointRadius: 4,
          pointHoverRadius: 6,
          pointBackgroundColor: 'rgb(239, 68, 68)',
          pointBorderColor: '#fff',
          pointBorderWidth: 2
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
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
                return `Cumulative: ${formatCurrency(context.parsed.y, currentCurrency)}`;
              },
              title: function(context) {
                const index = context[0].dataIndex;
                const date = new Date(cumulativeData[index].date);
                return `Day ${date.getDate()}`;
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
              text: `Cumulative Amount (${currentCurrency})`,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
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

  function createWeeklyChart() {
    if (!weeklyChartCanvas || weeklyData.length === 0) return;

    // Destroy existing chart if it exists
    if (weeklyChartInstance) {
      weeklyChartInstance.destroy();
      weeklyChartInstance = null;
    }

    const ctx = weeklyChartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
    const labels = weeklyData.map(item => item.week);
    const expenseData = weeklyData.map(item => item.expense);
    const incomeData = weeklyData.map(item => item.income);

    weeklyChartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: labels,
        datasets: [
          {
            label: 'Income',
            data: incomeData,
            backgroundColor: 'rgba(16, 185, 129, 0.8)',
            borderColor: 'rgba(16, 185, 129, 1)',
            borderWidth: 2,
            borderRadius: 4
          },
          {
            label: 'Expenses',
            data: expenseData,
            backgroundColor: 'rgba(239, 68, 68, 0.8)',
            borderColor: 'rgba(239, 68, 68, 1)',
            borderWidth: 2,
            borderRadius: 4
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
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
                return `${context.dataset.label}: ${formatCurrency(context.parsed.y, currentCurrency)}`;
              }
            }
          }
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Week',
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
              text: `Amount (${currentCurrency})`,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
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

  function createBudgetChart() {
    if (!budgetChartCanvas || budgets.length === 0) return;

    // Destroy existing chart if it exists
    if (budgetChartInstance) {
      budgetChartInstance.destroy();
      budgetChartInstance = null;
    }

    const ctx = budgetChartCanvas.getContext('2d');
    const currentCurrency = $currency || 'IDR';
    
    // Create expense map for quick lookup
    const expenseMap = new Map();
    expensesByCategory.forEach(item => {
      expenseMap.set(item.categoryId, item.total);
    });

    // Prepare chart data
    const chartData = budgets.map(budget => {
      const spent = expenseMap.get(budget.categoryId) || 0;
      return {
        category: budget.categoryName,
        budget: budget.amount,
        spent: spent
      };
    });

    const labels = chartData.map(item => item.category);
    const budgetData = chartData.map(item => item.budget);
    const spentData = chartData.map(item => item.spent);

    // Color based on percentage used
    const backgroundColors = chartData.map(item => {
      const percentage = item.budget > 0 ? (item.spent / item.budget) * 100 : 0;
      if (percentage >= 100) return 'rgba(239, 68, 68, 0.8)';
      if (percentage >= 75) return 'rgba(245, 158, 11, 0.8)';
      if (percentage >= 50) return 'rgba(251, 191, 36, 0.8)';
      return 'rgba(16, 185, 129, 0.8)';
    });

    budgetChartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: labels,
        datasets: [
          {
            label: 'Budget',
            data: budgetData,
            backgroundColor: 'rgba(99, 102, 241, 0.3)',
            borderColor: 'rgba(99, 102, 241, 1)',
            borderWidth: 2,
            borderRadius: 4
          },
          {
            label: 'Spent',
            data: spentData,
            backgroundColor: backgroundColors,
            borderColor: backgroundColors.map(color => color.replace('0.8', '1')),
            borderWidth: 2,
            borderRadius: 4
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: true,
            position: 'top'
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                return `${context.dataset.label}: ${formatCurrency(context.parsed.y, currentCurrency)}`;
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
              text: `Amount (${currentCurrency})`,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
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

        {#if cumulativeData.length > 0}
          <div class="chart-section">
            <h3>Cumulative Spending Trend</h3>
            <p class="chart-description">Track your spending accumulation throughout the month</p>
            <div class="chart-container">
              <canvas bind:this={cumulativeChartCanvas}></canvas>
            </div>
          </div>
        {/if}

        {#if weeklyData.length > 0}
          <div class="chart-section">
            <h3>Weekly Income vs Expenses</h3>
            <p class="chart-description">Compare your weekly income and expenses</p>
            <div class="chart-container">
              <canvas bind:this={weeklyChartCanvas}></canvas>
            </div>
          </div>
        {/if}

        {#if expensesByCategory.length > 0}
          <div class="category-list">
            <h4>Category Breakdown Details</h4>
            {#each expensesByCategory.sort((a, b) => b.total - a.total) as item}
              {@const total = expensesByCategory.reduce((sum, cat) => sum + cat.total, 0)}
              {@const percentage = ((item.total / total) * 100).toFixed(1)}
              <div class="category-item">
                <div class="category-info">
                  <span class="category-name">{item.category}</span>
                  <span class="category-percentage">{percentage}%</span>
                </div>
                <span class="category-amount">{formatCurrency(item.total)}</span>
              </div>
            {/each}
          </div>
        {:else if dailyExpenses.length === 0}
          <div class="no-data">No expenses data available</div>
        {/if}

        {#if budgets.length > 0}
          <div class="chart-section">
            <div class="space-xl"></div>
            <h3>Budget vs Actual Spending</h3>
            <p class="chart-description">Compare your budgeted amounts with actual spending</p>
            <div class="chart-container">
              <canvas bind:this={budgetChartCanvas}></canvas>
            </div>
          </div>
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
    max-width: 800px;
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
    margin-bottom: 0.5rem;
    font-weight: 600;
  }

  .chart-description {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1rem;
    font-style: italic;
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
    border: 1px solid var(--border);
  }

  .category-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .category-name {
    font-weight: 500;
    color: var(--text-primary);
    text-transform: capitalize;
  }

  .category-percentage {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .category-amount {
    font-weight: 600;
    color: var(--primary-color);
    font-size: 1.125rem;
  }
</style>

