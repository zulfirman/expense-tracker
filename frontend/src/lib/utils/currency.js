import { get } from 'svelte/store';
import { currency } from '../stores/currency';

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

/**
 * Format currency value based on user's currency preference
 * This function can be used in both Svelte components and regular JS files
 * @param {number|string} amount - Value to format
 * @param {string} [forceCurrency] - Optional currency code to force (for reactive contexts, use the store directly)
 * @returns {string} Formatted currency string
 */
export function formatCurrency(amount, forceCurrency = null) {
  if (!amount && amount !== 0) return '';
  
  // Get current currency from store or use forced currency
  let currentCurrency = forceCurrency;
  if (!currentCurrency) {
    try {
      currentCurrency = get(currency) || 'IDR';
    } catch (e) {
      // If store is not available (SSR), default to IDR
      currentCurrency = 'IDR';
    }
  }
  
  // Convert string to number if needed
  let numericValue;
  if (typeof amount === 'string') {
    numericValue = amount.replace(/\D/g, '');
    if (!numericValue) return '';
    numericValue = parseFloat(numericValue);
  } else {
    numericValue = amount;
  }
  
  if (isNaN(numericValue)) return '';
  
  const symbol = CURRENCY_SYMBOLS[currentCurrency] || 'Rp.';
  const locale = CURRENCY_LOCALES[currentCurrency] || 'id-ID';
  
  return symbol + ' ' + new Intl.NumberFormat(locale, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(numericValue);
}

