<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import axios from 'axios';
  import Swal from 'sweetalert2';
  import { auth } from '$lib/stores/auth';
  import { accounts } from '$lib/stores/accounts';

  let name = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  let loading = false;

  onMount(() => {
    // Redirect if already logged in
    if ($auth.isAuthenticated) {
      goto('/expenses');
    }
  });

  async function handleSignup() {
    if (!name || !email || !password) {
      Swal.fire({
        icon: 'warning',
        title: 'Missing Fields',
        text: 'Please fill in all fields',
        zIndex: 9999
      });
      return;
    }

    if (password.length < 6) {
      Swal.fire({
        icon: 'warning',
        title: 'Password Too Short',
        text: 'Password must be at least 6 characters',
        zIndex: 9999
      });
      return;
    }

    if (password !== confirmPassword) {
      Swal.fire({
        icon: 'warning',
        title: `Passwords Don't Match`,
        text: 'Please make sure both passwords are the same',
        zIndex: 9999
      });
      return;
    }

    loading = true;
    try {
      const response = await axios.post('/api/auth/signup', {
        name,
        email,
        password
      });

      // Signup is always a fresh start - clear any existing accounts
      // clearExisting=true means start fresh
      auth.login(response.data.user, response.data.token, true);
      
      Swal.fire({
        icon: 'success',
        title: 'Account Created!',
        text: 'Welcome to Expenses Tracker',
        timer: 2000,
        showConfirmButton: false,
        zIndex: 9999
      });

      goto('/expenses');
    } catch (error) {
      Swal.fire({
        icon: 'error',
        title: 'Signup Failed',
        text: error.response?.data?.message || 'Failed to create account',
        zIndex: 9999
      });
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter') {
      handleSignup();
    }
  }
</script>

<div class="auth-page">
  <div class="auth-container">
    <h1>Sign Up</h1>
    <p class="subtitle">Create your account</p>

    <form on:submit|preventDefault={handleSignup}>
      <div class="form-group">
        <label for="name">Name</label>
        <input
          id="name"
          type="text"
          bind:value={name}
          placeholder="Your name"
          class="form-input"
          on:keydown={handleKeyDown}
          disabled={loading}
          required
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
          placeholder="At least 6 characters"
          class="form-input"
          on:keydown={handleKeyDown}
          disabled={loading}
          required
          minlength="6"
        />
      </div>

      <div class="form-group">
        <label for="confirmPassword">Confirm Password</label>
        <input
          id="confirmPassword"
          type="password"
          bind:value={confirmPassword}
          placeholder="Confirm your password"
          class="form-input"
          on:keydown={handleKeyDown}
          disabled={loading}
          required
        />
      </div>

      <button type="submit" class="btn btn-primary" disabled={loading}>
        {#if loading}
          <span class="spinner"></span> Creating account...
        {:else}
          Sign Up
        {/if}
      </button>
    </form>

    <div class="auth-footer">
      <p>Already have an account? <a href="/login">Login</a></p>
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

