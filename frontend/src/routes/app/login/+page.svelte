<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import InputExpenses from "$components/InputExpenses.svelte";
  import { requirePublicWithSleep } from '$lib/utils/authSleep';
  import { workspace } from '$lib/stores/workspace';
  import { get } from 'svelte/store';
  import WorkspaceModal from '$lib/components/WorkspaceModal.svelte';

  let email = '';
  let password = '';
  let loading = false;
  let showWorkspaceModal = false;

  async function handleLogin() {
      if (loading) return; // Prevent double submission

      if (!email || !password) {
          setTimeout(() => {
              Swal.fire({
                  icon: 'warning',
                  title: 'Missing Fields',
                  text: 'Please enter both email and password',

              });
          }, 50);
          return;
      }

      loading = true;
      try {
          const response = await axios.post('/api/apps/auth/login', {
              email,
              password
          });

          // Check if this is a fresh login (no existing accounts)
          // If no existing accounts, treat as fresh login and clear any stale data
          const hasExistingAccounts = localStorage.getItem('accounts');
          const shouldClear = !hasExistingAccounts || JSON.parse(hasExistingAccounts || '[]').length === 0;

          // Add account to accounts store (handles multiple accounts)
          // clearExisting=true means start fresh, false means add to existing
          auth.login(response.data.user, response.data.token, response.data.refreshToken, shouldClear);

           // Check workspace state BEFORE redirecting
           await workspace.init();
           const wsState = get(workspace);
           
           // Check if first signin is not completed and has 0 workspace
           // Show workspace modal if needed (check every time to prevent empty workspace)
           if (!response.data.user.firstSigninCompleted && wsState.list.length === 0) {
               // No workspaces and first signin not completed - MUST create one
               // Redirect first, then show modal
               goto('/app/expenses');
               setTimeout(() => {
                   showWorkspaceModal = true;
               }, 300);
           } else if (wsState.list.length === 0) {
               // Edge case: first signin completed but no workspaces - still show modal
               goto('/app/expenses');
               setTimeout(() => {
                   showWorkspaceModal = true;
               }, 300);
           } else {
               // Has workspaces, proceed normally
               goto('/app/expenses');
               setTimeout(() => {
                   Swal.fire({
                       icon: 'success',
                       title: 'Welcome Back!',
                       timer: 1500,
                       showConfirmButton: false,
                   });
               }, 50);
           }
      } catch (error) {
          setTimeout(() => {
              Swal.fire({
                  icon: 'error',
                  title: 'Login Failed',
                  text: error.response?.data?.message || 'Invalid email or password',

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
        handleLogin();
      }, 50);
    }
  }
</script>

<div class="max-w-md mx-auto">
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body">
      <h1 class="card-title text-3xl justify-center mb-2">Login</h1>
      <p class="text-center text-base-content/70 mb-6">Sign in to your account</p>

      <form on:submit|preventDefault={handleLogin}>
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

        <fieldset class="fieldset mb-6">
          <legend class="fieldset-legend">Password</legend>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Enter your password"
            class="input input-bordered w-full border-2"
            on:keydown={handleKeyDown}
            disabled={loading}
            required
          />
        </fieldset>

        <button type="submit" class="btn btn-primary w-full" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Logging in...
          {:else}
            Login
          {/if}
        </button>
      </form>

      <div class="mt-6 text-center">
        <p class="text-sm text-base-content/70">
          Don't have an account?
          <a href="/app/signup" class="link link-primary font-medium">Sign up</a>
        </p>
        {#if $auth.isAuthenticated}
          <p class="mt-4">
            <a href="/app/expenses" class="link link-primary font-medium">
              Continue as {auth.user?.name || 'Current User'}
            </a>
          </p>
        {/if}
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
