<script>
  import { onMount } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';

  let categories = [];
  let loading = false;
  let showAddForm = false;
  let showEditForm = false;
  let editingCategory = null;
  let categoryName = '';
  let isActive = true;

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      return;
    }
    await loadCategories();
  });

  async function loadCategories() {
    loading = true;
    try {
      const response = await api.get('/categories');
      categories = response.data || [];
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
    isActive = true;
    showAddForm = true;
    showEditForm = false;
  }

  function openEditForm(category) {
    editingCategory = category;
    categoryName = category.name;
    isActive = category.isActive;
    showEditForm = true;
    showAddForm = false;
  }

  function closeForms() {
    showAddForm = false;
    showEditForm = false;
    editingCategory = null;
    categoryName = '';
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
        name: categoryName.trim()
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
      await api.put(`/categories/${id}`, {
        name: name,
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

<div class="categories-page">
  <div class="header-section">
    <h1>My Categories</h1>
    <button class="btn btn-primary" style="max-width: 50%" on:click={openAddForm} disabled={loading}>
      + Add Category
    </button>
  </div>

  {#if loading && categories.length === 0}
    <div class="loading">Loading categories...</div>
  {:else if categories.length === 0}
    <div class="empty-state">
      <p>No categories yet. Create your first category to get started!</p>
    </div>
  {:else}
    <div class="categories-list">
      {#each categories as category}
        <div class="category-card">
          <div class="category-info">
            <h3>{category.name}</h3>
            <span class="category-slug">{category.slug}</span>
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
            <button class="btn-icon" on:click={() => openEditForm(category)} title="Edit">
              ✏️
            </button>
            <button class="btn-icon" on:click={() => handleDelete(category)} title="Delete">
              🗑️
            </button>
          </div>
        </div>
      {/each}
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
        <div class="form-group">
          <label for="category-name">Category Name</label>
          <input
            id="category-name"
            type="text"
            bind:value={categoryName}
            placeholder="e.g., Groceries"
            class="form-input"
            on:keydown={(e) => { if (e.key === 'Enter') handleCreate(); }}
            disabled={loading}
          />
        </div>
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
        <div class="form-group">
          <label for="edit-category-name">Category Name</label>
          <input
            id="edit-category-name"
            type="text"
            bind:value={categoryName}
            placeholder="e.g., Groceries"
            class="form-input"
            on:keydown={(e) => { if (e.key === 'Enter') handleUpdate(); }}
            disabled={loading}
          />
        </div>
        <div class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={isActive} disabled={loading} />
            <span>Active</span>
          </label>
        </div>
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
    gap: 1rem;
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

  .category-slug {
    font-size: 0.875rem;
    color: var(--text-secondary);
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

  .btn-icon {
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.25rem;
    opacity: 0.7;
    transition: opacity 0.2s;
  }

  .btn-icon:hover {
    opacity: 1;
  }

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
    max-width: 500px;
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
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 0;
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
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }

  .form-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
    background: var(--background);
    color: var(--text-primary);
  }

  .form-input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .button-group {
    display: flex;
    gap: 1rem;
  }

  .btn {
    flex: 1;
    padding: 0.875rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background-color: #4338ca;
  }

  .btn-secondary {
    background-color: var(--surface);
    color: var(--text-primary);
    border: 1px solid var(--border);
  }

  .btn-secondary:hover:not(:disabled) {
    background-color: var(--background);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 0.6s linear infinite;
    margin-right: 0.5rem;
    vertical-align: middle;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>


