<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import '$lib/styles/shared.css';

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
  <div class="header-row">
    <h1>Change Password</h1>
    <button class="btn ghost" on:click={() => goto('/preferences')}>Back</button>
  </div>

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

  .header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
  }

  .password-form {
    background: var(--surface);
    border-radius: 0.75rem;
    padding: 1.5rem;
    border: 1px solid var(--border);
  }

  .button-group {
    margin-top: 2rem;
  }
</style>



