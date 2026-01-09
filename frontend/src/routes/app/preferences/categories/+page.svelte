<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageCode from '$lib/components/PageCode.svelte';
  import Toggle from '$lib/components/Toggle.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';

  let categories = [];
  let incomeCategories = [];
  let expenseCategories = [];
  let loading = false;
  let showAddForm = false;
  let showEditForm = false;
  let editingCategory = null;
  let categoryName = '';
  let categoryType = 'expense';
  let isActive = true;
  let activeTab = 'expense'; // 'income' or 'expense'
  let draggedCategory = null;
  let draggedOverIndex = null;
  
  $: pageCode = getPageCode($page.url.pathname);

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;
    await loadCategories();
  });

  async function loadCategories() {
    loading = true;
    try {
      const response = await api.get('/categories');
      categories = response.data || [];
      // Sort by sequence, then by name
      incomeCategories = categories
        .filter(cat => cat.type === 'income')
        .sort((a, b) => {
          const seqA = a.sequence || 0;
          const seqB = b.sequence || 0;
          if (seqA !== seqB) return seqA - seqB;
          return a.name.localeCompare(b.name);
        });
      expenseCategories = categories
        .filter(cat => cat.type === 'expense')
        .sort((a, b) => {
          const seqA = a.sequence || 0;
          const seqB = b.sequence || 0;
          if (seqA !== seqB) return seqA - seqB;
          return a.name.localeCompare(b.name);
        });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to load categories',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function openAddForm() {
    categoryName = '';
    categoryType = activeTab; // Set type based on active tab
    isActive = true;
    showAddForm = true;
    showEditForm = false;
  }

  function openEditForm(category) {
    editingCategory = category;
    categoryName = category.name;
    categoryType = category.type || 'expense';
    isActive = category.isActive;
    showEditForm = true;
    showAddForm = false;
  }

  function closeForms() {
    showAddForm = false;
    showEditForm = false;
    editingCategory = null;
    categoryName = '';
    categoryType = 'expense';
  }

  async function handleCreate() {
    if (loading) return; // Prevent double submission
    
    if (!categoryName.trim()) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Name',
          text: 'Please enter a category name',
          zIndex: 9999
        });
      }, 50);
      return;
    }

    loading = true;
    try {
      await api.post('/categories', {
        name: categoryName.trim(),
        type: categoryType,
        isActive: isActive
      });
      
      await loadCategories();
      closeForms();
      
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Category Created',
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to create category',
          zIndex: 9999
        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  async function handleUpdate() {
    if (loading) return; // Prevent double submission
    
    if (!categoryName.trim()) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Name',
          text: 'Please enter a category name',
          zIndex: 9999
        });
      }, 50);
      return;
    }

    loading = true;
    try {
      await api.put(`/categories/${editingCategory.id}`, {
        name: categoryName.trim(),
        type: categoryType,
        isActive: isActive
      });
      
      await loadCategories();
      closeForms();
      
      setTimeout(() => {
        Swal.fire({
          icon: 'success',
          title: 'Category Updated',
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      }, 50);
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to update category',
          zIndex: 9999
        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  async function handleDelete(category) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'Delete Category?',
      text: `Are you sure you want to delete "${category.name}"?`,
      showCancelButton: true,
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      zIndex: 9999
    });

    if (result.isConfirmed) {
      loading = true;
      try {
        await api.delete(`/categories/${category.id}`);
        await loadCategories();
        
        setTimeout(() => {
          Swal.fire({
            icon: 'success',
            title: 'Category Deleted',
            timer: 1500,
            showConfirmButton: false,
            zIndex: 9999
          });
        }, 50);
      } catch (error) {
        setTimeout(() => {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: error.response?.data?.message || 'Failed to delete category',
            zIndex: 9999
          });
        }, 50);
      } finally {
        loading = false;
      }
    }
  }

  function toggleActive(category) {
    if (loading) return; // Prevent double submission
    const newActiveState = !category.isActive;
    // Update local state immediately for better UX
    category.isActive = newActiveState;
    handleUpdateCategory(category.id, category.name, newActiveState);
  }

  async function handleUpdateCategory(id, name, active) {
    if (loading) return; // Prevent double submission
    loading = true;
    try {
      const category = categories.find(cat => cat.id === id);
      await api.put(`/categories/${id}`, {
        name: name,
        type: category?.type || 'expense',
        isActive: active
      });
      await loadCategories();
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to update category',
          zIndex: 9999
        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  function handleDragStart(category, index) {
    draggedCategory = { ...category, index };
  }

  function handleDragOver(e, index) {
    e.preventDefault();
    draggedOverIndex = index;
  }

  function handleDragLeave() {
    draggedOverIndex = null;
  }

  async function handleDrop(e, targetIndex, targetType) {
    e.preventDefault();
    draggedOverIndex = null;
    
    if (!draggedCategory || draggedCategory.index === targetIndex) {
      draggedCategory = null;
      return;
    }

    // Get categories of the same type
    const sameTypeCategories = categories.filter(cat => cat.type === targetType);
    const sortedCategories = [...sameTypeCategories].sort((a, b) => (a.sequence || 0) - (b.sequence || 0));
    
    // Remove dragged item from array
    const draggedItem = sortedCategories.find(cat => cat.id === draggedCategory.id);
    if (!draggedItem) {
      draggedCategory = null;
      return;
    }
    
    sortedCategories.splice(draggedCategory.index, 1);
    sortedCategories.splice(targetIndex, 0, draggedItem);
    
    // Update sequences
    const updates = sortedCategories.map((cat, idx) => ({
      id: cat.id,
      sequence: idx + 1
    }));

    loading = true;
    try {
      await api.put('/categories/sequence', { categories: updates });
      await loadCategories();
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to reorder categories',
        zIndex: 9999
      });
    } finally {
      loading = false;
      draggedCategory = null;
    }
  }
