export function getPageCode(pathname) {
  if (pathname === '/' || pathname === '/expenses') return 'PG-EXPENSES';
  if (pathname === '/income') return 'PG-INCOME';
  if (pathname === '/history') return 'PG-HISTORY';
  if (pathname === '/budget') return 'PG-BUDGET';
  if (pathname === '/app/login') return 'PG-LOGIN';
  if (pathname === '/app/signup') return 'PG-SIGNUP';
  if (pathname === '/profile') return 'PG-PROFILE';
  if (pathname === '/change-password') return 'PG-CHANGE-PW';
  if (pathname === '/preferences') return 'PG-PREF';
  if (pathname.startsWith('/preferences/categories')) return 'PG-PREF-CAT';
  if (pathname.startsWith('/preferences/currency')) return 'PG-PREF-CURR';
  if (pathname.startsWith('/preferences/quick-amounts')) return 'PG-PREF-QA';
  if (pathname.startsWith('/preferences/change-password')) return 'PG-PREF-PW';
  if (pathname.startsWith('/categories')) return 'PG-CATEGORIES';
  return `PG-${pathname.replace(/\//g, '-').toUpperCase() || 'ROOT'}`;
}



