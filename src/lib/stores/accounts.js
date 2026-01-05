import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createAccountsStore() {
  const { subscribe, set, update } = writable({
    accounts: [],
    currentAccountId: null
  });

  return {
    subscribe,
    addAccount: (user, token) => {
      if (!browser) return;
      
      update(state => {
        const account = {
          id: user.id,
          name: user.name,
          email: user.email,
          token: token,
          addedAt: new Date().toISOString()
        };

        // Check if account already exists
        const existingIndex = state.accounts.findIndex(acc => acc.id === user.id);
        if (existingIndex >= 0) {
          // Update existing account
          state.accounts[existingIndex] = account;
        } else {
          // Add new account
          state.accounts.push(account);
        }

        // Set as current account
        state.currentAccountId = user.id;

        // Save to localStorage
        localStorage.setItem('accounts', JSON.stringify(state.accounts));
        localStorage.setItem('currentAccountId', user.id.toString());
        localStorage.setItem('auth_token', token);
        localStorage.setItem('auth_user', JSON.stringify(user));

        return state;
      });
    },
    removeAccount: (accountId) => {
      if (!browser) return;
      
      update(state => {
        state.accounts = state.accounts.filter(acc => acc.id !== accountId);
        
        // If removed account was current, switch to first available or null
        if (state.currentAccountId === accountId) {
          state.currentAccountId = state.accounts.length > 0 ? state.accounts[0].id : null;
          if (state.currentAccountId) {
            const currentAccount = state.accounts.find(acc => acc.id === state.currentAccountId);
            localStorage.setItem('auth_token', currentAccount.token);
            localStorage.setItem('auth_user', JSON.stringify({ id: currentAccount.id, name: currentAccount.name, email: currentAccount.email }));
          } else {
            localStorage.removeItem('auth_token');
            localStorage.removeItem('auth_user');
          }
        }

        localStorage.setItem('accounts', JSON.stringify(state.accounts));
        localStorage.setItem('currentAccountId', state.currentAccountId ? state.currentAccountId.toString() : '');

        return state;
      });
    },
    switchAccount: (accountId) => {
      if (!browser) return;
      
      update(state => {
        const account = state.accounts.find(acc => acc.id === accountId);
        if (account) {
          state.currentAccountId = accountId;
          localStorage.setItem('currentAccountId', accountId.toString());
          localStorage.setItem('auth_token', account.token);
          localStorage.setItem('auth_user', JSON.stringify({ id: account.id, name: account.name, email: account.email }));
        }
        return state;
      });
    },
    init: () => {
      if (browser) {
        const accountsStr = localStorage.getItem('accounts');
        const currentAccountIdStr = localStorage.getItem('currentAccountId');
        
        if (accountsStr) {
          try {
            const accounts = JSON.parse(accountsStr);
            const currentAccountId = currentAccountIdStr ? parseInt(currentAccountIdStr) : (accounts.length > 0 ? accounts[0].id : null);
            
            set({
              accounts: accounts,
              currentAccountId: currentAccountId
            });

            // Set current account token if exists
            if (currentAccountId) {
              const currentAccount = accounts.find(acc => acc.id === currentAccountId);
              if (currentAccount) {
                localStorage.setItem('auth_token', currentAccount.token);
                localStorage.setItem('auth_user', JSON.stringify({ id: currentAccount.id, name: currentAccount.name, email: currentAccount.email }));
              }
            }
          } catch (e) {
            console.error('Failed to parse accounts:', e);
          }
        }
      }
    },
    getCurrentAccount: () => {
      let currentAccount = null;
      subscribe(state => {
        if (state.currentAccountId) {
          currentAccount = state.accounts.find(acc => acc.id === state.currentAccountId);
        }
      })();
      return currentAccount;
    },
    clearAccounts: () => {
      if (!browser) return;
      
      set({
        accounts: [],
        currentAccountId: null
      });
      
      localStorage.removeItem('accounts');
      localStorage.removeItem('currentAccountId');
    }
  };
}

export const accounts = createAccountsStore();


