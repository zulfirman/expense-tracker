import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { accounts } from './accounts';

function createAuthStore() {
  const { subscribe, set, update } = writable({
    user: null,
    token: null,
    refreshToken: null,
    isAuthenticated: false
  });

  return {
    subscribe,
    login: (user, token, refreshToken, clearExisting = false) => {
      if (browser) {
        // Clear existing accounts if requested (for fresh login/signup)
        if (clearExisting) {
          accounts.clearAccounts();
        }
        
        // Add to accounts store
        accounts.addAccount(user, token, refreshToken);
        
        localStorage.setItem('auth_token', token);
        localStorage.setItem('auth_refresh_token', refreshToken);
        localStorage.setItem('auth_user', JSON.stringify(user));
      }
      set({ user, token, refreshToken, isAuthenticated: true });
    },
    logout: () => {
      if (browser) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_refresh_token');
        localStorage.removeItem('auth_user');
        // Clear accounts when logging out (user wants fresh start)
        accounts.clearAccounts();
      }
      set({ user: null, token: null, refreshToken: null, isAuthenticated: false });
    },
    init: () => {
      if (browser) {
        // Initialize accounts first
        accounts.init();
        
        const token = localStorage.getItem('auth_token');
        const refreshToken = localStorage.getItem('auth_refresh_token');
        const userStr = localStorage.getItem('auth_user');
        if (token && refreshToken && userStr) {
          try {
            const user = JSON.parse(userStr);
            set({ user, token, refreshToken, isAuthenticated: true });
          } catch (e) {
            console.error('Failed to parse user data:', e);
            localStorage.removeItem('auth_token');
            localStorage.removeItem('auth_refresh_token');
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
    getRefreshToken: () => {
      if (browser) {
        return localStorage.getItem('auth_refresh_token');
      }
      return null;
    },
    setToken: (token, refreshToken) => {
      if (browser) {
        if (token) localStorage.setItem('auth_token', token);
        if (refreshToken) localStorage.setItem('auth_refresh_token', refreshToken);
        update(state => ({ ...state, token, refreshToken }));
      }
    },
    switchAccount: (accountId) => {
      accounts.switchAccount(accountId);
      const account = accounts.getCurrentAccount();
      if (account) {
        const user = { id: account.id, name: account.name, email: account.email };
        set({ user, token: account.token, refreshToken: account.refreshToken || null, isAuthenticated: true });
      }
    },
    refreshAccessToken: async function() {
      if (!browser) return false;
      
      const refreshToken = localStorage.getItem('auth_refresh_token');
      if (!refreshToken) {
        return false;
      }

      try {
        const response = await fetch('/api/auth/refresh', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ refreshToken })
        });

        if (!response.ok) {
          throw new Error('Failed to refresh token');
        }

        const data = await response.json();
        
        // Update tokens in store and localStorage
        if (browser) {
          localStorage.setItem('auth_token', data.token);
          localStorage.setItem('auth_refresh_token', data.refreshToken);
        }
        update(state => ({ ...state, token: data.token, refreshToken: data.refreshToken }));
        
        // Update accounts store
        const account = accounts.getCurrentAccount();
        if (account) {
          accounts.addAccount(
            { id: account.id, name: account.name, email: account.email },
            data.token,
            data.refreshToken
          );
        }
        
        return true;
      } catch (error) {
        console.error('Failed to refresh token:', error);
        // Logout on refresh failure
        if (browser) {
          localStorage.removeItem('auth_token');
          localStorage.removeItem('auth_refresh_token');
          localStorage.removeItem('auth_user');
          accounts.clearAccounts();
        }
        set({ user: null, token: null, refreshToken: null, isAuthenticated: false });
        return false;
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

