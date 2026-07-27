/**
 * Auth Utility Module
 * Manages token validation, session checks, and forced login redirection.
 */

export function isTokenValid() {
  const token = sessionStorage.getItem('token')
  if (!token) return false

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return false

    const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))
    if (payload.exp && payload.exp * 1000 < Date.now()) {
      return false // Token expired
    }
    return true
  } catch (e) {
    return false // Malformed token
  }
}

export function clearSessionAndForceLogin() {
  sessionStorage.clear()
  localStorage.clear()
  if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}
