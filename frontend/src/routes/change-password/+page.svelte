<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import api from '$lib/api';
  import Swal from 'sweetalert2';

  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  let loading = false;

  onMount(() => {
    if (!$auth.isAuthenticated) {
      goto('/login');
    }
  });

  async function handleSubmit() {
    if (!currentPassword || !newPassword || !confirmPassword) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Fields',
        text: 'Please fill in all fields',
        zIndex: 9999
      });
      return;
    }

    if (newPassword.length < 6) {
      setTimeout(()=>{
          Swal.fire({
              icon: 'warning',
              title: 'Invalid Password',
              text: 'Password must be at least 6 characters long',
              zIndex: 9999
          });
      },50);
      return;
    }

    if (newPassword !== confirmPassword) {
      Swal.fire({
        icon: 'warning',
        title: 'Password Mismatch',
        text: 'New password and confirm password do not match',
        zIndex: 9999
      });
      return;
    }

    loading = true;
    try {
      await api.put('/auth/password', {
        currentPassword,
        newPassword
      });

      Swal.fire({
        icon: 'success',
        title: 'Password Changed',
        text: 'Your password has been updated successfully',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });

      // Reset form
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Change Failed',
        text: error.response?.data?.message || 'Failed to change password',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !loading) {
      handleSubmit();
    }
  }
</script>

<div class="change-password-page">
  <h1>Change Password</h1>

  <div class="password-form">
    <div class="form-group">
      <label for="current-password">Current Password</label>
      <input
        id="current-password"
        type="password"
        bind:value={currentPassword}
        placeholder="Enter current password"
        class="form-input"
        on:keydown={handleKeyDown}
        disabled={loading}
      />
    </div>
    <div class="form-group">
      <label for="new-password">New Password</label>
      <input
        id="new-password"
        type="password"
        bind:value={newPassword}
        placeholder="Enter new password (min. 6 characters)"
        class="form-input"
        on:keydown={handleKeyDown}
        disabled={loading}
      />
    </div>
    <div class="form-group">
      <label for="confirm-password">Confirm New Password</label>
      <input
        id="confirm-password"
        type="password"
        bind:value={confirmPassword}
        placeholder="Confirm new password"
        class="form-input"
        on:keydown={handleKeyDown}
        disabled={loading}
      />
    </div>
    <div class="button-group">
      <button class="btn btn-primary" on:click={handleSubmit} disabled={loading}>
        {loading ? 'Changing...' : 'Change Password'}
      </button>
    </div>
  </div>
</div>

<style>
  .change-password-page {
    max-width: 600px;
    margin: 0 auto;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
  }

  .password-form {
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

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>

