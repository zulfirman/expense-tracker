<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';
  import InputExpenses from "../../components/InputExpenses.svelte";

  let email = '';
  let password = '';
  let loading = false;

  onMount(() => {
    // Redirect if already logged in
    if ($auth.isAuthenticated) {
      goto('/expenses');
    }
  });

  async function handleLogin() {
    if (!email || !password) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Fields',
        text: 'Please enter both email and password',
        zIndex: 9999
      });
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
      
      Swal.fire({
        icon: 'success',
        title: 'Welcome Back!',
        timer: 1500,
        showConfirmButton: false,
        zIndex: 9999
      });

      goto('/expenses');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Login Failed',
        text: error.response?.data?.message || 'Invalid email or password',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter') {
      handleLogin();
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
          <a href="/signup" class="link link-primary font-medium">Sign up</a>
        </p>
        {#if $auth.isAuthenticated}
          <p class="mt-4">
            <a href="/expenses" class="link link-primary font-medium">
              Continue as {auth.user?.name || 'Current User'}
            </a>
          </p>
        {/if}
      </div>
    </div>
  </div>
</div>
