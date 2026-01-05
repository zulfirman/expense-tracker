import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { accounts } from './accounts';

function createAuthStore() {
  const { subscribe, set, update } = writable({
    user: null,
    token: null,
    isAuthenticated: false
  });

  return {
    subscribe,
    login: (user, token, clearExisting = false) => {
      if (browser) {
        // Clear existing accounts if requested (for fresh login/signup)
        if (clearExisting) {
          accounts.clearAccounts();
        }
        
        // Add to accounts store
        accounts.addAccount(user, token);
        
        localStorage.setItem('auth_token', token);
        localStorage.setItem('auth_user', JSON.stringify(user));
      }
      set({ user, token, isAuthenticated: true });
    },
    logout: () => {
      if (browser) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
        // Clear accounts when logging out (user wants fresh start)
        accounts.clearAccounts();
      }
      set({ user: null, token: null, isAuthenticated: false });
    },
    init: () => {
      if (browser) {
        // Initialize accounts first
        accounts.init();
        
        const token = localStorage.getItem('auth_token');
        const userStr = localStorage.getItem('auth_user');
        if (token && userStr) {
          try {
            const user = JSON.parse(userStr);
            set({ user, token, isAuthenticated: true });
          } catch (e) {
            console.error('Failed to parse user data:', e);
            localStorage.removeItem('auth_token');
            localStorage.removeItem('auth_user');
          }
        }
      }
    },
    getToken: () => {
      if (browser) {
        return localStorage.getItem('auth_token');
      }
      return null;
    },
    switchAccount: (accountId) => {
      accounts.switchAccount(accountId);
      const account = accounts.getCurrentAccount();
      if (account) {
        const user = { id: account.id, name: account.name, email: account.email };
        set({ user, token: account.token, isAuthenticated: true });
      }
    }
  };
}

export const auth = createAuthStore();

// Derived store for current account info
export const currentAccount = derived(
  [auth, accounts],
  ([$auth, $accounts]) => {
    if ($accounts.currentAccountId) {
      return $accounts.accounts.find(acc => acc.id === $accounts.currentAccountId);
    }
    return null;
  }
);

