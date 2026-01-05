/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
  // Let client-side handle routing since auth is stored in localStorage
  // Server-side hooks can't access localStorage, so we handle auth checks client-side
  return resolve(event);
}
