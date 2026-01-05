import { writable } from 'svelte/store';

function createThemeStore() {
  const { subscribe, set, update } = writable('light');

  return {
    subscribe,
    setTheme: (theme) => {
      set(theme);
      if (typeof window !== 'undefined') {
        localStorage.setItem('theme', theme);
        applyTheme(theme);
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
        const theme = prefersDark ? 'dark' : 'light';
        set(theme);
        applyTheme(theme);
        
        // Listen for system theme changes
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
          if (!localStorage.getItem('theme')) {
            const newTheme = e.matches ? 'dark' : 'light';
            set(newTheme);
            applyTheme(newTheme);
          }
        });
      }
    },
    toggle: () => {
      update(theme => {
        const newTheme = theme === 'dark' ? 'light' : 'dark';
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
  if (theme === 'dark') {
    root.classList.add('dark');
  } else {
    root.classList.remove('dark');
  }
}

export const theme = createThemeStore();





