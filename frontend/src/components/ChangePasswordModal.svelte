<script>
  import { createEventDispatcher } from 'svelte';
  import api from '$lib/api';
  import Swal from 'sweetalert2';

  const dispatch = createEventDispatcher();

  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  let loading = false;

  function close() {
    dispatch('close');
    // Reset form
    currentPassword = '';
    newPassword = '';
    confirmPassword = '';
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      close();
    }
  }

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
      Swal.fire({
        icon: 'warning',
        title: 'Invalid Password',
        text: 'Password must be at least 6 characters long',
        zIndex: 9999
      });
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

      close();
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

<div class="modal-backdrop" on:click={handleBackdropClick}>
  <div class="modal-content" on:click|stopPropagation>
    <div class="modal-header">
      <h2>Change Password</h2>
      <button class="close-btn" on:click={close}>×</button>
    </div>
    <div class="modal-body">
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Current Password</legend>
        <input
          id="current-password"
          type="password"
          bind:value={currentPassword}
          placeholder="Enter current password"
          class="input input-bordered w-full border-2"
          on:keydown={handleKeyDown}
          disabled={loading}
        />
      </fieldset>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">New Password</legend>
        <input
          id="new-password"
          type="password"
          bind:value={newPassword}
          placeholder="Enter new password (min. 6 characters)"
          class="input input-bordered w-full border-2"
          on:keydown={handleKeyDown}
          disabled={loading}
        />
      </fieldset>
      <fieldset class="fieldset mb-3">
        <legend class="fieldset-legend">Confirm New Password</legend>
        <input
          id="confirm-password"
          type="password"
          bind:value={confirmPassword}
          placeholder="Confirm new password"
          class="input input-bordered w-full border-2"
          on:keydown={handleKeyDown}
          disabled={loading}
        />
      </fieldset>
      <div class="button-group">
        <button class="btn btn-secondary" on:click={close} disabled={loading}>
          Cancel
        </button>
        <button class="btn btn-primary" on:click={handleSubmit} disabled={loading}>
          {loading ? 'Changing...' : 'Change Password'}
        </button>
      </div>
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
</style>



