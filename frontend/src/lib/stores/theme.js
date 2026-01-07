import { writable } from 'svelte/store';

function createThemeStore() {
  const { subscribe, set, update } = writable('cupcake');

  return {
    subscribe,
    setTheme: (newTheme) => {
      set(newTheme);
      if (typeof window !== 'undefined') {
        localStorage.setItem('theme', newTheme);
        applyTheme(newTheme);
      }
    },
    init: () => {
      if (typeof window !== 'undefined') {
        // Check localStorage first
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
          set(savedTheme);
          applyTheme(savedTheme);
          return;
        }
        
        // Check system preference
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const initialTheme = prefersDark ? 'night' : 'cupcake';
        set(initialTheme);
        applyTheme(initialTheme);
        
        // Listen for system theme changes when user hasn't chosen explicitly
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
          if (!localStorage.getItem('theme')) {
            const newTheme = e.matches ? 'night' : 'cupcake';
            set(newTheme);
            applyTheme(newTheme);
          }
        });
      }
    },
    toggle: () => {
      update(current => {
        // Simple toggle between cupcake (light) and night (dark)
        const newTheme = current === 'night' ? 'cupcake' : 'night';
        if (typeof window !== 'undefined') {
          localStorage.setItem('theme', newTheme);
          applyTheme(newTheme);
        }
        return newTheme;
      });
    }
  };
}

function applyTheme(theme) {
  if (typeof document === 'undefined') return;
  
  const root = document.documentElement;
  // DaisyUI uses data-theme attribute
  root.setAttribute('data-theme', theme);
}

export const theme = createThemeStore();







