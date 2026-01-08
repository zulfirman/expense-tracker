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
      const response = await axios.post('/api/apps/auth/signup', {
        name,
        email,
        password
      });

      // Signup is always a fresh start - clear any existing accounts
      // clearExisting=true means start fresh
      auth.login(response.data.user, response.data.token, response.data.refreshToken, true);
      
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
          <a href="/login" class="link link-primary font-medium">Login</a>
        </p>
      </div>
    </div>
  </div>
</div>