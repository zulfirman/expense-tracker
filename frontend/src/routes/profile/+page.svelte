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
      reverseButtons: true,
      zIndex: 9999
    }).then((result) => {
      if (result.isConfirmed) {
        auth.logout();
        goto('/login');
      }
    });
  }
</script>

<div class="max-w-3xl mx-auto space-y-4">
  <div class="card bg-base-100 shadow-xl border border-base-300">
    <div class="card-body">
      <div class="flex items-center justify-between gap-2">
        <div>
          <p class="text-xs uppercase tracking-wide text-base-content/60">Account</p>
          <h1 class="text-2xl font-bold">Profile</h1>
        </div>
        <button class="btn btn-soft btn-sm" on:click={() => goto('/preferences')}>Preferences</button>
      </div>

      {#if loading}
        <div class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else}
        <div class="grid gap-4 md:grid-cols-2">
          <fieldset class="fieldset">
            <legend class="fieldset-legend">Name</legend>
            <input
              id="name"
              type="text"
              bind:value={name}
              placeholder="Your name"
              class="input input-bordered w-full border-2"
              disabled={saving}
            />
          </fieldset>

          <fieldset class="fieldset">
            <legend class="fieldset-legend">Email</legend>
            <input
              id="email"
              type="email"
              bind:value={email}
              placeholder="your@email.com"
              class="input input-bordered w-full border-2"
              disabled={saving}
            />
          </fieldset>
        </div>

        <div class="flex flex-col md:flex-row gap-3 mt-4">
          <button class="btn btn-primary w-full" on:click={handleSave} disabled={saving}>
            {#if saving}
              <span class="loading loading-spinner loading-sm mr-1"></span>
              Saving...
            {:else}
              Save Changes
            {/if}
          </button>
        </div>
        <div class="flex flex-col md:flex-row gap-3 mt-4">
          <button class="btn btn-error w-full" on:click={handleLogout} disabled={saving}>
            Logout
          </button>
          <div class="badge px-3 py-2 badge-accent">Accent</div>
        </div>
      {/if}
    </div>
  </div>
</div>



