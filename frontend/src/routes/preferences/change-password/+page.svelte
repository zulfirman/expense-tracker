<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import api from '$lib/api';
  import Swal from 'sweetalert2';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';

  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  let loading = false;
  
  $: pageCode = getPageCode($page.url.pathname);

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

<div class="max-w-md mx-auto space-y-4">
  <PageHeader
    title="Change Password"
    subtitle="Update the password for your account."
    pageCode={pageCode}
    actions={true}
  >
    <svelte:fragment slot="actions">
      <button class="btn btn-soft btn-sm" on:click={() => goto('/preferences')}>Back</button>
    </svelte:fragment>
  </PageHeader>

  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body space-y-4">
      <fieldset class="fieldset">
        <legend class="fieldset-legend">Current password</legend>
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

      <fieldset class="fieldset">
        <legend class="fieldset-legend">New password</legend>
        <input
          id="new-password"
          type="password"
          bind:value={newPassword}
          placeholder="At least 6 characters"
          class="input input-bordered w-full border-2"
          on:keydown={handleKeyDown}
          disabled={loading}
        />
      </fieldset>

      <fieldset class="fieldset">
        <legend class="fieldset-legend">Confirm new password</legend>
        <input
          id="confirm-password"
          type="password"
          bind:value={confirmPassword}
          placeholder="Re-enter new password"
          class="input input-bordered w-full border-2"
          on:keydown={handleKeyDown}
          disabled={loading}
        />
      </fieldset>

      <div class="pt-2">
        <button class="btn btn-primary w-full" on:click={handleSubmit} disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm mr-1"></span>
            Changing...
          {:else}
            Change password
          {/if}
        </button>
      </div>
    </div>
  </div>
</div>

