<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';

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
      const response = await axios.post('/api/auth/login', {
        email,
        password
      });

      // Check if this is a fresh login (no existing accounts)
      // If no existing accounts, treat as fresh login and clear any stale data
      const hasExistingAccounts = localStorage.getItem('accounts');
      const shouldClear = !hasExistingAccounts || JSON.parse(hasExistingAccounts || '[]').length === 0;

      // Add account to accounts store (handles multiple accounts)
      // clearExisting=true means start fresh, false means add to existing
      auth.login(response.data.user, response.data.token, shouldClear);
      
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

<div class="auth-page">
  <div class="auth-container">
    <h1>Login</h1>
    <p class="subtitle">Sign in to your account</p>

    <form on:submit|preventDefault={handleLogin}>
      <div class="form-group">
        <label for="email">Email</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          placeholder="your@email.com"
          class="form-input"
          on:keydown={handleKeyDown}
          disabled={loading}
          required
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          placeholder="Enter your password"
          class="form-input"
          on:keydown={handleKeyDown}
          disabled={loading}
          required
        />
      </div>

      <button type="submit" class="btn btn-primary" disabled={loading}>
        {#if loading}
          <span class="spinner"></span> Logging in...
        {:else}
          Login
        {/if}
      </button>
    </form>

    <div class="auth-footer">
      <p>Don't have an account? <a href="/signup">Sign up</a></p>
      {#if $auth.isAuthenticated}
        <p style="margin-top: 1rem;">
          <a href="/expenses" style="color: var(--primary-color); text-decoration: none; font-weight: 500;">
            Continue as {auth.user?.name || 'Current User'}
          </a>
        </p>
      {/if}
    </div>
  </div>
</div>

<style>
  .auth-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    background: var(--background);
  }

  .auth-container {
    width: 100%;
    max-width: 400px;
    background: var(--surface);
    border-radius: 1rem;
    padding: 2rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  h1 {
    font-size: 2rem;
    margin-bottom: 0.5rem;
    color: var(--text-primary);
    text-align: center;
  }

  .subtitle {
    text-align: center;
    color: var(--text-secondary);
    margin-bottom: 2rem;
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

  .btn {
    width: 100%;
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

  .btn-primary:disabled {
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

  .auth-footer {
    margin-top: 1.5rem;
    text-align: center;
  }

  .auth-footer p {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }

  .auth-footer a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
  }

  .auth-footer a:hover {
    text-decoration: underline;
  }
</style>

