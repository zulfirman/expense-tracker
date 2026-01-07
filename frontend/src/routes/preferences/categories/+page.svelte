<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';

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

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }
    await loadCategories();
  });

  async function loadCategories() {
    loading = true;
    try {
      const response = await api.get('/categories');
      categories = response.data || [];
      incomeCategories = categories.filter(cat => cat.type === 'income');
      expenseCategories = categories.filter(cat => cat.type === 'expense');
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
    if (!categoryName.trim()) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Name',
        text: 'Please enter a category name',
        zIndex: 9999
      });
      return;
    }

    loading = true;
    try {
      await api.post('/categories', {
        name: categoryName.trim(),
        type: categoryType
      });
      
      await loadCategories();
      closeForms();
      
      Swal.fire({
        icon: 'success',
        title: 'Category Created',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to create category',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  async function handleUpdate() {
    if (!categoryName.trim()) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Name',
        text: 'Please enter a category name',
        zIndex: 9999
      });
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
      
      Swal.fire({
        icon: 'success',
        title: 'Category Updated',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to update category',
        zIndex: 9999
      });
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
        
        Swal.fire({
          icon: 'success',
          title: 'Category Deleted',
          timer: 1500,
          showConfirmButton: false,
          zIndex: 9999
        });
      } catch (error) {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to delete category',
          zIndex: 9999
        });
      } finally {
        loading = false;
      }
    }
  }

  function toggleActive(category) {
    handleUpdateCategory(category.id, category.name, !category.isActive);
  }

  async function handleUpdateCategory(id, name, active) {
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
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to update category',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }
</script>

<div class="max-w-3xl mx-auto space-y-4">
  <div class="flex items-center justify-between gap-2">
    <div>
      <h1 class="text-2xl font-bold">My Categories</h1>
      <p class="text-sm text-base-content/70 mt-1">
        Manage your income and expense categories.
      </p>
    </div>
    <div class="flex gap-2">
      <button class="btn btn-soft btn-sm" on:click={() => goto('/preferences')}>Back</button>
      <button class="btn btn-primary btn-sm" on:click={openAddForm} disabled={loading}>
        + Add Category
      </button>
    </div>
  </div>

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
    <div class="tabs tabs-boxed w-fit">
      <button
        class="tab text-sm"
        class:tab-active={activeTab === 'income'}
        on:click={() => activeTab = 'income'}
        disabled={loading}
      >
        💰 Income ({incomeCategories.length})
      </button>
      <button
        class="tab text-sm"
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
          {#each incomeCategories as category}
    <div class="card bg-base-100 shadow-sm border-1">
              <div class="card-body flex flex-row items-center justify-between gap-3 py-3">
                <div>
                  <h3 class="font-semibold text-sm">{category.name}</h3>
                </div>
                <div class="flex items-center gap-2">
                  <input
                    type="checkbox"
                    class="toggle toggle-sm toggle-primary"
                    checked={category.isActive}
                    on:change={() => toggleActive(category)}
                    disabled={loading}
                  />
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
          {#each expenseCategories as category}
            <div class="card bg-base-100 shadow-sm border-1">
              <div class="card-body flex flex-row items-center justify-between gap-3 py-3">
                <div>
                  <h3 class="font-semibold text-sm">{category.name}</h3>
                </div>
                <div class="flex items-center gap-2">
                  <input
                    type="checkbox"
                    class="toggle toggle-sm toggle-primary"
                    checked={category.isActive}
                    on:change={() => toggleActive(category)}
                    disabled={loading}
                  />
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
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-semibold text-lg mb-2">Add Category</h3>
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
      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Category Name</legend>
        <input
          id="category-name"
          type="text"
          bind:value={categoryName}
          placeholder="e.g., Groceries"
          class="input input-bordered w-full border-2"
          on:keydown={(e) => { if (e.key === 'Enter') handleCreate(); }}
          disabled={loading}
        />
      </fieldset>
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
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-semibold text-lg mb-2">Edit Category</h3>
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
          on:keydown={(e) => { if (e.key === 'Enter') handleUpdate(); }}
          disabled={loading}
        />
      </fieldset>
      <fieldset class="fieldset mb-4">
        <legend class="fieldset-legend">Active</legend>
        <label class="cursor-pointer label justify-start gap-3">
          <input type="checkbox" class="toggle toggle-sm" bind:checked={isActive} disabled={loading} />
          <span class="label-text">Active</span>
        </label>
      </fieldset>
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

