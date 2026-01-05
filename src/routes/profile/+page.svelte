<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';

  let name = '';
  let email = '';
  let loading = false;
  let saving = false;

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadProfile();
  });

  async function loadProfile() {
    loading = true;
    try {
      const response = await api.get('/auth/profile');
      name = response.data.name;
      email = response.data.email;
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: error.response?.data?.message || 'Failed to load profile',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  async function handleSave() {
    if (!name || !email) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Fields',
        text: 'Please fill in all fields',
        zIndex: 9999
      });
      return;
    }

    saving = true;
    try {
      const response = await api.put('/auth/profile', {
        name,
        email
      });

      // Update auth store
      auth.login(response.data, $auth.token);
      
      Swal.fire({
        icon: 'success',
        title: 'Profile Updated',
        text: 'Your profile has been updated successfully',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Update Failed',
        text: error.response?.data?.message || 'Failed to update profile',
        zIndex: 9999
      });
    } finally {
      saving = false;
    }
  }

  function handleLogout() {
    Swal.fire({
      icon: 'question',
      title: 'Logout?',
      text: 'Are you sure you want to logout?',
      showCancelButton: true,
      confirmButtonText: 'Logout',
      cancelButtonText: 'Cancel',
      zIndex: 9999
    }).then((result) => {
      if (result.isConfirmed) {
        auth.logout();
        goto('/login');
      }
    });
  }
</script>

<div class="profile-page">
  <h1>Profile</h1>

  {#if loading}
    <div class="loading">Loading profile...</div>
  {:else}
    <div class="profile-form">
      <div class="form-group">
        <label for="name">Name</label>
        <input
          id="name"
          type="text"
          bind:value={name}
          placeholder="Your name"
          class="form-input"
          disabled={saving}
        />
      </div>

      <div class="form-group">
        <label for="email">Email</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          placeholder="your@email.com"
          class="form-input"
          disabled={saving}
        />
      </div>

      <div class="button-group">
        <button class="btn btn-primary" on:click={handleSave} disabled={saving}>
          {#if saving}
            <span class="spinner"></span> Saving...
          {:else}
            Save Changes
          {/if}
        </button>
        <button class="btn btn-danger" on:click={handleLogout} disabled={saving}>
          Logout
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .profile-page {
    max-width: 600px;
    margin: 0 auto;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: var(--text-secondary);
  }

  .profile-form {
    background: var(--surface);
    border-radius: 0.75rem;
    padding: 1.5rem;
    border: 1px solid var(--border);
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
    margin-top: 2rem;
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

  .btn-danger {
    background-color: var(--danger);
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background-color: #dc2626;
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

  @media (max-width: 768px) {
    .button-group {
      flex-direction: column;
    }
  }
</style>



