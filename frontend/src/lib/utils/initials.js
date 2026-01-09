/**
 * Get initials from a name
 * @param {string} name - Full name
 * @returns {string} - Initials (e.g., "ZF" for "Zul Firman")
 */
export function getInitials(name) {
  if (!name || typeof name !== 'string') return 'U';
  
  const parts = name.trim().split(/\s+/);
  if (parts.length === 0) return 'U';
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  
  // Get first letter of first name and first letter of last name
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

