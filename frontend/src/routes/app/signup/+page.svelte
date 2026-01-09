<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import { requirePublicWithSleep } from '$lib/utils/authSleep';
  import { workspace } from '$lib/stores/workspace';
  import { get } from 'svelte/store';
  import WorkspaceModal from '$lib/components/WorkspaceModal.svelte';

  let name = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  let loading = false;
  let showWorkspaceModal = false;

  async function handleSignup() {
    if (loading) return; // Prevent double submission
    
    if (!name || !email || !password) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Missing Fields',
          text: 'Please fill in all fields',

        });
      }, 50);
      return;
    }

    if (password.length < 6) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: 'Password Too Short',
          text: 'Password must be at least 6 characters',

        });
      }, 50);
      return;
    }

    if (password !== confirmPassword) {
      setTimeout(() => {
        Swal.fire({
          icon: 'warning',
          title: `Passwords Don't Match`,
          text: 'Please make sure both passwords are the same',

        });
      }, 50);
      return;
    }

    loading = true;
    try {
      const response = await axios.post('/api/apps/auth/signup', {
        name,
        email,
        password
      });

      // Signup is always a fresh start - clear any existing accounts
      // clearExisting=true means start fresh
      auth.login(response.data.user, response.data.token, response.data.refreshToken, true);
      
      // Check workspace state BEFORE redirecting
      await workspace.init();
      const wsState = get(workspace);
      
      // After signup, check if first signin is not completed and has 0 workspace
      // Show workspace modal if needed
      if (!response.data.user.firstSigninCompleted && wsState.list.length === 0) {
        // No workspaces and first signin not completed - MUST create one
        // Redirect first, then show modal
        goto('/app/expenses');
        setTimeout(() => {
          showWorkspaceModal = true;
        }, 300);
      } else {
        // Has workspaces or first signin completed, just show success message
        goto('/app/expenses');
        setTimeout(() => {
          Swal.fire({
            icon: 'success',
            title: 'Account Created!',
            text: 'Welcome to Expenses Tracker',
            timer: 2000,
            showConfirmButton: false,
          });
        }, 50);
      }
    } catch (error) {
      setTimeout(() => {
        Swal.fire({
          icon: 'error',
          title: 'Signup Failed',
          text: error.response?.data?.message || 'Failed to create account',

        });
      }, 50);
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !loading) {
      e.preventDefault();
      setTimeout(() => {
        handleSignup();
      }, 50);
    }
  }
</script>

<div class="max-w-md mx-auto">
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body">
      <h1 class="card-title text-3xl justify-center mb-2">Sign Up</h1>
      <p class="text-center text-base-content/70 mb-6">Create your account</p>

      <form on:submit|preventDefault={handleSignup}>
        <fieldset class="fieldset mb-4">
          <legend class="fieldset-legend">Name</legend>
          <input
            id="name"
            type="text"
            bind:value={name}
            placeholder="Your name"
            class="input input-bordered w-full border-2"
            on:keydown={handleKeyDown}
            disabled={loading}
            required
          />
        </fieldset>

        <fieldset class="fieldset mb-4">
          <legend class="fieldset-legend">Email</legend>
          <input
            id="email"
            type="email"
            bind:value={email}
            placeholder="your@email.com"
            class="input input-bordered w-full border-2"
            on:keydown={handleKeyDown}
            disabled={loading}
            required
          />
        </fieldset>

        <fieldset class="fieldset mb-4">
          <legend class="fieldset-legend">Password</legend>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="At least 6 characters"
            class="input input-bordered w-full border-2"
            on:keydown={handleKeyDown}
            disabled={loading}
            required
            minlength="6"
          />
        </fieldset>

        <fieldset class="fieldset mb-6">
          <legend class="fieldset-legend">Confirm Password</legend>
          <input
            id="confirmPassword"
            type="password"
            bind:value={confirmPassword}
            placeholder="Confirm your password"
            class="input input-bordered w-full border-2"
            on:keydown={handleKeyDown}
            disabled={loading}
            required
          />
        </fieldset>

        <button type="submit" class="btn btn-primary w-full" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Creating account...
          {:else}
            Sign Up
          {/if}
        </button>
      </form>

      <div class="mt-6 text-center">
        <p class="text-sm text-base-content/70">
          Already have an account?
          <a href="/app/login" class="link link-primary font-medium">Login</a>
        </p>
      </div>
    </div>
  </div>
</div>

<WorkspaceModal
  bind:open={showWorkspaceModal}
  required={true}
  on:saved={() => {
    showWorkspaceModal = false;
    goto('/app/expenses');
  }}
/>