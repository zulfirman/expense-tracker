<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
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
      incomeCategories = categories.filter(cat => cat.type === 'income');
      expenseCategories = categories.filter(cat => cat.type === 'expense');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to load categories',

      });
    } finally {
      loading = false;
    }
  }

  function openAddForm() {
    categoryName = '';
    categoryType = 'expense';
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

      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to create category',

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

      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to update category',

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

        });
      } catch (error) {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.response?.data?.message || 'Failed to delete category',

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

      });
    } finally {
      loading = false;
    }
  }
</script>

<div class="categories-page">
  <div class="header-section">
    <h1>My Categories</h1>
    <div class="header-actions">
      <button class="btn ghost" on:click={() => goto('/preferences')}>Back</button>
      <button class="btn primary" style="max-width: 50%" on:click={openAddForm} disabled={loading}>
        + Add Category
      </button>
    </div>
  </div>

  {#if loading && categories.length === 0}
    <div class="loading">Loading categories...</div>
  {:else if categories.length === 0}
    <div class="empty-state">
      <p>No categories yet. Create your first category to get started!</p>
    </div>
  {:else}
    <div class="categories-list">
      {#if incomeCategories.length > 0}
        <div class="category-section">
          <h2 class="section-title">💰 Income Categories</h2>
          {#each incomeCategories as category}
            <div class="category-card">
              <div class="category-info">
                <h3>{category.name}</h3>
              </div>
              <div class="category-actions">
                <label class="toggle-switch">
                  <input
                    type="checkbox"
                    checked={category.isActive}
                    on:change={() => toggleActive(category)}
                    disabled={loading}
                  />
                  <span class="toggle-slider"></span>
                </label>
                <button class="btn-text btn-edit" on:click={() => openEditForm(category)}>Edit</button>
                <button class="btn-text btn-delete" on:click={() => handleDelete(category)}>Delete</button>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      {#if expenseCategories.length > 0}
        <div class="category-section">
          <h2 class="section-title">💸 Expense Categories</h2>
          {#each expenseCategories as category}
            <div class="category-card">
              <div class="category-info">
                <h3>{category.name}</h3>
              </div>
              <div class="category-actions">
                <label class="toggle-switch">
                  <input
                    type="checkbox"
                    checked={category.isActive}
                    on:change={() => toggleActive(category)}
                    disabled={loading}
                  />
                  <span class="toggle-slider"></span>
                </label>
                <button class="btn-text btn-edit" on:click={() => openEditForm(category)}>Edit</button>
                <button class="btn-text btn-delete" on:click={() => handleDelete(category)}>Delete</button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</div>

{#if showAddForm}
  <div class="modal-backdrop" on:click={closeForms}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Add Category</h2>
        <button class="close-btn" on:click={closeForms}>×</button>
      </div>
      <div class="modal-body">
        <fieldset class="fieldset mb-4">
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
        <div class="button-group">
          <button class="btn btn-secondary" on:click={closeForms} disabled={loading}>Cancel</button>
          <button class="btn btn-primary" on:click={handleCreate} disabled={loading}>
            {#if loading}
              <span class="spinner"></span> Creating...
            {:else}
              Create
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

{#if showEditForm}
  <div class="modal-backdrop" on:click={closeForms}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Edit Category</h2>
        <button class="close-btn" on:click={closeForms}>×</button>
      </div>
      <div class="modal-body">
        <fieldset class="fieldset mb-4">
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
        <fieldset class="fieldset mb-4">
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
          <label class="cursor-pointer flex items-center gap-2">
            <input type="checkbox" bind:checked={isActive} disabled={loading} class="toggle toggle-sm" />
            <span>Active</span>
          </label>
        </fieldset>
        <div class="button-group">
          <button class="btn btn-secondary" on:click={closeForms} disabled={loading}>Cancel</button>
          <button class="btn btn-primary" on:click={handleUpdate} disabled={loading}>
            {#if loading}
              <span class="spinner"></span> Updating...
            {:else}
              Update
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .categories-page {
    max-width: 800px;
    margin: 0 auto;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
  }

  h1 {
    font-size: 1.5rem;
    color: var(--text-primary);
  }

  .loading, .empty-state {
    text-align: center;
    padding: 3rem;
    color: var(--text-secondary);
  }

  .categories-list {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .category-section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid var(--border);
  }

  .category-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    padding: 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: box-shadow 0.2s;
  }

  .category-card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .category-info h3 {
    font-size: 1.125rem;
    margin-bottom: 0.25rem;
    color: var(--text-primary);
  }

  .category-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .toggle-switch {
    position: relative;
    display: inline-block;
    width: 50px;
    height: 24px;
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--border);
    transition: 0.3s;
    border-radius: 24px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  .toggle-switch input:checked + .toggle-slider {
    background-color: var(--primary-color);
  }

  .toggle-switch input:checked + .toggle-slider:before {
    transform: translateX(26px);
  }

  .category-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

</style>


