import { get } from 'svelte/store';
import { goto } from '$app/navigation';
import { auth } from '$lib/stores/auth';

// Small utility sleep, default 1ms as per user's preferred pattern
export function sleep(ms = 1) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Ensure the user is authenticated.
 * Uses a tiny sleep so that the auth store has time to hydrate on first load.
 *
 * @param {string} redirectTo - Where to send unauthenticated users.
 * @returns {Promise<boolean>} - true if authenticated, false if redirected.
 */
export async function requireAuthWithSleep(redirectTo = '/app/login') {
  await sleep(1);

  const state = get(auth);
  if (!state?.isAuthenticated) {
    goto(redirectTo);
    return false;
  }

  return true;
}

/**
 * Ensure the page is only accessible to unauthenticated users (e.g. login/signup).
 * Redirects authenticated users away (usually to the main app).
 *
 * @param {string} redirectTo - Where to send authenticated users.
 * @returns {Promise<boolean>} - true if still public, false if redirected.
 */
export async function requirePublicWithSleep(redirectTo = '/app/expenses') {
  await sleep(1);

  const state = get(auth);
  if (state?.isAuthenticated) {
    goto(redirectTo);
    return false;
  }

  return true;
}


