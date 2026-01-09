import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { auth } from './auth';
import api from '$lib/api';

// Currency symbols mapping
const CURRENCY_SYMBOLS = {
  IDR: 'Rp.',
  USD: '$',
  EUR: '€',
  JPY: '¥'
};

// Currency locales for number formatting
const CURRENCY_LOCALES = {
  IDR: 'id-ID',
  USD: 'en-US',
  EUR: 'de-DE',
  JPY: 'ja-JP'
};

function createCurrencyStore() {
  const { subscribe, set, update } = writable('IDR');

  return {
    subscribe,
    setCurrency: async (currency) => {
      if (!browser) return;
      
      try {
        await api.put('/auth/currency', { currency });
        set(currency);
        localStorage.setItem('currency', currency);
      } catch (error) {
        console.error('Failed to update currency:', error);
        throw error;
      }
    },
    init: async () => {
      if (!browser) return;
      
      // First, load from localStorage immediately (synchronous) to avoid showing default
      const savedCurrency = localStorage.getItem('currency');
      if (savedCurrency) {
        set(savedCurrency);
      }
      
      // Then, try to get currency from user profile and update if different
      try {
        const response = await api.get('/auth/profile');
        const currency = response.data.currency || 'IDR';
        set(currency);
        localStorage.setItem('currency', currency);
      } catch (error) {
        // If API fails and no localStorage, use default
        if (!savedCurrency) {
          set('IDR');
          localStorage.setItem('currency', 'IDR');
        }
      }
    },
    formatCurrency: (amount, showSymbols = true) => {
      let currency = 'IDR';
      subscribe(value => { currency = value; })();
      
      if (!amount && amount !== 0) return '';
      
      let numericValue;
      if (typeof amount === 'string') {
        numericValue = amount.replace(/\D/g, '');
        if (!numericValue) return '';
        numericValue = parseFloat(numericValue);
      } else {
        numericValue = amount;
      }
      
      if (isNaN(numericValue)) return '';
      
      let symbol = CURRENCY_SYMBOLS[currency] || 'Rp.';
      const locale = CURRENCY_LOCALES[currency] || 'id-ID';

      if (showSymbols){
          symbol = ""
      }
      return symbol + ' ' + new Intl.NumberFormat(locale, {
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
      }).format(numericValue);
    }
  };
}

export const currency = createCurrencyStore();

