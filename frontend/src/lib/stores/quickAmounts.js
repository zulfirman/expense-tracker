import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';
import api from '$lib/api';
import { auth } from './auth';
import { currency as currencyStore } from './currency';

/**
 * Currency-aware default quick amounts
 */
const DEFAULTS = {
  IDR: [10000, 25000, 50000, 100000, 200000],
  USD: [5, 10, 20, 50, 100],
  EUR: [5, 10, 20, 50, 100],
  JPY: [500, 1000, 2000, 5000, 10000]
};

function getDefaultsByCurrency(code) {
  return DEFAULTS[code] || DEFAULTS.IDR;
}

function storageKey(userId) {
  return `quickAmounts_${userId || 'guest'}`;
}

function createQuickAmountsStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    /**
     * Initialize quick amounts.
     * Tries API first; falls back to localStorage and currency defaults.
     */
    init: async () => {
      if (!browser) return;

      const userId = get(auth)?.user?.id;
      const currentCurrency = get(currencyStore) || 'IDR';

      // 1) Try API
      try {
        const resp = await api.get('/quick-amounts');
        if (Array.isArray(resp.data) && resp.data.length > 0) {
          set(resp.data);
          localStorage.setItem(storageKey(userId), JSON.stringify(resp.data));
          return;
        }
      } catch (err) {
        // ignore and fallback
      }

      // 2) LocalStorage
      try {
        const saved = localStorage.getItem(storageKey(userId));
        if (saved) {
          const parsed = JSON.parse(saved);
          if (Array.isArray(parsed) && parsed.length > 0) {
            set(parsed);
            return;
          }
        }
      } catch (err) {
        // ignore
      }

      // 3) Currency defaults
      const defaults = getDefaultsByCurrency(currentCurrency);
      set(defaults);
      localStorage.setItem(storageKey(userId), JSON.stringify(defaults));
    },

    /**
     * Save quick amounts (CRUD single payload).
     * Tries API; falls back to localStorage.
     */
    setAmounts: async (amounts) => {
      if (!browser) return;

      const sanitized = (amounts || [])
        .map((v) => Number(v))
        .filter((v) => !isNaN(v) && v > 0);

      set(sanitized);

      const userId = get(auth)?.user?.id;
      localStorage.setItem(storageKey(userId), JSON.stringify(sanitized));

      try {
        await api.put('/quick-amounts', { amounts: sanitized });
      } catch (err) {
        // silent fallback to localStorage
        console.warn('Quick amounts saved locally (API unavailable)');
      }
    },

    getDefaultsByCurrency,
  };
}

export const quickAmounts = createQuickAmountsStore();