</script>

<div class="max-w-3xl mx-auto space-y-4">
  <PageHeader
    title="My Categories"
    subtitle="Manage your income and expense categories. Drag and drop categories to reorder them."
    pageCode={pageCode}
    actions={true}
  >
    <svelte:fragment slot="actions">
      <button class="btn btn-soft btn-sm" on:click={() => goto('/app/preferences')}>Back</button>
      <button class="btn btn-primary btn-sm" on:click={openAddForm} disabled={loading}>
        + Add Category
      </button>
    </svelte:fragment>
  </PageHeader>

  {#if loading && categories.length === 0}
    <div class="flex justify-center py-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if categories.length === 0}
    <div class="alert alert-info">
      <span>No categories yet. Create your first category to get started!</span>
    </div>
  {:else}
    <!-- Tabs -->
    <div class="tabs tabs-boxed w-full sm:w-fit">
      <button
        class="tab text-xs sm:text-sm flex-1 sm:flex-none"
        class:tab-active={activeTab === 'income'}
        on:click={() => activeTab = 'income'}
        disabled={loading}
      >
        💰 Income ({incomeCategories.length})
      </button>
      <button
        class="tab text-xs sm:text-sm flex-1 sm:flex-none"
        class:tab-active={activeTab === 'expense'}
        on:click={() => activeTab = 'expense'}
        disabled={loading}
      >
        💸 Expense ({expenseCategories.length})
      </button>
    </div>

    <!-- Tab Content -->
    <div class="space-y-3 mt-3">
      {#if activeTab === 'income'}
        {#if incomeCategories.length === 0}
          <div class="alert alert-soft">
            <span>No income categories yet. Create your first income category!</span>
          </div>
        {:else}
          {#each incomeCategories as category, index}
            <div
              class="card bg-base-100 shadow-sm border-1 cursor-move"
              class:opacity-50={draggedCategory?.id === category.id}
              class:border-primary={draggedOverIndex === index}
              draggable={!loading}
              on:dragstart={() => handleDragStart(category, index)}
              on:dragover={(e) => handleDragOver(e, index)}
              on:dragleave={handleDragLeave}
              on:drop={(e) => handleDrop(e, index, 'income')}
            >
              <div class="card-body flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-3">
                <div class="flex items-center gap-2 flex-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-base-content/40">
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                  </svg>
                  <h3 class="font-semibold text-sm">{category.name}</h3>
                </div>
                <div class="flex items-center gap-2 flex-wrap">
                  <button class="btn btn-xs btn-soft" on:click={() => openEditForm(category)}>Edit</button>
                  <button class="btn btn-xs btn-error" on:click={() => handleDelete(category)}>Delete</button>
                </div>

              </div>
            </div>
          {/each}
        {/if}
      {:else}
        {#if expenseCategories.length === 0}
          <div class="alert alert-soft">
            <span>No expense categories yet. Create your first expense category!</span>
          </div>
        {:else}
          {#each expenseCategories as category, index}
            <div
              class="card bg-base-100 shadow-sm border-1 cursor-move"
              class:opacity-50={draggedCategory?.id === category.id}
              class:border-primary={draggedOverIndex === index}
              draggable={!loading}
              on:dragstart={() => handleDragStart(category, index)}
              on:dragover={(e) => handleDragOver(e, index)}
              on:dragleave={handleDragLeave}
              on:drop={(e) => handleDrop(e, index, 'expense')}
            >
              <div class="card-body flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-3">
                <div class="flex items-center gap-2 flex-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-base-content/40">
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                  </svg>
                  <h3 class="font-semibold text-sm">{category.name}</h3>
                </div>
                <div class="flex items-center gap-2 flex-wrap">
                  <label class="label cursor-pointer gap-2">
                    <Toggle
                      size="sm"
                      color="primary"
                      checked={category.isActive}
                      disabled={loading}
                      on:change={() => toggleActive(category)}
                    />
                    <span class="label-text text-xs sm:text-sm">Active</span>
                  </label>
                  <button class="btn btn-xs btn-soft" on:click={() => openEditForm(category)}>Edit</button>
                  <button class="btn btn-xs btn-error" on:click={() => handleDelete(category)}>Delete</button>
                </div>
              </div>
            </div>
          {/each}
        {/if}
      {/if}
    </div>
  {/if}
</div>

{#if showAddForm}
  <div class="modal modal-open z-[2100]">
    <div class="modal-box w-11/12 max-w-md">
      <div class="flex items-center gap-2 mb-2">
        <span class="text-xs font-mono text-base-content/30">{pageCode}</span>
        <h3 class="font-semibold text-lg">Add Category</h3>
      </div>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Type</legend>
        <select
          id="category-type"
          bind:value={categoryType}
          class="select select-bordered w-full border-2"
          disabled={loading}
        >
          <option value="expense">Expense</option>
          <option value="income">Income</option>
        </select>
      </fieldset>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Category Name</legend>
        <input
          id="category-name"
          type="text"
          bind:value={categoryName}
          placeholder="e.g., Groceries"
          class="input input-bordered w-full border-2"
          on:keydown={(e) => { if (e.key === 'Enter' && !loading) { e.preventDefault(); setTimeout(() => handleCreate(), 50); } }}
          disabled={loading}
        />
      </fieldset>
      <label class="label cursor-pointer gap-3 mb-4">
        <Toggle
          size="md"
          color="primary"
          checked={isActive}
          disabled={loading}
          on:change={(e) => isActive = e.detail}
        />
        <span class="label-text">Active</span>
      </label>
      <div class="modal-action">
        <button class="btn btn-soft" on:click={closeForms} disabled={loading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleCreate} disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm mr-1"></span>
            Creating...
          {:else}
            Create
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

{#if showEditForm}
  <div class="modal modal-open z-[2100]">
    <div class="modal-box w-11/12 max-w-md">
      <div class="flex items-center gap-2 mb-2">
        <span class="text-xs font-mono text-base-content/30">{pageCode}</span>
        <h3 class="font-semibold text-lg">Edit Category</h3>
      </div>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Type</legend>
        <select
          id="edit-category-type"
          bind:value={categoryType}
          class="select select-bordered w-full border-2"
          disabled={loading}
        >
          <option value="expense">Expense</option>
          <option value="income">Income</option>
        </select>
      </fieldset>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Category Name</legend>
        <input
          id="edit-category-name"
          type="text"
          bind:value={categoryName}
          placeholder="e.g., Groceries"
          class="input input-bordered w-full border-2"
          on:keydown={(e) => { if (e.key === 'Enter' && !loading) { e.preventDefault(); setTimeout(() => handleUpdate(), 50); } }}
          disabled={loading}
        />
      </fieldset>
      <label class="label cursor-pointer gap-3 mb-4">
        <Toggle
          size="md"
          color="primary"
          checked={isActive}
          disabled={loading}
          on:change={(e) => isActive = e.detail}
        />
        <span class="label-text">Active</span>
      </label>
      <div class="modal-action">
        <button class="btn btn-soft" on:click={closeForms} disabled={loading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleUpdate} disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm mr-1"></span>
            Updating...
          {:else}
            Update
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

