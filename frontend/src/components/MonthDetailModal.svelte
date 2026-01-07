<script>
import { createEventDispatcher, onMount, onDestroy, afterUpdate } from 'svelte';
  import api from '$lib/api';
  import { Chart, registerables } from 'chart.js';
  import { formatCurrency } from '$lib/utils/currency';
  import { currency } from '$lib/stores/currency';
import { theme } from '$lib/stores/theme';

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
$: currentTheme = $theme || document.documentElement?.dataset?.theme || 'light';
$: isDarkTheme = ['dark', 'night', 'forest', 'business', 'dracula', 'sunset', 'dim', 'black'].some(t => currentTheme?.toLowerCase().includes(t));

function hexToRgba(hex, alpha = 1) {
  if (!hex) return `rgba(99,102,241,${alpha})`;
  const normalized = hex.replace('#', '');
  const bigint = parseInt(normalized.length === 3 ? normalized.split('').map(c => c + c).join('') : normalized, 16);
  const r = (bigint >> 16) & 255;
  const g = (bigint >> 8) & 255;
  const b = bigint & 255;
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

function getPalette() {
  return {
    text: isDarkTheme ? '#e5e7eb' : '#1f2937',
    grid: isDarkTheme ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)',
    surface: isDarkTheme ? '#0f172a' : '#ffffff',
    primary: '#6366f1',
    success: '#10b981',
    danger: '#ef4444',
    info: '#3b82f6',
    warning: '#f59e0b',
    accent: '#a855f7'
  };
}

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
    const palette = getPalette();
    
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
            borderColor: palette.success,
            backgroundColor: hexToRgba(palette.success, 0.16),
            borderWidth: 2,
            fill: false,
            tension: 0.4,
            pointRadius: 4,
            pointHoverRadius: 6,
            pointBackgroundColor: palette.success,
            pointBorderColor: palette.surface,
            pointBorderWidth: 2
          },
          {
            label: 'Expenses',
            data: expenseData,
            borderColor: palette.danger,
            backgroundColor: hexToRgba(palette.danger, 0.16),
            borderWidth: 2,
            fill: false,
            tension: 0.4,
            pointRadius: 4,
            pointHoverRadius: 6,
            pointBackgroundColor: palette.danger,
            pointBorderColor: palette.surface,
            pointBorderWidth: 2
          },
          {
            label: 'Net (Income - Expense)',
            data: netData,
            borderColor: palette.info,
            backgroundColor: hexToRgba(palette.info, 0.2),
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointRadius: 5,
            pointHoverRadius: 7,
            pointBackgroundColor: palette.info,
            pointBorderColor: palette.surface,
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
              color: palette.text,
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text
            },
            grid: {
              display: true,
              color: palette.grid
            }
          },
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: `Amount (${currentCurrency})`,
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text,
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
              }
            },
            grid: {
              display: true,
              color: palette.grid
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
    const palette = getPalette();
    
    // Sort categories by total (descending)
    const sortedCategories = [...expensesByCategory].sort((a, b) => b.total - a.total);

    const colors = [
      palette.primary,
      palette.info,
      palette.accent,
      palette.warning,
      palette.success,
      palette.danger
    ];

    categoryChartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: sortedCategories.map(item => item.category),
        datasets: [{
          label: 'Expenses by Category',
          data: sortedCategories.map(item => item.total),
          backgroundColor: sortedCategories.map((_, index) => hexToRgba(colors[index % colors.length], 0.75)),
          borderColor: sortedCategories.map((_, index) => colors[index % colors.length]),
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text,
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
              }
            },
            grid: {
              display: true,
              color: palette.grid
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
    const palette = getPalette();
    
    // Sort categories by total (descending)
    const sortedCategories = [...expensesByCategory].sort((a, b) => b.total - a.total);

    const colors = [
      palette.primary,
      palette.accent,
      palette.info,
      palette.warning,
      palette.success,
      palette.danger,
      '#0ea5e9',
      '#f97316',
      '#14b8a6',
      '#a3e635'
    ];

    const total = sortedCategories.reduce((sum, item) => sum + item.total, 0);

    pieChartInstance = new Chart(ctx, {
      type: 'pie',
      data: {
        labels: sortedCategories.map(item => item.category),
        datasets: [{
          data: sortedCategories.map(item => item.total),
          backgroundColor: sortedCategories.map((_, index) => hexToRgba(colors[index % colors.length], 0.78)),
          borderColor: sortedCategories.map((_, index) => colors[index % colors.length]),
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
              color: palette.text,
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
    const palette = getPalette();
    
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
          borderColor: palette.danger,
          backgroundColor: hexToRgba(palette.danger, 0.14),
          borderWidth: 3,
          fill: true,
          tension: 0.4,
          pointRadius: 4,
          pointHoverRadius: 6,
          pointBackgroundColor: palette.danger,
          pointBorderColor: palette.surface,
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
              color: palette.text,
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text
            },
            grid: {
              display: true,
              color: palette.grid
            }
          },
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: `Cumulative Amount (${currentCurrency})`,
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text,
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
              }
            },
            grid: {
              display: true,
              color: palette.grid
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
    const palette = getPalette();
    
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
            backgroundColor: hexToRgba(palette.success, 0.8),
            borderColor: palette.success,
            borderWidth: 2,
            borderRadius: 4
          },
          {
            label: 'Expenses',
            data: expenseData,
            backgroundColor: hexToRgba(palette.danger, 0.82),
            borderColor: palette.danger,
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
              color: palette.text,
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text,
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
              }
            },
            grid: {
              display: true,
              color: palette.grid
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
    const palette = getPalette();
    
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
            backgroundColor: hexToRgba(palette.primary, 0.3),
            borderColor: palette.primary,
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
            position: 'top',
            labels: {
              color: palette.text
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
              text: 'Category',
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text
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
              color: palette.text,
              font: {
                size: 12,
                weight: 'bold'
              }
            },
            ticks: {
              color: palette.text,
              callback: function(value) {
                return formatCurrency(value, currentCurrency);
              }
            },
            grid: {
              display: true,
              color: palette.grid
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

<div class="modal modal-open z-[2100]" on:click={handleBackdropClick}>
  <div class="modal-box w-11/12 max-w-5xl" on:click|stopPropagation>
    <div class="flex items-start justify-between gap-3 mb-4">
      <div>
        <p class="text-xs uppercase tracking-wide text-base-content/60">Insights</p>
        <h2 class="text-2xl font-bold">Advanced Analytics</h2>
        <p class="text-sm text-base-content/70">{formatMonthYear(month.month)}</p>
      </div>
      <button class="btn btn-ghost btn-sm" on:click={close}>✕</button>
    </div>

    <div class="max-h-[75vh] overflow-y-auto pr-1 space-y-4">
      {#if loading}
        <div class="flex justify-center py-10">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else if dailyExpenses.length > 0 || expensesByCategory.length > 0}
        {#if dailyExpenses.length > 0}
          <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body space-y-3">
              <h3 class="font-semibold text-lg">Daily Income vs Expenses</h3>
              <div class="rounded-lg border border-base-300 bg-base-200/40 p-3">
                <div class="relative h-[320px]">
                  <canvas bind:this={chartCanvas}></canvas>
                </div>
              </div>
            </div>
          </div>
        {/if}

        {#if cumulativeData.length > 0}
          <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body space-y-2">
              <div>
                <h3 class="font-semibold text-lg">Cumulative Spending Trend</h3>
                <p class="text-sm text-base-content/70">Track your spending accumulation throughout the month</p>
              </div>
              <div class="rounded-lg border border-base-300 bg-base-200/40 p-3">
                <div class="relative h-[320px]">
                  <canvas bind:this={cumulativeChartCanvas}></canvas>
                </div>
              </div>
            </div>
          </div>
        {/if}

        {#if weeklyData.length > 0}
          <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body space-y-2">
              <div>
                <h3 class="font-semibold text-lg">Weekly Income vs Expenses</h3>
                <p class="text-sm text-base-content/70">Compare your weekly income and expenses</p>
              </div>
              <div class="rounded-lg border border-base-300 bg-base-200/40 p-3">
                <div class="relative h-[320px]">
                  <canvas bind:this={weeklyChartCanvas}></canvas>
                </div>
              </div>
            </div>
          </div>
        {/if}

        {#if expensesByCategory.length > 0}
          <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body space-y-3">
              <h3 class="font-semibold text-lg">Category Breakdown</h3>
              <div class="space-y-2">
                {#each expensesByCategory.sort((a, b) => b.total - a.total) as item}
                  {@const total = expensesByCategory.reduce((sum, cat) => sum + cat.total, 0)}
                  {@const percentage = ((item.total / total) * 100).toFixed(1)}
                  <div class="flex items-center justify-between p-3 rounded-lg border border-base-200 bg-base-200/40">
                    <div>
                      <p class="font-semibold text-base">{item.category}</p>
                      <p class="text-xs text-base-content/70">{percentage}% of spend</p>
                    </div>
                    <div class="font-semibold">{formatCurrency(item.total)}</div>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        {:else if dailyExpenses.length === 0}
          <div class="alert alert-info">
            <span>No expenses data available</span>
          </div>
        {/if}

        {#if budgets.length > 0}
          <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body space-y-2">
              <div>
                <h3 class="font-semibold text-lg">Budget vs Actual Spending</h3>
                <p class="text-sm text-base-content/70">Compare your budgeted amounts with actual spending</p>
              </div>
              <div class="rounded-lg border border-base-300 bg-base-200/40 p-3">
                <div class="relative h-[320px]">
                  <canvas bind:this={budgetChartCanvas}></canvas>
                </div>
              </div>
            </div>
          </div>
        {/if}
      {:else}
        <div class="alert alert-info">
          <span>No data available</span>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .chart-container canvas {
    max-width: 100%;
  }
</style>

